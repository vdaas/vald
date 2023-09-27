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
	"github.com/vdaas/vald/internal/core/algorithm/ngt"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/http/metrics"
	"github.com/vdaas/vald/internal/strings"
	"gonum.org/v1/hdf5"
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

	vectors, _, _ := load(os.Getenv("DATA_PATH"))
	log.Infof("# of vectors: %v", len(vectors))
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
	sleep(ctx, time.Second*5, time.Minute*4, func() {
		output("waiting for start")
	}, func() {
		runtime.GC()
		output("gc")
		time.Sleep(time.Minute)
		output("starting")
	})

	ids := make([]uint, len(vectors))
	run(ctx, false, path, len(vectors[0]), vectors, ids, time.Hour*1, output)
	sleep(ctx, time.Second*5, time.Minute*4, func() {
		output("waiting for next")
	}, func() {
		runtime.GC()
		output("gc")
		time.Sleep(time.Minute)
		output("starting")
	})
	// run(ctx, true, path, len(vectors[0]), nil, nil, 0, output)
	// sleep(ctx, time.Second*5, time.Minute*4, func() {
	// 	output("waiting for next")
	// }, func() {
	// 	runtime.GC()
	// 	output("gc")
	// 	time.Sleep(time.Minute)
	// 	output("starting")
	// })
	// run(ctx, true, path, len(vectors[0]), vectors, ids, time.Hour*2, output)

	ids = ids[:0:0]
	ids = nil
	vectors = vectors[:0:0]
	vectors = nil
	sleep(ctx, time.Second*5, time.Minute*5, func() {
		output("waiting for gc")
	}, func() {
		runtime.GC()
		output("gc")
	})
	sleep(ctx, time.Second*5, time.Minute*5, func() {
		output("waiting for gc")
	}, func() {
		runtime.GC()
		output("gc")
	})
	sleep(ctx, time.Second*5, time.Minute*5, func() {
		output("finalizing")
	}, func() {
		cancel()
		wg.Wait()
	})
}

func run(ctx context.Context, load bool, path string, dim int, vectors [][]float32, ids []uint, dur time.Duration, output func(header string)) {
	var n ngt.NGT
	if load {
		n, _ = ngt.Load(
			ngt.WithDimension(dim),
			ngt.WithDefaultPoolSize(8),
			ngt.WithObjectType(ngt.Float),
			ngt.WithDistanceType(ngt.L2),
		)
	} else {
		n, _ = ngt.New(
			ngt.WithDimension(dim),
			ngt.WithDefaultPoolSize(8),
			ngt.WithObjectType(ngt.Float),
			ngt.WithDistanceType(ngt.L2),
		)
	}

	if vectors != nil {
		var (
			i      int
			vector []float32
			id     uint
			err    error
		)
		if ids == nil {
			ids = make([]uint, len(vectors))
		} else if load {
			for _, id = range ids {
				_ = n.Remove(id)
			}
			output("remove")
		}
		sleep(ctx, 0, dur, func() {
			for i, vector = range vectors {
				id, err = n.Insert(vector)
				if err != nil {
					log.Fatal(err)
				}
				ids[i] = id
			}
			output("insert")
			if err = n.CreateIndex(8); err != nil {
				log.Fatal(err)
			}
			output("create index")
			for _, id = range ids {
				if err = n.Remove(id); err != nil {
					log.Fatal(err)
				}
			}
			output("remove")
		}, func() {
			for _, vector = range vectors {
				_, err = n.Insert(vector)
				if err != nil {
					log.Fatal(err)
				}
			}
			output("insert")
			if err = n.CreateIndex(8); err != nil {
				log.Fatal(err)
			}
			output("create index")
			if err = n.SaveIndex(); err != nil {
				log.Fatal(err)
			}
			output("save index")
		})
	}
	sleep(ctx, time.Second*5, time.Minute*10, func() {
		output("finalizing")
	}, func() {
		n.Close()
		n = nil
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

// load function loads training and test vector from hdf file. The size of ids is same to the number of training data.
// Each id, which is an element of ids, will be set a random number.
func load(path string) (train, test [][]float32, err error) {
	var f *hdf5.File
	f, err = hdf5.OpenFile(path, hdf5.F_ACC_RDONLY)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	// readFn function reads vectors of the hierarchy with the given the name.
	readFn := func(name string) ([][]float32, error) {
		// Opens and returns a named Dataset.
		// The returned dataset must be closed by the user when it is no longer needed.
		d, err := f.OpenDataset(name)
		if err != nil {
			return nil, err
		}
		defer d.Close()

		// Space returns an identifier for a copy of the dataspace for a dataset.
		sp := d.Space()
		defer sp.Close()

		// SimpleExtentDims returns dataspace dimension size and maximum size.
		dims, _, _ := sp.SimpleExtentDims()
		row, dim := int(dims[0]), int(dims[1])

		// Gets the stored vector. All are represented as one-dimensional arrays.
		// The type of the slice depends on your dataset.
		// For fashion-mnist-784-euclidean.hdf5, the datatype is float32.
		vec := make([]float32, sp.SimpleExtentNPoints())
		if err := d.Read(&vec); err != nil {
			return nil, err
		}

		// Converts a one-dimensional array to a two-dimensional array.
		// Use the `dim` variable as a separator.
		vecs := make([][]float32, row)
		for i := 0; i < row; i++ {
			vecs[i] = make([]float32, dim)
			for j := 0; j < dim; j++ {
				vecs[i][j] = float32(vec[i*dim+j])
			}
		}

		return vecs, nil
	}

	// Gets vector of `train` hierarchy.
	train, err = readFn("train")
	if err != nil {
		return nil, nil, err
	}

	// Gets vector of `test` hierarchy.
	test, err = readFn("test")
	if err != nil {
		return nil, nil, err
	}

	return
}
