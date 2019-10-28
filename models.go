package godaq

import "errors"

const (
	ModelMId 		 = 1
	ModelSId 		 = 2
	ModelNId 		 = 3
	ModelEM08ABBRId  = 10
	ModelTP04ARId    = 11
	ModelTP04ABId    = 12
	ModelEM08RRLLId  = 13
	ModelEM08LLLBId  = 14
	ModelEM08LLLLId  = 15
	ModelEM08LLARId  = 16
	ModelEM08ABRR2Id = 17
)

type HwFeatures struct {
	Name                              string
	NPIOs, NLeds                      uint
	NInputs, NOutputs, NHiddenOutputs uint
	NCalibRegs                        uint
	DacTypes                          []uint
	AdcTypes                          []uint
}

type HwModel interface {
	GetFeatures() HwFeatures
	GetCalibIndex(isOutput, diffMode, secondStage bool, n, gainId, modeInput uint) (uint, error)
}

var hwModels = make(map[uint8]HwModel)

func registerModel(model uint8, hw HwModel) error {
	if _, exists := hwModels[model]; exists {
		return errors.New("Hardware model already registered!")
	}
	hwModels[model] = hw
	return nil
}
// OPENDAQ-M MODEL
type ModelBase struct {
	HwFeatures
}

func (m *ModelBase) GetFeatures() HwFeatures {
	return m.HwFeatures
}

func (m *ModelBase) GetCalibIndex(isOutput, diffMode, secondStage bool, n, gainId, modeInput uint) (uint, error) {
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
// OPENDAQ-S MODEL
func newModelS() *ModelBase {
	nInputs := uint(8)
	nOutputs := uint(1)

	return &ModelBase{HwFeatures{
		Name:       "OpenDAQ S",
		NLeds:      1,
		NPIOs:      6,
		NInputs:    nInputs,
		NOutputs:   nOutputs,
		NCalibRegs: nOutputs + 2*nInputs,
		DacTypes:   []uint{OutputSId},
		AdcTypes:   []uint{InputSId, InputSId, InputSId, InputSId, InputSId, InputSId, InputSId, InputSId}}}
}
// OPENDAQ-N MODEL
func newModelN() *ModelBase {
	nInputs := uint(8)
	nOutputs := uint(1)
	gains := Inputtypes[InputNId].GetFeatures().gains

	return &ModelBase{HwFeatures{
		Name:       "OpenDAQ N",
		NLeds:      1,
		NPIOs:      6,
		NInputs:    nInputs,
		NOutputs:   nOutputs,
		NCalibRegs: nOutputs + 2*(nInputs+uint(len(gains))),
		DacTypes:   []uint{OutputMId},
		AdcTypes:   []uint{InputNId, InputNId, InputNId, InputNId, InputNId, InputNId, InputNId, InputNId}}}
}
// ABRR model has 2 static/tacho inputs, 2 static inputs and 4 relays
func newModelTP08ABRR() *ModelBase {
	nInputs := uint(4)
	nOutputs := uint(2) // internal tacho bias DAC references
	return &ModelBase{HwFeatures{
		Name:           "EM08S-ABRR",
		NLeds:          nInputs,
		NPIOs:          4,
		NInputs:        nInputs,
		NHiddenOutputs: nOutputs,
		NCalibRegs:     uint(nOutputs + 2*nInputs),
		DacTypes:   	[]uint{OutputTId, OutputTId},
		AdcTypes:   	[]uint{InputAId, InputAId, InputAId, InputAId}}}
}
// AR model has 2 static/tacho inputs and 2 relays
func newModelTP04AR() *ModelBase {
	nInputs := uint(2)
	nOutputs := uint(2) // internal tacho bias DAC references
	return &ModelBase{HwFeatures{
		Name:           "TP04AR",
		NLeds:          nInputs,
		NPIOs:          2,
		NInputs:        nInputs,
		NHiddenOutputs: nOutputs,
		NCalibRegs:     uint(nOutputs + 2*nInputs),
		DacTypes:   	[]uint{OutputTId, OutputTId},
		AdcTypes:   	[]uint{InputAId, InputAId}}}
}
// AB model has 2 static/tacho inputs and 2 static inputs
func newModelTP04AB() *ModelBase {
	nInputs := uint(4)
	nOutputs := uint(2) // internal tacho bias DAC references
	return &ModelBase{HwFeatures{
		Name:           "TP04AB",
		NLeds:          nInputs,
		NPIOs:          0,
		NInputs:        nInputs,
		NHiddenOutputs: nOutputs,
		NCalibRegs:     uint(nOutputs + 2*nInputs),
		DacTypes:   	[]uint{OutputTId, OutputTId},
		AdcTypes:   	[]uint{InputAId, InputAId, InputAId, InputAId}}}
}
// RRLL model has 4 relays and 4 current outputs
// In this case, Dac.VMin and Dac.VMax represent output values in mA
func newModelEM08RRLL() *ModelBase {
	nOutputs := uint(4)
	return &ModelBase{HwFeatures{
		Name:       "EM08C-RRLL",
		NLeds:      0,
		NPIOs:      4,
		NInputs:    0,
		NOutputs:   nOutputs,
		NCalibRegs: nOutputs,
		DacTypes:   	[]uint{OutputLId, OutputLId, OutputLId, OutputLId}}}
}
// LLLB model has 2 analog inputs and 6 current outputs
// In this case, Dac.VMin and Dac.VMax represent output values in mA
func newModelEM08LLLB() *ModelBase {
	nInputs := uint(2)
	nOutputs := uint(6)
	return &ModelBase{HwFeatures{
		Name:       "EM08C-LLLB",
		NLeds:      nInputs,
		NPIOs:      0,
		NInputs:    nInputs,
		NOutputs:   nOutputs,
		NCalibRegs: uint(nOutputs + 2*nInputs),
		DacTypes:   []uint{OutputLId, OutputLId, OutputLId, OutputLId, OutputLId, OutputLId},
		AdcTypes:   []uint{InputAId, InputAId}}}
}

// LLLL model has 8 current outputs
// In this case, Dac.VMin and Dac.VMax represent output values in mA
func newModelEM08LLLL() *ModelBase {
	nOutputs := uint(8)
	return &ModelBase{HwFeatures{
		Name:       "EM08C-LLLL",
		NLeds:      0,
		NPIOs:      0,
		NInputs:    0,
		NOutputs:   nOutputs,
		NCalibRegs: nOutputs,
		DacTypes:   []uint{OutputLId, OutputLId, OutputLId, OutputLId, OutputLId, OutputLId, OutputLId, OutputLId}}}
}
// LLAR model has 4 current outputs, 2 relays, 2 static inputs
// In this case, Dac.VMin and Dac.VMax represent output values in mA
type ModelM struct {
	HwFeatures
}
func newModelM() *ModelM {
	nInputs := uint(8)
	nOutputs := uint(1)
	gains := Inputtypes[InputMId].GetFeatures().gains

	return &ModelM{HwFeatures{
		Name:       "OpenDAQ M",
		NLeds:      1,
		NPIOs:      6,
		NInputs:    nInputs,
		NOutputs:   nOutputs,
		NCalibRegs: nOutputs + nInputs + uint(len(gains)),
		DacTypes:   []uint{OutputMId},
		AdcTypes:   []uint{InputMId, InputMId, InputMId, InputMId, InputMId, InputMId, InputMId, InputMId}}}
}

func (m *ModelM) GetFeatures() HwFeatures {
	return m.HwFeatures
}

func (m *ModelM) GetCalibIndex(isOutput, diffMode, secondStage bool, n, gainId, modeInput uint) (uint, error) {
	input_id := uint8(m.GetFeatures().AdcTypes[0])
	gains := Inputtypes[input_id].GetFeatures().gains
	if isOutput {
		if n < 1 || n > m.NOutputs {
			return 0, ErrInvalidOutput
		}
		return n - 1, nil
	}
	if secondStage {
		if gainId >= uint(len(gains)) {
			return 0, ErrInvalidGainID
		}
		return m.NOutputs + m.NInputs + gainId, nil
	}
	if n < 1 || n > m.NInputs {
		return 0, ErrInvalidInput
	}
	return m.NOutputs + n - 1, nil
}
// LLAR model has 4 current outputs, 2 relays, 2 static inputs
// In this case, Dac.VMin and Dac.VMax represent output values in mA
type ModelEM08LLAR struct {
	HwFeatures
}

func newModelEM08LLAR() *ModelEM08LLAR {
	nInputs := uint(2)
	nOutputs := uint(4)
	return &ModelEM08LLAR{HwFeatures{
		Name:       "EM08C-LLAR",
		NLeds:      nInputs,
		NPIOs:      2,
		NInputs:    nInputs,
		NOutputs:   nOutputs,
		NCalibRegs: uint(nOutputs + 2*nInputs),
		DacTypes:   []uint{OutputLId, OutputLId, OutputLId, OutputLId},
		AdcTypes:   []uint{InputASId, InputASId}}}
}

func (m *ModelEM08LLAR) GetFeatures() HwFeatures {
	return m.HwFeatures
}

func (m *ModelEM08LLAR) GetCalibIndex(isOutput, diffMode, secondStage bool, n, gainId, modeInput uint) (uint, error) {
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
	return m.NOutputs + m.NHiddenOutputs + n - 1 + 2 * m.NInputs * modeInput, nil
}
// ABRR2 model has 2 static/tacho inputs, 2 static inputs and 4 relays.
// new version of ABRR with shunt resistors for loop current
type ModelEM08ABRR2 struct {
	HwFeatures
}
func newModelEM08ABRR2() *ModelEM08ABRR2 {
	nInputs := uint(4)
	nOutputs := uint(2) // internal tacho bias DAC references
	return &ModelEM08ABRR2{HwFeatures{
		Name:           "EM08S-ABRR",
		NLeds:          nInputs,
		NPIOs:          4,
		NInputs:        nInputs,
		NHiddenOutputs: nOutputs,
		NCalibRegs:     uint(nOutputs + 2*nInputs),
		DacTypes:   	[]uint{OutputTId, OutputTId},
		AdcTypes:   	[]uint{InputASId, InputASId, InputASId, InputASId}}}
}

func (m *ModelEM08ABRR2) GetFeatures() HwFeatures {
	return m.HwFeatures
}

func (m *ModelEM08ABRR2) GetCalibIndex(isOutput, diffMode, secondStage bool, n, gainId, modeInput uint) (uint, error) {
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
	return m.NOutputs + m.NHiddenOutputs + n - 1 + 2 * m.NInputs * modeInput, nil
}

func init() {
	registerModel(ModelMId, newModelM())
	registerModel(ModelNId, newModelN())
	registerModel(ModelSId, newModelS())
	registerModel(ModelEM08ABBRId, newModelTP08ABRR())
	registerModel(ModelTP04ARId, newModelTP04AR())
	registerModel(ModelTP04ABId, newModelTP04AB())
	registerModel(ModelEM08RRLLId, newModelEM08RRLL())
	registerModel(ModelEM08LLLBId, newModelEM08LLLB())
	registerModel(ModelEM08LLLLId, newModelEM08LLAR())
	registerModel(ModelEM08ABRR2Id, newModelEM08ABRR2())
}