//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

import "github.com/vdaas/vald/internal/strings"

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
	for i := len(str); i > 0; i-- {
		switch str[:i] {
		case DEBUG.String(), "DEB", "DEBG", "DB", "DBG", "D":
			return DEBUG
		case INFO.String(), "IFO", "INF", "IF", "IN", "I":
			return INFO
		case WARN.String(), "WARNING", "WAR", "WRN", "WN", "W":
			return WARN
		case ERROR.String(), "ERROR", "ERRO", "ER", "ERR", "E":
			return ERROR
		case FATAL.String(), "FATA", "FAT", "FT", "FL", "F":
			return FATAL
		}
	}
	return Unknown
}
