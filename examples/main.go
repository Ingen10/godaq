package main

import (
	"fmt"
	"log"

	"github.com/opendaq/godaq"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	devices, err := godaq.ListDevicePorts()
	checkErr(err)
	if len(devices) == 0 {
		log.Fatal("No devices found")
	}

	fmt.Println(devices)
	daq, err := godaq.New(devices[0].Port)
	checkErr(err)
	defer daq.Close()

	for i := uint(1); i <= daq.NPIOs; i++ {
		checkErr(daq.SetPIODir(i, true))
		checkErr(daq.SetPIO(i, false))
	}

	model, version, _, err := daq.GetInfo()
	checkErr(err)
	fmt.Println("model:", model, "version:", version)

	for i := uint(1); i <= daq.NLeds; i++ {
		checkErr(daq.SetLED(i, godaq.RED))
	}

	//for i := uint(1); i <= daq.NOutputs; i++ {
	checkErr(daq.SetAnalog(3, -1))
	checkErr(daq.SetAnalog(4, 1))
	//}

	checkErr(daq.ConfigureADC(1, 0, 1, 1))

	for i := 0; i < 10; i++ {
		val, err := daq.ReadAnalog()
		checkErr(err)
		fmt.Println(val)
	}
}
