package uarm

import (
	"fmt"
	"strconv"
	"strings"
)

func (arm *Arm) GetCurrentPosXYZ() (x, y, z float64, err error) {
	response, err := arm.SendRaw("P2220")
	if err != nil {
		arm.Logf("INFO", "GetCurrentPos Error: %v\n", err)
		err = fmt.Errorf("failed to send raw P2220 command: %v", err)
		return
	}

	// TODO: Validate this response?
	arm.Logf("INFO", "Response: %v", response)

	// TODO: Need to do more validation here to avoid panic
	p := strings.Split(response, " ")
	x, _ = strconv.ParseFloat(p[1][1:], 64)
	y, _ = strconv.ParseFloat(p[2][1:], 64)
	z, _ = strconv.ParseFloat(p[3][1:], 64)
	return
}
