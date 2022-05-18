package circuitbreaker

import "time"

// state is custome type for circuit breaker state.
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

type stater interface {
	state() state
	onEntry()
	onSuccess()
	onFailuer()
}

type closedState struct{}

func (cs *closedState) state() state {
	return stateClosed
}
func (cs *closedState) onEntry()   {}
func (cs *closedState) onSuccess() {}
func (cs *closedState) onFailuer() {}

type openState struct{}

func (cs *openState) state() state {
	return stateClosed
}
func (cs *openState) onEntry()   {}
func (cs *openState) onSuccess() {}
func (cs *openState) onFailuer() {}

type halfOpenState struct {
	timeoutDur time.Duration
}

func (cs *halfOpenState) state() state {
	return stateClosed
}
func (cs *halfOpenState) onEntry()   {}
func (cs *halfOpenState) onSuccess() {}
func (cs *halfOpenState) onFailuer() {}

var (
	_ stater = (*closedState)(nil)
	_ stater = (*openState)(nil)
	_ stater = (*halfOpenState)(nil)
)
