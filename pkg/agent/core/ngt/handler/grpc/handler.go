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

package grpc

import (
	"fmt"
	"reflect"
	"unicode"
	"unicode/utf8"

	agent "github.com/vdaas/vald/apis/grpc/v1/agent/core"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/agent/core/ngt/service"
)

type Server interface {
	agent.AgentServer
	vald.Server
}

type server struct {
	agent.UnimplementedAgentServer
	vald.UnimplementedValdServer
	ngt               service.NGT
	eg                errgroup.Group
	name              string
	ip                string
	streamConcurrency int
}

const (
	apiName         = "vald/agent/core/ngt"
	ngtResourceType = "vald/internal/core/algorithm"
)

// resourceInfo builds the errdetails.ResourceInfo shared by every handler error
// path: the resourceType is caller-specific, the resource name is always this
// server's apiName/name/ip triple.
func (s *server) resourceInfo(resourceType string) *errdetails.ResourceInfo {
	return &errdetails.ResourceInfo{
		ResourceType: resourceType,
		ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
	}
}

func New(opts ...Option) (Server, error) {
	s := new(server)

	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(s); err != nil {
			werr := errors.ErrOptionFailed(err, reflect.ValueOf(opt))

			e := new(errors.ErrCriticalOption)
			if errors.As(err, &e) {
				log.Error(werr)
				return nil, werr
			}
			log.Warn(werr)
		}
	}
	return s, nil
}

func (s *server) newLocations(uuids ...string) (locs *payload.Object_Locations) {
	if len(uuids) == 0 {
		return nil
	}
	locs = &payload.Object_Locations{
		Locations: make([]*payload.Object_Location, 0, len(uuids)),
	}
	for _, uuid := range uuids {
		locs.Locations = append(locs.GetLocations(), &payload.Object_Location{
			Name: s.name,
			Uuid: uuid,
			Ips:  []string{s.ip},
		})
	}
	return locs
}

// errMatcher returns a func reporting the same as errors.Is(err, target) for
// targets that can match only by message, as the freshly built per-ID errors
// (errors.ErrObjectIDNotFound(id) etc.) do. It folds the message of every node
// of err's tree once, so classifying a batch of N IDs costs O(N) instead of N
// errors.Is calls that each rebuild the whole joined message (O(N^2)).
func errMatcher(err error) func(target error) bool {
	msgs := make(map[string]struct{})
	var walk func(error)
	walk = func(e error) {
		for e != nil {
			msgs[foldMessage(e.Error())] = struct{}{}
			//nolint:errorlint // walks the error tree layer by layer like errors.Is; errors.As would skip layers.
			switch x := e.(type) {
			case interface{ Unwrap() error }:
				e = x.Unwrap()
			case interface{ Unwrap() []error }:
				for _, c := range x.Unwrap() {
					walk(c)
				}
				return
			default:
				return
			}
		}
	}
	walk(err)
	return func(target error) bool {
		_, ok := msgs[foldMessage(target.Error())]
		return ok
	}
}

// foldMessage maps every rune to the smallest rune of its simple case-folding
// orbit, so two messages fold to the same string iff strings.EqualFold (used
// by errors.Is) reports them equal.
func foldMessage(s string) string {
	return strings.Map(func(r rune) rune {
		if r < utf8.RuneSelf {
			// the upper-case ASCII letter is the smallest rune of its orbit.
			if 'a' <= r && r <= 'z' {
				return r - 'a' + 'A'
			}
			return r
		}
		m := r
		for f := unicode.SimpleFold(r); f != r; f = unicode.SimpleFold(f) {
			m = min(m, f)
		}
		return m
	}, s)
}

func (s *server) newLocation(uuid string) *payload.Object_Location {
	return &payload.Object_Location{
		Name: s.name,
		Uuid: uuid,
		Ips:  []string{s.ip},
	}
}
