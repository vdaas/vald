//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

package main

import (
	"context"
	"encoding/gob"
	"flag"
	"fmt"
	"io/fs"
	"os"

	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/pkg/agent/core/ngt/service/kvs"
)

var (
	format               = flag.String("format", "csv", "file format(csv,tsv)")
	kvsFileName          = flag.String("file", "ngt-meta.kvsdb", "kvsdb file name")
	kvsTimestampFileName = flag.String("timestamp-file", "ngt-timestamp.kvsdb", "kvsdb timestamp file name")
	path                 = flag.String("path", ".", "kvsdb file path")
)

func main() {
	flag.Parse()
	log.Init()

	kvsdb := kvs.New()

	// value
	m := make(map[string]uint32)
	gob.Register(map[string]uint32{})
	var f *os.File
	defer f.Close()
	f, _ = file.Open(
		file.Join(*path, *kvsFileName),
		os.O_RDONLY|os.O_SYNC,
		fs.ModePerm,
	)
	_ = gob.NewDecoder(f).Decode(&m)

	// timestamp
	mt := make(map[string]int64)
	gob.Register(map[string]int64{})
	var ft *os.File
	defer ft.Close()
	ft, _ = file.Open(
		file.Join(*path, *kvsTimestampFileName),
		os.O_RDONLY|os.O_SYNC,
		fs.ModePerm,
	)
	_ = gob.NewDecoder(ft).Decode(&mt)

	// kvs load
	for k, id := range m {
		if ts, ok := mt[k]; ok {
			kvsdb.Set(k, id, ts)
		} else {
			kvsdb.Set(k, id, 0)
		}
	}

	// print
	kvsdb.Range(context.TODO(), func(uuid string, oid uint32, ts int64) bool {
		if *format == "csv" {
			fmt.Printf("%s,%d,%d\n", uuid, oid, ts)
		} else if *format == "tsv" {
			fmt.Printf("%s\t%d\t%d\n", uuid, oid, ts)
		} else {
			fmt.Println(uuid, oid, ts)
		}
		return true
	})
}
