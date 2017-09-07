package main

import (
	"log"
	"time"

	"github.com/tarm/serial"
	"harrisonhjones.com/uarm"
)

func main() {
	log.Println("Opening serial port")
	c := &serial.Config{
		Name:        "COM4",
		Baud:        115200,
		ReadTimeout: time.Millisecond * 100,
	}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	arm, _ := uarm.New(s)
	time.Sleep(time.Millisecond * 5000)

	/*response := arm.MustSendRaw("P2234")
	log.Printf("Response: %s\n", response)*/
	rate := 5000
	for {
		arm.MoveTo(25, 0, 0, rate)
		arm.MoveTo(25, 5, 25, rate)
	}
}
