package log

import "strings"

type logLevel uint8

const (
	DEBUG logLevel = iota
	INFO
	DEBUGL
	WARN
	ERROR
	FATAL
)

func (lv logLevel) String() string {
	switch lv {
	case DEBUG:
		return "debug"
	case INFO:
		return "info"
	case WARN:
		return "warn"
	case ERROR:
		return "error"
	case FATAL:
		return "fatal"
	}
	return "unknown"
}

func level(lv string) logLevel {
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
