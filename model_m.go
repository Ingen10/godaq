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
	return &ModelM{HwFeatures{
		Name:       "OpenDAQ M",
		NLeds:      1,
		NPIOs:      6,
		NInputs:    8,
		NOutputs:   1,
		NCalibRegs: uint(1 + len(adcGainsM)),

		Adc: ADC{Bits: 16, Signed: true, VMin: -4.096, VMax: 4.096,
			Invert: true, Gains: adcGainsM},
		Dac: DAC{Bits: 14, Signed: false, VMin: -4.096, VMax: 4.096},
	}}
}

func (m *ModelM) GetFeatures() HwFeatures {
	return m.HwFeatures
}

func (m *ModelM) GetCalibIndex(isOutput bool, n, gainId uint, diffMode bool) (uint, error) {
	if isOutput {
		return 0, nil
	}
	if gainId >= uint(len(m.Adc.Gains)) {
		return 0, ErrInvalidGainID
	}
	return gainId + 1, nil
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
