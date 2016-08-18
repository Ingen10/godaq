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

func testLeds(daq *godaq.OpenDAQ) {
	for j := 0; j <= 2000; j++ {
		for color := godaq.OFF; color <= godaq.YELLOW; color++ {
			for i := uint(1); i <= daq.NLeds; i++ {
				checkErr(daq.SetLED(i, color))
			}
			time.Sleep(time.Millisecond * 10)
		}
		fmt.Println(j)
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

	//testLeds(daq)

	for i := uint(1); i <= daq.NOutputs; i++ {
		checkErr(daq.SetAnalog(i, 20.0))
	}

	checkErr(daq.ConfigureADC(2, 0, 0, 1))

	for i := 0; i < 20; i++ {
		val, err := daq.ReadAnalog()
		checkErr(err)
		fmt.Println(val)
	}
}
