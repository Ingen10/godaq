// Copyright 2016 The Godaq Authors. All rights reserved
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package godaq

import (
	"encoding/binary"
	"errors"
	"io"
	"sync"
	"time"

	"github.com/tarm/serial"
	try "gopkg.in/matryer/try.v1"
)

type Color uint8

const (
	OFF Color = iota
	GREEN
	RED
	YELLOW
)

var (
	ErrUnknownModel  = errors.New("Unknown device model number")
	ErrInvalidLed    = errors.New("Invalid LED number")
	ErrInvalidInput  = errors.New("Invalid input number")
	ErrInvalidOutput = errors.New("Invalid output number")
	ErrInvalidPIO    = errors.New("Invalid PIO number")
	ErrInvalidGainID = errors.New("Invalid gain ID")
)

type Calib struct {
	Gain   float32 // Gain calibration (-1 to 1)
	Offset float32 // Offset calibraton in ADUs
}

type HwFeatures struct {
	Name                              string
	NPIOs, NLeds                      uint
	NInputs, NOutputs, NHiddenOutputs uint
	NCalibRegs                        uint
	Dac                               DAC
	Adc                               ADC
}

type HwModel interface {
	GetFeatures() HwFeatures
	GetCalibIndex(isOutput, diffMode, secondStage bool, n, gainId uint) (uint, error)
	CheckValidInputs(pos, neg uint) error
}

var hwModels = make(map[uint8]HwModel)

func registerModel(model uint8, hw HwModel) error {
	if _, exists := hwModels[model]; exists {
		return errors.New("Hardware model already registered!")
	}
	hwModels[model] = hw
	return nil
}

func boolToByte(val bool) byte {
	if val {
		return 1
	}
	return 0
}

type OpenDAQ struct {
	ser *serial.Port
	HwFeatures
	hw    HwModel
	calib []Calib
	sync.Mutex

	// Input state (needed for converting ADC values to volts)
	gainId   uint
	posInput uint
	diffMode bool
}

func New(port string) (*OpenDAQ, error) {
	var err error
	daq := OpenDAQ{}
	daq.posInput = 1 // 0 is not a valid default for posInput

	// Setup and open the serial port
	serCfg := &serial.Config{Name: port, Baud: 115200, ReadTimeout: time.Millisecond * 100}
	daq.ser, err = serial.OpenPort(serCfg)
	if err != nil {
		return nil, err
	}
	time.Sleep(1500 * time.Millisecond)

	// Obtain the device model number
	model, _, _, err := daq.GetInfo()
	if err != nil {
		return nil, err
	}
	hw, ok := hwModels[model]
	if !ok {
		return nil, ErrUnknownModel
	}
	daq.hw = hw
	daq.HwFeatures = hw.GetFeatures()

	// Read the calibration registers from the device
	daq.calib = make([]Calib, daq.NCalibRegs)
	for i := range daq.calib {
		if daq.calib[i], err = daq.readCalib(uint8(i)); err != nil {
			return nil, err
		}
	}
	return &daq, nil
}

func (daq *OpenDAQ) Close() error {
	return daq.ser.Close()
}

// Send a comand and returns its response
func (daq *OpenDAQ) sendCommand(command *Message, respLen int) (r io.Reader, err error) {
	daq.Lock()
	defer daq.Unlock()
	// Retry the command up to 8 times
	err = try.Do(func(attempt int) (bool, error) {
		var e error
		r, e = sendCommand(daq.ser, command, respLen)
		if e != nil {
			daq.ser.Flush()
		}
		return attempt < 8, e
	})
	return
}

// Return the calibration values for a given input or output.
// The gain ID and the input mode (single-ended or differential) are needed.
// Different device models use different calibration schemas.
func (daq *OpenDAQ) GetCalib(isOutput, diffMode, secondStage bool, n, gainId uint) Calib {
	idx, err := daq.hw.GetCalibIndex(isOutput, diffMode, secondStage, n, gainId)
	if err != nil {
		return Calib{1, 0}
	}
	return daq.calib[idx]
}

// Convert a voltage to a DAC value given the number of the output
func (daq *OpenDAQ) voltsToDac(v float32, n uint) int {
	// TODO: add caching?
	cal := daq.GetCalib(true, false, false, n, 0)
	return daq.Dac.FromVolts(v, cal)
}

// Convert an ADC value to volts
func (daq *OpenDAQ) adcToVolts(raw int) float32 {
	// TODO: add caching?
	cal1 := daq.GetCalib(false, daq.diffMode, false, daq.posInput, daq.gainId)
	cal2 := daq.GetCalib(false, daq.diffMode, true, daq.posInput, daq.gainId)
	return daq.Adc.ToVolts(raw, daq.gainId, cal1, cal2)
}

func (daq *OpenDAQ) GetInfo() (model, version uint8, serial uint32, err error) {
	var buf io.Reader
	buf, err = daq.sendCommand(&Message{Number: 39}, 6)
	if err != nil {
		return
	}
	var info = struct {
		Model, Version uint8
		Serial         uint32
	}{}
	binary.Read(buf, binary.BigEndian, &info)
	model = info.Model
	version = info.Version

	if (model >= ModelEM08ABBRId){
		serial = fmt.Sprintf("8%02d%04d", model, info.Serial)
	} else {
		serial = info.Serial
	}
	return
}

// Read the calibration register stored at index nReg
func (daq *OpenDAQ) readCalib(nReg uint8) (Calib, error) {
	buf, err := daq.sendCommand(&Message{36, []byte{nReg}}, 5)
	if err != nil {
		return Calib{1, 0}, err
	}
	var ret = struct {
		_    uint8
		Gain int16
		Offs int16
	}{}
	binary.Read(buf, binary.BigEndian, &ret)
	//TODO: refactor this
	if uint(nReg) < daq.NOutputs+daq.NHiddenOutputs {
		return Calib{1. + float32(ret.Gain)/(1<<16), float32(ret.Offs) / (1 << 16)}, nil
	}
	return Calib{1. + float32(ret.Gain)/(1<<16), float32(ret.Offs) / (1 << 5)}, nil
}

func (daq *OpenDAQ) SetLED(n uint, c Color) error {
	if n < 1 || n > daq.NLeds {
		return ErrInvalidLed
	}
	if c > 3 {
		return errors.New("Invalid LED color")
	}
	_, err := daq.sendCommand(&Message{18, []byte{byte(c), byte(n)}}, 2)
	return err
}

func (daq *OpenDAQ) ConfigureADC(posInput, negInput, gainId uint, nSamples uint8) error {
	if err := daq.hw.CheckValidInputs(posInput, negInput); err != nil {
		return err
	}
	if gainId >= uint(len(daq.Adc.Gains)) {
		return ErrInvalidGainID
	}
	daq.posInput = posInput
	daq.gainId = gainId
	daq.diffMode = false
	if negInput != 0 {
		daq.diffMode = true
	}
	_, err := daq.sendCommand(&Message{2, []byte{byte(posInput), byte(negInput),
		byte(gainId), nSamples}}, 6)
	return err
}

// Read a raw value from the ADC
func (daq *OpenDAQ) ReadADC() (int16, error) {
	buf, err := daq.sendCommand(&Message{Number: 1}, 2)
	if err != nil {
		return 0, err
	}
	var val int16
	binary.Read(buf, binary.BigEndian, &val)
	return val, nil
}

// Read a value in volts from the ADC
func (daq *OpenDAQ) ReadAnalog() (float32, error) {
	val, err := daq.ReadADC()
	if err != nil {
		return 0, err
	}
	return daq.adcToVolts(int(val)), nil
}

// Set the raw value of the DAC at output n
func (daq *OpenDAQ) SetDAC(n uint, val int) error {
	if n < 1 || n > (daq.NOutputs+daq.NHiddenOutputs) {
		return ErrInvalidOutput
	}
	out := toBytes(int16(val))
	out = append(out, byte(n))
	_, err := daq.sendCommand(&Message{13, out}, 3)
	return err
}

// Set the voltage at output n
func (daq *OpenDAQ) SetAnalog(n uint, val float32) error {
	return daq.SetDAC(n, daq.voltsToDac(val, n))
}

func (daq *OpenDAQ) SetPIO(n uint, value bool) error {
	if n < 1 || n > daq.NPIOs {
		return ErrInvalidPIO
	}
	val := boolToByte(value)
	_, err := daq.sendCommand(&Message{3, []byte{byte(n), val}}, 2)
	return err
}

func (daq *OpenDAQ) SetPIODir(n uint, out bool) error {
	if n < 1 || n > daq.NPIOs {
		return ErrInvalidPIO
	}
	dir := boolToByte(out)
	_, err := daq.sendCommand(&Message{5, []byte{byte(n), dir}}, 2)
	return err
}
