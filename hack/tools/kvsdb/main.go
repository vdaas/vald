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
	"encoding/csv"
	"encoding/gob"
	"flag"
	"io/fs"
	"os"
	"strconv"
	"unsafe"

	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/log"
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

	// print
	var s [][]string
	w := csv.NewWriter(os.Stdout)
	defer w.Flush()
	switch *format {
	case "csv":
	case "tsv":
		w.Comma = '\t'
	default:
		w.Comma = ' '
	}
	s = append(s, []string{"uuid", "oid", "timestamp"})
	for k, id := range m {
		if ts, ok := mt[k]; ok {
			s = append(s, []string{k, strconv.FormatUint(uint64(id), 10), strconv.FormatInt(ts, 10)})
		} else {
			s = append(s, []string{k, strconv.FormatUint(uint64(id), 10), "0"})
		}
		if len(s)*int(unsafe.Sizeof("")) > 4e+6 {
			w.WriteAll(s)
			s = nil
		}
	}
	w.WriteAll(s)
}
