package godaq

import "errors"

const (
	InputAId  = 1
	InputASId  = 2
	InputMId   = 3
	InputSId   = 4
	InputNId   = 5
	InputPId   = 6
)

type InputFeatures struct {
	type_str   string
	bits       uint
	vmin       float32
	vmax       float32
	gains      []float32
	inputmodes []uint
	unit       []string
	signed     bool
	invert     bool
}

type InputBase struct{
	InputFeatures
}

type InputAS struct {
	InputFeatures
}

type InputP struct {
	InputFeatures
}

type InputModel interface {
	RawToUnits(raw int, gainId, modeInput uint, cal1, cal2 Calib) (float32, string)
	GetFeatures() InputFeatures
}

var Inputtypes = make(map[uint8]InputModel)

func registerInput(input_id uint8, it InputModel) error {
	if _, exists := Inputtypes[input_id]; exists {
		return errors.New("Input type already registered!")
	}
	Inputtypes[input_id] = it
	return nil
}

func newInputA() *InputBase {
	return &InputBase{InputFeatures{
		type_str:   "INPUT_TYPE_A",
		bits:		16,
		vmin:		-24,
		vmax:		24,
		gains:		[]float32{1, 2, 4, 5, 8, 10, 16, 32},
		inputmodes: []uint{0},
		unit:		[]string{"V"},
		signed:     true}}
	}

func newInputAS() *InputAS {
	return &InputAS{InputFeatures{
		type_str:   "INPUT_TYPE_AS",
		bits:		16,
		vmin:		-24,
		vmax:		24,
		gains:		[]float32{1, 2, 4, 5, 8, 10, 16, 32},
		inputmodes: []uint{0, 1},
		unit:		[]string{"V", "mA"},
		signed:     true}}
}

func newInputM() *InputBase {
	return &InputBase{InputFeatures{
		type_str:   "INPUT_TYPE_M",
		bits:		16,
		vmin:		-4.096,
		vmax:		4.096,
		gains:		[]float32{1.0/3, 1, 2, 10, 100},
		inputmodes: []uint{0, 5, 6, 7, 8, 25},
		unit:		[]string{"V"},
		signed:     true}}
}

func newInputS() *InputBase {
	return &InputBase{InputFeatures{
		type_str:   "INPUT_TYPE_S",
		bits:		16,
		vmin:		-12.0,
		vmax:		12.0,
		gains:		[]float32{1, 2, 4, 5, 8, 10, 16, 20},
		inputmodes: []uint{0, 1, 2, 3, 4, 5, 6, 7, 8}, //MIRAR COMO PONER RANGE
		unit:		[]string{"V"},
		signed:     true}}
}

func newInputN() *InputBase {
	return &InputBase{InputFeatures{
		type_str:   "INPUT_TYPE_N",
		bits:		16,
		vmin:		-12.288,
		vmax:		12.288,
		gains:		[]float32{1, 2, 4, 5, 8, 10, 16, 32},
		inputmodes: []uint{0, 1, 2, 3, 4, 5, 6, 7, 8}, //MIRAR COMO PONER RANGE
		unit:		[]string{"V"},
		signed:     true}}
}

func newInputP() *InputP {
	return &InputP{InputFeatures{type_str:   "INPUT_TYPE_P",
		bits:		16,
		vmin:		-24,
		vmax:		24,
		gains:		[]float32{1, 2, 4, 5, 8, 10, 16, 32},
		inputmodes: []uint{0},
		unit:		[]string{"ohm"},
		signed:     true}}
}
// GetFeatures returns the input features struct.
func (it *InputBase) GetFeatures() InputFeatures {
	return it.InputFeatures
}
// RawToUnits Converts an ADC value to a specific unit
// gainId: Gain index
// modeInput: Inpuct lecture mode (0: V)
// cal1: pre-PGA calibration values
// cal2: post-PGA calibration values
func (it *InputBase) RawToUnits(raw int, gainId, modeInput uint, cal1, cal2 Calib) (float32, string) {
	input_feat := it.GetFeatures()
	baseOffs := 0
	if !input_feat.signed {
		baseOffs = 1 << (input_feat.bits) / 2
	}

	max := 1 << input_feat.bits
	adcGain := float32(max) / (input_feat.vmax - input_feat.vmin)
	pgaGain := input_feat.gains[gainId]
	offset := cal1.Offset + cal2.Offset*pgaGain
	gain := adcGain * pgaGain * cal1.Gain * cal2.Gain

	v := (float32(raw-baseOffs) - offset) / gain
	if input_feat.invert {
		return -v, input_feat.unit[modeInput]
	}
	return v, input_feat.unit[modeInput]
}
// GetFeatures returns the input features struct.
func (it *InputAS) GetFeatures() InputFeatures {
	return it.InputFeatures
}
// RawToUnits Converts an ADC value to a specific unit
// gainId: Gain index
// modeInput: Inpuct lecture mode (0: V, 1: mA)
// cal1: pre-PGA calibration values
// cal2: post-PGA calibration values
func (it *InputAS) RawToUnits(raw int, gainId, modeInput uint, cal1, cal2 Calib) (float32, string) {
	input_feat := it.GetFeatures()
	baseOffs := 0
	if !input_feat.signed {
		baseOffs = 1 << (input_feat.bits) / 2
	}

	max := 1 << input_feat.bits
	adcGain := float32(max) / (input_feat.vmax - input_feat.vmin)
	pgaGain := input_feat.gains[gainId]
	offset := cal1.Offset + cal2.Offset*pgaGain
	gain := adcGain * pgaGain * cal1.Gain * cal2.Gain

	v := (float32(raw-baseOffs) - offset) / gain
	if input_feat.invert {
		return -v, input_feat.unit[modeInput]
	}
	if modeInput == 1 {
		v *= 10
	}
	return v, input_feat.unit[modeInput]
}
// GetFeatures returns the input features struct.
func (it *InputP) GetFeatures() InputFeatures {
	return it.InputFeatures
}
// RawToUnits Converts an ADC value to a specific unit
// gainId: Gain index
// modeInput: Inpuct lecture mode (0: ohmnios)
// cal1: pre-PGA calibration values
// cal2: post-PGA calibration values
func (it *InputP) RawToUnits(raw int, gainId, modeInput uint, cal1, cal2 Calib) (float32, string) {
	input_feat := it.GetFeatures()
	baseOffs := 0
	if !input_feat.signed {
		baseOffs = 1 << (input_feat.bits) / 2
	}

	max := 1 << input_feat.bits
	adcGain := float32(max) / (input_feat.vmax - input_feat.vmin)
	pgaGain := input_feat.gains[gainId]
	offset := cal1.Offset + cal2.Offset*pgaGain
	gain := adcGain * pgaGain * cal1.Gain * cal2.Gain

	v := (float32(raw-baseOffs) - offset) / gain
	if input_feat.invert {
		return -v, input_feat.unit[modeInput]
	}
	return v, input_feat.unit[modeInput]
}

func init() {
	// Register this imputs
	registerInput(InputAId, newInputA())
	registerInput(InputASId, newInputAS())
	registerInput(InputMId, newInputM())
	registerInput(InputSId, newInputS())
	registerInput(InputNId, newInputN())
	registerInput(InputPId, newInputP())
}