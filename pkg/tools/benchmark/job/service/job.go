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
	"os"
	"reflect"
	"syscall"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/test/data/hdf5"
)

type Job interface {
	PreStart(context.Context) error
	Start(context.Context) (<-chan error, error)
	Stop(context.Context) error
}

type jobType int

const (
	USERDEFINED jobType = iota
	SEARCH
)

func (jt jobType) String() string {
	switch jt {
	case USERDEFINED:
		return "userdefined"
	case SEARCH:
		return "search"
	}
	return ""
}

type job struct {
	eg           errgroup.Group
	dimension    int
	dataset      *config.BenchmarkDataset
	jobType      jobType
	jobFunc      func(context.Context, chan error) error
	insertConfig *config.InsertConfig
	updateConfig *config.UpdateConfig
	upsertConfig *config.UpsertConfig
	searchConfig *config.SearchConfig
	removeConfig *config.RemoveConfig
	client       vald.Client
	hdf5         hdf5.Data
}

func New(opts ...Option) (Job, error) {
	j := new(job)
	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(j); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}
	if j.jobFunc == nil {
		switch j.jobType {
		case USERDEFINED:
			opt := WithJobFunc(j.jobFunc)
			err := opt(j)
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		case SEARCH:
			j.jobFunc = j.search
		}
	} else if j.jobType != USERDEFINED {
		log.Warnf("[benchmark job] userdefined jobFunc is set but jobType is set %s", j.jobType.String())
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

	j.eg.Go(func() (err error) {
		defer func() {
			p, perr := os.FindProcess(os.Getpid())
			if perr != nil {
				log.Error(perr)
				return
			}
			if err != nil {
				select {
				case <-ctx.Done():
					ech <- errors.Wrap(err, ctx.Err().Error())
				case ech <- err:
				}
			}
			if err := p.Signal(syscall.SIGTERM); err != nil {
				log.Error(err)
			}
		}()
		err = j.jobFunc(ctx, ech)
		if err != nil {
			log.Errorf("[benchmark job] failed to job: %v", err)
		}
		return
	})

	return ech, nil
}

func (j *job) Stop(ctx context.Context) (err error) {
	err = j.client.Stop(ctx)
	return
}

func calcRecall(linearRes, searchRes []*payload.Object_Distance) (recall float64) {
	if len(linearRes) == 0 || len(searchRes) == 0 {
		return
	}
	linearIds := map[string]struct{}{}
	for _, v := range linearRes {
		linearIds[v.Id] = struct{}{}
	}
	for _, v := range searchRes {
		if _, ok := linearIds[v.Id]; ok {
			recall++
		}
	}
	return recall / float64(len(linearRes))
}

func genVec(data [][]float32, cfg *config.BenchmarkDataset) [][]float32 {
	start := cfg.Range.Start
	end := cfg.Range.End
	if (end - start) < cfg.Indexes {
		end = cfg.Indexes
	}
	num := end - start + 1
	if len(data) < num {
		num = len(data)
		end = start + num + 1
	}
	vectors := data[start : end+1]
	return vectors
}
