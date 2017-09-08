# go-uarm
Go package for interacting with the uFactory uArm

## Install

To install locally do a `go get harrisonhjones.com/uarm`

## Example

This examples assumes that the uArm is on COM4. You may need to change it to support your system.

```golang
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

	rate := 5000
	for {
		arm.MoveXYZ(25, 0, 0, rate, false)
		arm.MoveXYZ(25, 5, 25, rate, false)
	}
}
```

## Commands

| Command                                   | Supported                             | Documented                | Description  |
| -------------                             |:-------------:                        |:-------------:            | ----- |
| New(conn io.ReadWriter)                   | :heavy_check_mark: Yes                | :x: No                    | Creates a connection to the uArm. Starts monitoring for events |
| SendRaw(cmd string)                       | :heavy_check_mark: Yes                | :x: No                    | Sends a raw command to the uArm |
| MoveXYZ(x, y, z, rate int, relative bool) | :heavy_exclamation_mark: (Untested)   | :heavy_check_mark: Yes    | Moves to either an absolute XYZ position or by a relative XYZ amount |
| GetCurrentPosXYZ()                        | :heavy_check_mark: Yes                | :x: No                    | Gets the current XYZ position of the arm |