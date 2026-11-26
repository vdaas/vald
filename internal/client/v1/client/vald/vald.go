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

package vald

import (
	"context"
	"reflect"
	"slices"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/observability/trace"
)

type Client interface {
	vald.Client
	GRPCClient() grpc.Client
	Start(context.Context) (<-chan error, error)
	Stop(context.Context) error
}

type client struct {
	c       grpc.Client
	vcCache ConnCache
	addrs   []string
}

type singleClient struct {
	vc vald.Client
}

const (
	apiName = "vald/internal/client/v1/client/vald"
)

func New(opts ...Option) (Client, error) {
	c := new(client)
	for _, opt := range append(defaultOptions, opts...) {
		err := opt(c)
		if err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}
	if c.c == nil {
		if c.addrs == nil {
			return nil, errors.ErrGRPCTargetAddrNotFound
		}
		c.c = grpc.New("Vald Client", grpc.WithAddrs(c.addrs...))
	}
	return c, nil
}

var nilSingleClient = &singleClient{} //nolint:gochecknoglobals // shared zero-allocation sentinel returned when ClientConn is nil

func NewFromConn(cc *grpc.ClientConn) Client {
	if cc == nil {
		return nilSingleClient
	}
	return &singleClient{vc: vald.NewValdClient(cc)}
}

func (c *client) Start(ctx context.Context) (<-chan error, error) {
	return c.c.StartConnectionMonitor(ctx)
}

func (c *client) Stop(ctx context.Context) error {
	return c.c.Close(ctx)
}

func (c *client) GRPCClient() grpc.Client {
	return c.c
}

func mergeCallOpts(copts, opts []grpc.CallOption) []grpc.CallOption {
	switch {
	case len(opts) == 0:
		return copts
	case len(copts) == 0:
		return opts
	default:
		return slices.Concat(copts, opts)
	}
}

func roundRobinCall[Req, Res any](
	ctx context.Context,
	c *client,
	rpcName string,
	in Req,
	fn func(vald.Client, context.Context, Req, ...grpc.CallOption) (Res, error),
	opts ...grpc.CallOption,
) (Res, error) {
	if trace.Enabled() {
		var span trace.Span
		ctx, span = trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+rpcName), apiName+"/"+rpcName)
		defer trace.End(span)
	}
	return grpc.RoundRobin(ctx, c.c, func(ictx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (Res, error) {
		return fn(c.vcCache.Get(conn), ictx, in, mergeCallOpts(copts, opts)...)
	})
}

func roundRobinStream[Res any](
	ctx context.Context,
	c *client,
	rpcName string,
	fn func(vald.Client, context.Context, ...grpc.CallOption) (Res, error),
	opts ...grpc.CallOption,
) (Res, error) {
	if trace.Enabled() {
		var span trace.Span
		ctx, span = trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+rpcName), apiName+"/"+rpcName)
		defer trace.End(span)
	}
	// the stream is opened with the method ctx: the per-attempt ctx is cancelled by backoff.Do on return.
	return grpc.RoundRobin(ctx, c.c, func(_ context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (Res, error) {
		return fn(c.vcCache.Get(conn), ctx, mergeCallOpts(copts, opts)...)
	})
}

func (c *client) Exists(
	ctx context.Context, in *payload.Object_ID, opts ...grpc.CallOption,
) (*payload.Object_ID, error) {
	return roundRobinCall(ctx, c, vald.ExistsRPCName, in, vald.Client.Exists, opts...)
}

func (c *client) Search(
	ctx context.Context, in *payload.Search_Request, opts ...grpc.CallOption,
) (*payload.Search_Response, error) {
	return roundRobinCall(ctx, c, vald.SearchRPCName, in, vald.Client.Search, opts...)
}

func (c *client) SearchByID(
	ctx context.Context, in *payload.Search_IDRequest, opts ...grpc.CallOption,
) (*payload.Search_Response, error) {
	return roundRobinCall(ctx, c, vald.SearchByIDRPCName, in, vald.Client.SearchByID, opts...)
}

func (c *client) StreamSearch(
	ctx context.Context, opts ...grpc.CallOption,
) (vald.Search_StreamSearchClient, error) {
	return roundRobinStream(ctx, c, vald.StreamSearchRPCName, vald.Client.StreamSearch, opts...)
}

func (c *client) StreamSearchByID(
	ctx context.Context, opts ...grpc.CallOption,
) (vald.Search_StreamSearchByIDClient, error) {
	return roundRobinStream(ctx, c, vald.StreamSearchByIDRPCName, vald.Client.StreamSearchByID, opts...)
}

func (c *client) MultiSearch(
	ctx context.Context, in *payload.Search_MultiRequest, opts ...grpc.CallOption,
) (*payload.Search_Responses, error) {
	return roundRobinCall(ctx, c, vald.MultiSearchRPCName, in, vald.Client.MultiSearch, opts...)
}

func (c *client) MultiSearchByID(
	ctx context.Context, in *payload.Search_MultiIDRequest, opts ...grpc.CallOption,
) (*payload.Search_Responses, error) {
	return roundRobinCall(ctx, c, vald.MultiSearchByIDRPCName, in, vald.Client.MultiSearchByID, opts...)
}

func (c *client) LinearSearch(
	ctx context.Context, in *payload.Search_Request, opts ...grpc.CallOption,
) (*payload.Search_Response, error) {
	return roundRobinCall(ctx, c, vald.LinearSearchRPCName, in, vald.Client.LinearSearch, opts...)
}

func (c *client) LinearSearchByID(
	ctx context.Context, in *payload.Search_IDRequest, opts ...grpc.CallOption,
) (*payload.Search_Response, error) {
	return roundRobinCall(ctx, c, vald.LinearSearchByIDRPCName, in, vald.Client.LinearSearchByID, opts...)
}

func (c *client) StreamLinearSearch(
	ctx context.Context, opts ...grpc.CallOption,
) (vald.Search_StreamLinearSearchClient, error) {
	return roundRobinStream(ctx, c, vald.StreamLinearSearchRPCName, vald.Client.StreamLinearSearch, opts...)
}

func (c *client) StreamLinearSearchByID(
	ctx context.Context, opts ...grpc.CallOption,
) (vald.Search_StreamLinearSearchByIDClient, error) {
	return roundRobinStream(ctx, c, vald.StreamLinearSearchByIDRPCName, vald.Client.StreamLinearSearchByID, opts...)
}

func (c *client) MultiLinearSearch(
	ctx context.Context, in *payload.Search_MultiRequest, opts ...grpc.CallOption,
) (*payload.Search_Responses, error) {
	return roundRobinCall(ctx, c, vald.MultiLinearSearchRPCName, in, vald.Client.MultiLinearSearch, opts...)
}

func (c *client) MultiLinearSearchByID(
	ctx context.Context, in *payload.Search_MultiIDRequest, opts ...grpc.CallOption,
) (*payload.Search_Responses, error) {
	return roundRobinCall(ctx, c, vald.MultiLinearSearchByIDRPCName, in, vald.Client.MultiLinearSearchByID, opts...)
}

func (c *client) Insert(
	ctx context.Context, in *payload.Insert_Request, opts ...grpc.CallOption,
) (*payload.Object_Location, error) {
	return roundRobinCall(ctx, c, vald.InsertRPCName, in, vald.Client.Insert, opts...)
}

func (c *client) StreamInsert(
	ctx context.Context, opts ...grpc.CallOption,
) (vald.Insert_StreamInsertClient, error) {
	return roundRobinStream(ctx, c, vald.StreamInsertRPCName, vald.Client.StreamInsert, opts...)
}

func (c *client) MultiInsert(
	ctx context.Context, in *payload.Insert_MultiRequest, opts ...grpc.CallOption,
) (*payload.Object_Locations, error) {
	return roundRobinCall(ctx, c, vald.MultiInsertRPCName, in, vald.Client.MultiInsert, opts...)
}

func (c *client) Update(
	ctx context.Context, in *payload.Update_Request, opts ...grpc.CallOption,
) (*payload.Object_Location, error) {
	return roundRobinCall(ctx, c, vald.UpdateRPCName, in, vald.Client.Update, opts...)
}

func (c *client) StreamUpdate(
	ctx context.Context, opts ...grpc.CallOption,
) (vald.Update_StreamUpdateClient, error) {
	return roundRobinStream(ctx, c, vald.StreamUpdateRPCName, vald.Client.StreamUpdate, opts...)
}

func (c *client) MultiUpdate(
	ctx context.Context, in *payload.Update_MultiRequest, opts ...grpc.CallOption,
) (*payload.Object_Locations, error) {
	return roundRobinCall(ctx, c, vald.MultiUpdateRPCName, in, vald.Client.MultiUpdate, opts...)
}

func (c *client) UpdateTimestamp(
	ctx context.Context, in *payload.Update_TimestampRequest, opts ...grpc.CallOption,
) (*payload.Object_Location, error) {
	return roundRobinCall(ctx, c, vald.UpdateTimestampRPCName, in, vald.Client.UpdateTimestamp, opts...)
}

func (c *client) Upsert(
	ctx context.Context, in *payload.Upsert_Request, opts ...grpc.CallOption,
) (*payload.Object_Location, error) {
	return roundRobinCall(ctx, c, vald.UpsertRPCName, in, vald.Client.Upsert, opts...)
}

func (c *client) StreamUpsert(
	ctx context.Context, opts ...grpc.CallOption,
) (vald.Upsert_StreamUpsertClient, error) {
	return roundRobinStream(ctx, c, vald.StreamUpsertRPCName, vald.Client.StreamUpsert, opts...)
}

func (c *client) MultiUpsert(
	ctx context.Context, in *payload.Upsert_MultiRequest, opts ...grpc.CallOption,
) (*payload.Object_Locations, error) {
	return roundRobinCall(ctx, c, vald.MultiUpsertRPCName, in, vald.Client.MultiUpsert, opts...)
}

func (c *client) Remove(
	ctx context.Context, in *payload.Remove_Request, opts ...grpc.CallOption,
) (*payload.Object_Location, error) {
	return roundRobinCall(ctx, c, vald.RemoveRPCName, in, vald.Client.Remove, opts...)
}

func (c *client) StreamRemove(
	ctx context.Context, opts ...grpc.CallOption,
) (vald.Remove_StreamRemoveClient, error) {
	return roundRobinStream(ctx, c, vald.StreamRemoveRPCName, vald.Client.StreamRemove, opts...)
}

func (c *client) MultiRemove(
	ctx context.Context, in *payload.Remove_MultiRequest, opts ...grpc.CallOption,
) (*payload.Object_Locations, error) {
	return roundRobinCall(ctx, c, vald.MultiRemoveRPCName, in, vald.Client.MultiRemove, opts...)
}

func (c *client) Flush(
	ctx context.Context, in *payload.Flush_Request, opts ...grpc.CallOption,
) (*payload.Info_Index_Count, error) {
	return roundRobinCall(ctx, c, vald.FlushRPCName, in, vald.Client.Flush, opts...)
}

func (c *client) RemoveByTimestamp(
	ctx context.Context, in *payload.Remove_TimestampRequest, opts ...grpc.CallOption,
) (*payload.Object_Locations, error) {
	return roundRobinCall(ctx, c, vald.RemoveByTimestampRPCName, in, vald.Client.RemoveByTimestamp, opts...)
}

func (c *client) GetObject(
	ctx context.Context, in *payload.Object_VectorRequest, opts ...grpc.CallOption,
) (*payload.Object_Vector, error) {
	return roundRobinCall(ctx, c, vald.GetObjectRPCName, in, vald.Client.GetObject, opts...)
}

func (c *client) StreamGetObject(
	ctx context.Context, opts ...grpc.CallOption,
) (vald.Object_StreamGetObjectClient, error) {
	return roundRobinStream(ctx, c, vald.StreamGetObjectRPCName, vald.Client.StreamGetObject, opts...)
}

func (c *client) StreamListObject(
	ctx context.Context, in *payload.Object_List_Request, opts ...grpc.CallOption,
) (vald.Object_StreamListObjectClient, error) {
	if trace.Enabled() {
		var span trace.Span
		ctx, span = trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.StreamListObjectRPCName), apiName+"/"+vald.StreamListObjectRPCName)
		defer trace.End(span)
	}
	// the stream is opened with the method ctx: the per-attempt ctx is cancelled by backoff.Do on return.
	return grpc.RoundRobin(ctx, c.c, func(_ context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (vald.Object_StreamListObjectClient, error) {
		return c.vcCache.Get(conn).StreamListObject(ctx, in, mergeCallOpts(copts, opts)...)
	})
}

func (c *client) IndexInfo(
	ctx context.Context, in *payload.Empty, opts ...grpc.CallOption,
) (*payload.Info_Index_Count, error) {
	return roundRobinCall(ctx, c, vald.IndexInfoRPCName, in, vald.Client.IndexInfo, opts...)
}

func (c *client) IndexDetail(
	ctx context.Context, in *payload.Empty, opts ...grpc.CallOption,
) (*payload.Info_Index_Detail, error) {
	return roundRobinCall(ctx, c, vald.IndexDetailRPCName, in, vald.Client.IndexDetail, opts...)
}

func (c *client) IndexStatistics(
	ctx context.Context, in *payload.Empty, opts ...grpc.CallOption,
) (*payload.Info_Index_Statistics, error) {
	return roundRobinCall(ctx, c, vald.IndexStatisticsRPCName, in, vald.Client.IndexStatistics, opts...)
}

func (c *client) IndexStatisticsDetail(
	ctx context.Context, in *payload.Empty, opts ...grpc.CallOption,
) (*payload.Info_Index_StatisticsDetail, error) {
	return roundRobinCall(ctx, c, vald.IndexStatisticsDetailRPCName, in, vald.Client.IndexStatisticsDetail, opts...)
}

func (c *client) IndexProperty(
	ctx context.Context, in *payload.Empty, opts ...grpc.CallOption,
) (*payload.Info_Index_PropertyDetail, error) {
	return roundRobinCall(ctx, c, vald.IndexPropertyRPCName, in, vald.Client.IndexProperty, opts...)
}

func (c *client) GetTimestamp(
	ctx context.Context, in *payload.Object_TimestampRequest, opts ...grpc.CallOption,
) (*payload.Object_Timestamp, error) {
	return roundRobinCall(ctx, c, vald.GetTimestampRPCName, in, vald.Client.GetTimestamp, opts...)
}

func (*singleClient) Start(context.Context) (<-chan error, error) {
	return nil, nil
}

func (*singleClient) Stop(context.Context) error {
	return nil
}

func (*singleClient) GRPCClient() grpc.Client {
	return nil
}

func singleCall[Req, Res any](
	ctx context.Context,
	vc vald.Client,
	rpcName string,
	in Req,
	fn func(vald.Client, context.Context, Req, ...grpc.CallOption) (Res, error),
	opts ...grpc.CallOption,
) (res Res, err error) {
	if vc == nil {
		return res, errors.ErrGRPCClientConnNotFound("*")
	}
	if !trace.Enabled() {
		return fn(vc, ctx, in, opts...)
	}
	return singleCallSlow(ctx, vc, rpcName, in, fn, opts...)
}

func singleCallSlow[Req, Res any](
	ctx context.Context,
	vc vald.Client,
	rpcName string,
	in Req,
	fn func(vald.Client, context.Context, Req, ...grpc.CallOption) (Res, error),
	opts ...grpc.CallOption,
) (Res, error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+rpcName), apiName+"/"+rpcName)
	defer trace.End(span)
	return fn(vc, ctx, in, opts...)
}

func singleStream[Res any](
	ctx context.Context,
	vc vald.Client,
	rpcName string,
	fn func(vald.Client, context.Context, ...grpc.CallOption) (Res, error),
	opts ...grpc.CallOption,
) (res Res, err error) {
	if vc == nil {
		return res, errors.ErrGRPCClientConnNotFound("*")
	}
	if !trace.Enabled() {
		return fn(vc, ctx, opts...)
	}
	return singleStreamSlow(ctx, vc, rpcName, fn, opts...)
}

func singleStreamSlow[Res any](
	ctx context.Context,
	vc vald.Client,
	rpcName string,
	fn func(vald.Client, context.Context, ...grpc.CallOption) (Res, error),
	opts ...grpc.CallOption,
) (Res, error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+rpcName), apiName+"/"+rpcName)
	defer trace.End(span)
	return fn(vc, ctx, opts...)
}

func (c *singleClient) Exists(
	ctx context.Context, in *payload.Object_ID, opts ...grpc.CallOption,
) (*payload.Object_ID, error) {
	return singleCall(ctx, c.vc, vald.ExistsRPCName, in, vald.Client.Exists, opts...)
}

func (c *singleClient) Search(
	ctx context.Context, in *payload.Search_Request, opts ...grpc.CallOption,
) (*payload.Search_Response, error) {
	return singleCall(ctx, c.vc, vald.SearchRPCName, in, vald.Client.Search, opts...)
}

func (c *singleClient) SearchByID(
	ctx context.Context, in *payload.Search_IDRequest, opts ...grpc.CallOption,
) (*payload.Search_Response, error) {
	return singleCall(ctx, c.vc, vald.SearchByIDRPCName, in, vald.Client.SearchByID, opts...)
}

func (c *singleClient) StreamSearch(
	ctx context.Context, opts ...grpc.CallOption,
) (vald.Search_StreamSearchClient, error) {
	return singleStream(ctx, c.vc, vald.StreamSearchRPCName, vald.Client.StreamSearch, opts...)
}

func (c *singleClient) StreamSearchByID(
	ctx context.Context, opts ...grpc.CallOption,
) (vald.Search_StreamSearchByIDClient, error) {
	return singleStream(ctx, c.vc, vald.StreamSearchByIDRPCName, vald.Client.StreamSearchByID, opts...)
}

func (c *singleClient) MultiSearch(
	ctx context.Context, in *payload.Search_MultiRequest, opts ...grpc.CallOption,
) (*payload.Search_Responses, error) {
	return singleCall(ctx, c.vc, vald.MultiSearchRPCName, in, vald.Client.MultiSearch, opts...)
}

func (c *singleClient) MultiSearchByID(
	ctx context.Context, in *payload.Search_MultiIDRequest, opts ...grpc.CallOption,
) (*payload.Search_Responses, error) {
	return singleCall(ctx, c.vc, vald.MultiSearchByIDRPCName, in, vald.Client.MultiSearchByID, opts...)
}

func (c *singleClient) LinearSearch(
	ctx context.Context, in *payload.Search_Request, opts ...grpc.CallOption,
) (*payload.Search_Response, error) {
	return singleCall(ctx, c.vc, vald.LinearSearchRPCName, in, vald.Client.LinearSearch, opts...)
}

func (c *singleClient) LinearSearchByID(
	ctx context.Context, in *payload.Search_IDRequest, opts ...grpc.CallOption,
) (*payload.Search_Response, error) {
	return singleCall(ctx, c.vc, vald.LinearSearchByIDRPCName, in, vald.Client.LinearSearchByID, opts...)
}

func (c *singleClient) StreamLinearSearch(
	ctx context.Context, opts ...grpc.CallOption,
) (vald.Search_StreamLinearSearchClient, error) {
	return singleStream(ctx, c.vc, vald.StreamLinearSearchRPCName, vald.Client.StreamLinearSearch, opts...)
}

func (c *singleClient) StreamLinearSearchByID(
	ctx context.Context, opts ...grpc.CallOption,
) (vald.Search_StreamLinearSearchByIDClient, error) {
	return singleStream(ctx, c.vc, vald.StreamLinearSearchByIDRPCName, vald.Client.StreamLinearSearchByID, opts...)
}

func (c *singleClient) MultiLinearSearch(
	ctx context.Context, in *payload.Search_MultiRequest, opts ...grpc.CallOption,
) (*payload.Search_Responses, error) {
	return singleCall(ctx, c.vc, vald.MultiLinearSearchRPCName, in, vald.Client.MultiLinearSearch, opts...)
}

func (c *singleClient) MultiLinearSearchByID(
	ctx context.Context, in *payload.Search_MultiIDRequest, opts ...grpc.CallOption,
) (*payload.Search_Responses, error) {
	return singleCall(ctx, c.vc, vald.MultiLinearSearchByIDRPCName, in, vald.Client.MultiLinearSearchByID, opts...)
}

func (c *singleClient) Insert(
	ctx context.Context, in *payload.Insert_Request, opts ...grpc.CallOption,
) (*payload.Object_Location, error) {
	return singleCall(ctx, c.vc, vald.InsertRPCName, in, vald.Client.Insert, opts...)
}

func (c *singleClient) StreamInsert(
	ctx context.Context, opts ...grpc.CallOption,
) (vald.Insert_StreamInsertClient, error) {
	return singleStream(ctx, c.vc, vald.StreamInsertRPCName, vald.Client.StreamInsert, opts...)
}

func (c *singleClient) MultiInsert(
	ctx context.Context, in *payload.Insert_MultiRequest, opts ...grpc.CallOption,
) (*payload.Object_Locations, error) {
	return singleCall(ctx, c.vc, vald.MultiInsertRPCName, in, vald.Client.MultiInsert, opts...)
}

func (c *singleClient) Update(
	ctx context.Context, in *payload.Update_Request, opts ...grpc.CallOption,
) (*payload.Object_Location, error) {
	return singleCall(ctx, c.vc, vald.UpdateRPCName, in, vald.Client.Update, opts...)
}

func (c *singleClient) UpdateTimestamp(
	ctx context.Context, in *payload.Update_TimestampRequest, opts ...grpc.CallOption,
) (*payload.Object_Location, error) {
	return singleCall(ctx, c.vc, vald.UpdateTimestampRPCName, in, vald.Client.UpdateTimestamp, opts...)
}

func (c *singleClient) StreamUpdate(
	ctx context.Context, opts ...grpc.CallOption,
) (vald.Update_StreamUpdateClient, error) {
	return singleStream(ctx, c.vc, vald.StreamUpdateRPCName, vald.Client.StreamUpdate, opts...)
}

func (c *singleClient) MultiUpdate(
	ctx context.Context, in *payload.Update_MultiRequest, opts ...grpc.CallOption,
) (*payload.Object_Locations, error) {
	return singleCall(ctx, c.vc, vald.MultiUpdateRPCName, in, vald.Client.MultiUpdate, opts...)
}

func (c *singleClient) Upsert(
	ctx context.Context, in *payload.Upsert_Request, opts ...grpc.CallOption,
) (*payload.Object_Location, error) {
	return singleCall(ctx, c.vc, vald.UpsertRPCName, in, vald.Client.Upsert, opts...)
}

func (c *singleClient) StreamUpsert(
	ctx context.Context, opts ...grpc.CallOption,
) (vald.Upsert_StreamUpsertClient, error) {
	return singleStream(ctx, c.vc, vald.StreamUpsertRPCName, vald.Client.StreamUpsert, opts...)
}

func (c *singleClient) MultiUpsert(
	ctx context.Context, in *payload.Upsert_MultiRequest, opts ...grpc.CallOption,
) (*payload.Object_Locations, error) {
	return singleCall(ctx, c.vc, vald.MultiUpsertRPCName, in, vald.Client.MultiUpsert, opts...)
}

func (c *singleClient) Remove(
	ctx context.Context, in *payload.Remove_Request, opts ...grpc.CallOption,
) (*payload.Object_Location, error) {
	return singleCall(ctx, c.vc, vald.RemoveRPCName, in, vald.Client.Remove, opts...)
}

func (c *singleClient) StreamRemove(
	ctx context.Context, opts ...grpc.CallOption,
) (vald.Remove_StreamRemoveClient, error) {
	return singleStream(ctx, c.vc, vald.StreamRemoveRPCName, vald.Client.StreamRemove, opts...)
}

func (c *singleClient) MultiRemove(
	ctx context.Context, in *payload.Remove_MultiRequest, opts ...grpc.CallOption,
) (*payload.Object_Locations, error) {
	return singleCall(ctx, c.vc, vald.MultiRemoveRPCName, in, vald.Client.MultiRemove, opts...)
}

func (c *singleClient) Flush(
	ctx context.Context, in *payload.Flush_Request, opts ...grpc.CallOption,
) (*payload.Info_Index_Count, error) {
	return singleCall(ctx, c.vc, vald.FlushRPCName, in, vald.Client.Flush, opts...)
}

func (c *singleClient) RemoveByTimestamp(
	ctx context.Context, in *payload.Remove_TimestampRequest, opts ...grpc.CallOption,
) (*payload.Object_Locations, error) {
	return singleCall(ctx, c.vc, vald.RemoveByTimestampRPCName, in, vald.Client.RemoveByTimestamp, opts...)
}

func (c *singleClient) GetObject(
	ctx context.Context, in *payload.Object_VectorRequest, opts ...grpc.CallOption,
) (*payload.Object_Vector, error) {
	return singleCall(ctx, c.vc, vald.GetObjectRPCName, in, vald.Client.GetObject, opts...)
}

func (c *singleClient) StreamGetObject(
	ctx context.Context, opts ...grpc.CallOption,
) (vald.Object_StreamGetObjectClient, error) {
	return singleStream(ctx, c.vc, vald.StreamGetObjectRPCName, vald.Client.StreamGetObject, opts...)
}

func (c *singleClient) StreamListObject(
	ctx context.Context, in *payload.Object_List_Request, opts ...grpc.CallOption,
) (vald.Object_StreamListObjectClient, error) {
	return singleCall(ctx, c.vc, vald.StreamListObjectRPCName, in, vald.Client.StreamListObject, opts...)
}

func (c *singleClient) IndexInfo(
	ctx context.Context, in *payload.Empty, opts ...grpc.CallOption,
) (*payload.Info_Index_Count, error) {
	return singleCall(ctx, c.vc, vald.IndexInfoRPCName, in, vald.Client.IndexInfo, opts...)
}

func (c *singleClient) IndexDetail(
	ctx context.Context, in *payload.Empty, opts ...grpc.CallOption,
) (*payload.Info_Index_Detail, error) {
	return singleCall(ctx, c.vc, vald.IndexDetailRPCName, in, vald.Client.IndexDetail, opts...)
}

func (c *singleClient) IndexStatistics(
	ctx context.Context, in *payload.Empty, opts ...grpc.CallOption,
) (*payload.Info_Index_Statistics, error) {
	return singleCall(ctx, c.vc, vald.IndexStatisticsRPCName, in, vald.Client.IndexStatistics, opts...)
}

func (c *singleClient) IndexStatisticsDetail(
	ctx context.Context, in *payload.Empty, opts ...grpc.CallOption,
) (*payload.Info_Index_StatisticsDetail, error) {
	return singleCall(ctx, c.vc, vald.IndexStatisticsDetailRPCName, in, vald.Client.IndexStatisticsDetail, opts...)
}

func (c *singleClient) IndexProperty(
	ctx context.Context, in *payload.Empty, opts ...grpc.CallOption,
) (*payload.Info_Index_PropertyDetail, error) {
	return singleCall(ctx, c.vc, vald.IndexPropertyRPCName, in, vald.Client.IndexProperty, opts...)
}

func (c *singleClient) GetTimestamp(
	ctx context.Context, in *payload.Object_TimestampRequest, opts ...grpc.CallOption,
) (*payload.Object_Timestamp, error) {
	return singleCall(ctx, c.vc, vald.GetTimestampRPCName, in, vald.Client.GetTimestamp, opts...)
}
