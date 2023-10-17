//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

// Package ngt provides implementation of Go API for https://github.com/yahoojapan/NGT
package main

/*
#cgo LDFLAGS: -lhdf5 -lngt
#include <NGT/Capi.h>
#include <hdf5/serial/hdf5.h>
#include <stdlib.h>

float* vectors;
size_t size;
size_t dim;

void load(const char* filename) {
	hid_t file_id = H5Fopen(filename, H5F_ACC_RDONLY, H5P_DEFAULT);
    hid_t dataset_id = H5Dopen(file_id, "train", H5P_DEFAULT);
    hid_t space_id = H5Dget_space(dataset_id);
    int ndims = H5Sget_simple_extent_ndims(space_id);
    hsize_t *dims = (hsize_t *)malloc(sizeof(hsize_t) * ndims);
    H5Sget_simple_extent_dims(space_id, dims, NULL);

	size = dims[0];
	dim = dims[1];

	vectors = (float *)malloc(size * dim * sizeof(float));
    H5Dread(dataset_id, H5T_NATIVE_FLOAT, H5S_ALL, space_id, H5P_DEFAULT, vectors);

    H5Sclose(space_id);
    H5Dclose(dataset_id);
    H5Fclose(file_id);
}

float* vector(int i) {
	return &vectors[i * dim];
}

void free_vector() {
	free(vectors);
}
*/
import "C"

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/vdaas/vald/internal/conv"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/http/metrics"
	"github.com/vdaas/vald/internal/strings"
)

func main() {
	var (
		buf           []byte
		err           error
		lines, fields []string
		line          string
		vmpeak,
		vmsize,
		vmdata,
		vmrss,
		vmhwm,
		vmstack,
		vmswap,
		vmexe,
		vmlib,
		vmlock,
		vmpin,
		vmpte,
		gc,
		save,
		cls float64

		pfile  = fmt.Sprintf("/proc/%d/status", os.Getpid())
		zero   = float64(0.0)
		format = "%s\t" + strings.TrimSuffix(strings.Repeat("%.2f\t", 42), "\t")
	)
	output := func(header string) {
		buf, err = os.ReadFile(pfile)
		if err != nil {
			log.Fatal(err)
		}
		lines = strings.Split(conv.Btoa(buf), "\n")
		for _, line = range lines {
			fields = strings.Fields(line)

			switch {
			case strings.HasPrefix(line, "VmPeak"):
				f, err := strconv.ParseFloat(fields[1], 64)
				if err == nil {
					vmpeak = f
				}
			case strings.HasPrefix(line, "VmSize"):
				f, err := strconv.ParseFloat(fields[1], 64)
				if err == nil {
					vmsize = f
				}
			case strings.HasPrefix(line, "VmHWM"):
				f, err := strconv.ParseFloat(fields[1], 64)
				if err == nil {
					vmhwm = f
				}
			case strings.HasPrefix(line, "VmRSS"):
				f, err := strconv.ParseFloat(fields[1], 64)
				if err == nil {
					vmrss = f
				}
			case strings.HasPrefix(line, "VmData"):
				f, err := strconv.ParseFloat(fields[1], 64)
				if err == nil {
					vmdata = f
				}
			case strings.HasPrefix(line, "VmStk"):
				f, err := strconv.ParseFloat(fields[1], 64)
				if err == nil {
					vmstack = f
				}
			case strings.HasPrefix(line, "VmExe"):
				f, err := strconv.ParseFloat(fields[1], 64)
				if err == nil {
					vmexe = f
				}
			case strings.HasPrefix(line, "VmLck"):
				f, err := strconv.ParseFloat(fields[1], 64)
				if err == nil {
					vmlock = f
				}
			case strings.HasPrefix(line, "VmLib"):
				f, err := strconv.ParseFloat(fields[1], 64)
				if err == nil {
					vmlib = f
				}
			case strings.HasPrefix(line, "VmPTE"):
				f, err := strconv.ParseFloat(fields[1], 64)
				if err == nil {
					vmpte = f
				}
			case strings.HasPrefix(line, "VmSwap"):
				f, err := strconv.ParseFloat(fields[1], 64)
				if err == nil {
					vmswap = f
				}
			case strings.HasPrefix(line, "VmPin"):
				f, err := strconv.ParseFloat(fields[1], 64)
				if err == nil {
					vmpin = f
				}
			}
			fields = fields[:0:0]
			fields = nil
		}
		buf = buf[:0:0]
		buf = nil
		lines = lines[:0:0]
		lines = nil
		switch {
		case strings.Contains(header, "gc"):
			gc = vmpeak
			save = zero
			cls = zero
		case strings.Contains(header, "save"):
			save = vmpeak
			gc = zero
			cls = zero
		case strings.Contains(header, "close"):
			cls = vmpeak
			gc = zero
			save = zero
		default:
			gc = zero
			save = zero
			cls = zero
		}

		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		metrics := []interface{}{
			header,
			vmpeak,
			vmsize,
			vmdata,
			vmrss,
			vmhwm,
			vmstack,
			vmswap,
			vmexe,
			vmlib,
			vmlock,
			vmpin,
			vmpte,
			float64(m.Alloc) / 1024.0,
			float64(m.BuckHashSys),
			float64(m.Frees),
			float64(m.GCSys) / 1024.0,
			float64(m.HeapAlloc) / 1024.0,
			float64(m.HeapIdle) / 1024.0,
			float64(m.HeapInuse) / 1024.0,
			float64(m.HeapObjects),
			float64(m.HeapReleased) / 1024.0,
			float64(m.HeapSys) / 1024.0,
			float64(m.HeapIdle-m.HeapReleased) / 1024.0,
			float64(m.Lookups),
			float64(m.MCacheInuse),
			float64(m.MCacheSys),
			float64(m.MSpanInuse) / 1024.0,
			float64(m.MSpanSys),
			float64(m.Mallocs),
			float64(m.Mallocs - m.Frees),
			float64(m.NextGC),
			float64(m.NumForcedGC),
			float64(m.NumGC),
			float64(m.OtherSys),
			float64(m.PauseTotalNs) / 1024.0,
			float64(m.StackInuse),
			float64(m.StackSys),
			float64(m.Sys) / 1024.0,
			float64(m.TotalAlloc) / 1024.0,
			gc,
			save,
			cls,
		}
		log.Infof(format, metrics...)
		switch {
		case strings.Contains(header, "gc"),
			strings.Contains(header, "save"),
			strings.Contains(header, "close"):
			log.Infof(format, metrics...)
		}
	}
	defer output("end")
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		srv := &http.Server{
			Addr:    "0.0.0.0:6060",
			Handler: metrics.NewPProfHandler(),
		}
		go srv.ListenAndServe()
		<-ctx.Done()
		srv.Shutdown(context.Background())
	}()

	C.load(C.CString(os.Getenv("DATA_PATH")))
	log.Infof("# of vectors: %v", C.size)
	log.Info(strings.Join([]string{
		"Operation",
		"VmPeak",
		"VmSize",
		"VmData",
		"VmRSS",
		"VmHWM",
		"VmStack",
		"VmSwap",
		"VmEXE",
		"VmLib",
		"VmLock",
		"VmPin",
		"VmPTE",
		"Alloc",
		"BuckHashSys",
		"Frees",
		"GCSys",
		"HeapAlloc",
		"HeapIdle",
		"HeapInuse",
		"HeapObjects",
		"HeapReleased",
		"HeapSys",
		"HeapWillReturn",
		"Lookups",
		"MCacheInuse",
		"MCacheSys",
		"MSpanInuse",
		"MSpanSys",
		"Mallocs",
		"LiveObjects",
		"NextGC",
		"NumForcedGC",
		"NumGC",
		"OtherSys",
		"PauseTotalNs",
		"StackInuse",
		"StackSys",
		"Sys",
		"TotalAlloc",
		"GC",
		"Save",
		"Close",
	}, "\t"))
	output("start")
	path, _ := file.MkdirTemp("")
	sleep(ctx, time.Second*5, time.Minute*1, func() {
		output("waiting for start")
	}, func() {
		runtime.GC()
		output("gc")
		time.Sleep(time.Minute)
		output("starting")
	})

	run(ctx, false, path, time.Minute*10, output)
	sleep(ctx, time.Second*5, time.Minute*4, func() {
		output("waiting for next")
	}, func() {
		runtime.GC()
		output("gc")
		time.Sleep(time.Minute)
		output("starting")
	})
	run(ctx, true, path, 0, output)
	// sleep(ctx, time.Second*5, time.Minute*4, func() {
	// 	output("waiting for next")
	// }, func() {
	// 	runtime.GC()
	// 	output("gc")
	// 	time.Sleep(time.Minute)
	// 	output("starting")
	// })
	// run(ctx, true, path, len(vectors[0]), vectors, ids, time.Hour*2, output)

	// ids = ids[:0:0]
	// ids = nil
	// vectors = vectors[:0:0]
	// vectors = nil
	sleep(ctx, time.Second*5, time.Minute*1, func() {
		output("waiting for gc")
	}, func() {
		runtime.GC()
		output("gc")
	})
	sleep(ctx, time.Second*5, time.Minute*1, func() {
		output("waiting for gc")
	}, func() {
		runtime.GC()
		output("gc")
	})
	sleep(ctx, time.Second*5, time.Minute*1, func() {
		output("finalizing")
	}, func() {
		cancel()
		wg.Wait()
	})
}

func run(ctx context.Context, load bool, path string, dur time.Duration, output func(header string)) {
	dim := C.uint32_t(C.dim)
	nerr := C.ngt_create_error_object()
	nprop := C.ngt_create_property(nerr)
	C.ngt_set_property_dimension(nprop, C.int32_t(dim), nerr)
	C.ngt_set_property_distance_type_l2(nprop, nerr)
	index := C.ngt_create_graph_and_tree_in_memory(nprop, nerr)
	C.ngt_destroy_property(nprop)

	if C.vectors != nil {
		ids := make([]C.ObjectID, int(C.size))

		sleep(ctx, 0, dur, func() {
			for i := 0; i < int(C.dim); i++ {
				ids[i] = C.ngt_insert_index_as_float(index, C.vector(C.int(i)), dim, nerr)
			}
			output("insert")

			C.ngt_create_index(index, 8, nerr)
			output("create index")
			for _, id := range ids {
				C.ngt_remove_index(index, id, nerr)
			}
			output("remove")
		}, func() {
			for i := 0; i < int(C.dim); i++ {
				ids[i] = C.ngt_insert_index_as_float(index, C.vector(C.int(i)), dim, nerr)
			}
			output("insert")

			C.ngt_create_index(index, 8, nerr)
			output("create index")
			for _, id := range ids {
				C.ngt_remove_index(index, id, nerr)
			}
			output("remove")
		})
	}
	sleep(ctx, time.Second*5, time.Minute*10, func() {
		output("finalizing")
	}, func() {
		C.ngt_destroy_error_object(nerr)
		C.ngt_close_index(index)
		output("close")
	})
}

func sleep(ctx context.Context, duration, limit time.Duration, fn, efn func()) {
	if limit == 0 {
		fn()
		efn()
		return
	}
	defer efn()
	end := time.NewTimer(limit)
	defer end.Stop()
	if duration == 0 {
		for {
			select {
			case <-ctx.Done():
				return
			case <-end.C:
				return
			default:
				fn()
			}
		}
		return
	}
	ticker := time.NewTicker(duration)
	defer ticker.Stop()
	for range ticker.C {
		select {
		case <-ctx.Done():
			return
		case <-end.C:
			return
		default:
			fn()
		}
	}
	return
}
