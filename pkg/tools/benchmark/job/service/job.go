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

// Package service manages the main logic of benchmark job.
package service

import (
	"context"
	"reflect"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/test/data/hdf5"
)

type Job interface {
	PreStart(context.Context) error
	Start(context.Context) (<-chan error, error)
}

type jobType int

const (
	SEARCH jobType = iota
)

type job struct {
	eg        errgroup.Group
	jobType   jobType
	dimension int
	iter      int
	num       uint32
	minNum    uint32
	radius    float64
	epsilon   float64
	timeout   time.Duration
	client    vald.Client
	hdf5      hdf5.Data
}

func New(opts ...Option) (Job, error) { j := new(job)
	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(j); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}
	return j, nil
}

func (j *job) PreStart(ctx context.Context) error {
	log.Infof("[benchmark job] start download dataset of %s", j.hdf5.GetName().String())
	if err := j.hdf5.Download(); err != nil {
		return err
	}
	log.Infof("[benchmark job] success download dataset of %s", j.hdf5.GetName().String())
	log.Infof("[benchmark job] start load dataset of %s", j.hdf5.GetName().String())
	if err := j.hdf5.Read(); err != nil {
		return err
	}
	log.Infof("[benchmark job] success load dataset of %s", j.hdf5.GetName().String())
	return nil
}

func (j *job) Start(ctx context.Context) (<-chan error, error) {
	ech := make(chan error, 3)
	cech, err := j.client.Start(ctx)
	if err != nil {
		log.Error("[benchmark job] failed to start connection monitor")
		return nil, err
	}
	j.eg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return nil
			case ech <- <-cech:
			}
		}
	})

	switch j.jobType {
	case SEARCH:
		err := search(ctx, j, ech)
		if err != nil {
			return ech, err
		}
	}
	return ech, nil
}

func calcRecall(linearRes, searchRes []*payload.Object_Distance) float64 {
	if len(linearRes) == 0 || len(searchRes) == 0 {
		return 0
	}
	linearIds := make([]string, len(linearRes))
	for i, v := range linearRes {
		linearIds[i] = v.Id
	}
	cnt := 0
	for _, v := range searchRes {
		if contains(v.Id, linearIds) {
			cnt++
		}
	}
	return float64(cnt / len(linearRes))
}

func contains(target string, arr []string) bool {
	for _, v := range arr {
		if v == target {
			return true
		}
	}
	return false
}
