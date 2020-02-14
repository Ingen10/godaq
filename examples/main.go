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
	fmt.Println("N DAC:", daq.NOutputs)
	fmt.Println("N ADCs:", daq.NInputs)
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

	// Configure the ADC: read from input 1 with gainID=1,
	// average 10 samples each time
	checkErr(daq.ConfigureADC(1, 0, 1, 10))

	// Read 20 samples
	for i := 0; i < 10; i++ {
		val, err := daq.ReadAnalog()
		checkErr(err)
		fmt.Println(val)
		time.Sleep(100 * time.Millisecond)
		//read_val, err := daq.ReadPIO(4)
		checkErr(err)
		//fmt.Println("PIO:", read_val)
		time.Sleep(500 * time.Millisecond)
	}
	err = daq.SetPortDir(0)
	checkErr(err)
	err = daq.SetPort(42)
	checkErr(err)
}
