package godaq

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestBaseOutputbitRange(t *testing.T) {
	op := newOutputS()
	lower, upper := op.bitRange()
	assert.EqualValues(t, -32768, lower)
	assert.EqualValues(t, 32767, upper)
}
func TestBaseOutputclampValue(t *testing.T) {
	op := newOutputS()
	assert.Equal(t, -32768, op.clampValue(-200000))
	assert.Equal(t, 32767, op.clampValue(86000))
	assert.Equal(t, 200, op.clampValue(200))
}
func TestBaseOutputFromVolts(t *testing.T) {
	op := newOutputS()
	assert.Equal(t, 0, op.FromVolts(0.0, Calib{1, 0}))
	assert.Equal(t, 16384, op.FromVolts(2.048, Calib{1, 0}))
	assert.Equal(t, 32767, op.FromVolts(4.096, Calib{1, 0}))
	assert.Equal(t, 32767, op.FromVolts(10.0, Calib{1, 0}))

}