package operation

import (
	"context"
	"io"
	"reflect"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/agent/core"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type client struct {
	host string
	port int
}

type Dataset struct {
	Train     [][]float32
	Test      [][]float32
	Neighbors [][]int
}

type Client interface {
	Search(t *testing.T, ctx context.Context, ds Dataset) error
	SearchByID(t *testing.T, ctx context.Context, ds Dataset) error
	Insert(t *testing.T, ctx context.Context, ds Dataset) error
	Update(t *testing.T, ctx context.Context, ds Dataset) error
	Remove(t *testing.T, ctx context.Context, ds Dataset) error
	GetObject(t *testing.T, ctx context.Context, ds Dataset) error
	CreateIndex(t *testing.T, ctx context.Context) error
	SaveIndex(t *testing.T, ctx context.Context) error
	IndexInfo(t *testing.T, ctx context.Context) (*payload.Info_Index_Count, error)
}

func New(host string, port int) (Client, error) {
	return &client{
		host: host,
		port: port,
	}, nil
}

func (c *client) Search(t *testing.T, ctx context.Context, ds Dataset) error {
	t.Log("search operation started")

	client, err := c.getClient(ctx)
	if err != nil {
		return err
	}

	sc, err := client.StreamSearch(ctx)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			res, err := sc.Recv()
			if err == io.EOF {
				return
			}

			if err != nil {
				st, serr := status.FromError(err)
				t.Errorf("error: %v\tserror: %v\tcode: %s\tdetails: %s\tmessage: %s\tstatus-error: %s\tproto: %s", err, serr, st.Code().String(), errdetails.Serialize(st.Details()), st.Message(), st.Err().Error(), errdetails.Serialize(st.Proto()))
				continue
			}

			resp := res.GetResponse()
			if resp == nil {
				err := res.GetStatus()
				if err != nil {
					t.Errorf("an error returned:\tcode: %d\tmessage: %s\tdetails: %s", err.GetCode(), err.GetMessage(), errdetails.Serialize(err.GetDetails()))
					continue
				}

				t.Error("returned response is nil")
				continue
			}

			topKIDs := make([]string, 0, len(resp.GetResults()))
			for _, d := range resp.GetResults() {
				topKIDs = append(topKIDs, d.Id)
			}

			if len(topKIDs) == 0 {
				t.Errorf("empty result is returned for test ID %s: %#v", resp.GetRequestId(), topKIDs)
				continue
			}

			idx, err := strconv.Atoi(resp.GetRequestId())
			if err != nil {
				t.Errorf("an error occurred while converting RequestId into int: %s", err)
				continue
			}

			t.Logf("results: %d, recall: %f", len(topKIDs), c.recall(topKIDs, ds.Neighbors[idx][:len(topKIDs)]))
		}
	}()

	for i := 0; i < len(ds.Test); i++ {
		id := strconv.Itoa(i)
		err := sc.Send(&payload.Search_Request{
			Vector: ds.Test[i],
			Config: &payload.Search_Config{
				RequestId: id,
				Num:       100,
				Radius:    -1.0,
				Epsilon:   0.01,
				Timeout:   3000000000,
			},
		})
		if err != nil {
			return err
		}
	}

	sc.CloseSend()

	wg.Wait()

	t.Log("search operation finished")

	return nil
}

func (c *client) SearchByID(t *testing.T, ctx context.Context, ds Dataset) error {
	t.Log("searchByID operation started")

	client, err := c.getClient(ctx)
	if err != nil {
		return err
	}

	sc, err := client.StreamSearchByID(ctx)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			res, err := sc.Recv()
			if err == io.EOF {
				return
			}

			if err != nil {
				st, serr := status.FromError(err)
				t.Errorf("error: %v\tserror: %v\tcode: %s\tdetails: %s\tmessage: %s\tstatus-error: %s\tproto: %s", err, serr, st.Code().String(), errdetails.Serialize(st.Details()), st.Message(), st.Err().Error(), errdetails.Serialize(st.Proto()))
				continue
			}

			resp := res.GetResponse()
			if resp == nil {
				err := res.GetStatus()
				if err != nil {
					t.Errorf("an error returned:\tcode: %d\tmessage: %s\tdetails: %s", err.GetCode(), err.GetMessage(), errdetails.Serialize(err.GetDetails()))
					continue
				}

				t.Error("returned response is nil")
			}

			topKIDs := make([]string, 0, len(resp.GetResults()))
			for _, d := range resp.GetResults() {
				topKIDs = append(topKIDs, d.Id)
			}

			if len(topKIDs) == 0 {
				t.Errorf("empty result is returned: %#v", topKIDs)
			}
		}
	}()

	for i := 0; i < len(ds.Train); i++ {
		id := strconv.Itoa(i)
		err := sc.Send(&payload.Search_IDRequest{
			Id: id,
			Config: &payload.Search_Config{
				RequestId: id,
				Num:       100,
				Radius:    -1.0,
				Epsilon:   0.01,
				Timeout:   3000000000,
			},
		})
		if err != nil {
			return err
		}
	}

	sc.CloseSend()

	wg.Wait()

	t.Log("searchByID operation finished")

	return nil
}

func (c *client) Insert(t *testing.T, ctx context.Context, ds Dataset) error {
	t.Log("insert operation started")

	client, err := c.getClient(ctx)
	if err != nil {
		return err
	}

	sc, err := client.StreamInsert(ctx)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			res, err := sc.Recv()
			if err == io.EOF {
				return
			}

			if err != nil {
				st, serr := status.FromError(err)
				t.Errorf("error: %v\tserror: %v\tcode: %s\tdetails: %s\tmessage: %s\tstatus-error: %s\tproto: %s", err, serr, st.Code().String(), errdetails.Serialize(st.Details()), st.Message(), st.Err().Error(), errdetails.Serialize(st.Proto()))
				continue
			}

			loc := res.GetLocation()
			if loc == nil {
				err := res.GetStatus()
				if err != nil {
					t.Errorf("an error returned:\tcode: %d\tmessage: %s\tdetails: %s", err.GetCode(), err.GetMessage(), errdetails.Serialize(err.GetDetails()))
					continue
				}

				t.Error("returned loc is nil")
				continue
			}

			t.Logf("returned loc: %s", loc)
		}
	}()

	for i := 0; i < len(ds.Train); i++ {
		id := strconv.Itoa(i)
		err := sc.Send(&payload.Insert_Request{
			Vector: &payload.Object_Vector{
				Id:     id,
				Vector: ds.Train[i],
			},
			Config: &payload.Insert_Config{
				SkipStrictExistCheck: false,
			},
		})
		if err != nil {
			return err
		}
	}

	sc.CloseSend()

	wg.Wait()

	t.Log("insert operation finished")

	return nil
}

func (c *client) Update(t *testing.T, ctx context.Context, ds Dataset) error {
	t.Log("update operation started")

	client, err := c.getClient(ctx)
	if err != nil {
		return err
	}

	sc, err := client.StreamUpdate(ctx)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			res, err := sc.Recv()
			if err == io.EOF {
				return
			}

			if err != nil {
				st, serr := status.FromError(err)
				t.Errorf("error: %v\tserror: %v\tcode: %s\tdetails: %s\tmessage: %s\tstatus-error: %s\tproto: %s", err, serr, st.Code().String(), errdetails.Serialize(st.Details()), st.Message(), st.Err().Error(), errdetails.Serialize(st.Proto()))
				continue
			}

			loc := res.GetLocation()
			if loc == nil {
				err := res.GetStatus()
				if err != nil {
					t.Errorf("an error returned:\tcode: %d\tmessage: %s\tdetails: %s", err.GetCode(), err.GetMessage(), errdetails.Serialize(err.GetDetails()))
					continue
				}

				t.Error("returned loc is nil")
			}

			t.Logf("returned: %s", loc)
		}
	}()

	for i := 0; i < len(ds.Train); i++ {
		id := strconv.Itoa(i)
		v := ds.Train[i]
		err := sc.Send(&payload.Update_Request{
			Vector: &payload.Object_Vector{
				Id:     id,
				Vector: append(v[1:], v[0]),
			},
			Config: &payload.Update_Config{
				SkipStrictExistCheck: false,
			},
		})
		if err != nil {
			return err
		}
	}

	sc.CloseSend()

	wg.Wait()

	t.Log("update operation finished")

	return nil
}

func (c *client) Remove(t *testing.T, ctx context.Context, ds Dataset) error {
	t.Log("remove operation started")

	client, err := c.getClient(ctx)
	if err != nil {
		return err
	}

	sc, err := client.StreamRemove(ctx)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			res, err := sc.Recv()
			if err == io.EOF {
				return
			}

			if err != nil {
				st, serr := status.FromError(err)
				t.Errorf("error: %v\tserror: %v\tcode: %s\tdetails: %s\tmessage: %s\tstatus-error: %s\tproto: %s", err, serr, st.Code().String(), errdetails.Serialize(st.Details()), st.Message(), st.Err().Error(), errdetails.Serialize(st.Proto()))
				continue
			}

			loc := res.GetLocation()
			if loc == nil {
				err := res.GetStatus()
				if err != nil {
					t.Errorf("an error returned:\tcode: %d\tmessage: %s\tdetails: %s", err.GetCode(), err.GetMessage(), errdetails.Serialize(err.GetDetails()))
					continue
				}

				t.Error("returned response is nil")
				continue
			}

			t.Logf("returned: %s", loc)
		}
	}()

	for i := 0; i < len(ds.Train); i++ {
		id := strconv.Itoa(i)
		err := sc.Send(&payload.Remove_Request{
			Id: &payload.Object_ID{
				Id: id,
			},
			Config: &payload.Remove_Config{
				SkipStrictExistCheck: false,
			},
		})
		if err != nil {
			return err
		}
	}

	sc.CloseSend()

	wg.Wait()

	t.Log("remove operation finished")

	return nil
}

func (c *client) GetObject(t *testing.T, ctx context.Context, ds Dataset) error {
	t.Log("getObject operation started")

	client, err := c.getClient(ctx)
	if err != nil {
		return err
	}

	sc, err := client.StreamGetObject(ctx)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			res, err := sc.Recv()
			if err == io.EOF {
				return
			}

			if err != nil {
				st, serr := status.FromError(err)
				t.Errorf("error: %v\tserror: %v\tcode: %s\tdetails: %s\tmessage: %s\tstatus-error: %s\tproto: %s", err, serr, st.Code().String(), errdetails.Serialize(st.Details()), st.Message(), st.Err().Error(), errdetails.Serialize(st.Proto()))
				continue
			}

			resp := res.GetVector()
			if resp == nil {
				err := res.GetStatus()
				if err != nil {
					t.Errorf("an error returned:\tcode: %d\tmessage: %s\tdetails: %s", err.GetCode(), err.GetMessage(), errdetails.Serialize(err.GetDetails()))
					continue
				}

				t.Error("returned response is nil")
				continue
			}

			idx, err := strconv.Atoi(resp.GetId())
			if err != nil {
				t.Errorf("an error occurred while converting Id into int: %s", err)
				continue
			}

			if !reflect.DeepEqual(res.GetVector().GetVector(), ds.Train[idx]) {
				t.Errorf(
					"got: %#v, expected: %#v",
					res.GetVector().GetVector(),
					ds.Train[idx],
				)
			}
		}
	}()

	for i := 0; i < len(ds.Train); i++ {
		id := strconv.Itoa(i)
		err := sc.Send(&payload.Object_VectorRequest{
			Id: &payload.Object_ID{
				Id: id,
			},
		})
		if err != nil {
			return err
		}
	}

	sc.CloseSend()

	wg.Wait()

	t.Log("getObject operation finished")

	return nil
}

func (c *client) CreateIndex(t *testing.T, ctx context.Context) error {
	client, err := c.getAgentClient(ctx)
	if err != nil {
		return err
	}

	_, err = client.CreateIndex(ctx, &payload.Control_CreateIndexRequest{
		PoolSize: 10000,
	})

	return err
}

func (c *client) SaveIndex(t *testing.T, ctx context.Context) error {
	client, err := c.getAgentClient(ctx)
	if err != nil {
		return err
	}

	_, err = client.SaveIndex(ctx, &payload.Empty{})

	return err
}

func (c *client) IndexInfo(t *testing.T, ctx context.Context) (*payload.Info_Index_Count, error) {
	client, err := c.getAgentClient(ctx)
	if err != nil {
		return nil, err
	}

	return client.IndexInfo(ctx, &payload.Empty{})
}

func (c *client) getGRPCConn(ctx context.Context) (*grpc.ClientConn, error) {
	return grpc.DialContext(
		ctx,
		c.host+":"+strconv.Itoa(c.port),
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(
			keepalive.ClientParameters{
				Time:                time.Second,
				Timeout:             5 * time.Second,
				PermitWithoutStream: true,
			},
		),
	)
}

func (c *client) getClient(ctx context.Context) (vald.Client, error) {
	conn, err := c.getGRPCConn(ctx)
	if err != nil {
		return nil, err
	}

	return vald.NewValdClient(conn), nil
}

func (c *client) getAgentClient(ctx context.Context) (core.AgentClient, error) {
	conn, err := c.getGRPCConn(ctx)
	if err != nil {
		return nil, err
	}

	return core.NewAgentClient(conn), nil
}

func (c *client) recall(results []string, neighbors []int) (recall float64) {
	ns := map[string]struct{}{}
	for _, n := range neighbors {
		ns[strconv.Itoa(n)] = struct{}{}
	}

	for _, r := range results {
		if _, ok := ns[r]; ok {
			recall++
		}
	}

	return recall / float64(len(neighbors))
}
