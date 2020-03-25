package strategy

import "github.com/cockroachdb/errors"

func wrapErrors(errs []error) (wrapped error) {
	for _, err := range errs {
		if err != nil {
			if wrapped == nil {
				wrapped = err
			} else {
				wrapped = errors.Wrap(err, wrapped.Error())
			}
		}
	}
	return
}
