//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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

// Package grpc provides grpc server logic
package grpc

import (
	"reflect"

	agent "github.com/vdaas/vald/apis/grpc/v1/agent/core"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/pkg/agent/core/ngt/service"
)

type Server interface {
	agent.AgentServer
	vald.Server
}

type server struct {
	name              string
	ip                string
	ngt               service.NGT
	eg                errgroup.Group
	streamConcurrency int
	agent.UnimplementedAgentServer
	vald.UnimplementedValdServer
}

const (
	apiName         = "vald/agent/core/ngt"
	ngtResourceType = "vald/internal/core/algorithm"
)

var errNGT = new(errors.NGTError)

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
	locs = payload.Object_LocationsFromVTPool()
	if cap(locs.GetLocations()) < len(uuids) {
		locs.Locations = make([]*payload.Object_Location, 0, len(uuids))
	}
	for _, uuid := range uuids {
		loc := payload.Object_LocationFromVTPool()
		loc.Name = s.name
		loc.Uuid = uuid
		loc.Ips = []string{s.ip}
		locs.Locations = append(locs.GetLocations(), loc)
	}
	return locs
}

func (s *server) newLocation(uuid string) *payload.Object_Location {
	locs := s.newLocations(uuid)
	if locs != nil && locs.GetLocations() != nil && len(locs.GetLocations()) > 0 {
		return locs.Locations[0]
	}
	return nil
}
