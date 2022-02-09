//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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
package zap

import (
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log/format"
	"github.com/vdaas/vald/internal/log/level"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	detailsKey = "details"

	defaultLevel = zapcore.DebugLevel
)

var (
	zapcore_NewConsoleEncoder = zapcore.NewConsoleEncoder
	zapcore_NewJSONEncoder    = zapcore.NewJSONEncoder
)

type logger struct {
	format format.Format
	level  level.Level

	enableCaller bool

	logger *zap.Logger
	sugar  *zap.SugaredLogger
}

// New returns a new logger instance.
func New(opts ...Option) (*logger, error) {
	l := new(logger)
	for _, opt := range append(defaultOpts, opts...) {
		opt(l)
	}

	err := l.initialize("stdout", "stderr")
	if err != nil {
		return nil, err
	}

	return l, nil
}

func (l *logger) initialize(sinkPath, errSinkPath string) (err error) {
	sink, closeOut, err := zap.Open(sinkPath)
	if err != nil {
		return err
	}

	errSink, _, err := zap.Open(errSinkPath)
	if err != nil {
		closeOut()
		return err
	}

	core := zapcore.NewCore(
		toZapEncoder(l.format),
		sink,
		toZapLevel(l.level),
	)

	opts := []zap.Option{
		zap.ErrorOutput(errSink),
	}

	if l.enableCaller {
		opts = append(opts, zap.AddCaller())
	}

	l.logger = zap.New(core, opts...)

	l.sugar = l.logger.Sugar()

	return nil
}

func (l *logger) Close() error {
	err := l.logger.Sync()
	if err != nil {
		return errors.Wrap(l.sugar.Sync(), err.Error())
	}

	return l.sugar.Sync()
}

func toZapLevel(lv level.Level) zapcore.Level {
	switch lv {
	case level.DEBUG:
		return zapcore.DebugLevel
	case level.INFO:
		return zapcore.InfoLevel
	case level.WARN:
		return zapcore.WarnLevel
	case level.ERROR:
		return zapcore.ErrorLevel
	case level.FATAL:
		return zapcore.FatalLevel
	case level.Unknown:
		fallthrough
	default:
		return defaultLevel
	}
}

func toZapEncoder(fmt format.Format) zapcore.Encoder {
	cfg := zap.NewProductionEncoderConfig()

	switch fmt {
	case format.RAW:
		return zapcore_NewConsoleEncoder(cfg)
	case format.JSON:
		return zapcore_NewJSONEncoder(cfg)
	case format.Unknown:
		fallthrough
	default:
		return zapcore_NewJSONEncoder(cfg)
	}
}

func (l *logger) Debug(vals ...interface{}) {
	l.sugar.Debug(vals...)
}

func (l *logger) Debugf(format string, vals ...interface{}) {
	l.sugar.Debugf(format, vals...)
}

func (l *logger) Debugd(msg string, details ...interface{}) {
	if len(details) == 1 {
		l.logger.Debug(msg, zap.Any(detailsKey, details[0]))
		return
	}

	l.logger.Debug(msg, zap.Any(detailsKey, details))
}

func (l *logger) Info(vals ...interface{}) {
	l.sugar.Info(vals...)
}

func (l *logger) Infof(format string, vals ...interface{}) {
	l.sugar.Infof(format, vals...)
}

func (l *logger) Infod(msg string, details ...interface{}) {
	if len(details) == 1 {
		l.logger.Info(msg, zap.Any(detailsKey, details[0]))
		return
	}

	l.logger.Info(msg, zap.Any(detailsKey, details))
}

func (l *logger) Warn(vals ...interface{}) {
	l.sugar.Warn(vals...)
}

func (l *logger) Warnf(format string, vals ...interface{}) {
	l.sugar.Warnf(format, vals...)
}

func (l *logger) Warnd(msg string, details ...interface{}) {
	if len(details) == 1 {
		l.logger.Warn(msg, zap.Any(detailsKey, details[0]))
		return
	}

	l.logger.Warn(msg, zap.Any(detailsKey, details))
}

func (l *logger) Error(vals ...interface{}) {
	l.sugar.Error(vals...)
}

func (l *logger) Errorf(format string, vals ...interface{}) {
	l.sugar.Errorf(format, vals...)
}

func (l *logger) Errord(msg string, details ...interface{}) {
	if len(details) == 1 {
		l.logger.Error(msg, zap.Any(detailsKey, details[0]))
		return
	}

	l.logger.Error(msg, zap.Any(detailsKey, details))
}

func (l *logger) Fatal(vals ...interface{}) {
	l.sugar.Fatal(vals...)
}

func (l *logger) Fatalf(format string, vals ...interface{}) {
	l.sugar.Fatalf(format, vals...)
}

func (l *logger) Fatald(msg string, details ...interface{}) {
	if len(details) == 1 {
		l.logger.Fatal(msg, zap.Any(detailsKey, details[0]))
		return
	}

	l.logger.Fatal(msg, zap.Any(detailsKey, details))
}
