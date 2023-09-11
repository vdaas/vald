package bbolt_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vdaas/vald/internal/db/kvs/bbolt"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

func TestGetSetClose(t *testing.T) {
	t.Parallel()

	tempdir := t.TempDir()
	tmpfile := filepath.Join(tempdir, "test.db")
	b, err := bbolt.New(tmpfile, "", nil)
	require.NoError(t, err)

	err = b.Set("key", []byte("value"))
	require.NoError(t, err)

	val, ok, err := b.Get("key")
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, []byte("value"), val)

	val, ok, err = b.Get("no exist key")
	require.NoError(t, err)
	require.False(t, ok)
	require.Nil(t, val)

	err = b.Close(false)
	require.NoError(t, err)

	b, err = bbolt.New(tmpfile, "", nil)
	require.NoError(t, err)

	// recover from the file
	val, ok, err = b.Get("key")
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, []byte("value"), val)

	err = b.Close(true)
	require.NoError(t, err)

	// now the file is deleted
	_, err = os.Stat(tmpfile)
	require.True(t, os.IsNotExist(err))
}

func TestAsyncSet(t *testing.T) {
	t.Parallel()

	tempdir := t.TempDir()
	tmpfile := filepath.Join(tempdir, "test.db")
	b, err := bbolt.New(tmpfile, "", nil)
	require.NoError(t, err)

	kv := map[string]string{
		"key1": "val1",
		"key2": "val2",
		"key3": "val3",
		"key4": "val4",
		"key5": "val5",
	}

	eg, _ := errgroup.New(context.Background())
	for k, v := range kv {
		b.AsyncSet(&eg, k, []byte(v))
	}

	// wait until all set is done
	eg.Wait()

	for k := range kv {
		_, ok, err := b.Get(k)
		require.NoError(t, err)
		require.True(t, ok)
	}

	err = b.Close(true)
	require.NoError(t, err)
}
