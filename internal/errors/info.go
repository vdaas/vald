package errors

var (
	// ErrFailedToInitInfo represents an error to initialize info.
	ErrFailedToInitInfo = func(err error) error {
		return Wrapf(err, "failed to init info")
	}
)
