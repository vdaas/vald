package strategy

import "github.com/vdaas/vald/internal/errors"

func wrapErrors(errs []error) (wrapped error) {
	for _, err := range errs {
		if err != nil {
			if wrapped == nil {
				wrapped = err
			} else {
				wrapped = errors.Wrap(wrapped, err.Error())
			}
		}
	}
	return
}
