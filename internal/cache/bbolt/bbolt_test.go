package bbolt_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vdaas/vald/internal/cache/bbolt"
)

func TestBbolt(t *testing.T) {
	tempdir := t.TempDir()
	tmpfile := filepath.Join(tempdir, "test.db")
	b, err := bbolt.New(tmpfile)
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

	err = b.Close()
	require.NoError(t, err)
}

func TestSetBatch(t *testing.T) {
	tempdir := t.TempDir()
	tmpfile := filepath.Join(tempdir, "test.db")
	b, err := bbolt.New(tmpfile)
	require.NoError(t, err)

	kv := map[string]struct{}{
		"key1": {},
		"key2": {},
		"key3": {},
		"key4": {},
		"key5": {},
	}

	err = b.SetBatch(kv)
	require.NoError(t, err)

	for k := range kv {
		_, ok, err := b.Get(k)
		require.NoError(t, err)
		require.True(t, ok)
	}

	err = b.Close()
	require.NoError(t, err)
}
