package logging

import (
	"fmt"
)

type Logger struct {
	logBuffer		chan Log
	endSig			chan bool
}

var globalLogger = Logger{ logBuffer: make(chan Log, 50) }

func StartLogging() {
	go globalLogger.startLogging()
}

func EndLogging() {
	globalLogger.endSig <- true
}

func LogInfo(area, msg string) {
	log := Log{ Level: "Info", Area: area, Message: msg }
	globalLogger.logBuffer <- log
}

func LogError(area, msg string) {
	log := Log{ Level: "Error", Area: area, Message: msg }
	globalLogger.logBuffer <- log
}

func (this *Logger) writeLog(log Log) {
	fmt.Printf("Level: %v\t|\t Area: %v\t|\t Info: %v\n", log.Level, log.Area, log.Message)
}

func (this *Logger) startLogging() {
	PollLoop:
	for {
		select {
			case log := <- this.logBuffer:
				this.writeLog(log)
			case <- this.endSig:
				break PollLoop
			default:
		}
	}
}