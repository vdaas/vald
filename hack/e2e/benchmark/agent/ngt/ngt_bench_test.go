package ngt

import (
	"context"
	"github.com/vdaas/vald/apis/grpc/payload"
	"io"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/agent"
	"github.com/vdaas/vald/hack/e2e/benchmark/internal"
	"github.com/vdaas/vald/internal/log"
	"google.golang.org/grpc"
)

var (
	train [][]float64
	test [][]float64
	ids []string
)

func Setup(tb testing.TB) (context.Context, agent.AgentClient, func()) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)

	tb.Helper()
	conn, err := grpc.DialContext(ctx, "localhost:8082", grpc.WithInsecure())
	if err != nil {
		tb.Fatal(err)
	}
	return ctx, agent.NewAgentClient(conn), func(){
		defer cancel()
	}
}

func BenchmarkAgentNGTSequential(b *testing.B) {
	ctx, client, teardown := Setup(b)
	defer teardown()

	b.Run("Insert", func(bb *testing.B) {
		bb.ResetTimer()
		var t time.Duration
		ids, t = internal.Insert(bb, train, func(id string, vector []float64) error {
			_, err := client.Insert(ctx, &payload.Object_Vector{
				Id: &payload.Object_ID{
					Id: id,
				},
				Vector: vector,
			})
			return err
		})
		bb.Logf("duration: %v", t)
	})

	b.Run("CreateIndex", func(bb *testing.B) {
		bb.ResetTimer()
		t := internal.CreateIndex(bb, func() error {
			_, err := client.CreateIndex(ctx, &payload.Controll_CreateIndexRequest{
				PoolSize: 1000,
			})
			return err
		})
		bb.Logf("duration: %v", t)
	})

	b.Run("Search", func(bb *testing.B) {
		bb.ResetTimer()
		t := internal.Search(bb, test, func(vector []float64) error {
			_, err := client.Search(ctx, &payload.Search_Request{
				Vector: &payload.Object_Vector{
					Vector: vector,
				},
				Config: &payload.Search_Config{
					Num:     10,
					Radius:  -1,
					Epsilon: 0.01,
				},
			})
			return err
		})
		bb.Logf("duration: %v", t)
	})

	b.Run("Remove", func(bb *testing.B) {
		bb.ResetTimer()
		t := internal.Remove(bb, ids[:len(ids)/10], func(id string) error {
			_, err := client.Remove(ctx, &payload.Object_ID{
				Id: id,
			})
			return err
		})
		bb.Logf("duration: %v", t)
	})
}

func BenchmarkAgentNGTStream(b *testing.B) {
	ctx, client, teardown := Setup(b)
	defer teardown()

	b.Run("Insert", func(bb *testing.B) {
		st, err := client.StreamInsert(ctx)
		if err != nil {
			bb.Error(err)
		}
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go func() {
			start := time.Now()
			defer func() {
				bb.Logf("recieve duration: %v", time.Now().Sub(start))
				wg.Done()
			}()
			for {
				_, err := st.Recv()
				if err == io.EOF {
					break
				} else if err != nil {
					bb.Error(err)
				}
			}
		}()
		var t time.Duration
		bb.ResetTimer()
		ids, t = internal.Insert(bb, train, func(id string, vector []float64) error {
			err := st.Send(&payload.Object_Vector{
				Id: &payload.Object_ID{
					Id: id,
				},
				Vector: vector,
			})
			return err
		})
		if err := st.CloseSend(); err != nil {
			bb.Error(err)
		}
		bb.Logf("send duration: %v", t)
		wg.Wait()
	})

	b.Run("CreateIndex", func(bb *testing.B) {
		bb.ResetTimer()
		t := internal.CreateIndex(bb, func() error {
			_, err := client.CreateIndex(ctx, &payload.Controll_CreateIndexRequest{
				PoolSize: 1000,
			})
			return err
		})
		bb.Logf("duration: %v", t)
	})

	b.Run("Search", func(bb * testing.B) {
		st, err := client.StreamSearch(ctx)
		if err != nil {
			bb.Error(err)
		}
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go func() {
			start := time.Now()
			defer func() {
				bb.Logf("recieve duration: %v", time.Now().Sub(start))
				wg.Done()
			}()
			for {
				_, err := st.Recv()
				if err == io.EOF {
					break
				} else if err != nil {
					bb.Error(err)
				}
			}
		}()
		bb.ResetTimer()
		t := internal.Search(bb, test, func(vector []float64) error {
			err := st.Send(&payload.Search_Request{
				Vector: &payload.Object_Vector{
					Vector: vector,
				},
				Config: &payload.Search_Config{
					Num:     10,
					Radius:  -1,
					Epsilon: 0.01,
				},
			})
			return err
		})

		bb.Logf("send duration: %v", t)
		wg.Wait()
	})

	b.Run("Remove", func(bb * testing.B) {
		st, err := client.StreamRemove(ctx)
		if err != nil {
			bb.Error(err)
		}
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go func() {
			start := time.Now()
			defer func() {
				bb.Logf("recieve duration: %v", time.Now().Sub(start))
				wg.Done()
			}()
			for {
				_, err := st.Recv()
				if err == io.EOF {
					break
				} else if err != nil {
					bb.Error(err)
				}
			}
		}()
		bb.ResetTimer()
		t := internal.Remove(bb, ids[:len(ids)/10], func(id string) error {
			err := st.Send(&payload.Object_ID{
				Id: id,
			})
			return err
		})

		bb.Logf("send duration: %v", t)
		wg.Wait()
	})
}

func TestMain(m *testing.M) {
	log.Init(log.DefaultGlg())

	datasetName := "../../assets/dataset/fashion-mnist-784-euclidean.hdf5"
	var err error
	log.Info("start load dataset")
	train, test, err = internal.Load(datasetName)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("finish load dataset")

	ret := m.Run()
	os.Exit(ret)
}
