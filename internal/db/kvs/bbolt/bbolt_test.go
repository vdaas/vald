// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

	err = b.Set([]byte("key"), []byte("value"))
	require.NoError(t, err)

	val, ok, err := b.Get([]byte("key"))
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, []byte("value"), val)

	val, ok, err = b.Get([]byte("no exist key"))
	require.NoError(t, err)
	require.False(t, ok)
	require.Nil(t, val)

	err = b.Close(false)
	require.NoError(t, err)

	b, err = bbolt.New(tmpfile, "", nil)
	require.NoError(t, err)

	// recover from the file
	val, ok, err = b.Get([]byte("key"))
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
		b.AsyncSet(&eg, []byte(k), []byte(v))
	}

	// wait until all set is done
	eg.Wait()

	for k := range kv {
		_, ok, err := b.Get([]byte(k))
		require.NoError(t, err)
		require.True(t, ok)
	}

	err = b.Close(true)
	require.NoError(t, err)
}
