//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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
	l, ok := map[string]Level{
		DEBUG.String():       DEBUG,
		DEBUG.String() + "S": DEBUG,
		"D":                  DEBUG,
		"DB":                 DEBUG,
		"DBG":                DEBUG,
		"DEB":                DEBUG,
		"DEBG":               DEBUG,
		INFO.String():        INFO,
		INFO.String() + "S":  INFO,
		"I":                  INFO,
		"IF":                 INFO,
		"IFO":                INFO,
		"IN":                 INFO,
		"INF":                INFO,
		WARN.String():        WARN,
		WARN.String() + "S":  WARN,
		"W":                  WARN,
		"WAR":                WARN,
		"WARNING":            WARN,
		"WN":                 WARN,
		"WRN":                WARN,
		ERROR.String():       ERROR,
		ERROR.String() + "S": ERROR,
		"E":                  ERROR,
		"ER":                 ERROR,
		"ERR":                ERROR,
		"ERRO":               ERROR,
		FATAL.String():       FATAL,
		FATAL.String() + "S": FATAL,
		"F":                  FATAL,
		"FAT":                FATAL,
		"FATA":               FATAL,
		"FL":                 FATAL,
		"FT":                 FATAL,
	}[strings.ToUpper(str)]
	if ok {
		return l
	}
	return Unknown
}
