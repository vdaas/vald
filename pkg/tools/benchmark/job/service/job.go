//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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

// Package service manages the main logic of benchmark job.
package service

import (
	"context"
	"os"
	"reflect"
	"strconv"
	"syscall"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s/client"
	v1 "github.com/vdaas/vald/internal/k8s/vald/benchmark/api/v1"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/rand"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/internal/test/data/hdf5"
	"github.com/vdaas/vald/internal/timeutil/rate"
)

type Job interface {
	PreStart(context.Context) error
	Start(context.Context) (<-chan error, error)
	Stop(context.Context) error
}

type jobType int

const (
	USERDEFINED jobType = iota
	INSERT
	SEARCH
	UPDATE
	UPSERT
	REMOVE
	GETOBJECT
	EXISTS
)

func (jt jobType) String() string {
	switch jt {
	case USERDEFINED:
		return "userdefined"
	case INSERT:
		return "insert"
	case SEARCH:
		return "search"
	case UPDATE:
		return "update"
	case UPSERT:
		return "upsert"
	case REMOVE:
		return "remove"
	case GETOBJECT:
		return "getobject"
	case EXISTS:
		return "exists"
	}
	return ""
}

type job struct {
	eg                 errgroup.Group
	dataset            *config.BenchmarkDataset
	jobType            jobType
	jobFunc            func(context.Context, chan error) error
	insertConfig       *config.InsertConfig
	updateConfig       *config.UpdateConfig
	upsertConfig       *config.UpsertConfig
	searchConfig       *config.SearchConfig
	removeConfig       *config.RemoveConfig
	objectConfig       *config.ObjectConfig
	client             vald.Client
	hdf5               hdf5.Data
	beforeJobName      string
	beforeJobNamespace string
	k8sClient          client.Client
	beforeJobDur       time.Duration
	limiter            rate.Limiter
	rps                int
	concurrencyLimit   int
	timeout            time.Duration
	timestamp          int64
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
		case INSERT:
			j.jobFunc = j.insert
			if j.insertConfig == nil {
				return nil, errors.NewErrInvalidOption("insert config", j.insertConfig)
			}
			ts, err := strconv.Atoi(j.insertConfig.Timestamp)
			if err != nil {
				log.Warn("[benchmark job]: ", errors.NewErrInvalidOption("insert config timestamp", j.insert, err).Error())
			} else {
				j.timestamp = int64(ts)
			}
		case SEARCH:
			j.jobFunc = j.search
			if j.searchConfig == nil {
				return nil, errors.NewErrInvalidOption("search config", j.searchConfig)
			}
			to, err := time.ParseDuration(j.searchConfig.Timeout)
			if err != nil {
				log.Warn("[benchmark job]: ", errors.NewErrInvalidOption("search config timeout", j.searchConfig.Timeout, err).Error())
			} else {
				j.timeout = to
			}
		case UPDATE:
			j.jobFunc = j.update
			if j.updateConfig == nil {
				return nil, errors.NewErrInvalidOption("update config", j.updateConfig)
			}
			ts, err := strconv.Atoi(j.updateConfig.Timestamp)
			if err != nil {
				log.Warn("[benchmark job]: ", errors.NewErrInvalidOption("update config timestamp", j.updateConfig.Timestamp, err).Error())
			} else {
				j.timestamp = int64(ts)
			}
		case UPSERT:
			j.jobFunc = j.upsert
			if j.upsertConfig == nil {
				return nil, errors.NewErrInvalidOption("upsert config", j.insertConfig)
			}
			ts, err := strconv.Atoi(j.upsertConfig.Timestamp)
			if err != nil {
				log.Warn("[benchmark job]: ", errors.NewErrInvalidOption("upsert config timestamp", j.upsertConfig.Timestamp, err).Error())
			} else {
				j.timestamp = int64(ts)
			}
		case REMOVE:
			j.jobFunc = j.remove
			if j.removeConfig == nil {
				return nil, errors.NewErrInvalidOption("insert config", j.insertConfig)
			}
			ts, err := strconv.Atoi(j.removeConfig.Timestamp)
			if err != nil {
				log.Warn("[benchmark job]: ", errors.NewErrInvalidOption("remove config timestamp", j.removeConfig.Timestamp, err).Error())
			} else {
				j.timestamp = int64(ts)
			}
		case GETOBJECT:
			j.jobFunc = j.getObject
			if j.objectConfig == nil {
				log.Warnf("[benchmark job] No get object config is set: %v", j.objectConfig)
			}
		case EXISTS:
			j.jobFunc = j.exists
		}
	} else if j.jobType != USERDEFINED {
		log.Warnf("[benchmark job] userdefined jobFunc is set but jobType is set %s", j.jobType.String())
	}
	if j.rps > 0 {
		j.limiter = rate.NewLimiter(j.rps)
	}
	// If (Range.End - Range.Start) is smaller than Indexes, Indexes are prioritized based on Range.Start.
	if (j.dataset.Range.End - j.dataset.Range.Start + 1) < j.dataset.Indexes {
		j.dataset.Range.End = j.dataset.Range.Start + j.dataset.Indexes
	}

	return j, nil
}

func (j *job) PreStart(ctx context.Context) error {
	if j.jobType != GETOBJECT && j.jobType != EXISTS && j.jobType != REMOVE {
		log.Infof("[benchmark job] start download dataset of %s", j.hdf5.GetName().String())
		if err := j.hdf5.Download(j.dataset.URL); err != nil {
			return err
		}
		log.Infof("[benchmark job] success download dataset of %s", j.hdf5.GetName().String())
		log.Infof("[benchmark job] start load dataset of %s", j.hdf5.GetName().String())
		var key hdf5.Hdf5Key
		switch j.dataset.Group {
		case "train":
			key = hdf5.Train
		case "test":
			key = hdf5.Test
		case "neighbors":
			key = hdf5.Neighors
		default:
		}
		if err := j.hdf5.Read(key); err != nil {
			return err
		}
		log.Infof("[benchmark job] success load dataset of %s", j.hdf5.GetName().String())
	}
	// Wait for beforeJob completed if exists
	if len(j.beforeJobName) != 0 {
		var jobResource v1.ValdBenchmarkJob
		log.Info("[benchmark job] check before benchjob is completed or not...")
		j.eg.Go(safety.RecoverFunc(func() error {
			dt := time.NewTicker(j.beforeJobDur)
			defer dt.Stop()
			for {
				select {
				case <-ctx.Done():
					return nil
				case <-dt.C:
					err := j.k8sClient.Get(ctx, j.beforeJobName, j.beforeJobNamespace, &jobResource)
					if err != nil {
						return err
					}
					if jobResource.Status == v1.BenchmarkJobCompleted {
						log.Infof("[benchmark job ] before job (%s) is completed, job service will start soon.", j.beforeJobName)
						return nil
					}
					log.Infof("[benchmark job] before job (%s/%s) is not completed...", j.beforeJobName, jobResource.Status)
				}
			}
		}))
		if err := j.eg.Wait(); err != nil {
			return err
		}
	}
	return nil
}

func (j *job) Start(ctx context.Context) (<-chan error, error) {
	ech := make(chan error, 3)
	cech, err := j.client.Start(ctx)
	if err != nil {
		log.Error("[benchmark job] failed to start connection monitor")
		close(ech)
		return nil, err
	}
	j.eg.Go(func() error {
		defer close(ech)
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
					ech <- errors.Join(err, ctx.Err())
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

func calcRecall(linearRes, searchRes *payload.Search_Response) (recall float64) {
	if linearRes == nil || searchRes == nil {
		return
	}
	lres := linearRes.Results
	sres := searchRes.Results
	if len(lres) == 0 || len(sres) == 0 {
		return
	}
	linearIds := map[string]struct{}{}
	for _, v := range lres {
		linearIds[v.Id] = struct{}{}
	}
	for _, v := range sres {
		if _, ok := linearIds[v.Id]; ok {
			recall++
		}
	}
	return recall / float64(len(lres))
}

// TODO: apply many object type.
func addNoiseToVec(oVec []float32) []float32 {
	noise := rand.Float32()
	vec := oVec
	idx := rand.LimitedUint32(uint64(len(oVec) - 1))
	vec[idx] += noise
	return vec
}
