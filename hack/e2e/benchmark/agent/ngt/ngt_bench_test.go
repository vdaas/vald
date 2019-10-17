package ngt

import (
	"context"
	"io"
	"sync"
	"testing"

	"github.com/vdaas/vald/apis/grpc/agent"
	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/hack/e2e/benchmark/internal"
	"github.com/vdaas/vald/internal/log"
	"google.golang.org/grpc"
)

var (
	train [][]float64
	test [][]float64
	ids []string
	client agent.AgentClient
	searchConfig = &payload.Search_Config{
		Num:     10,
		Radius:  -1,
		Epsilon: 0.01,
	}
)

func init() {
	log.Init(log.DefaultGlg())

	conn, err := grpc.Dial( "localhost:8082", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	client = agent.NewAgentClient(conn)

	datasetName := "../../assets/dataset/fashion-mnist-784-euclidean.hdf5"
	train, test, err = internal.Load(datasetName)
	if err != nil {
		log.Fatal(err)
	}

	ids = internal.CreateIDs(len(train))
}

func BenchmarkAgentNGTSequentialInsert(b *testing.B) {
	ids = internal.Insert(b, ids, train, func(id string, vector []float64) error {
		_, err := client.Insert(context.Background(), &payload.Object_Vector{
			Id: &payload.Object_ID{
				Id: id,
			},
			Vector: vector,
		})
		return err
	})
}

func BenchmarkAgentNGTSequentialCreateIndex(b *testing.B) {
	internal.CreateIndex(b, func() error {
		_, err := client.CreateIndex(context.Background(), &payload.Controll_CreateIndexRequest{
			PoolSize: 10000,
		})
		return err
	})
}

func BenchmarkAgentNGTSequentialSearch(b *testing.B) {
	internal.Search(b, test, func(vector []float64) error {
		_, err := client.Search(context.Background(), &payload.Search_Request{
			Vector: &payload.Object_Vector{
				Vector: vector,
			},
			Config: searchConfig,
		})
		return err
	})
}

func BenchmarkAgentNGTSequentialRemove(b *testing.B) {
	internal.Remove(b, ids, func(id string) error {
		_, err := client.Remove(context.Background(), &payload.Object_ID{
			Id: id,
		})
		return err
	})
}

func BenchmarkAgentNGTStreamInsert(b *testing.B) {
	st, err := client.StreamInsert(context.Background())
	if err != nil {
		b.Error(err)
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(st agent.Agent_StreamInsertClient, wg *sync.WaitGroup) {
		defer wg.Done()
		for {
			_, err := st.Recv()
			if err == io.EOF {
				break
			} else if err != nil {
				b.Error(err)
			}
		}
	}(st, wg)
	ids = internal.Insert(b, ids, train, func(id string, vector []float64) error {
		err := st.Send(&payload.Object_Vector{
			Id: &payload.Object_ID{
				Id: id,
			},
			Vector: vector,
		})
		if err == io.EOF {
			return nil
		}
		return err
	})
	if err := st.CloseSend(); err != nil {
		b.Error(err)
	}

	wg.Wait()
}

func BenchmarkAgentNGTStreamCreateIndex(b *testing.B) {
	internal.CreateIndex(b, func() error {
		_, err := client.CreateIndex(context.Background(), &payload.Controll_CreateIndexRequest{
			PoolSize: 10000,
		})
		return err
	})
}

func BenchmarkAgentNGTStreamSearch(b *testing.B) {
	st, err := client.StreamSearch(context.Background())
	if err != nil {
		b.Error(err)
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(st agent.Agent_StreamSearchClient, wg *sync.WaitGroup) {
		defer wg.Done()
		for {
			_, err := st.Recv()
			if err == io.EOF {
				break
			} else if err != nil {
				b.Error(err)
			}
		}
	}(st, wg)
	internal.Search(b, test, func(vector []float64) error {
		err := st.Send(&payload.Search_Request{
			Vector: &payload.Object_Vector{
				Vector: vector,
			},
			Config: searchConfig,
		})
		if err == io.EOF {
			return nil
			}
		return err
	})
	if err := st.CloseSend(); err != nil {
		b.Error(err)
	}

	wg.Wait()
}

func BenchmarkAgentNGTStreamRemove(b *testing.B) {
	st, err := client.StreamRemove(context.Background())
	if err != nil {
		b.Error(err)
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(st agent.Agent_StreamRemoveClient, wg *sync.WaitGroup) {
		defer wg.Done()
		for {
			_, err := st.Recv()
			if err == io.EOF {
				break
			} else if err != nil {
				b.Error(err)
			}
		}
	}(st, wg)
	internal.Remove(b, ids, func(id string) error {
		err := st.Send(&payload.Object_ID{
			Id: id,
		})
		if err == io.EOF {
			return nil
		}
		return err
	})
	if err := st.CloseSend(); err != nil {
		b.Error(err)
	}

	wg.Wait()
}
