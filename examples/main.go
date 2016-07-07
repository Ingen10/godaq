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

	info, err := daq.GetInfo()
	checkErr(err)
	fmt.Println("model:", info.Model, "version:", info.Version)

	checkErr(daq.SetLED(1, godaq.GREEN))
	checkErr(daq.SetAnalog(1, 4.0))
	checkErr(daq.ConfigureADC(1, 0, 1, 1))

	for i := 0; i < 20; i++ {
		val, err := daq.ReadAnalog()
		checkErr(err)
		fmt.Println(val)
		time.Sleep(time.Millisecond * 100)
	}
}
