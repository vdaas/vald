package internal

import (
	"testing"
	"time"

	"github.com/kpango/fuid"
	"github.com/vdaas/vald/internal/log"
)

func createIDs(n int) []string {
	ids := make([]string, 0, n)
	for i := 0; i < n; i++ {
		ids = append(ids, fuid.String())
	}
	return ids
}

func Insert(tb testing.TB, dataset [][]float64, insert func(string, []float64) error) ([]string, time.Duration) {
	tb.Helper()
	ids := createIDs(len(dataset))
	start := time.Now()
	for i, vector := range dataset {
		if err := insert(ids[i], vector); err != nil {
			log.Error(err)
		}
	}
	return ids, time.Now().Sub(start)
}

func CreateIndex(tb testing.TB, createIndex func() error) time.Duration {
	tb.Helper()
	start := time.Now()
	if err := createIndex(); err != nil {
		log.Error(err)
	}
	return time.Now().Sub(start)
}

func Search(tb testing.TB, dataset [][]float64, search func([]float64) error) time.Duration {
	tb.Helper()
	start := time.Now()
	for _, vector := range dataset {
		if err := search(vector); err != nil {
			log.Error(err)
		}
	}
	return time.Now().Sub(start)
}

func Remove(tb testing.TB, dataset []string, remove func(id string) error) time.Duration {
	tb.Helper()

	start := time.Now()
	for _, id := range dataset {
		if err := remove(id); err != nil {
			log.Error(err)
		}
	}
	return time.Now().Sub(start)
}