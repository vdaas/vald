package circuitbreaker

type state uint32

const (
	stateUnknown state = iota
	stateClosed
	stateOpen
	stateHalfOpen
)

func (s state) String() string {
	switch s {
	case stateClosed:
		return "closed"
	case stateOpen:
		return "open"
	case stateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}
