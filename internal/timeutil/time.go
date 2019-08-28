package timeutil

import (
	"time"

	"github.com/vdaas/vald/internal/errors"
)

// ParseTime parses string to time.Duration.
func Parse(t string) (time.Duration, error) {
	if t == "" {
		return 0, nil
	}
	dur, err := time.ParseDuration(t)
	if err != nil {
		return 0, errors.Wrap(err, errors.ErrTimeoutParseFailed(t).Error())
	}
	return dur, nil
}
