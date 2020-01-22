package glg

import "strings"

type level uint8

const (
	DEBUG level = iota
	INFO
	DEBUGL
	WARN
	ERROR
	FATAL
)

func (lv level) String() string {
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

func toLevel(lv string) level {
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
