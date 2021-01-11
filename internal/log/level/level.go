//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package level

import "strings"

type Level uint8

const (
	Unknown Level = iota

	// DEBUG is debug log level.
	DEBUG

	// INFO is info log level.
	INFO

	// WARN is warning log level.
	WARN

	// ERRO is error log level.
	ERROR

	// FATAL is fatal log level.
	FATAL
)

func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	}
	return "Unknown"
}

func Atol(str string) Level {
	str = strings.ToUpper(str)
	switch str {
	case DEBUG.String():
		return DEBUG
	case INFO.String():
		return INFO
	case WARN.String():
		return WARN
	case ERROR.String():
		return ERROR
	case FATAL.String():
		return FATAL
	}
	switch {
	case len(str) < 3:
		return Unknown
	case str[:3] == DEBUG.String()[:3]:
		return DEBUG
	case str[:3] == INFO.String()[:3]:
		return INFO
	case str[:3] == WARN.String()[:3]:
		return WARN
	case str[:3] == ERROR.String()[:3]:
		return ERROR
	case str[:3] == FATAL.String()[:3]:
		return FATAL
	}

	return Unknown
}
