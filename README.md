Golang bindings for OpenDAQ
===========================

[OpenDAQ](http://www.open-daq.com) is an open source data acquisition device.

The hardware uses a VCP (Virtual COM Port) USB interface. The communication
protocol is described [here](http://opendaq.readthedocs.io/en/latest/serial-protocol.html).

**Note: This project is in an early stage. Only a small subset of the commands
is supported.**


Installation
------------

Install godaq and its dependencies to your `src` directory:

	go get -v github.com/opendaq/godaq

Updating:

	go get -u -v github.com/opendaq/godaq


Usage example
-------------

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/opendaq/godaq"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	daq, err := godaq.New("/dev/ttyUSB0")
	checkErr(err)
	defer daq.Close()

	model, version, serial, err := daq.GetInfo()
	checkErr(err)
	fmt.Println("model:", model, "version:", version, "serial:", serial)

	checkErr(daq.SetLED(1, godaq.RED))

	// Set the output voltage to 2 V
	checkErr(daq.SetAnalog(1, 2.0))

	// Configure the ADC: read from input 1 with gainID=1,
	// average 10 samples each time
	checkErr(daq.ConfigureADC(1, 0, 1, 10))

	// Read 20 samples
	for i := 0; i < 20; i++ {
		val, err := daq.ReadAnalog()
		checkErr(err)
		fmt.Println(val)
		time.Sleep(100*time.Millisecond)
	}
}
```
