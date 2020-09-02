package errors

var (
	// ErrInvalidOption returns invalid option error.
	ErrInvalidOption = func(name string, val interface{}) error {
		if val == nil {
			return Errorf("invalid option. name: %s, val: nil", name)
		}
		return Errorf("invalid option. name: %s, val: %#v", name, val)
	}
)
