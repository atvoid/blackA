package logging

import (
	"time"
)

type Log struct {
	TimeStamp time.Time
	Level     string
	Area      string
	Message   string
	Verbosity int
}
