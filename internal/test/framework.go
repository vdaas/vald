//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Package test provides functions for general testing use
package test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/conv"
	"github.com/vdaas/vald/internal/encoding/json"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/test/goleak"
)

type Case[T, A any] struct {
	Want       Result[T]
	Args       A
	BeforeFunc BeforeFunc[A]
	AfterFunc  AfterFunc[T, A]
	CheckFunc  CheckFunc[T]
	Name       string
}

type Result[T any] struct {
	Val T
	Err error
}

type (
	BeforeFunc[A any]   func(context.Context, *testing.T, A) A
	AfterFunc[T, A any] func(context.Context, *testing.T, A, T, error) error
	CheckFunc[T any]    func(tt *testing.T, want Result[T], got Result[T]) error
	Do[T, A any]        func(*testing.T, A) (T, error)
)

func DefaultCheck[T any](tt *testing.T, want Result[T], got Result[T]) error {
	tt.Helper()
	if !errors.Is(got.Err, want.Err) {
		return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", got.Err, want.Err)
	}
	if !reflect.DeepEqual(got.Val, want.Val) {
		gb, err := json.Marshal(got.Val)
		gs := conv.Btoa(gb)
		if err != nil || gb == nil {
			gs = fmt.Sprintf("%#v", got.Val)
		}

		wb, err := json.Marshal(want.Val)
		ws := conv.Btoa(wb)
		if err != nil || wb == nil {
			ws = fmt.Sprintf("%#v", want.Val)
		}
		return errors.Errorf("got: \"%s\",\n\t\t\t\twant: \"%s\"", gs, ws)
	}
	return nil
}

func Run[T, A any](ctx context.Context, t *testing.T, do Do[T, A], tests ...Case[T, A]) error {
	t.Helper()
	ech := make(chan error, len(tests))
	defer close(ech)
	for _, tc := range tests {
		select {
		case err := <-ech:
			if err != nil {
				return err
			}
		case <-ctx.Done():
			err := ctx.Err()
			return err
		default:
			test := tc
			t.Run(test.Name, func(tt *testing.T) {
				tt.Helper()
				err := safety.RecoverFunc(func() error {
					defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
					args := test.Args
					if test.BeforeFunc != nil {
						args = test.BeforeFunc(ctx, tt, args)
					}
					checkFunc := test.CheckFunc
					if checkFunc == nil {
						checkFunc = DefaultCheck
					}
					got, err := do(tt, args)
					if err = checkFunc(tt, test.Want, Result[T]{
						Val: got,
						Err: err,
					}); err != nil {
						return err
					}
					if test.AfterFunc != nil {
						err = test.AfterFunc(ctx, tt, args, got, err)
						if err != nil {
							return err
						}
					}
					return nil
				})()
				if err != nil {
					select {
					case ech <- err:
					case <-ctx.Done():
						err := ctx.Err()
						tt.Error(err)
					}
				}
			})
		}
	}
	select {
	case err := <-ech:
		return err
	case <-ctx.Done():
		err := ctx.Err()
		return err
	default:
		return nil
	}
}
