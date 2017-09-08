package uarm

import (
	"fmt"
	"math"
	"sync"
	"time"
)

// MoveXYZ moves the uArm to a specified point (mm) at a given rate (mm/min) if relative is false.
// It moves the uArm by a specified amount (mm) at a given rate (mm/min) if relative if true
// It returns the response from the arm or any error encountered.
// This function is largely untested. As of 9/8/2017 the firmware on the uArm can return early, before the
// arm reaches the desired position. As such, this function does not return until it has confirmed the arm has moved
// into position.
// TODO: Return on error
// TODO: Validate response and return on response error
func (arm *Arm) MoveXYZ(x, y, z, rate int, relative bool) (string, error) {
	cmd := "G0" // Absolute movement
	if relative {
		cmd = "G2204"
	}

	cmdSent := false
	respChan := make(chan string)
	respStr := ""
	respLock := sync.RWMutex{}

	// Wait for the response on another thread and update the response placeholder with a lock
	go func() {
		// TODO: Validate here
		msg := <-respChan
		respLock.Lock()
		defer respLock.Unlock()
		respStr = msg
		return
	}()

	i := 0
	for {
		i++
		// TODO: Handle errors
		// Step 1 - Get delta between current position and desired position
		delta, _ := arm.getXYZDelta(x, y, z, rate)
		if delta < 1 {
			arm.Logf("INFO", "it took %d attempts to move\n", i)
			respLock.RLock()
			defer respLock.RUnlock()
			return respStr, nil
		}
		// Step 2 - Calculate travel time & start a timer
		tt, _ := arm.getTravelTimeToPosition(delta, rate)
		timer := time.NewTimer(tt)

		// Step 3 - (once) Start movement and push the response onto the response channel
		if !cmdSent {
			go func() {
				r, err := arm.SendRaw(fmt.Sprintf("%s X%d Y%d Z%d F%d", cmd, x, y, z, rate))
				if err != nil {
					arm.Logf("INFO", "WARNING: SendRaw error: %v", err)
				}
				respChan <- r
			}()
			cmdSent = true
		}

		// Step 4 - Wait for the timeout before repeating
		<-timer.C
	}

}

func (arm *Arm) getXYZDelta(dx, dy, dz, rate int) (delta float64, err error) {
	cx, cy, cz, err := arm.GetCurrentPosXYZ()
	if err != nil {
		err = fmt.Errorf("unable to read current position: %v", err)
		return
	}
	delta = math.Abs(cx-float64(dx)) + math.Abs(cy-float64(dy)) + math.Abs(cz-float64(dz))
	arm.Logf("INFO", "Calculated delta: %v mm\n", delta)
	return
}

func (arm *Arm) getTravelTimeToPosition(delta float64, rate int) (d time.Duration, err error) {
	d = time.Millisecond * time.Duration(math.Max(minTravelTime, math.Ceil(delta/float64(rate)*msInMin*arm.movementSafetyFactor)))
	arm.Logf("INFO", "Calculated travel time: %v\n", d)
	return
}
