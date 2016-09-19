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

const (
	ModelSId  = 2
	sNInputs  = 8
	sNOutputs = 1
)

var adcGainsS = []float32{1, 2, 4, 5, 8, 10, 16, 20}

type ModelS struct {
	HwFeatures
}

func NewModelS() *ModelS {
	return &ModelS{HwFeatures{
		Name:       "OpenDAQ S",
		NLeds:      1,
		NPIOs:      6,
		NInputs:    sNInputs,
		NOutputs:   sNOutputs,
		NCalibRegs: uint(sNOutputs + 2*(sNInputs+len(adcGainsS))),

		Adc: ADC{Bits: 14, Signed: true, VMin: -4.096, VMax: 4.096, Gains: adcGainsS},
		// The DAC has 12 bits, but the firmware transforms the values
		Dac: DAC{Bits: 16, Signed: false, VMin: 0.0, VMax: 4.096},
	}}
}

func (m *ModelS) GetFeatures() HwFeatures {
	return m.HwFeatures
}

func (m *ModelS) GetCalibIndex(isOutput, diffMode, secondStage bool, n, gainId uint) (uint, error) {
	if isOutput {
		if n < 1 || n > m.NOutputs {
			return 0, ErrInvalidOutput
		}
		return n - 1, nil
	}

	var index uint
	if secondStage {
		if gainId >= uint(len(m.Adc.Gains)) {
			return 0, ErrInvalidGainID
		}
		index = m.NOutputs + m.NInputs + gainId
	} else {
		if n < 1 || n > m.NInputs {
			return 0, ErrInvalidInput
		}
		index = m.NOutputs + n - 1
	}
	if diffMode {
		index += m.NInputs + uint(len(m.Adc.Gains))
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
	registerModel(ModelSId, NewModelS())
}
