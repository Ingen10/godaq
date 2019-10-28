package godaq

import (
	"testing"
	"github.com/stretchr/testify/assert"
)
func TestBaseInputRawToUnits(t *testing.T) {
	ip := newInputA()
	assert.Equal(t, "INPUT_TYPE_A", ip.type_str)
	assert.EqualValues(t, 16, ip.bits)
	// gain x1
	value, units := ip.RawToUnits(-2048, 0, 0, Calib{1, 0}, Calib{1, 0})
	assert.EqualValues(t, float32(-1.5), value)
	assert.Equal(t, units, "V")
	value, _ = ip.RawToUnits(0, 0, 0, Calib{1, 0}, Calib{1, 0})
	assert.EqualValues(t, float32(0.0), value)
	value, _ = ip.RawToUnits(2048, 0, 0, Calib{1, 0}, Calib{1, 0})
	assert.EqualValues(t, float32(1.5), value)
	// gain x4
	value, units = ip.RawToUnits(-2048, 2, 0, Calib{1, 0}, Calib{1, 0})
	assert.EqualValues(t, float32(-0.375), value)
	assert.Equal(t, units, "V")
	value, _ = ip.RawToUnits(0, 2, 0, Calib{1, 0}, Calib{1, 0})
	assert.EqualValues(t, float32(0.0), value)
	value, _ = ip.RawToUnits(2048, 2, 0, Calib{1, 0}, Calib{1, 0})
	assert.EqualValues(t, float32(0.375), value)
}
func TestASInputRawToUnits(t *testing.T) {
	ip := newInputAS()
	assert.Equal(t, "INPUT_TYPE_AS", ip.type_str)
	assert.EqualValues(t, 16, ip.bits)
	// MODE V (0)
	// gain x1
	value, units := ip.RawToUnits(-2048, 0, 0, Calib{1, 0}, Calib{1, 0})
	assert.EqualValues(t, float32(-1.5), value)
	assert.Equal(t, units, "V")
	value, _ = ip.RawToUnits(0, 0, 0, Calib{1, 0}, Calib{1, 0})
	assert.EqualValues(t, float32(0.0), value)
	value, _ = ip.RawToUnits(2048, 0, 0, Calib{1, 0}, Calib{1, 0})
	assert.EqualValues(t, float32(1.5), value)
	// gain x4
	value, units = ip.RawToUnits(-2048, 2, 0, Calib{1, 0}, Calib{1, 0})
	assert.EqualValues(t, float32(-0.375), value)
	assert.Equal(t, units, "V")
	value, _ = ip.RawToUnits(0, 2, 0, Calib{1, 0}, Calib{1, 0})
	assert.EqualValues(t, float32(0.0), value)
	value, _ = ip.RawToUnits(2048, 2, 0, Calib{1, 0}, Calib{1, 0})
	assert.EqualValues(t, float32(0.375), value)
	// MODE CURRENT (1)
	value, units = ip.RawToUnits(2048, 0, 1, Calib{1, 0}, Calib{1, 0})
	assert.EqualValues(t, float32(15.0), value)
	assert.Equal(t, units, "mA")
}
func TestPInputRawToUnits(t *testing.T) {
	ip := newInputP()
	assert.Equal(t, "INPUT_TYPE_P", ip.type_str)
	assert.EqualValues(t, 16, ip.bits)
}
