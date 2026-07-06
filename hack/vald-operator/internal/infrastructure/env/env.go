package env

import (
	"os"
	"strconv"
)

func GetEnv(key string, fallback ...string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	if len(fallback) > 0 {
		return fallback[0]
	}
	return ""
}

func GetEnvInt(key string, fallback ...string) int {
	res, err := strconv.Atoi(GetEnv(key, fallback...))
	if err != nil {
		return 0
	}
	return res
}

func GetEnvBool(key string, fallback ...string) bool {
	res, err := strconv.ParseBool(GetEnv(key, fallback...))
	if err != nil {
		return false
	}
	return res
}

func GetEnvInt64(key string, fallback ...string) int64 {
	res, err := strconv.ParseInt(GetEnv(key, fallback...), 10, 64)
	if err != nil {
		return 0
	}
	return res
}

func GetEnvFloat64(key string, fallback ...string) float64 {
	res, err := strconv.ParseFloat(GetEnv(key, fallback...), 64)
	if err != nil {
		return 0
	}
	return res
}
