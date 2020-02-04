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
package log

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log/format"
	"github.com/vdaas/vald/internal/log/level"
	loggertype "github.com/vdaas/vald/internal/log/logger_type"
	"github.com/vdaas/vald/internal/log/mock"
)

func TestWithLogger(t *testing.T) {
	type test struct {
		name      string
		logger    Logger
		checkFunc func(Option) error
	}

	tests := []test{
		func() test {
			logger := new(mock.Logger)

			return test{
				name:   "set success when logger is not nil",
				logger: logger,
				checkFunc: func(opt Option) error {
					option := new(option)
					opt(option)

					if !reflect.DeepEqual(option.logger, logger) {
						return errors.New("invalid params was set")
					}

					return nil
				},
			}
		}(),

		func() test {
			logger := new(mock.Logger)

			return test{
				name:   "returns nothing when logger is nil",
				logger: nil,
				checkFunc: func(opt Option) error {
					option := &option{
						logger: logger,
					}
					opt(option)

					if !reflect.DeepEqual(option.logger, logger) {
						return errors.New("invalid params was set")
					}

					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithLogger(tt.logger)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithLoggerType(t *testing.T) {
	type test struct {
		name      string
		str       string
		checkFunc func(Option) error
	}

	tests := []test{
		{
			name: "set success when str is not empty",
			str:  loggertype.GLG.String(),
			checkFunc: func(opt Option) error {
				option := new(option)
				opt(option)

				if option.loggerType != loggertype.GLG {
					return errors.New("invalid params was set")
				}
				return nil
			},
		},

		{
			name: "returns nothing when str is empty",
			checkFunc: func(opt Option) error {
				option := &option{
					loggerType: loggertype.ZAP,
				}
				opt(option)

				if option.loggerType != loggertype.ZAP {
					return errors.New("invalid params was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithLoggerType(tt.str)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithLevel(t *testing.T) {
	type test struct {
		name      string
		str       string
		checkFunc func(Option) error
	}

	tests := []test{
		{
			name: "set success when str is not empty",
			str:  level.DEBUG.String(),
			checkFunc: func(opt Option) error {
				option := new(option)
				opt(option)

				if option.level != level.DEBUG {
					return errors.New("invalid params was set")
				}
				return nil
			},
		},

		{
			name: "returns nothing when str is empty",
			checkFunc: func(opt Option) error {
				option := &option{
					level: level.ERROR,
				}
				opt(option)

				if option.level != level.ERROR {
					return errors.New("invalid params was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithLevel(tt.str)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithFormat(t *testing.T) {
	type test struct {
		name      string
		str       string
		checkFunc func(Option) error
	}

	tests := []test{
		{
			name: "set success when str is not empty",
			str:  format.JSON.String(),
			checkFunc: func(opt Option) error {
				option := new(option)
				opt(option)

				if option.format != format.JSON {
					return errors.New("invalid params was set")
				}
				return nil
			},
		},

		{
			name: "returns nothing when str is empty",
			checkFunc: func(opt Option) error {
				option := &option{
					format: format.JSON,
				}
				opt(option)

				if option.format != format.JSON {
					return errors.New("invalid params was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithFormat(tt.str)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}
