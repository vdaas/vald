package errors

var (
	ErrInvalidAgentResourceType = func(art string) error {
		return Errorf("invalid agent resource type: %s", art)
	}
)
