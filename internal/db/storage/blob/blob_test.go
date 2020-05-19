//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package blob

import (
	"context"
	"testing"

	"github.com/vdaas/vald/internal/db/storage/blob/s3"
)

const (
	endpoint        = ""
	region          = ""
	accessKey       = ""
	secretAccessKey = ""
	bucketURL       = ""
)

func TestS3Write(t *testing.T) {
	opener, err := s3.NewSession(
		s3.WithEndpoint(endpoint),
		s3.WithRegion(region),
		s3.WithAccessKey(accessKey),
		s3.WithSecretAccessKey(secretAccessKey),
	).URLOpener()
	if err != nil {
		t.Fatalf("opener initialize failed: %s", err)
	}

	bucket, err := NewBucket(
		WithBucketURLOpener(opener),
		WithBucketURL(bucketURL),
	)
	if err != nil {
		t.Fatalf("bucket initialize failed: %s", err)
	}

	ctx := context.Background()

	err = bucket.Open(ctx)
	if err != nil {
		t.Fatalf("bucket open failed: %s", err)
	}

	defer func() {
		err = bucket.Close()
		if err != nil {
			t.Fatalf("bucket close failed: %s", err)
		}
	}()

	w, err := bucket.Writer(ctx, "writer-test.txt")
	if err != nil {
		t.Fatalf("fetch writer failed: %s", err)
	}
	defer func() {
		err = w.Close()
		if err != nil {
			t.Fatalf("writer close failed: %s", err)
		}
	}()

	_, err = w.Write([]byte("Hello from blob world!"))
	if err != nil {
		t.Fatalf("write failed: %s", err)
	}
}

func TestS3Read(t *testing.T) {
	opener, err := s3.NewSession(
		s3.WithEndpoint(endpoint),
		s3.WithRegion(region),
		s3.WithAccessKey(accessKey),
		s3.WithSecretAccessKey(secretAccessKey),
	).URLOpener()
	if err != nil {
		t.Fatalf("opener initialize failed: %s", err)
	}

	bucket, err := NewBucket(
		WithBucketURLOpener(opener),
		WithBucketURL(bucketURL),
	)
	if err != nil {
		t.Fatalf("bucket initialize failed: %s", err)
	}

	ctx := context.Background()

	err = bucket.Open(ctx)
	if err != nil {
		t.Fatalf("bucket open failed: %s", err)
	}

	defer func() {
		err = bucket.Close()
		if err != nil {
			t.Fatalf("bucket close failed: %s", err)
		}
	}()

	r, err := bucket.Reader(ctx, "writer-test.txt")
	if err != nil {
		t.Fatalf("fetch reader failed: %s", err)
	}
	defer func() {
		err = r.Close()
		if err != nil {
			t.Fatalf("reader close failed: %s", err)
		}
	}()

	rbuf := make([]byte, 16)
	_, err = r.Read(rbuf)
	if err != nil {
		t.Fatalf("read failed: %s", err)
	}

	t.Logf("read: %s", string(rbuf))
}
