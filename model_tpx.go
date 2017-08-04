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
	ModelTP08ABBRId = 10
	ModelTP04ARId   = 11
	ModelTP04ABId   = 12
)

var (
	adcGainsTP04 = []float32{1, 2, 4, 5, 8, 10, 16, 32}
)

type ModelTPX struct {
	HwFeatures
}

func newModelTP08ABBR() *ModelTPX {
	nInputs := uint(4)
	nOutputs := uint(2)
	return &ModelTPX{HwFeatures{
		Name:       "TP08ABBR",
		NLeds:      8,
		NPIOs:      4,
		NInputs:    nInputs,
		NOutputs:   nOutputs,
		NCalibRegs: uint(nOutputs + 2*nInputs),

		Adc: ADC{Bits: 16, Signed: true, VMin: -24.0, VMax: 24.0, Gains: adcGainsTP04},
		Dac: DAC{Bits: 16, Signed: true, VMin: -24.0, VMax: 24.0},
	}}
}

func newModelTP04AR() *ModelTPX {
	nInputs := uint(2)
	nOutputs := uint(2)
	return &ModelTPX{HwFeatures{
		Name:       "TP04AR",
		NLeds:      2,
		NPIOs:      2,
		NInputs:    nInputs,
		NOutputs:   nOutputs,
		NCalibRegs: uint(nOutputs + 2*nInputs),

		Adc: ADC{Bits: 16, Signed: true, VMin: -24.0, VMax: 24.0, Gains: adcGainsTP04},
		Dac: DAC{Bits: 16, Signed: true, VMin: -24.0, VMax: 24.0},
	}}
}

func newModelTP04AB() *ModelTPX {
	nInputs := uint(4)
	nOutputs := uint(2)
	return &ModelTPX{HwFeatures{
		Name:       "TP04AB",
		NLeds:      4,
		NPIOs:      0,
		NInputs:    nInputs,
		NOutputs:   nOutputs,
		NCalibRegs: uint(nOutputs + 2*nInputs),

		Adc: ADC{Bits: 16, Signed: true, VMin: -24.0, VMax: 24.0, Gains: adcGainsTP04},
		Dac: DAC{Bits: 16, Signed: true, VMin: -24.0, VMax: 24.0},
	}}
}

func (m *ModelTPX) GetFeatures() HwFeatures {
	return m.HwFeatures
}

func (m *ModelTPX) GetCalibIndex(isOutput, diffMode, secondStage bool, n, gainId uint) (uint, error) {
	if isOutput {
		if n < 1 || n > m.NOutputs {
			return 0, ErrInvalidOutput
		}
		return n - 1, nil
	}
	if n < 1 || n > m.NInputs {
		return 0, ErrInvalidInput
	}

	if secondStage {
		return m.NOutputs + m.NInputs + n - 1, nil
	}
	return m.NOutputs + n - 1, nil
}

func (m *ModelTPX) CheckValidInputs(pos, neg uint) error {
	if pos < 1 || pos > m.NInputs {
		return ErrInvalidInput
	}
	if neg != 0 {
		return ErrInvalidInput
	}
	return nil
}

func init() {
	// Register this models
	registerModel(ModelTP08ABBRId, newModelTP08ABBR())
	registerModel(ModelTP04ARId, newModelTP04AR())
	registerModel(ModelTP04ABId, newModelTP04AB())
}
