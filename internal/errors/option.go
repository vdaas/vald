package errors

var (
	// ErrNilObject represents error for nil Object.
	ErrNilObject = New("empty errorgroup")
)

var (
	// ErrEmptyString represents error for empty string.
	ErrEmptyString = New("empty string")
	// ErrInvalidNumber represents error for invalid number.
	ErrInvalidNumber = func(num int64) error {
		return Errorf("invalid number: %d", num)
	}
)
