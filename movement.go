package uarm

import (
	"fmt"
	"math"
	"time"
)

func (arm *Arm) MoveTo(x, y, z, rate int) error {
	start := make(chan bool, 0)
	response := make(chan string)

	// TODO: Move this inside of the for loop
	go func() {
		<-start
		r, err := arm.SendRaw(fmt.Sprintf("G0 X%d Y%d Z%d F%d", x, y, z, rate))
		if err != nil {
			arm.Logf("INFO", "WARNING: SendRaw error: %v", err)
		}
		response <- r
	}()

	i := 0
	for {
		i++
		// TODO: Implement
		// TODO: Handle errors
		// Step 1 - Get delta between current position and desired position
		delta, _ := arm.getPositionDelta(x, y, z, rate)
		if delta < 1 {
			arm.Logf("INFO", "it took %d attempts to move\n", i)
			return nil
		}
		// Step 2 - Calculate travel time & start a timer
		tt, _ := arm.getTravelTimeToPosition(delta, rate)
		timer := time.NewTimer(tt)

		// Step 3 - (once) Start movement
		select {
		case start <- true:
		default:
		}

		// Step 4 - On timeout go to step 1
		<-timer.C
		// TODO: Do something with the response / validate it
		go func() { <-response }() // Drain the response channel
	}

}

func (arm *Arm) getPositionDelta(x, y, z, rate int) (delta float64, err error) {
	cx, cy, cz, err := arm.GetCurrentPos()
	if err != nil {
		err = fmt.Errorf("unable to read current position: %v", err)
		return
	}
	delta = math.Abs(cx-float64(x)) + math.Abs(cy-float64(y)) + math.Abs(cz-float64(z))
	arm.Logf("INFO", "Calculated delta: %v mm\n", delta)
	return
}

func (arm *Arm) getTravelTimeToPosition(delta float64, rate int) (d time.Duration, err error) {
	d = time.Millisecond * time.Duration(math.Max(minTravelTime, math.Ceil(delta/float64(rate)*msInMin*arm.movementSafetyFactor)))
	arm.Logf("INFO", "Calculated travel time: %v\n", d)
	return
}
