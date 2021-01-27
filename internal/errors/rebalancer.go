package errors

var (
	ErrInvalidAgentResourceType = func(art string) error {
		return Errorf("invalid agent resource type: %s", art)
	}

	ErrJobTemplateNotFound = func() error {
		return New("job template not found")
	}

	ErrFailedToDecodeJobTemplate = func(err error) error {
		return Wrap(err, "failed to decode job template")
	}
)
