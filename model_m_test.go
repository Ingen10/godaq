package godaq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMCalibIndex(t *testing.T) {
	hw := NewModelM()
	idx, err := hw.GetCalibIndex(true, false, false, 1, 0)
	assert.EqualValues(t, 0, idx)
	assert.Nil(t, err)

	for i := uint(1); i <= hw.NInputs; i++ {
		idx, err := hw.GetCalibIndex(false, false, false, i, 0)
		assert.EqualValues(t, i, idx)
		assert.Nil(t, err)
	}

	for i := uint(0); i < uint(len(hw.Adc.Gains)); i++ {
		idx, err := hw.GetCalibIndex(false, false, true, 1, i)
		assert.EqualValues(t, 9+i, idx)
		assert.Nil(t, err)

		idx, err = hw.GetCalibIndex(false, true, true, i+1, i)
		assert.EqualValues(t, 9+i, idx)
		assert.Nil(t, err)
	}
}

func TestMValidInputs(t *testing.T) {
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
