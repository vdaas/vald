// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package usecase

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	iconfig "github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/index/operator/config"
)

// fakeOperator is a minimal service.Operator stand-in: it structurally
// satisfies the interface without requiring an import of pkg/index/operator/service.
type fakeOperator struct {
	startErr error
	ch       chan error
}

func (f *fakeOperator) Start(context.Context) (<-chan error, error) {
	if f.startErr != nil {
		return nil, f.startErr
	}
	return f.ch, nil
}

// fakeServer is a minimal starter.Server (servers.Listener) stand-in.
type fakeServer struct {
	shutdownErr    error
	ch             chan error
	shutdownCalled bool
}

func (f *fakeServer) ListenAndServe(context.Context) <-chan error {
	return f.ch
}

func (f *fakeServer) Shutdown(context.Context) error {
	f.shutdownCalled = true
	return f.shutdownErr
}

// fakeObservability is a minimal observability.Observability stand-in.
type fakeObservability struct {
	preStartErr    error
	stopErr        error
	startCh        chan error
	preStartCalled bool
	startCalled    bool
	stopCalled     bool
}

func (f *fakeObservability) PreStart(context.Context) error {
	f.preStartCalled = true
	return f.preStartErr
}

func (f *fakeObservability) Start(context.Context) <-chan error {
	f.startCalled = true
	return f.startCh
}

func (f *fakeObservability) Stop(context.Context) error {
	f.stopCalled = true
	return f.stopErr
}

func mustNewGroup(t *testing.T) errgroup.Group {
	t.Helper()
	eg, _ := errgroup.New(t.Context())
	return eg
}

func validOperatorConfig() *config.Data {
	return &config.Data{
		Operator: &iconfig.IndexOperator{
			Namespace:                         "default",
			AgentName:                         "vald-agent",
			RotatorName:                       "vald-readreplica-rotate",
			TargetReadReplicaIDAnnotationsKey: "vald.vdaas.org/target-read-replica-id",
			RotationJobConcurrency:            1,
		},
	}
}

// TestNew only exercises the paths reachable without a live (or fake) Kubernetes
// API server: option-application failures inside service.New, and the
// controller-runtime manager construction failure when no cluster is reachable.
// The success path is out of scope for a unit test here; see the omissions
// reported alongside this test file.
func TestNew(t *testing.T) {
	t.Run("returns the critical option error surfaced by service.New", func(tt *testing.T) {
		tt.Parallel()

		cfg := validOperatorConfig()
		cfg.Operator.RotationJobConcurrency = 0

		r, err := New(cfg)
		require.Nil(tt, r)

		var critical *errors.ErrCriticalOption
		require.ErrorAs(tt, err, &critical)
	})

	t.Run("returns an error when no kubernetes cluster config is reachable", func(tt *testing.T) {
		// Forces client/config.GetConfig() (used by internal/k8s.New) to fail
		// deterministically, regardless of the ambient environment.
		tt.Setenv("KUBECONFIG", filepath.Join(tt.TempDir(), "does-not-exist.yaml"))

		r, err := New(validOperatorConfig())
		require.Error(tt, err)
		require.Nil(tt, r)
	})
}

func Test_run_PreStart(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("prestart failed")

	tests := []struct {
		observability *fakeObservability
		wantErr       error
		name          string
	}{
		{
			name: "returns nil when observability is nil",
		},
		{
			name:          "delegates to observability.PreStart and returns nil on success",
			observability: &fakeObservability{},
		},
		{
			name:          "forwards the error returned by observability.PreStart",
			observability: &fakeObservability{preStartErr: wantErr},
			wantErr:       wantErr,
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()

			r := &run{}
			if test.observability != nil {
				r.observability = test.observability
			}

			err := r.PreStart(tt.Context())
			if test.wantErr != nil {
				require.ErrorIs(tt, err, test.wantErr)
			} else {
				require.NoError(tt, err)
			}
			if test.observability != nil {
				require.True(tt, test.observability.preStartCalled)
			}
		})
	}
}

func Test_run_PostStop(t *testing.T) {
	t.Parallel()

	r := &run{}
	require.NoError(t, r.PostStop(t.Context()))
}

func Test_run_Stop(t *testing.T) {
	t.Parallel()

	obsErr := errors.New("observability stop failed")
	srvErr := errors.New("server shutdown failed")

	tests := []struct {
		observability *fakeObservability
		server        *fakeServer
		name          string
		wantErrs      []error
	}{
		{
			name: "returns nil when observability and server are both nil",
		},
		{
			name:          "returns the observability error when only observability fails",
			observability: &fakeObservability{stopErr: obsErr},
			wantErrs:      []error{obsErr},
		},
		{
			name:     "returns the server error when only server fails",
			server:   &fakeServer{shutdownErr: srvErr},
			wantErrs: []error{srvErr},
		},
		{
			name:          "joins both errors when observability and server fail",
			observability: &fakeObservability{stopErr: obsErr},
			server:        &fakeServer{shutdownErr: srvErr},
			wantErrs:      []error{obsErr, srvErr},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()

			r := &run{}
			if test.observability != nil {
				r.observability = test.observability
			}
			if test.server != nil {
				r.server = test.server
			}

			err := r.Stop(tt.Context())
			if len(test.wantErrs) == 0 {
				require.NoError(tt, err)
			} else {
				for _, wantErr := range test.wantErrs {
					require.ErrorIs(tt, err, wantErr)
				}
			}
			if test.observability != nil {
				require.True(tt, test.observability.stopCalled)
			}
			if test.server != nil {
				require.True(tt, test.server.shutdownCalled)
			}
		})
	}
}

func Test_run_Start(t *testing.T) {
	t.Parallel()

	const waitTimeout = 2 * time.Second

	t.Run("propagates the operator start error and closes the channel", func(tt *testing.T) {
		tt.Parallel()

		wantErr := errors.New("operator start failed")
		r := &run{
			eg:       mustNewGroup(tt),
			operator: &fakeOperator{startErr: wantErr},
			server:   &fakeServer{},
		}

		ctx, cancel := context.WithTimeout(tt.Context(), waitTimeout*2)
		defer cancel()

		ech, err := r.Start(ctx)
		require.NoError(tt, err)

		select {
		case got, ok := <-ech:
			require.True(tt, ok)
			require.ErrorIs(tt, got, wantErr)
		case <-time.After(waitTimeout):
			tt.Fatal("timed out waiting for the operator start error")
		}

		select {
		case _, ok := <-ech:
			require.False(tt, ok, "the channel must be closed after the operator start error")
		case <-time.After(waitTimeout):
			tt.Fatal("timed out waiting for the channel to close")
		}
	})

	t.Run("forwards operator, server and observability loop errors until ctx is canceled", func(tt *testing.T) {
		tt.Parallel()

		opErr := errors.New("operator loop error")
		srvErr := errors.New("server loop error")
		obsErr := errors.New("observability loop error")

		dech := make(chan error, 1)
		sech := make(chan error, 1)
		oech := make(chan error, 1)

		r := &run{
			eg:            mustNewGroup(tt),
			operator:      &fakeOperator{ch: dech},
			server:        &fakeServer{ch: sech},
			observability: &fakeObservability{startCh: oech},
		}

		ctx, cancel := context.WithCancel(tt.Context())
		defer cancel()

		ech, err := r.Start(ctx)
		require.NoError(tt, err)

		dech <- opErr
		select {
		case got := <-ech:
			require.ErrorIs(tt, got, opErr)
		case <-time.After(waitTimeout):
			tt.Fatal("timed out waiting for the operator loop error")
		}

		sech <- srvErr
		select {
		case got := <-ech:
			require.ErrorIs(tt, got, srvErr)
		case <-time.After(waitTimeout):
			tt.Fatal("timed out waiting for the server loop error")
		}

		oech <- obsErr
		select {
		case got := <-ech:
			require.ErrorIs(tt, got, obsErr)
		case <-time.After(waitTimeout):
			tt.Fatal("timed out waiting for the observability loop error")
		}

		cancel()
		select {
		case _, ok := <-ech:
			require.False(tt, ok, "the channel must be closed once ctx is canceled")
		case <-time.After(waitTimeout):
			tt.Fatal("timed out waiting for the channel to close after ctx cancellation")
		}
	})
}

// NOT IMPLEMENTED BELOW
//
// func Test_run_PreStop(t *testing.T) {
// 	type args struct {
// 		in0 context.Context
// 	}
// 	type fields struct {
// 		eg            errgroup.Group
// 		cfg           *config.Data
// 		observability observability.Observability
// 		server        starter.Server
// 		operator      service.Operator
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           in0:nil,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           cfg:nil,
// 		           observability:nil,
// 		           server:nil,
// 		           operator:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           in0:nil,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           cfg:nil,
// 		           observability:nil,
// 		           server:nil,
// 		           operator:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			r := &run{
// 				eg:            test.fields.eg,
// 				cfg:           test.fields.cfg,
// 				observability: test.fields.observability,
// 				server:        test.fields.server,
// 				operator:      test.fields.operator,
// 			}
//
// 			err := r.PreStop(test.args.in0)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
