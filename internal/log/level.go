package log

import "strings"

type Level uint8

const (
	DEBUG Level = iota
	INFO
	DEBUGL
	WARN
	ERROR
	FATAL
)

func (lv Level) String() string {
	switch lv {
	case DEBUG:
		return "Debug"
	case INFO:
		return "Info"
	case WARN:
		return "Warn"
	case ERROR:
		return "Error"
	case FATAL:
		return "Fatal"
	}
	return "Unknown"
}

func ToLevel(lv string) Level {
	switch strings.ToLower(lv) {
	case "debug":
		return DEBUG
	case "info":
		return INFO
	case "warn":
		return WARN
	case "error":
		return ERROR
	case "fatal":
		return FATAL
	}
	return INFO
}
