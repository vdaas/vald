//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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

// Package search manages the main logic of search job.
package search

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/test/data/hdf5"
)

type SearchJob interface {
	PreStart(context.Context) error
	Start(context.Context) (<-chan error, error)
}

type searchJob struct {
	eg        errgroup.Group
	dimension int
	num       uint32
	minNum    uint32
	radius    float64
	epsilon   float64
	timeout   string
	client    vald.Client
	hdf5      hdf5.Data
}

func New(opts ...Option) (SearchJob, error) {
	s := new(searchJob)
	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(s); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}
	return s, nil
}

func (s *searchJob) PreStart(ctx context.Context) error {
	log.Infof("[bench: pre start search job] start download dataset of %#v", s.hdf5.GetName())
	if err := s.hdf5.Download(); err != nil {
		return err
	}
	log.Infof("[bench: pre start search job] success download dataset of %#v", s.hdf5.GetName())
	return nil
}

func (s *searchJob) Start(ctx context.Context) (<-chan error, error) {
	ech := make(chan error, 3)
	cech, err := s.client.Start(ctx)
	if err != nil {
		log.Error("[bench: search job] failed to start connection monitor")
		return nil, err
	}
	dur, err := time.ParseDuration(s.timeout)
	if err != nil {
		log.Error("[bench: search job] failed to timeout setting")
		return nil, err
	}

	s.eg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return nil
			case ech <- <-cech:
			}
		}
	})

	// TODO: gen search query
	var vec []float32
	for i := 0; i < s.dimension; i++ {
		vec = append(vec, float32(i)*0.1)
	}
	scfg := &payload.Search_Config{
		Num:     s.num,
		MinNum:  s.minNum,
		Radius:  float32(s.radius),
		Epsilon: float32(s.epsilon),
		Timeout: dur.Microseconds(),
	}
	s.eg.Go(func() (err error) {
		log.Infof("[bench: search job] Start search bench")
		lres, err := s.client.LinearSearch(ctx, &payload.Search_Request{
			Vector: vec,
			Config: scfg,
		})
		log.Infof("LinearSearch Result:\n %#v\n", lres)
		if err != nil {
			select {
			case <-ctx.Done():
				if err != context.Canceled {
					ech <- errors.Wrap(err, ctx.Err().Error())
				} else {
					ech <- err
				}
			case ech <- err:
			}
			return err
		}
		bres := testing.Benchmark(func(b *testing.B) {
			b.ResetTimer()
			start := time.Now()
			sres, err := s.client.Search(ctx, &payload.Search_Request{
				Vector: vec,
				Config: scfg,
			})
			if err != nil {
				select {
				case <-ctx.Done():
					if err != context.Canceled {
						ech <- errors.Wrap(err, ctx.Err().Error())
					} else {
						ech <- err
					}
				case ech <- err:
				}
			}
			latency := time.Since(start)
			recall := calcRecall(lres.Results, sres.Results)
			b.ReportMetric(recall, "recall")
			b.ReportMetric(float64(latency.Microseconds()), "latency")
		})
		// TODO: send metrics to the Prometeus
		log.Infof("[bench: search job] Finish search bench: \n%#v\n", bres)
		return nil
	})
	return ech, nil
}

func calcRecall(lres, sres []*payload.Object_Distance) (recall float64) {
	if len(lres) == 0 || len(sres) == 0 {
		return
	}
	sIds := make([]string, len(sres))
	for i, v := range sres {
		sIds[i] = v.Id
	}
	cnt := 0
	for _, v := range lres {
		if contains(v.Id, sIds) {
			cnt++
		}
	}
	return float64(cnt / len(lres))
}

func contains(tgt string, arr []string) bool {
	for _, v := range arr {
		if v == tgt {
			return true
		}
	}
	return false
}
