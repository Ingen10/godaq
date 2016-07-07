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
	"time"

	"github.com/tarm/serial"
)

type Color uint8

const (
	OFF Color = iota
	GREEN
	RED
	YELLOW
)

const (
	// Hardware models
	modelM = 1
	modelS = 2

	// Voltage ranges of the ADC (with gain factor = 1)
	adcVMinM = -4.096
	adcVMaxM = 4.096
	adcVMinS = -4.096
	adcVMaxS = 4.096

	// Voltage ranges of the DAC
	dacVMinM = -4.096
	dacVMaxM = 4.096
	dacVMinS = 0.0
	dacVMaxS = 4.096

	// Number of bits of the ADC
	adcBitsM = 16
	adcBitsS = 14

	// Number of bits of the DAC
	dacBitsM = 14
	dacBitsS = 16 // The DAC has only 12 bits, but the firmware scales the value
)

var (
	// Gain factors
	adcGainsM = []float32{1.0 / 3, 1, 2, 10, 100}
	adcGainsS = []float32{1, 2, 4, 5, 8, 10, 16, 20}

	ErrInvalidModel  = errors.New("Invalid device model number")
	ErrInvalidOutput = errors.New("Invalid output number")
	ErrInvalidLed    = errors.New("Invalid LED number")
	ErrInvalidInput  = errors.New("Invalid input number")
	ErrInvalidGainID = errors.New("Invalid gain ID")
)

type Info struct {
	Model, Version uint8
	SerialNumber   uint32
}

type Calib struct {
	Gain   float32 // Gain calibration (-1 to 1)
	Offset float32 // Offset calibraton in ADUs
}

type OpenDAQ struct {
	ser *serial.Port

	model             uint8
	nPIOs, nLeds      uint
	nInputs, nOutputs uint
	nCalibRegs        uint
	adcBits, dacBits  uint
	adcVMin, adcVMax  float32
	dacVMin, dacVMax  float32
	adcGains          []float32
	calib             []Calib

	// Input state (needed for converting ADC values to volts)
	gainId   uint
	posInput uint
	diffMode bool
}

// Return the range of an integer given the number of bits
func bitRange(bits uint, signed bool) (int, int) {
	if signed {
		return -(1 << (bits - 1)), 1<<(bits-1) - 1
	}
	return 0, 1<<bits - 1
}

// Limit an integer value within the representable range, given the
// number of bits
func clampToBitRange(value int, bits uint, signed bool) int {
	lower, upper := bitRange(bits, signed)
	if value < lower {
		return lower
	} else if value > upper {
		return upper
	}
	return value
}

func New(port string) (*OpenDAQ, error) {
	var err error
	daq := OpenDAQ{}

	// Setup and open the serial port
	serCfg := &serial.Config{Name: port, Baud: 115200, ReadTimeout: time.Second}
	daq.ser, err = serial.OpenPort(serCfg)
	if err != nil {
		return nil, err
	}
	time.Sleep(2 * time.Second)

	// Obtain the device model number
	info, err := daq.GetInfo()
	if err != nil {
		return nil, err
	}
	daq.model = info.Model
	daq.posInput = 1 // 0 is not a valid default for posInput

	switch daq.model {
	case modelM:
		daq.nLeds = 1
		daq.nPIOs = 6
		daq.nInputs = 8
		daq.nOutputs = 1
		daq.adcBits = adcBitsM
		daq.dacBits = dacBitsM
		daq.adcVMin = adcVMinM
		daq.adcVMax = adcVMaxM
		daq.dacVMin = dacVMinM
		daq.dacVMax = dacVMaxM
		daq.adcGains = adcGainsM
		daq.calib = make([]Calib, 1+len(adcGainsM))
	case modelS:
		daq.nLeds = 1
		daq.nPIOs = 6
		daq.nInputs = 8
		daq.nOutputs = 1
		daq.adcBits = adcBitsS
		daq.dacBits = dacBitsS
		daq.adcVMin = adcVMinS
		daq.adcVMax = adcVMaxS
		daq.dacVMin = dacVMinS
		daq.dacVMax = dacVMaxS
		daq.adcGains = adcGainsS
		daq.calib = make([]Calib, 1+2*len(adcGainsS))
	default:
		return nil, ErrInvalidModel
	}

	// Read the calibration registers from the device
	for i := range daq.calib {
		if daq.calib[i], err = daq.readCalib(uint8(i)); err != nil {
			return nil, err
		}
	}
	return &daq, nil
}

// Check the validity of  a combination of positive and negative inputs
func (daq *OpenDAQ) checkValidInputs(posInput, negInput uint) error {
	if posInput < 1 || posInput > daq.nInputs {
		return ErrInvalidInput
	}
	if negInput == 0 {
		return nil
	}

	switch daq.model {
	case modelM:
		if negInput == 25 {
			break
		}
		if negInput < 5 || negInput > 8 {
			return ErrInvalidInput
		}
	case modelS:
		if posInput%2 == 0 && negInput != posInput-1 {
			return ErrInvalidInput
		} else if posInput%2 == 1 && negInput != posInput+1 {
			return ErrInvalidInput
		}
		return nil
	default:
		return ErrInvalidModel
	}
	return nil
}

func (daq *OpenDAQ) Close() error {
	return daq.ser.Close()
}

// Send a comand and returns its response
func (daq *OpenDAQ) sendCommand(command *Message, respLen int) (io.Reader, error) {
	return sendCommand(daq.ser, command, respLen)
}

// Return the calibration values of an output given its number
func (daq *OpenDAQ) GetDACCalib(n uint) Calib {
	if n < 1 || n > daq.nOutputs {
		return Calib{1, 0}
	}
	return daq.calib[int(n-1)]
}

// Return the calibration values of the ADC, given the positive input,
// the gain ID and the input mode (single-ended or differential).
// Different device models use different calibration schemas.
func (daq *OpenDAQ) GetADCCalib(posInput, gainId uint, diffMode bool) (cal Calib) {
	switch daq.model {
	case modelM:
		if gainId >= uint(len(daq.adcGains)) {
			cal = Calib{1, 0}
		}
		cal = daq.calib[int(daq.nOutputs+gainId)]
		cal.Gain = -cal.Gain // in model M the input value is inverted!
	case modelS:
		if posInput < 1 || posInput > daq.nInputs {
			cal = Calib{1, 0}
		}
		index := daq.nOutputs + posInput - 1
		if diffMode {
			index += daq.nInputs
		}
		cal = daq.calib[int(index)]
	}
	return
}

// Convert a voltage to a DAC value given the number of the output.
// The DAC value is always positive (from 0 to 2^nbits - 1), even though
// the output voltage can be negative
func (daq *OpenDAQ) voltsToDac(v float32, n uint) int {
	// TODO: add caching?
	cal := daq.GetDACCalib(n)
	baseGain := float32(int(1)<<daq.dacBits) / (daq.dacVMax - daq.dacVMin)
	baseOffs := 0
	if daq.dacVMin < 0 {
		baseOffs = 1 << (daq.dacBits - 1)
	}
	val := int(v*baseGain*cal.Gain+cal.Offset) + baseOffs
	return clampToBitRange(val, daq.dacBits, false)
}

// Convert an ADC value to volts
func (daq *OpenDAQ) adcToVolts(raw int) float32 {
	// TODO: add caching?
	cal := daq.GetADCCalib(daq.posInput, daq.gainId, daq.diffMode)
	baseGain := daq.adcGains[daq.gainId] * float32(int(1)<<daq.adcBits) /
		(daq.adcVMax - daq.adcVMin)
	return (float32(raw) - cal.Offset) / (baseGain * cal.Gain)
}

func (daq *OpenDAQ) GetInfo() (*Info, error) {
	buf, err := daq.sendCommand(&Message{Number: 39}, 6)
	if err != nil {
		return nil, err
	}
	var info Info
	binary.Read(buf, binary.BigEndian, &info)
	return &info, nil
}

// Read the calibration register stored at index n
func (daq *OpenDAQ) readCalib(n uint8) (Calib, error) {
	buf, err := daq.sendCommand(&Message{36, []byte{n}}, 5)
	if err != nil {
		return Calib{1, 0}, err
	}
	var ret = struct {
		_    uint8
		Gain int16
		Offs int16
	}{}
	binary.Read(buf, binary.BigEndian, &ret)
	return Calib{1. + float32(ret.Gain)/1e5, float32(ret.Offs)}, nil
}

func (daq *OpenDAQ) SetLED(n uint, c Color) error {
	if n < 1 || n > daq.nLeds {
		return ErrInvalidLed
	}
	if c > 3 {
		return errors.New("Invalid LED color")
	}
	_, err := daq.sendCommand(&Message{18, []byte{byte(c)}}, 1)
	return err
}

func (daq *OpenDAQ) ConfigureADC(posInput, negInput, gainId uint, nSamples uint8) error {
	if err := daq.checkValidInputs(posInput, negInput); err != nil {
		return err
	}
	if gainId >= uint(len(daq.adcGains)) {
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
	if n < 1 || n > daq.nOutputs {
		return ErrInvalidOutput
	}
	_, err := daq.sendCommand(&Message{24, toBytes(int16(val))}, 2)
	return err
}

// Set the voltage at output n
func (daq *OpenDAQ) SetAnalog(n uint, val float32) error {
	return daq.SetDAC(n, daq.voltsToDac(val, n))
}
