package persistent_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vdaas/vald/internal/cache/persistent"
	"github.com/vdaas/vald/internal/sync"
)

func TestPersistentCache(t *testing.T) {
	base := t.TempDir()
	pc, err := persistent.NewPCache(base)
	require.NoError(t, err)

	len := 4096

	for i := 0; i < len; i++ {
		err := pc.Set(fmt.Sprint(i), struct{}{})
		require.NoError(t, err)
	}

	for i := 0; i < len; i++ {
		_, ok, err := pc.Get(fmt.Sprint(i))
		require.NoError(t, err)
		require.True(t, ok, fmt.Sprintf("i: %d", i))
	}

	for i := 0; i < len; i++ {
		err := pc.Delete(fmt.Sprint(i))
		require.NoError(t, err)
	}

	for i := 0; i < len; i++ {
		_, ok, err := pc.Get(fmt.Sprint(i))
		require.NoError(t, err)
		require.False(t, ok, fmt.Sprintf("i: %d", i))
	}

	err = pc.Close()
	require.NoError(t, err)
}

func TestPersistentCacheConcurrent(t *testing.T) {
	base := t.TempDir()
	pc, err := persistent.NewPCache(base)
	require.NoError(t, err)

	len := 4096

	var wg sync.WaitGroup
	for i := 0; i < len; i++ {
		wg.Add(1)
		go func(key int) {
			defer wg.Done()
			err := pc.Set(fmt.Sprint(key), struct{}{})
			require.NoError(t, err)
		}(i)
	}

	wg.Wait()

	for i := 0; i < len; i++ {
		wg.Add(1)
		go func(key int) {
			defer wg.Done()
			_, ok, err := pc.Get(fmt.Sprint(key))
			require.NoError(t, err)
			require.True(t, ok, fmt.Sprintf("i: %d", key))
		}(i)
	}

	wg.Wait()

	for i := 0; i < len; i++ {
		wg.Add(1)
		go func(key int) {
			defer wg.Done()
			err := pc.Delete(fmt.Sprint(key))
			require.NoError(t, err)
		}(i)
	}

	wg.Wait()

	for i := 0; i < len; i++ {
		wg.Add(1)
		go func(key int) {
			defer wg.Done()
			_, ok, err := pc.Get(fmt.Sprint(key))
			require.NoError(t, err)
			require.False(t, ok, fmt.Sprintf("i: %d", key))
		}(i)
	}

	wg.Wait()

	err = pc.Close()
	require.NoError(t, err)
}
