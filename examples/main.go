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
	for j := 0; j <= 20; j++ {
		for color := godaq.OFF; color <= godaq.YELLOW; color++ {
			for i := uint(1); i <= daq.NLeds; i++ {
				checkErr(daq.SetLED(i, color))
			}
			time.Sleep(time.Millisecond * 100)
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

	model, version, _, err := daq.GetInfo()
	checkErr(err)
	fmt.Println("model:", model, "version:", version)

	fmt.Println("\nDAC calib:")
	for i := uint(1); i <= daq.NOutputs; i++ {
		calib := daq.GetCalib(true, false, false, i, 0)
		fmt.Println(calib)
	}
	fmt.Println("\nADC calib:")
	for i := uint(1); i <= daq.NInputs; i++ {
		calib := daq.GetCalib(false, false, false, i, 0)
		fmt.Printf("Calib %d (1st stage): %v\n", i, calib)
		calib = daq.GetCalib(false, false, true, i, 0)
		fmt.Printf("Calib %d (2nd stage): %v\n", i, calib)
	}

	for i := uint(1); i <= daq.NPIOs; i++ {
		checkErr(daq.SetPIODir(i, true))
		checkErr(daq.SetPIO(i, false))
	}

	//testLeds(daq)

	for i := uint(1); i <= daq.NOutputs; i++ {
		checkErr(daq.SetAnalog(i, 2.0))
	}

	checkErr(daq.ConfigureADC(1, 0, 1, 10))

	fmt.Println("\nAnalog readings:")
	for i := 0; i < 8; i++ {
		val, err := daq.ReadAnalog()
		checkErr(err)
		fmt.Println(val)
	}
}
