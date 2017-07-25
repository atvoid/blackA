package logging

import (
	"fmt"
)

const (
	LOGLEVEL_SIMPLE = 0
	LOGLEVEL_NORMAL = 1
	LOGLEVEL_DETAIL = 2
)

type Logger struct {
	logBuffer chan Log
	endSig    chan bool
}

var globalLogger = Logger{logBuffer: make(chan Log, 50)}
var globalLogVerbosity = LOGLEVEL_SIMPLE

func StartLogging(logLevel int) {
	globalLogVerbosity = logLevel
	go globalLogger.startLogging()
}

func EndLogging() {
	globalLogger.endSig <- true
}

func LogInfo(area, msg string) {
	log := Log{Level: "Info", Area: area, Message: msg, Verbosity: LOGLEVEL_SIMPLE}
	globalLogger.logBuffer <- log
}

func LogInfo_Normal(area, msg string) {
	log := Log{Level: "Info", Area: area, Message: msg, Verbosity: LOGLEVEL_NORMAL}
	globalLogger.logBuffer <- log
}

func LogInfo_Detail(area, msg string) {
	log := Log{Level: "Info", Area: area, Message: msg, Verbosity: LOGLEVEL_DETAIL}
	globalLogger.logBuffer <- log
}

func LogError(area, msg string) {
	log := Log{Level: "Error", Area: area, Message: msg}
	globalLogger.logBuffer <- log
}

func (this *Logger) writeLog(log Log) {
	fmt.Printf("Level: %v\t|\t Area: %v\t|\t Msg: %v\n", log.Level, log.Area, log.Message)
}

func (this *Logger) filterLog(log *Log) bool {
	if log.Level == "Info" && log.Verbosity > globalLogVerbosity {
		return false
	}
	return true
}

func (this *Logger) startLogging() {
PollLoop:
	for {
		select {
		case log := <-this.logBuffer:
			if this.filterLog(&log) {
				this.writeLog(log)
			}
		case <-this.endSig:
			break PollLoop
		default:
		}
	}
}
