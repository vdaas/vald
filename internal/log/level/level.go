//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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

	// DEBG is debug log level
	DEBG

	// INFO is info log level
	INFO

	// WARN is warning log level
	WARN

	// ERR is error log level
	ERR

	// CRIT is critical error log level
	CRIT

	// FATAL is fatal log level
	FATAL
)

func (l Level) String() string {
	switch l {
	case DEBG:
		return "DEBG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERR:
		return "ERR"
	case CRIT:
		return "CRIT"
	case FATAL:
		return "FATAL"
	}
	return "unknown"
}

func Atol(str string) Level {
	switch strings.ToLower(str) {
	case "debug":
		return DEBG
	case "info":
		return INFO
	case "warn":
		return WARN
	case "error":
		return ERR
	case "crit":
		return CRIT
	case "fatal":
		return FATAL
	}
	return Unknown
}
