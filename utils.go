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

import (
	"errors"
	"io/ioutil"
	"runtime"
	"strings"
)

// List all available USB-serial ports (Linux only)
func ListPorts() ([]string, error) {
	if runtime.GOOS != "linux" {
		return nil, errors.New("Not supported OS")
	}
	files, err := ioutil.ReadDir("/dev")
	if err != nil {
		return nil, err
	}
	var list []string
	for _, file := range files {
		n := file.Name()
		if strings.HasPrefix(n, "ttyUSB") || strings.HasPrefix(n, "ttyACM") {
			list = append(list, "/dev/"+n)
		}
	}
	return list, nil
}

type DevicePort struct {
	Model uint8
	Port  string
}

func ListDevicePorts() ([]DevicePort, error) {
	ports, err := ListPorts()
	if err != nil {
		return nil, err
	}
	var list []DevicePort
	for _, port := range ports {
		if dev, err := New(port); err == nil {
			if model, _, _, err := dev.GetInfo(); err == nil {
				list = append(list, DevicePort{model, port})
			}
			dev.Close()
		}
	}
	return list, nil
}
