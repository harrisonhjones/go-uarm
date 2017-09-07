package uarm

import (
	"io"
	"log"
	"time"
)

type Arm struct {
	conn                 io.ReadWriter
	cmdIndex             chan int
	responses            [100]chan string
	readTimeout          time.Duration
	logger               *log.Logger
	movementSafetyFactor float64
}

type ReadTimeout struct {
	msg string
}

func (e *ReadTimeout) Error() string { return e.msg }
