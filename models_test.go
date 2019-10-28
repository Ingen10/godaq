package godaq

import (
	"testing"
	"github.com/stretchr/testify/assert"
)
func TestBaseModelCalibIndex(t *testing.T) {
	hw:= newModelTP04AB()
	assert.Equal(t, "TP04AB", hw.Name)
	assert.EqualValues(t, 10, hw.NCalibRegs)

	nOuts := hw.NOutputs + hw.NHiddenOutputs

	for i := uint(0); i < hw.NOutputs; i++ {
		idx, err := hw.GetCalibIndex(true, false, false, i+1, 0, 0)
		assert.NoError(t, err)
		assert.EqualValues(t, i, idx)
	}
	for i := uint(0); i < hw.NInputs; i++ {
		// first stage calibration slots
		idx, err := hw.GetCalibIndex(false, false, false, i+1, 0, 0)
		assert.EqualValues(t, nOuts+i, idx)
		assert.Nil(t, err)
		// second stage calibration slots
		idx, err = hw.GetCalibIndex(false, false, true, i+1, 0, 0)
		assert.EqualValues(t, nOuts+hw.NInputs+i, idx)
		assert.Nil(t, err)
	}
}
func TestBaseModelGetFeatures(t *testing.T) {
	hw:= newModelM()
	features := hw.GetFeatures()
	assert.Equal(t, "OpenDAQ M", features.Name)
	assert.EqualValues(t, 1, hw.NLeds)
	assert.EqualValues(t, 6, hw.NPIOs)
	assert.EqualValues(t, 8, hw.NInputs)
	assert.EqualValues(t, 1, hw.NOutputs)
	assert.EqualValues(t, 8, hw.NInputs)
	assert.EqualValues(t, 14, hw.NCalibRegs)
	assert.EqualValues(t, []uint{OutputMId}, hw.DacTypes)
	assert.EqualValues(t, []uint{InputMId, InputMId, InputMId, InputMId, InputMId, InputMId, InputMId, InputMId},
		 hw.AdcTypes)
}
func TestABRR2ModelCalibIndex(t *testing.T) {
	hw:= newModelEM08ABRR2()
	assert.Equal(t, "EM08S-ABRR", hw.Name)
	assert.EqualValues(t, 10, hw.NCalibRegs)
	nOuts := hw.NOutputs + hw.NHiddenOutputs
	for i := uint(0); i < hw.NOutputs; i++ {
		idx, err := hw.GetCalibIndex(true, false, false, i+1, 0, 0)
		assert.NoError(t, err)
		assert.EqualValues(t, i, idx)
	}
	for i := uint(0); i < hw.NInputs; i++ {
		// second stage calibration slots
		idx, err := hw.GetCalibIndex(false, false, true, i+1, 0, 0)
		assert.EqualValues(t, nOuts+hw.NInputs+i, idx)
		assert.Nil(t, err)
		// 2 input modes
		for j:= uint(0); j < 2; j++ {
			// first stage calibration slots
			idx, err = hw.GetCalibIndex(false, false, false, i+1, 0, j)
			assert.EqualValues(t, nOuts+i+2*hw.NInputs*j, idx)
			assert.Nil(t, err)
		}
	}
}
func TestLLARModelCalibIndex(t *testing.T) {
	hw:= newModelEM08LLAR()
	assert.Equal(t, "EM08C-LLAR", hw.Name)
	assert.EqualValues(t, 8, hw.NCalibRegs)
	nOuts := hw.NOutputs + hw.NHiddenOutputs
	for i := uint(0); i < hw.NOutputs; i++ {
		idx, err := hw.GetCalibIndex(true, false, false, i+1, 0, 0)
		assert.NoError(t, err)
		assert.EqualValues(t, i, idx)
	}
	for i := uint(0); i < hw.NInputs; i++ {
		// second stage calibration slots
		idx, err := hw.GetCalibIndex(false, false, true, i+1, 0, 0)
		assert.EqualValues(t, nOuts+hw.NInputs+i, idx)
		assert.Nil(t, err)
		// 2 input modes
		for j:= uint(0); j < 2; j++ {
			// first stage calibration slots
			idx, err = hw.GetCalibIndex(false, false, false, i+1, 0, j)
			assert.EqualValues(t, nOuts+i+2*hw.NInputs*j, idx)
			assert.Nil(t, err)
		}
	}
}