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

const modelS = 2

var adcGainsS = []float32{1, 2, 4, 5, 8, 10, 16, 20}

type ModelS struct {
	HwFeatures
}

func NewModelS() *ModelS {
	return &ModelS{HwFeatures{
		NLeds:      1,
		NPIOs:      6,
		NInputs:    8,
		NOutputs:   1,
		NCalibRegs: 1 + 8*2,

		Adc: ADC{Bits: 14, Signed: true, VMin: -4.096, VMax: 4.096, Gains: adcGainsS},
		// The DAC has 12 bits, but the firmware transforms the values
		Dac: DAC{Bits: 16, Signed: false, VMin: 0.0, VMax: 4.096},
	}}
}

func (m *ModelS) GetFeatures() HwFeatures {
	return m.HwFeatures
}

func (m *ModelS) GetCalibIndex(isOutput bool, n, gainId uint, diffMode bool) (uint, error) {
	if isOutput {
		return 0, nil
	}
	if n < 1 || n > m.NInputs {
		return 0, ErrInvalidInput
	}
	index := m.NOutputs + n
	if diffMode {
		index += m.NInputs
	}
	return index, nil
}

func (m *ModelS) CheckValidInputs(pos, neg uint) error {
	if pos < 1 || pos > m.NInputs {
		return ErrInvalidInput
	}
	if neg != 0 && neg != 25 && (neg < 5 || neg > 8) {
		return ErrInvalidInput
	}
	return nil
}

func init() {
	// Register this model
	registerModel(modelS, NewModelS())
}
