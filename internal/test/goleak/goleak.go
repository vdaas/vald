// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package goleak

import "go.uber.org/goleak"

type (
	Option   = goleak.Option
	TestingT = goleak.TestingT
	TestingM = goleak.TestingM
)

var (
	defaultGoleakOptions = []goleak.Option{
		// ignore conflict with testing.T.Parallel()
		goleak.IgnoreTopFunction("testing.(*testContext).waitParallel"),
		goleak.IgnoreTopFunction("github.com/kpango/fastime.(*fastime).StartTimerD.func1"),
		goleak.IgnoreTopFunction("internal/poll.runtime_pollWait"),
		goleak.IgnoreTopFunction("net._C2func_getaddrinfo"),
		goleak.IgnoreTopFunction("go.opencensus.io/stats/view.(*worker).start"),
	}

	IgnoreTopFunction = goleak.IgnoreTopFunction
	IgnoreCurrent     = goleak.IgnoreCurrent
)

func VerifyNone(t TestingT, options ...Option) {
	goleak.VerifyNone(t, append(options, defaultGoleakOptions...)...)
}

func VerifyTestMain(m TestingM, options ...Option) {
	goleak.VerifyTestMain(m, append(options, defaultGoleakOptions...)...)
}
