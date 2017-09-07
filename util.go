package uarm

import (
	"fmt"
	"strings"
	"time"
)

func (arm *Arm) drainTimer(t *time.Timer, index int) {
	// TODO: Consider removing index param
	if !t.Stop() {
		<-t.C
	}
	arm.Logf("INFO", "%d - timer drained", index)
}

func (arm *Arm) Logf(level, format string, args ...interface{}) {
	if arm.logger != nil {
		arm.logger.Printf("[%s] %s\n", strings.ToUpper(level), fmt.Sprintf(strings.TrimSpace(format), args))
	}
}
