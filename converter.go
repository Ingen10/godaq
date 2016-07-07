// Copyright 2016 The Godaq Authors. All rights reserved
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package godaq

import "math"

func roundInt(f float32) int {
	return int(math.Floor(float64(f) + .5))
}

// Digital-to-analog converter
type DAC struct {
	Bits       uint
	Signed     bool
	Invert     bool
	VMin, VMax float32
}

// Return the range of an integer given the number of bits
func (dac *DAC) bitRange() (int, int) {
	if dac.Signed {
		return -(1 << (dac.Bits - 1)), 1<<(dac.Bits-1) - 1
	}
	return 0, 1<<dac.Bits - 1
}

// Limit an integer value within the representable range
func (dac *DAC) clampValue(value int) int {
	lower, upper := dac.bitRange()
	if value < lower {
		return lower
	} else if value > upper {
		return upper
	}
	return value
}

// Convert a voltage to a DAC value
func (dac *DAC) FromVolts(v float32, cal Calib) int {
	baseGain := float32(int(1)<<dac.Bits) / (dac.VMax - dac.VMin)
	if dac.Invert {
		baseGain = -baseGain
	}
	baseOffs := float32(0)
	if !dac.Signed {
		baseOffs = -dac.VMin * baseGain
	}
	val := roundInt(v*baseGain*cal.Gain + cal.Offset + baseOffs)
	return dac.clampValue(val)
}

// Analog-to-digital converter
type ADC struct {
	Bits       uint
	Signed     bool
	Invert     bool
	VMin, VMax float32
	Gains      []float32
}

// Convert an ADC value to volts
func (adc *ADC) ToVolts(raw int, gainId uint, cal Calib) float32 {
	baseGain := adc.Gains[gainId] * float32(int(1)<<adc.Bits) / (adc.VMax - adc.VMin)
	if adc.Invert {
		baseGain = -baseGain
	}
	baseOffs := 0
	if !adc.Signed {
		baseOffs = 1 << (adc.Bits) / 2
	}
	return (float32(raw-baseOffs) - cal.Offset) / (baseGain * cal.Gain)
}