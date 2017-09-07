package uarm

import (
	"io"
	"log"
	"time"
)

const (
	maxMovementRate             = 100 // mm/min // TODO: Implement
	minTravelTime               = 500 // ms
	msInMin                     = 60000
	defaultMovementSafetyFactor = 1.5
	numOfResponseChannels       = 100 // Arbitrary
)

func New(conn io.ReadWriter) (arm *Arm, err error) {
	arm = &Arm{
		conn:                 conn,
		cmdIndex:             make(chan int, numOfResponseChannels),
		movementSafetyFactor: defaultMovementSafetyFactor,
	}

	for i := 1; i < numOfResponseChannels; i++ {
		arm.cmdIndex <- i
		arm.responses[i] = make(chan string, 1)
	}

	go arm.monitor()
	return
}

func (arm *Arm) SetLogger(logger *log.Logger) {
	// TODO: Validate
	arm.logger = logger
}

func (arm *Arm) SetReadTimeout(d time.Duration) {
	// TODO: Validate
	arm.readTimeout = d
}

func (arm *Arm) SetMovementSafetyFactor(f float64) {
	// TODO: Validate
	arm.movementSafetyFactor = f
}
