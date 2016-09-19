package godaq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBitRange(t *testing.T) {
	dac := DAC{Bits: 12, Signed: false}
	lower, upper := dac.bitRange()
	assert.Equal(t, 0, lower)
	assert.Equal(t, 4095, upper)

	dac = DAC{Bits: 16, Signed: true}
	lower, upper = dac.bitRange()
	assert.Equal(t, -32768, lower)
	assert.Equal(t, 32767, upper)
}

func TestClampValue(t *testing.T) {
	dac := DAC{Bits: 12, Signed: false}
	assert.Equal(t, 0, dac.clampValue(-200))
	assert.Equal(t, 4095, dac.clampValue(8600))
	assert.Equal(t, 200, dac.clampValue(200))

	dac = DAC{Bits: 16, Signed: true}
	assert.Equal(t, -32768, dac.clampValue(-200000))
	assert.Equal(t, 32767, dac.clampValue(86000))
	assert.Equal(t, 200, dac.clampValue(200))
}

func TestFromVoltsSigned(t *testing.T) {
	dac := DAC{Bits: 12, Signed: true, VMin: -4.096, VMax: 4.096}
	assert.Equal(t, -2048, dac.FromVolts(-10, Calib{1, 0}))
	assert.Equal(t, -2048, dac.FromVolts(-4.096, Calib{1, 0}))
	assert.Equal(t, 0, dac.FromVolts(0.0, Calib{1, 0}))
	assert.Equal(t, 2047, dac.FromVolts(4.096, Calib{1, 0}))
	assert.Equal(t, 2047, dac.FromVolts(10.0, Calib{1, 0}))
}

func TestFromVoltsUnsigned(t *testing.T) {
	dac := DAC{Bits: 12, VMin: 0, VMax: 4.096}
	assert.Equal(t, 0, dac.FromVolts(0.0, Calib{1, 0}))
	assert.Equal(t, 2048, dac.FromVolts(2.048, Calib{1, 0}))
	assert.Equal(t, 4095, dac.FromVolts(4.095, Calib{1, 0}))

	dac = DAC{Bits: 12, VMin: -4.096, VMax: 4.096}
	assert.Equal(t, 0, dac.FromVolts(-10, Calib{1, 0}))
	assert.Equal(t, 0, dac.FromVolts(-4.096, Calib{1, 0}))
	assert.Equal(t, 2048, dac.FromVolts(0.0, Calib{1, 0}))
	assert.Equal(t, 3072, dac.FromVolts(2.048, Calib{1, 0}))
	assert.Equal(t, 4095, dac.FromVolts(4.096, Calib{1, 0}))
	assert.Equal(t, 4095, dac.FromVolts(10.0, Calib{1, 0}))
}

func TestToVoltsSigned(t *testing.T) {
	gains := []float32{1, 2, 4, 8}
	adc := ADC{Bits: 12, Signed: true, VMin: -4.096, VMax: 4.096, Gains: gains}

	// gain x1
	assert.Equal(t, -4.096, adc.ToVolts(-2048, 0, Calib{1, 0}, Calib{1, 0}))
	assert.Equal(t, 0.0, adc.ToVolts(0, 0, Calib{1, 0}, Calib{1, 0}))
	assert.Equal(t, 4.096, adc.ToVolts(2048, 0, Calib{1, 0}, Calib{1, 0}))

	// gain x4
	assert.Equal(t, -1.024, adc.ToVolts(-2048, 2, Calib{1, 0}, Calib{1, 0}))
	assert.Equal(t, 0.0, adc.ToVolts(0, 2, Calib{1, 0}, Calib{1, 0}))
	assert.Equal(t, 1.024, adc.ToVolts(2048, 2, Calib{1, 0}, Calib{1, 0}))
}

func TestToVoltsUnSigned(t *testing.T) {
	gains := []float32{1, 2, 4, 8}
	adc := ADC{Bits: 12, VMin: -4.096, VMax: 4.096, Gains: gains}

	// gain x1
	assert.Equal(t, -4.096, adc.ToVolts(0, 0, Calib{1, 0}, Calib{1, 0}))
	assert.Equal(t, 0.0, adc.ToVolts(2048, 0, Calib{1, 0}, Calib{1, 0}))
	assert.Equal(t, 4.096, adc.ToVolts(4096, 0, Calib{1, 0}, Calib{1, 0}))

	// gain x4
	assert.Equal(t, -1.024, adc.ToVolts(0, 2, Calib{1, 0}, Calib{1, 0}))
	assert.Equal(t, 0.0, adc.ToVolts(2048, 2, Calib{1, 0}, Calib{1, 0}))
	assert.Equal(t, 1.024, adc.ToVolts(4096, 2, Calib{1, 0}, Calib{1, 0}))
}
