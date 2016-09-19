package godaq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTP4XCalibIndex(t *testing.T) {
	hw := newModelTP4X()

	for i := uint(0); i < hw.NOutputs; i++ {
		idx, err := hw.GetCalibIndex(true, false, false, i+1, 0)
		assert.NoError(t, err)
		assert.EqualValues(t, i, idx)
	}

	for i := uint(0); i < hw.NInputs; i++ {
		idx, err := hw.GetCalibIndex(false, false, false, i+1, 0)
		assert.EqualValues(t, 4+i, idx)
		assert.Nil(t, err)
	}

	for i := uint(0); i < uint(len(hw.Adc.Gains)); i++ {
		idx, err := hw.GetCalibIndex(false, false, true, 1, i)
		assert.EqualValues(t, 8+i, idx)
		assert.Nil(t, err)

		idx, err = hw.GetCalibIndex(false, true, true, i+1, i)
		assert.EqualValues(t, 8+i, idx)
		assert.Nil(t, err)
	}
}
