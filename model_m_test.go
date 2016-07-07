package godaq

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCalibIndex(t *testing.T) {
	hw := NewModelM()
	idx, err := hw.GetCalibIndex(true, 0, 0, false)
	assert.EqualValues(t, 0, idx)
	assert.Nil(t, err)

	for i := uint(0); i < hw.NInputs; i++ {
		idx, err := hw.GetCalibIndex(false, i, 0, false)
		assert.EqualValues(t, 1, idx)
		assert.Nil(t, err)
	}

	for i := uint(0); i < uint(len(hw.Adc.Gains)); i++ {
		idx, err := hw.GetCalibIndex(false, 0, i, false)
		assert.EqualValues(t, 1+i, idx)
		assert.Nil(t, err)

		idx, err = hw.GetCalibIndex(false, i, i, true)
		assert.EqualValues(t, 1+i, idx)
		assert.Nil(t, err)
	}
}

func TestValidInputs(t *testing.T) {
	validNegInputs := []uint{0, 5, 6, 7, 8, 25}
	hw := NewModelM()

	var negInputs []uint
	for i := uint(0); i < 32; i++ {
		if err := hw.CheckValidInputs(1, i); err == nil {
			negInputs = append(negInputs, i)
		}
	}
	assert.Equal(t, validNegInputs, negInputs)
}
