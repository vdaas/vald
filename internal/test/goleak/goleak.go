package goleak

import "go.uber.org/goleak"

type (
	Option = goleak.Option
)

var (
	defaultGoleakOptions = []goleak.Option{
		goleak.IgnoreTopFunction("testing.(*testContext).waitParallel"),
	}

	IgnoreTopFunction = goleak.IgnoreTopFunction
	IgnoreCurrent     = goleak.IgnoreCurrent
)

func VerifyNone(t goleak.TestingT, options ...goleak.Option) {
	goleak.VerifyNone(t, append(options, defaultGoleakOptions...)...)
}
