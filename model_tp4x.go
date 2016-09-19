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
	ModelTP4XId  = 10
	tp4xNInputs  = 4
	tp4xNOutputs = 4
	tp4xNPIOs    = 4
	tp4xNLeds    = 4
)

var adcGainsTP4X = []float32{1, 2, 4, 8, 16, 32, 64, 128}

type ModelTP4X struct {
	HwFeatures
}

func newModelTP4X() *ModelTP4X {
	return &ModelTP4X{HwFeatures{
		Name:       "TP4X",
		NLeds:      tp4xNLeds,
		NPIOs:      tp4xNPIOs,
		NInputs:    tp4xNInputs,
		NOutputs:   tp4xNOutputs,
		NCalibRegs: uint(tp4xNOutputs + tp4xNInputs + len(adcGainsTP4X)),

		Adc: ADC{Bits: 16, Signed: true, VMin: -23.75, VMax: 23.75, Gains: adcGainsTP4X},
		// The DAC has 12 bits, but the firmware transforms the values
		Dac: DAC{Bits: 16, Signed: true, VMin: -23.75, VMax: 23.75},
	}}
}

func (m *ModelTP4X) GetFeatures() HwFeatures {
	return m.HwFeatures
}

func (m *ModelTP4X) GetCalibIndex(isOutput, diffMode, secondStage bool, n, gainId uint) (uint, error) {
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
