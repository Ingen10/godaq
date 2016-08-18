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

const ModelTP4XId = 10

var adcGainsTP4X = []float32{1, 2, 4, 8, 16, 32, 64, 128}

type ModelTP4X struct {
	HwFeatures
}

func newModelTP4X() *ModelTP4X {
	return &ModelTP4X{HwFeatures{
		Name:       "TP4X",
		NLeds:      4,
		NPIOs:      4,
		NInputs:    4,
		NOutputs:   4,
		NCalibRegs: 1,

		Adc: ADC{Bits: 16, Signed: true, VMin: -23.75, VMax: 23.75, Gains: adcGainsTP4X},
		// The DAC has 12 bits, but the firmware transforms the values
		Dac: DAC{Bits: 16, Signed: true, VMin: -23.75, VMax: 23.75},
	}}
}

func (m *ModelTP4X) GetFeatures() HwFeatures {
	return m.HwFeatures
}

func (m *ModelTP4X) GetCalibIndex(isOutput bool, n, gainId uint, diffMode bool) (uint, error) {
	return 0, nil
}

func (m *ModelTP4X) CheckValidInputs(pos, neg uint) error {
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
	registerModel(ModelTP4XId, newModelTP4X())
}
