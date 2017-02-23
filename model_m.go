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

const ModelMId = 1

var adcGainsM = []float32{1.0 / 3, 1, 2, 10, 100}

type ModelM struct {
	HwFeatures
}

func NewModelM() *ModelM {
	nInputs := uint(8)
	nOutputs := uint(1)

	return &ModelM{HwFeatures{
		Name:       "OpenDAQ M",
		NLeds:      1,
		NPIOs:      6,
		NInputs:    nInputs,
		NOutputs:   nOutputs,
		NCalibRegs: nOutputs + nInputs + uint(len(adcGainsM)),

		Adc: ADC{Bits: 16, Signed: true, VMin: -4.096, VMax: 4.096,
			Invert: true, Gains: adcGainsM},
		Dac: DAC{Bits: 16, Signed: true, VMin: -4.096, VMax: 4.096},
	}}
}

func (m *ModelM) GetFeatures() HwFeatures {
	return m.HwFeatures
}

// Get the index of a calibration register.
// Each register contains a pair of calibration values: a gain and and offset.
//
// isOutput: Obtain the calibration values of an output
// diffMode: Some models have different calibration values depending on the input mode (single-ended or differential)
// secondStage: The inputs with a PGA need two calibration registers. One is applied before the PGA and the other
// is applied after the PGA (second stage)
func (m *ModelM) GetCalibIndex(isOutput, diffMode, secondStage bool, n, gainId uint) (uint, error) {
	if isOutput {
		if n < 1 || n > m.NOutputs {
			return 0, ErrInvalidOutput
		}
		return n - 1, nil
	}
	if secondStage {
		if gainId >= uint(len(m.Adc.Gains)) {
			return 0, ErrInvalidGainID
		}
		return m.NOutputs + m.NInputs + gainId, nil
	}

	if n < 1 || n > m.NInputs {
		return 0, ErrInvalidInput
	}
	return m.NOutputs + n - 1, nil
}

func (m *ModelM) CheckValidInputs(pos, neg uint) error {
	if pos < 1 || pos > m.NInputs {
		return ErrInvalidInput
	}
	if neg != 0 && neg != 25 && (neg < 5 || neg > 8) {
		return ErrInvalidInput
	}
	return nil
}

func init() {
	registerModel(ModelMId, NewModelM())
}
