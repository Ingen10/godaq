package godaq

import ("errors"
		"math")

const (
	OutputMId  = 1
	OutputSId  = 2
	OutputTId   = 3
	OutputLId   = 4
)

type OutputFeatures struct {
	type_str   string
	bits       uint
	vmin       float32
	vmax       float32
	unit       string
	signed     bool
	invert     bool
}

type OutputBase struct{
	OutputFeatures
}

type OutputModel interface {
	bitRange() (int, int)
	clampValue(value int) int
	FromVolts(v float32, cal Calib) int
}

var Outputtypes = make(map[uint8]OutputModel)

func registerOutput(output_id uint8, ot OutputModel) error {
	if _, exists := Outputtypes[output_id]; exists {
		return errors.New("Input type already registered!")
	}
	Outputtypes[output_id] = ot
	return nil
}

func newOutputM() *OutputBase { 
	return &OutputBase{OutputFeatures{
		type_str:   "OUTPUT_TYPE_M",
		bits:		16,
		vmin:		-4.096,
		vmax:		4.096,
		unit:		"V",
		signed:     true}}
}

func newOutputS() *OutputBase {
	return &OutputBase{OutputFeatures{
		type_str:   "OUTPUT_TYPE_S",
		bits:		16,
		vmin:		0,
		vmax:		4.096,
		unit:		"V",
		signed:     true}}
}

func newOutputT() *OutputBase {
	return &OutputBase{OutputFeatures{
		type_str:   "OUTPUT_TYPE_T",
		bits:		16,
		vmin:		-24,
		vmax:		24,
		unit:		"V",
		signed:     true}}
}

func newOutputL() *OutputBase {
	return &OutputBase{OutputFeatures{
		type_str:   "OUTPUT_TYPE_L",
		bits:		16,
		vmin:		0,
		vmax:		40.96,
		unit:		"mA",
		signed:     true}}
}
// GetFeatures returns the output features struct.
func (ot *OutputBase) GetFeatures() OutputFeatures {
	return ot.OutputFeatures
}
func roundInt(f float32) int {
	return int(math.Floor(float64(f) + .5))
}
// bitRange returns the range of an integer given the number of bits
func (ot *OutputBase) bitRange() (int, int) {
	out_feat := ot.GetFeatures()
	if out_feat.signed {
		return -(1 << (out_feat.bits - 1)), 1<<(out_feat.bits-1) - 1
	}
	return 0, 1<<out_feat.bits - 1
}
// clampValue limits an integer value within the representable range
func (ot *OutputBase) clampValue(value int) int {
	lower, upper := ot.bitRange()
	if value < lower {
		return lower
	} else if value > upper {
		return upper
	}
	return value
}
// FromVolts converts a voltage to a DAC value
func (ot *OutputBase) FromVolts(v float32, cal Calib) int {
	out_feat := ot.GetFeatures()
	min, max := ot.bitRange()
	var baseGain float32
	if out_feat.signed {
		baseGain = out_feat.vmax / float32(max+1)
	} else {
		baseGain = (out_feat.vmax - out_feat.vmin) / float32(max-min+1)
	}
	if out_feat.invert {
		baseGain = -baseGain
	}
	val := roundInt((v - cal.Offset) / (baseGain * cal.Gain))
	if !out_feat.signed {
		val -= int(out_feat.vmin / baseGain)
	}
	return ot.clampValue(val)
}

func init() {
	// Register this outputs
	registerOutput(OutputMId, newOutputM())
	registerOutput(OutputSId, newOutputS())
	registerOutput(OutputTId, newOutputT())
	registerOutput(OutputLId, newOutputL())
}