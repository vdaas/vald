//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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

package cassandra

import (
	"context"
	"sync/atomic"

	"github.com/gocql/gocql"
)

type QueryObserver interface {
	gocql.QueryObserver
	CompletedQueryTotal() uint64
	QueryNsTotal() uint64
}

type BatchObserver interface {
	gocql.BatchObserver
}

type ConnectObserver interface {
	gocql.ConnectObserver
}

type FrameHeaderObserver interface {
	gocql.FrameHeaderObserver
}

type queryObserver struct {
	queryTotal   uint64 // completed query count total
	queryNsTotal uint64 // query time total
}

// NewQueryObserver returns a new QueryObserver instance.
func NewQueryObserver() QueryObserver {
	return &queryObserver{}
}

// ObserveQuery updates the member variables by passed ObservedQuery instance.
func (qo *queryObserver) ObserveQuery(ctx context.Context, q gocql.ObservedQuery) {
	atomic.AddUint64(&qo.queryTotal, 1)
	atomic.AddUint64(&qo.queryNsTotal, uint64(q.End.Sub(q.Start)))
}

// CompletedQueryTotal returns the cumulative number of completed query during uptime.
func (qo *queryObserver) CompletedQueryTotal() uint64 {
	return atomic.LoadUint64(&qo.queryTotal)
}

// QueryNsTotal returns the total consumed time by query operations in nanosecond.
func (qo *queryObserver) QueryNsTotal() uint64 {
	return atomic.LoadUint64(&qo.queryNsTotal)
}
