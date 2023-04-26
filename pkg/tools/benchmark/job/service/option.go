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

// Package service manages the main logic of benchmark job.
package service

import (
	"context"

	"github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s/client"
	"github.com/vdaas/vald/internal/test/data/hdf5"
	"github.com/vdaas/vald/internal/timeutil"
)

type Option func(j *job) error

var defaultOpts = []Option{
	// TODO: set default config for client
	WithDimension(748),
	WithBeforeJobDuration("30s"),
	WithRPS(100),
}

// WithDimension sets the vector's dimension for running benchmark job with dataset.
func WithDimension(dim int) Option {
	return func(j *job) error {
		if dim > 0 {
			j.dimension = dim
		}
		return nil
	}
}

// WithInsertConfig sets the insert API config for running insert request job.
func WithInsertConfig(c *config.InsertConfig) Option {
	return func(j *job) error {
		if c != nil {
			j.insertConfig = c
		}
		return nil
	}
}

// WithUpdateConfig sets the update API config for running update request job.
func WithUpdateConfig(c *config.UpdateConfig) Option {
	return func(j *job) error {
		if c != nil {
			j.updateConfig = c
		}
		return nil
	}
}

// WithUpsertConfig sets the upsert API config for running upsert request job.
func WithUpsertConfig(c *config.UpsertConfig) Option {
	return func(j *job) error {
		if c != nil {
			j.upsertConfig = c
		}
		return nil
	}
}

// WithSearchConfig sets the search API config for running search request job.
func WithSearchConfig(c *config.SearchConfig) Option {
	return func(j *job) error {
		if c != nil {
			j.searchConfig = c
		}
		return nil
	}
}

// WithRemoveConfig sets the remove API config for running remove request job.
func WithRemoveConfig(c *config.RemoveConfig) Option {
	return func(j *job) error {
		if c != nil {
			j.removeConfig = c
		}
		return nil
	}
}

// WithObjectConfig sets the get object API config for running get object request job.
func WithObjectConfig(c *config.ObjectConfig) Option {
	return func(j *job) error {
		if c != nil {
			j.objectConfig = c
		}
		return nil
	}
}

// WithValdClient sets the Vald client for sending request to the target Vald cluster.
func WithValdClient(c vald.Client) Option {
	return func(j *job) error {
		if c == nil {
			return errors.NewErrInvalidOption("client", c)
		}
		j.client = c
		return nil
	}
}

// WithErrGroup sets the errgroup to the job struct to handle errors.
func WithErrGroup(eg errgroup.Group) Option {
	return func(j *job) error {
		if eg == nil {
			return errors.NewErrInvalidOption("error group", eg)
		}
		j.eg = eg
		return nil
	}
}

// WithHdf5 sets the hdf5.Data which is used for benchmark job dataset.
func WithHdf5(d hdf5.Data) Option {
	return func(j *job) error {
		if d == nil {
			return errors.NewErrInvalidOption("hdf5", d)
		}
		j.hdf5 = d
		return nil
	}
}

// WithDataset sets the config.BenchmarkDataset including benchmakr dataset name, group name of hdf5.Data, the number of index, start range and end range.
func WithDataset(d *config.BenchmarkDataset) Option {
	return func(j *job) error {
		if d == nil {
			return errors.NewErrInvalidOption("dataset", d)
		}
		j.dataset = d
		return nil
	}
}

// WithJobTypeByString converts given string to JobType.
func WithJobTypeByString(t string) Option {
	var jt jobType
	switch t {
	case "userdefined":
		jt = USERDEFINED
	case "insert":
		jt = INSERT
	case "search":
		jt = SEARCH
	case "update":
		jt = UPDATE
	case "upsert":
		jt = UPSERT
	case "remove":
		jt = REMOVE
	case "getobject":
		jt = GETOBJECT
	case "exists":
		jt = EXISTS
	}
	return WithJobType(jt)
}

// WithJobType sets the jobType for running benchmark job.
func WithJobType(jt jobType) Option {
	return func(j *job) error {
		if len(jt.String()) == 0 {
			return errors.NewErrInvalidOption("jobType", jt.String())
		}
		j.jobType = jt
		return nil
	}
}

// WithJobFunc sets the job function.
func WithJobFunc(jf func(context.Context, chan error) error) Option {
	return func(j *job) error {
		if jf == nil {
			return errors.NewErrInvalidOption("jobFunc", jf)
		}
		j.jobFunc = jf
		return nil
	}
}

// WithBeforeJobName sets the beforeJobName which we should wait for until finish before running job.
func WithBeforeJobName(bjn string) Option {
	return func(j *job) error {
		if len(bjn) > 0 {
			j.beforeJobName = bjn
		}
		return nil
	}
}

// WithBeforeJobNamespace sets the beforeJobNamespace of the beforeJobName which we should wait for until finish before running job.
func WithBeforeJobNamespace(bjns string) Option {
	return func(j *job) error {
		if len(bjns) > 0 {
			j.beforeJobNamespace = bjns
		}
		return nil
	}
}

// WithBeforeJobDuration sets the duration for watching beforeJobName's status.
func WithBeforeJobDuration(dur string) Option {
	return func(j *job) error {
		if len(dur) == 0 {
			return nil
		}
		dur, err := timeutil.Parse(dur)
		if err != nil {
			return err
		}
		j.beforeJobDur = dur
		return nil
	}
}

// WithK8sClient binds the k8s client to the job struct which is used for get BenchmarkJobResource from Kubernetes API server.
func WithK8sClient(cli client.Client) Option {
	return func(j *job) error {
		if cli != nil {
			j.k8sClient = cli
		}
		return nil
	}
}

// WithRPS sets the rpc for sending request per seconds to the target Vald cluster.
func WithRPS(rps int) Option {
	return func(j *job) error {
		if rps > 0 {
			j.rps = rps
		}
		return nil
	}
}
