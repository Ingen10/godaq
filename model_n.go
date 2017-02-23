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

const ModelNId = 3

var adcGainsN = []float32{1, 2, 4, 5, 8, 10, 16, 32}

type ModelN struct {
	HwFeatures
}

func NewModelN() *ModelN {
	nInputs := uint(8)
	nOutputs := uint(1)

	return &ModelN{HwFeatures{
		Name:       "OpenDAQ N",
		NLeds:      1,
		NPIOs:      6,
		NInputs:    nInputs,
		NOutputs:   nOutputs,
		NCalibRegs: nOutputs + 2*(nInputs+uint(len(adcGainsN))),

		Adc: ADC{Bits: 16, Signed: true, VMin: -12.288, VMax: 12.288, Gains: adcGainsN},
		// The DAC has 12 bits, but the firmware transforms the values
		Dac: DAC{Bits: 16, Signed: true, VMin: -4.096, VMax: 4.096},
	}}
}

func (m *ModelN) GetFeatures() HwFeatures {
	return m.HwFeatures
}

func (m *ModelN) GetCalibIndex(isOutput, diffMode, secondStage bool, n, gainId uint) (uint, error) {
	if isOutput {
		if n < 1 || n > m.NOutputs {
			return 0, ErrInvalidOutput
		}
		return n - 1, nil
	}

	if n < 1 || n > m.NInputs {
		return 0, ErrInvalidInput
	}
	index := m.NOutputs + n - 1
	if secondStage {
		index += m.NInputs + n - 1
	}
	return index, nil
}

func (m *ModelN) CheckValidInputs(pos, neg uint) error {
	if pos < 1 || pos > m.NInputs {
		return ErrInvalidInput
	}
	if neg > 8 {
		return ErrInvalidInput
	}
	return nil
}

func init() {
	// Register this model
	registerModel(ModelNId, NewModelN())
}
