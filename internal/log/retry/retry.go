package retry

import (
	"reflect"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
)

func Out(
	fn func(vals ...interface{}) error,
	vals ...interface{},
) {
	Outf(func(format string, vals ...interface{}) error {
		return fn(vals...)
	}, "", vals...)
}

func Outf(
	fn func(format string, vals ...interface{}) error,
	format string, vals ...interface{},
) {
	if err := fn(format, vals...); err != nil {
		rv := reflect.ValueOf(fn)

		log.Warn(errors.ErrLoggingRetry(err, rv))

		err = fn(format, vals...)
		if err != nil {
			log.Error(errors.ErrLoggingFaild(err, rv))

			err = fn(format, vals...)
			if err != nil {
				panic(err)
			}
		}
	}
}
