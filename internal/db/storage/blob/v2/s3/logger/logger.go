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

package logger

import (
	"github.com/aws/smithy-go/logging"
	"github.com/vdaas/vald/internal/log"
)

type logger struct{}

// New returns aws wrapper logger function
func New() logging.Logger {
	return logger{}
}

func (l logger) Logf(classification logging.Classification, format string, v ...interface{}) {
	switch classification {
	case logging.Warn:
		log.Warnf(format, v...)
	case logging.Debug:
		log.Debugf(format, v...)
	default:
		log.Infof(format, v...)
	}
}
