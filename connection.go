package uarm

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func (arm *Arm) SendRaw(cmd string) (response string, err error) {
	// TODO: Consider implementing a rate limiter here
	index := <-arm.cmdIndex
	cmd = fmt.Sprintf("#%d %s\n", index, cmd)

	arm.Logf("INFO", "%d - sending command: %s", index, cmd)

	select {
	case response = <-arm.responses[index]:
		arm.Logf("WARN", "%d - Warning: Old response present in channel: %v\n", index, response)
	default:
	}

	timer := time.NewTimer(arm.readTimeout)
	_, err = arm.conn.Write([]byte(cmd))
	if err != nil {
		err = fmt.Errorf("unable to write command to serial port: %v", err)
		arm.cmdIndex <- index
		return
	}

	select {
	case <-timer.C:
		err = &ReadTimeout{msg: fmt.Sprintf("timeout of %v reached", arm.readTimeout)}
		arm.Logf("INFO", "%d - timeout: %v\n", index, err)
		go arm.returnIndex(index)
	case response = <-arm.responses[index]:
		arm.Logf("INFO", "%d - Got response: %s\n", index, response)
		go arm.drainTimer(timer, index)
		arm.cmdIndex <- index
	}
	return
}

func (arm *Arm) monitor() {
	for {
		scanner := bufio.NewScanner(arm.conn)
		for scanner.Scan() {
			line := scanner.Text()
			switch {
			case strings.HasPrefix(line, "#"):
				arm.Logf("INFO", "Outgoing command: %v\n", line)
			case strings.HasPrefix(line, "$"):
				arm.Logf("INFO", "Incoming command response: %v\n", line)

				i := strings.Index(line, " ")
				number, _ := strconv.ParseInt(line[1:i], 10, 32)
				response := strings.TrimSpace(line[i:])

				select {
				case arm.responses[number] <- response:
				default:
					arm.Logf("INFO", "Unable to send response \"%s\" to response channel #%d\n", response, number)
				}
			case strings.HasPrefix(line, "@"):
				arm.Logf("INFO", "Incoming event: %v\n", line)
			default:
			}
		}
		if err := scanner.Err(); err != nil {
			arm.Logf("ERROR", "reading standard input: %v", err)
		}
	}
}

func (arm *Arm) returnIndex(index int) {
	<-arm.responses[index]
	arm.cmdIndex <- index
}
