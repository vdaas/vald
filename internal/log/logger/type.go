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

package logger

import "strings"

type Type uint8

const (
	Unknown Type = iota
	GLG
	ZAP
	ZEROLOG
	LOGRUS
	KLOG
)

func (m Type) String() string {
	switch m {
	case GLG:
		return "glg"
	case ZAP:
		return "zap"
	case ZEROLOG:
		return "zerolog"
	case LOGRUS:
		return "logrus"
	case KLOG:
		return "klog"
	}
	return "unknown"
}

func Atot(str string) Type {
	switch strings.ToLower(str) {
	case "glg":
		return GLG
	case "zap":
		return ZAP
	case "zerolog":
		return ZEROLOG
	case "logrus":
		return LOGRUS
	case "klog":
		return KLOG
	}
	return Unknown
}
