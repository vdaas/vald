package internal

import (
	"testing"

	"github.com/kpango/fuid"
	"github.com/vdaas/vald/internal/log"
)

func CreateIDs(n int) []string {
	ids := make([]string, 0, n)
	for i := 0; i < n; i++ {
		ids = append(ids, fuid.String())
	}
	return ids
}

func Insert(b *testing.B, ids []string, dataset [][]float64, insert func(string, []float64) error) []string {
	b.Helper()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		idx := i % b.N
		if err := insert(ids[idx], dataset[idx]); err != nil {
			log.Error(err)
		}
	}
	return ids
}

func CreateIndex(b *testing.B, createIndex func() error) {
	b.Helper()
	b.ResetTimer()
	if err := createIndex(); err != nil {
		log.Error(err)
	}
}

func Search(b *testing.B, dataset [][]float64, search func([]float64) error) {
	b.Helper()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := search(dataset[i%b.N]); err != nil {
			log.Error(err)
		}
	}
}

func Remove(b *testing.B, dataset []string, remove func(id string) error) {
	b.Helper()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		idx := i % b.N
		if err := remove(dataset[idx]); err != nil {
			log.Error(err)
		}
	}
}