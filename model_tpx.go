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
	ModelEM08ABBRId = 10
	ModelTP04ARId   = 11
	ModelTP04ABId   = 12
	ModelEM08RRLLId = 13
	ModelEM08LLLBId = 14
	ModelEM08LLLLId = 15
)

var adcGainsTPX = []float32{1, 2, 4, 5, 8, 10, 16, 32}

type ModelTPX struct {
	HwFeatures
}

// AB model has 2 static/tacho inputs and 2 static inputs
func newModelTP04AB() *ModelTPX {
	nInputs := uint(4)
	nOutputs := uint(2) // internal tacho bias DAC references
	return &ModelTPX{HwFeatures{
		Name:           "TP04AB",
		NLeds:          nInputs,
		NPIOs:          0,
		NInputs:        nInputs,
		NHiddenOutputs: nOutputs,
		NCalibRegs:     uint(nOutputs + 2*nInputs),

		Adc: ADC{Bits: 16, Signed: true, VMin: -24.0, VMax: 24.0, Gains: adcGainsTPX},
		Dac: DAC{Bits: 16, Signed: true, VMin: -24.0, VMax: 24.0},
	}}
}

// AR model has 2 static/tacho inputs and 2 relays
func newModelTP04AR() *ModelTPX {
	nInputs := uint(2)
	nOutputs := uint(2) // internal tacho bias DAC references
	return &ModelTPX{HwFeatures{
		Name:           "TP04AR",
		NLeds:          nInputs,
		NPIOs:          2,
		NInputs:        nInputs,
		NHiddenOutputs: nOutputs,
		NCalibRegs:     uint(nOutputs + 2*nInputs),

		Adc: ADC{Bits: 16, Signed: true, VMin: -24.0, VMax: 24.0, Gains: adcGainsTPX},
		Dac: DAC{Bits: 16, Signed: true, VMin: -24.0, VMax: 24.0},
	}}
}

// ABRR model has 2 static/tacho inputs, 2 static inputs and 4 relays
func newModelEM08ABRR() *ModelTPX {
	nInputs := uint(4)
	nOutputs := uint(2) // internal tacho bias DAC references
	return &ModelTPX{HwFeatures{
		Name:           "EM08-ABRR",
		NLeds:          nInputs,
		NPIOs:          4,
		NInputs:        nInputs,
		NHiddenOutputs: nOutputs,
		NCalibRegs:     uint(nOutputs + 2*nInputs),

		Adc: ADC{Bits: 16, Signed: true, VMin: -24.0, VMax: 24.0, Gains: adcGainsTPX},
		Dac: DAC{Bits: 16, Signed: true, VMin: -24.0, VMax: 24.0},
	}}
}

// RRLL model has 4 relays and 4 current outputs
// In this case, Dac.VMin and Dac.VMax represent output values in mA
func newModelEM08RRLL() *ModelTPX {
	nOutputs := uint(4)
	return &ModelTPX{HwFeatures{
		Name:       "EM08-RRLL",
		NLeds:      0,
		NPIOs:      4,
		NInputs:    0,
		NOutputs:   nOutputs,
		NCalibRegs: nOutputs,

		Dac: DAC{Bits: 16, Signed: true, VMin: 0, VMax: 40.96},
	}}
}

// LLLB model has 2 analog inputs and 6 current outputs
// In this case, Dac.VMin and Dac.VMax represent output values in mA
func newModelEM08LLLB() *ModelTPX {
	nInputs := uint(2)
	nOutputs := uint(6)
	return &ModelTPX{HwFeatures{
		Name:       "EM08-LLLB",
		NLeds:      nInputs,
		NPIOs:      0,
		NInputs:    nInputs,
		NOutputs:   nOutputs,
		NCalibRegs: uint(nOutputs + 2*nInputs),

		Adc: ADC{Bits: 16, Signed: true, VMin: -24.0, VMax: 24.0, Gains: adcGainsTPX},
		Dac: DAC{Bits: 16, Signed: true, VMin: 0, VMax: 40.96},
	}}
}

// LLLL model has 8 current outputs
// In this case, Dac.VMin and Dac.VMax represent output values in mA
func newModelEM08LLLL() *ModelTPX {
	nOutputs := uint(8)
	return &ModelTPX{HwFeatures{
		Name:       "EM08-LLLL",
		NLeds:      0,
		NPIOs:      0,
		NInputs:    0,
		NOutputs:   nOutputs,
		NCalibRegs: nOutputs,

		Dac: DAC{Bits: 16, Signed: true, VMin: 0, VMax: 40.96},
	}}
}

func (m *ModelTPX) GetFeatures() HwFeatures {
	return m.HwFeatures
}

func (m *ModelTPX) GetCalibIndex(isOutput, diffMode, secondStage bool, n, gainId uint) (uint, error) {
	if isOutput {
		if n < 1 || n > (m.NOutputs+m.NHiddenOutputs) {
			return 0, ErrInvalidOutput
		}
		return n - 1, nil
	}
	if n < 1 || n > m.NInputs {
		return 0, ErrInvalidInput
	}

	if secondStage {
		return m.NOutputs + m.NHiddenOutputs + m.NInputs + n - 1, nil
	}
	return m.NOutputs + m.NHiddenOutputs + n - 1, nil
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
	registerModel(ModelTP04ABId, newModelTP04AB())
	registerModel(ModelTP04ARId, newModelTP04AR())
	registerModel(ModelEM08ABBRId, newModelEM08ABRR())
	registerModel(ModelEM08RRLLId, newModelEM08RRLL())
	registerModel(ModelEM08LLLBId, newModelEM08LLLB())
	registerModel(ModelEM08LLLLId, newModelEM08LLLL())
}
