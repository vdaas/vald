package agent

import (
	"context"

	"github.com/vdaas/vald/apis/grpc/agent"
	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/internal/client"
	igrpc "github.com/vdaas/vald/internal/net/grpc"
	"google.golang.org/grpc"
)

type Client interface {
	client.Client
	client.ObjectReader
	client.Indexer
}

type agentClient struct {
	streamConcurrency int
	agent.AgentClient
}

func New(ctx context.Context, addr string) (Client, error) {
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &agentClient{
		AgentClient: agent.NewAgentClient(conn),
	}, nil
}

func (c *agentClient) Exists(ctx context.Context, objectID *client.ObjectID) (*client.ObjectID, error) {
	return c.AgentClient.Exists(ctx, objectID)
}

func (c *agentClient) Search(ctx context.Context, searchRequest *client.SearchRequest) (*client.SearchResponse, error) {
	return c.AgentClient.Search(ctx, searchRequest)
}

func (c *agentClient) SearchByID(ctx context.Context, searchIDRequest *client.SearchIDRequest) (*client.SearchResponse, error) {
	return c.AgentClient.SearchByID(ctx, searchIDRequest)
}

func (c *agentClient) StreamSearch(ctx context.Context, newData func() *client.SearchRequest, f func(*client.SearchResponse, error)) (err error) {
	var st agent.Agent_StreamSearchClient

	st, err = c.AgentClient.StreamSearch(ctx)
	if err != nil {
		return err
	}

	defer func() {
		err = st.CloseSend()
	}()

	return igrpc.BidirectionalStreamClient(st, c.streamConcurrency, func() interface{} {
		return newData()
	}, func(res interface{}, err error) {
		f(res.(*client.SearchResponse), err)
	})
}

func (c *agentClient) StreamSearchByID(ctx context.Context, newData func() *client.SearchRequest, f func(*client.SearchResponse, error)) (err error) {
	var st agent.Agent_StreamSearchByIDClient

	st, err = c.AgentClient.StreamSearchByID(ctx)
	if err != nil {
		return err
	}

	defer func() {
		err = st.CloseSend()
	}()

	return igrpc.BidirectionalStreamClient(st, c.streamConcurrency, func() interface{} {
		return newData()
	}, func(res interface{}, err error) {
		f(res.(*client.SearchResponse), err)
	})
}

func (c *agentClient) Insert(ctx context.Context, objectVector *client.ObjectVector) error {
	_, err := c.AgentClient.Insert(ctx, objectVector)
	return err
}

func (c *agentClient) StreamInsert(ctx context.Context, newData func() *client.ObjectVector, f func(error)) (err error) {
	var st agent.Agent_StreamInsertClient

	st, err = c.AgentClient.StreamInsert(ctx)
	if err != nil {
		return err
	}

	defer func() {
		err = st.CloseSend()
	}()

	return igrpc.BidirectionalStreamClient(st, c.streamConcurrency, func() interface{} {
		return newData()
	}, func(_ interface{}, err error) {
		f(err)
	})
}

func (c *agentClient) MultiInsert(ctx context.Context, objectVectors *client.ObjectVectors) error {
	_, err := c.AgentClient.MultiInsert(ctx, objectVectors)
	return err
}

func (c *agentClient) Update(ctx context.Context, objectVector *client.ObjectVector) error {
	_, err := c.AgentClient.Update(ctx, objectVector)
	return err
}

func (c *agentClient) StreamUpdate(ctx context.Context, newData func() *client.ObjectVector, f func(error)) (err error) {
	var st agent.Agent_StreamUpdateClient

	st, err = c.AgentClient.StreamUpdate(ctx)
	if err != nil {
		return err
	}

	defer func() {
		err = st.CloseSend()
	}()

	return igrpc.BidirectionalStreamClient(st, c.streamConcurrency, func() interface{} {
		return newData()
	}, func(_ interface{}, err error) {
		f(err)
	})
}

func (c *agentClient) MultiUpdate(ctx context.Context, objectVectors *client.ObjectVectors) error {
	_, err := c.AgentClient.MultiUpdate(ctx, objectVectors)
	return err
}

func (c *agentClient) Remove(ctx context.Context, objectID *client.ObjectID) error {
	_, err := c.AgentClient.Remove(ctx, objectID)
	return err
}

func (c *agentClient) StreamRemove(ctx context.Context, newData func() *client.ObjectID, f func(error)) (err error) {
	var st agent.Agent_StreamRemoveClient

	st, err = c.AgentClient.StreamRemove(ctx)
	if err != nil {
		return err
	}

	defer func() {
		err = st.CloseSend()
	}()

	return igrpc.BidirectionalStreamClient(st, c.streamConcurrency, func() interface{} {
		return newData()
	}, func(_ interface{}, err error) {
		f(err)
	})
}

func (c *agentClient) MultiRemove(ctx context.Context, objectIDs *client.ObjectIDs) error {
	_, err := c.AgentClient.MultiRemove(ctx, objectIDs)
	return err
}

func (c *agentClient) GetObject(ctx context.Context, objectID *client.ObjectID) (*client.ObjectVector, error) {
	return c.AgentClient.GetObject(ctx, objectID)
}

func (c *agentClient) StreamGetObject(ctx context.Context, newData func() *client.ObjectID, f func(*client.ObjectVector, error)) (err error) {
	var st agent.Agent_StreamGetObjectClient

	st, err = c.AgentClient.StreamGetObject(ctx)
	if err != nil {
		return err
	}

	defer func() {
		err = st.CloseSend()
	}()

	return igrpc.BidirectionalStreamClient(st, c.streamConcurrency, func() interface{} {
		return newData()
	}, func(res interface{}, err error) {
		f(res.(*client.ObjectVector), err)
	})
}

func (c *agentClient) CreateIndex(ctx context.Context, controlCreateIndexRequest *client.ControlCreateIndexRequest) error {
	_, err := c.AgentClient.CreateIndex(ctx, controlCreateIndexRequest)
	return err
}

func (c *agentClient) SaveIndex(ctx context.Context) error {
	_, err := c.AgentClient.SaveIndex(ctx, new(payload.Empty))
	return err
}

func (c *agentClient) CreateAndSaveIndex(ctx context.Context, controlCreateIndexRequest *client.ControlCreateIndexRequest) error {
	_, err := c.AgentClient.CreateAndSaveIndex(ctx, controlCreateIndexRequest)
	return err
}

func (c *agentClient) IndexInfo(ctx context.Context) (*client.InfoIndex, error) {
	return c.AgentClient.IndexInfo(ctx, new(payload.Empty))
}
