package log

import "strings"

type logLevel uint8

const (
	INFO logLevel = iota
	DEBUGL
	WARN
	ERROR
	FATAL
)

func (lv logLevel) String() string {
	switch lv {
	case INFO:
		return "info"
	case DEBUGL:
		return "debug"
	case WARN:
		return "warn"
	case ERROR:
		return "error"
	case FATAL:
		return "fatal"
	}
	return "unknown"
}

func (lv logLevel) Level(ll string) logLevel {
	switch strings.ToLower(ll) {
	case "info":
		return INFO
	case "debug":
		return DEBUGL
	case "warn":
		return WARN
	case "error":
		return ERROR
	case "fatal":
		return FATAL
	}
	return INFO
}
