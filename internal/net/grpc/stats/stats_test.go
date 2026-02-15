//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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

package stats

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestRegister(t *testing.T) {
	t.Parallel()
	type args struct {
		srv *grpc.Server
	}
	type want struct{}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		func() test {
			srv := grpc.NewServer()
			return test{
				name: "success to register the stats server",
				args: args{
					srv: srv,
				},
				checkFunc: func(w want) error {
					info := srv.GetServiceInfo()
					if _, exists := info["rpc.v1.Stats"]; !exists {
						return errors.New("Stats service not registered")
					}
					return nil
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			Register(test.args.srv)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_ResourceStats(t *testing.T) {
	// Global variables are modified, so we cannot run in parallel
	// t.Parallel()
	type args struct {
		ctx context.Context
		req *payload.Empty
	}
	type want struct {
		stats *payload.Info_ResourceStats
		err   error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *payload.Info_ResourceStats, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, stats *payload.Info_ResourceStats, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if stats == nil && w.stats != nil {
			return errors.New("stats is nil but want is not nil")
		}
		if stats != nil && w.stats == nil {
			return errors.New("stats is not nil but want is nil")
		}
		return nil
	}
	tests := []test{
		func() test {
			return test{
				name: "success to get resource stats",
				args: args{
					ctx: context.Background(),
					req: &payload.Empty{},
				},
				checkFunc: func(w want, stats *payload.Info_ResourceStats, err error) error {
					if err != nil {
						return errors.Errorf("unexpected error: %v", err)
					}
					if stats == nil {
						return errors.New("stats should not be nil")
					}
					if stats.Name == "" {
						return errors.New("name should not be empty")
					}
					if stats.Ip == "" {
						return errors.New("ip should not be empty")
					}
					// CgroupStats might be nil depending on environment
					return nil
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			s := &server{}
			stats, err := s.ResourceStats(test.args.ctx, test.args.req)
			if err := test.checkFunc(test.want, stats, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_detectCgroupMode(t *testing.T) {
	// t.Parallel()
	type want struct {
		mode CgroupMode
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, CgroupMode) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, mode CgroupMode) error {
		if mode != w.mode {
			return errors.Errorf("got mode: %v, want: %v", mode, w.mode)
		}
		return nil
	}

	// Save original global variables
	origCgroupBasePath := cgroupBasePath
	origProcCgroupPath := procCgroupPath

	tests := []test{
		func() test {
			tmpDir := t.TempDir()

			return test{
				name: "detects cgroup v2 via cgroup.controllers",
				want: want{
					mode: CGV2,
				},
				beforeFunc: func() {
					cgroupBasePath = tmpDir
					_ = os.WriteFile(filepath.Join(tmpDir, "cgroup.controllers"), []byte(""), 0644)
				},
				afterFunc: func() {
					cgroupBasePath = origCgroupBasePath
				},
				checkFunc: func(w want, mode CgroupMode) error {
					if mode != CGV2 {
						return errors.Errorf("expected CGV2, got %v", mode)
					}
					return nil
				},
			}
		}(),
		func() test {
			tmpDir := t.TempDir()
			procFile := filepath.Join(tmpDir, "cgroup")

			return test{
				name: "detects cgroup v2 via proc file",
				want: want{
					mode: CGV2,
				},
				beforeFunc: func() {
					cgroupBasePath = tmpDir
					procCgroupPath = procFile
					// Write cgroup v2 entry
					_ = os.WriteFile(procFile, []byte("0::/foo/bar\n"), 0644)
				},
				afterFunc: func() {
					cgroupBasePath = origCgroupBasePath
					procCgroupPath = origProcCgroupPath
				},
				checkFunc: func(w want, mode CgroupMode) error {
					if mode != CGV2 {
						return errors.Errorf("expected CGV2, got %v", mode)
					}
					return nil
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			mode := detectCgroupMode()
			if err := test.checkFunc(test.want, mode); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_getCgroupV2Path(t *testing.T) {
	// t.Parallel()

	origCgroupBasePath := cgroupBasePath
	origProcCgroupPath := procCgroupPath

	type test struct {
		name       string
		setup      func(t *testing.T) (baseDir, procFile string)
		want       func(baseDir string) string
		cleanup    func()
	}

	tests := []test{
		{
			name: "cgroup namespace disabled (host path exposed)",
			setup: func(t *testing.T) (string, string) {
				tmpDir := t.TempDir()
				procFile := filepath.Join(tmpDir, "proc_cgroup")

				// Create sub-directory representing the cgroup
				subPath := "system.slice/docker-123.scope"
				fullPath := filepath.Join(tmpDir, subPath)
				if err := os.MkdirAll(fullPath, 0755); err != nil {
					t.Fatal(err)
				}

				// Create cgroup.controllers in the sub-directory
				if err := os.WriteFile(filepath.Join(fullPath, "cgroup.controllers"), []byte(""), 0644); err != nil {
					t.Fatal(err)
				}

				// Write proc file pointing to that subpath
				content := "0::/" + subPath + "\n"
				if err := os.WriteFile(procFile, []byte(content), 0644); err != nil {
					t.Fatal(err)
				}

				return tmpDir, procFile
			},
			want: func(baseDir string) string {
				return filepath.Join(baseDir, "system.slice/docker-123.scope")
			},
		},
		{
			name: "cgroup namespace enabled (root path)",
			setup: func(t *testing.T) (string, string) {
				tmpDir := t.TempDir()
				procFile := filepath.Join(tmpDir, "proc_cgroup")

				// Create cgroup.controllers in root
				if err := os.WriteFile(filepath.Join(tmpDir, "cgroup.controllers"), []byte(""), 0644); err != nil {
					t.Fatal(err)
				}

				// Write proc file pointing to root
				content := "0::/\n"
				if err := os.WriteFile(procFile, []byte(content), 0644); err != nil {
					t.Fatal(err)
				}

				return tmpDir, procFile
			},
			want: func(baseDir string) string {
				return baseDir
			},
		},
		{
			name: "fallback to base path if subpath not found",
			setup: func(t *testing.T) (string, string) {
				tmpDir := t.TempDir()
				procFile := filepath.Join(tmpDir, "proc_cgroup")

				// Create cgroup.controllers in root (so it's valid fallback)
				if err := os.WriteFile(filepath.Join(tmpDir, "cgroup.controllers"), []byte(""), 0644); err != nil {
					t.Fatal(err)
				}

				// Write proc file pointing to non-existent subpath
				content := "0::/non-existent\n"
				if err := os.WriteFile(procFile, []byte(content), 0644); err != nil {
					t.Fatal(err)
				}

				return tmpDir, procFile
			},
			want: func(baseDir string) string {
				return baseDir
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			baseDir, procFile := tc.setup(t)
			cgroupBasePath = baseDir
			procCgroupPath = procFile

			defer func() {
				cgroupBasePath = origCgroupBasePath
				procCgroupPath = origProcCgroupPath
			}()

			got := getCgroupV2Path()
			wantPath := tc.want(baseDir)

			if got != wantPath {
				t.Errorf("getCgroupV2Path() = %v, want %v", got, wantPath)
			}
		})
	}
}

func Test_calculateCpuUsageCores(t *testing.T) {
	t.Parallel()
	type args struct {
		m1        *CgroupMetrics
		m2        *CgroupMetrics
		deltaTime time.Duration
	}
	type want struct {
		stats *CgroupStats
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *CgroupStats) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, stats *CgroupStats) error {
		if stats == nil && w.stats != nil {
			return errors.New("stats is nil but want is not nil")
		}
		if stats != nil && w.stats == nil {
			return errors.New("stats is not nil but want is nil")
		}
		return nil
	}
	tests := []test{
		func() test {
			m1 := &CgroupMetrics{
				Mode:          CGV2,
				MemUsageBytes: 1000000,
				MemLimitBytes: 2000000,
				CPUUsageNano:  1000000000,
				CPUQuotaUs:    100000,
				CPUPeriodUs:   100000,
			}
			m2 := &CgroupMetrics{
				Mode:          CGV2,
				MemUsageBytes: 1500000,
				MemLimitBytes: 2000000,
				CPUUsageNano:  1100000000,
				CPUQuotaUs:    100000,
				CPUPeriodUs:   100000,
			}
			deltaTime := 1 * time.Second
			return test{
				name: "success to calculate cgroup stats",
				args: args{
					m1:        m1,
					m2:        m2,
					deltaTime: deltaTime,
				},
				checkFunc: func(w want, stats *CgroupStats) error {
					if stats == nil {
						return errors.New("stats should not be nil")
					}
					if stats.MemoryUsageBytes != 1500000 {
						return errors.Errorf("memory usage: got %d, want %d", stats.MemoryUsageBytes, 1500000)
					}
					if stats.MemoryLimitBytes != 2000000 {
						return errors.Errorf("memory limit: got %d, want %d", stats.MemoryLimitBytes, 2000000)
					}
					if stats.CPULimitCores != 1.0 {
						return errors.Errorf("cpu limit cores: got %f, want %f", stats.CPULimitCores, 1.0)
					}
					return nil
				},
			}
		}(),
		func() test {
			m1 := &CgroupMetrics{
				Mode:          CGV1,
				MemUsageBytes: 500000,
				MemLimitBytes: 1000000,
				CPUUsageNano:  500000000,
				CPUQuotaUs:    0,
				CPUPeriodUs:   100000,
			}
			m2 := &CgroupMetrics{
				Mode:          CGV1,
				MemUsageBytes: 600000,
				MemLimitBytes: 1000000,
				CPUUsageNano:  600000000,
				CPUQuotaUs:    0,
				CPUPeriodUs:   100000,
			}
			deltaTime := 1 * time.Second
			return test{
				name: "calculate stats with unlimited CPU quota",
				args: args{
					m1:        m1,
					m2:        m2,
					deltaTime: deltaTime,
				},
				checkFunc: func(w want, stats *CgroupStats) error {
					if stats == nil {
						return errors.New("stats should not be nil")
					}
					if stats.CPULimitCores != 0 {
						return errors.Errorf("cpu limit cores should be 0 for unlimited quota, got %f", stats.CPULimitCores)
					}
					return nil
				},
			}
		}(),
		func() test {
			m1 := &CgroupMetrics{
				Mode:          CGV2,
				MemUsageBytes: 1000000,
				MemLimitBytes: 0,
				CPUUsageNano:  1000000000,
				CPUQuotaUs:    0,
				CPUPeriodUs:   0,
			}
			m2 := &CgroupMetrics{
				Mode:          CGV2,
				MemUsageBytes: 1200000,
				MemLimitBytes: 0,
				CPUUsageNano:  900000000,
				CPUQuotaUs:    0,
				CPUPeriodUs:   0,
			}
			deltaTime := 1 * time.Second
			return test{
				name: "calculate stats with zero quota and negative CPU delta",
				args: args{
					m1:        m1,
					m2:        m2,
					deltaTime: deltaTime,
				},
				checkFunc: func(w want, stats *CgroupStats) error {
					if stats == nil {
						return errors.New("stats should not be nil")
					}
					if stats.CPULimitCores != 0 {
						return errors.Errorf("cpu limit cores should be 0 for zero quota, got %f", stats.CPULimitCores)
					}
					if stats.CPUUsageCores != 0 {
						return errors.Errorf("cpu usage cores should be 0 for negative delta, got %f", stats.CPUUsageCores)
					}
					return nil
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			stats := calculateCPUUsageCores(test.args.m1, test.args.m2, test.args.deltaTime)
			if err := test.checkFunc(test.want, &stats); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_readCgroupMetrics(t *testing.T) {
	// t.Parallel()
	type want struct {
		metrics *CgroupMetrics
		err     error
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, *CgroupMetrics, error) error
		beforeFunc func()
		afterFunc  func()
	}

	origCgroupBasePath := cgroupBasePath
	origProcCgroupPath := procCgroupPath

	defaultCheckFunc := func(w want, metrics *CgroupMetrics, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if metrics == nil && w.metrics != nil {
			return errors.New("metrics is nil but want is not nil")
		}
		if metrics != nil && w.metrics == nil {
			return errors.New("metrics is not nil but want is nil")
		}
		return nil
	}
	tests := []test{
		func() test {
			tmpDir := t.TempDir()
			procFile := filepath.Join(tmpDir, "cgroup")

			return test{
				name: "successfully reads cgroup v2 metrics",
				want: want{
					metrics: &CgroupMetrics{
						Mode:          CGV2,
						MemUsageBytes: 1024,
						MemLimitBytes: 2048,
						CPUUsageNano:  1000000,
						CPUQuotaUs:    1000,
						CPUPeriodUs:   1000,
					},
				},
				beforeFunc: func() {
					cgroupBasePath = tmpDir
					procCgroupPath = procFile

					// Setup cgroup v2 files in root
					_ = os.WriteFile(filepath.Join(tmpDir, "cgroup.controllers"), []byte("memory cpu"), 0644)
					_ = os.WriteFile(filepath.Join(tmpDir, "memory.current"), []byte("1024"), 0644)
					_ = os.WriteFile(filepath.Join(tmpDir, "memory.max"), []byte("2048"), 0644)
					_ = os.WriteFile(filepath.Join(tmpDir, "cpu.stat"), []byte("usage_usec 1000"), 0644)
					_ = os.WriteFile(filepath.Join(tmpDir, "cpu.max"), []byte("1000 1000"), 0644)

					_ = os.WriteFile(procFile, []byte("0::/\n"), 0644)
				},
				afterFunc: func() {
					cgroupBasePath = origCgroupBasePath
					procCgroupPath = origProcCgroupPath
				},
				checkFunc: func(w want, metrics *CgroupMetrics, err error) error {
					if err != nil {
						return errors.Errorf("unexpected error: %v", err)
					}
					if metrics == nil {
						return errors.New("metrics should not be nil")
					}
					if metrics.Mode != CGV2 {
						return errors.Errorf("expected CGV2, got %v", metrics.Mode)
					}
					if metrics.MemUsageBytes != w.metrics.MemUsageBytes {
						return errors.Errorf("mem usage: got %d, want %d", metrics.MemUsageBytes, w.metrics.MemUsageBytes)
					}
					return nil
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			metrics, err := readCgroupMetrics()
			if err := test.checkFunc(test.want, metrics, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_measureCgroupStats(t *testing.T) {
	// t.Parallel()
	type args struct {
		ctx context.Context
	}
	type want struct {
		stats *CgroupStats
		err   error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *CgroupStats, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	origCgroupBasePath := cgroupBasePath
	origProcCgroupPath := procCgroupPath

	defaultCheckFunc := func(w want, stats *CgroupStats, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if stats == nil && w.stats != nil {
			return errors.New("stats is nil but want is not nil")
		}
		if stats != nil && w.stats == nil {
			return errors.New("stats is not nil but want is nil")
		}
		return nil
	}
	tests := []test{
		func() test {
			tmpDir := t.TempDir()
			procFile := filepath.Join(tmpDir, "cgroup")

			return test{
				name: "successfully measures cgroup stats",
				args: args{
					ctx: context.Background(),
				},
				beforeFunc: func(a args) {
					cgroupBasePath = tmpDir
					procCgroupPath = procFile

					// Setup cgroup v2 files
					_ = os.WriteFile(filepath.Join(tmpDir, "cgroup.controllers"), []byte("memory cpu"), 0644)
					_ = os.WriteFile(filepath.Join(tmpDir, "memory.current"), []byte("1000"), 0644)
					_ = os.WriteFile(filepath.Join(tmpDir, "memory.max"), []byte("2000"), 0644)
					_ = os.WriteFile(filepath.Join(tmpDir, "cpu.stat"), []byte("usage_usec 1000"), 0644)
					_ = os.WriteFile(filepath.Join(tmpDir, "cpu.max"), []byte("1000 1000"), 0644)

					_ = os.WriteFile(procFile, []byte("0::/\n"), 0644)
				},
				afterFunc: func(a args) {
					cgroupBasePath = origCgroupBasePath
					procCgroupPath = origProcCgroupPath
				},
				checkFunc: func(w want, stats *CgroupStats, err error) error {
					if err != nil && !errors.Is(err, context.DeadlineExceeded) && !errors.Is(err, context.Canceled) {
						return errors.Errorf("unexpected error: %v", err)
					}
					if stats == nil {
						return errors.New("stats should not be nil")
					}
					if stats.MemoryUsageBytes == 0 {
						return errors.New("memory usage should be greater than 0")
					}
					return nil
				},
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			tmpDir := t.TempDir()
			procFile := filepath.Join(tmpDir, "cgroup")

			return test{
				name: "context canceled during measurement",
				args: args{
					ctx: ctx,
				},
				beforeFunc: func(a args) {
					cgroupBasePath = tmpDir
					procCgroupPath = procFile

					// Setup minimal files to pass first read
					_ = os.WriteFile(filepath.Join(tmpDir, "cgroup.controllers"), []byte("memory cpu"), 0644)
					_ = os.WriteFile(filepath.Join(tmpDir, "memory.current"), []byte("1000"), 0644)
					_ = os.WriteFile(filepath.Join(tmpDir, "memory.max"), []byte("2000"), 0644)
					_ = os.WriteFile(filepath.Join(tmpDir, "cpu.stat"), []byte("usage_usec 1000"), 0644)
					_ = os.WriteFile(filepath.Join(tmpDir, "cpu.max"), []byte("1000 1000"), 0644)

					_ = os.WriteFile(procFile, []byte("0::/\n"), 0644)
				},
				afterFunc: func(a args) {
					cgroupBasePath = origCgroupBasePath
					procCgroupPath = origProcCgroupPath
				},
				checkFunc: func(w want, stats *CgroupStats, err error) error {
					if err == nil {
						return errors.New("expected context cancellation error")
					}
					if stats != nil {
						return errors.New("stats should be nil when context is canceled")
					}
					return nil
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			stats, err := measureCgroupStats(test.args.ctx)
			if err := test.checkFunc(test.want, stats, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_readCgroupV2Metrics(t *testing.T) {
	// t.Parallel()
	type want struct {
		metrics *CgroupMetrics
		err     error
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, *CgroupMetrics, error) error
		beforeFunc func()
		afterFunc  func()
	}

	origCgroupBasePath := cgroupBasePath
	origProcCgroupPath := procCgroupPath

	defaultCheckFunc := func(w want, metrics *CgroupMetrics, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if metrics == nil && w.metrics != nil {
			return errors.New("metrics is nil but want is not nil")
		}
		if metrics != nil && w.metrics == nil {
			return errors.New("metrics is not nil but want is nil")
		}
		return nil
	}
	tests := []test{
		func() test {
			tmpDir := t.TempDir()
			procFile := filepath.Join(tmpDir, "cgroup")

			return test{
				name: "reads cgroup v2 metrics when available",
				want: want{
					metrics: &CgroupMetrics{
						Mode:          CGV2,
						MemUsageBytes: 123456,
						MemLimitBytes: 987654,
						CPUUsageNano:  123000000,
						CPUQuotaUs:    50000,
						CPUPeriodUs:   100000,
					},
				},
				beforeFunc: func() {
					cgroupBasePath = tmpDir
					procCgroupPath = procFile

					// Setup cgroup v2 files
					_ = os.WriteFile(filepath.Join(tmpDir, "cgroup.controllers"), []byte("memory cpu"), 0644)
					_ = os.WriteFile(filepath.Join(tmpDir, "memory.current"), []byte("123456"), 0644)
					_ = os.WriteFile(filepath.Join(tmpDir, "memory.max"), []byte("987654"), 0644)
					_ = os.WriteFile(filepath.Join(tmpDir, "cpu.stat"), []byte("usage_usec 123000\n"), 0644)
					_ = os.WriteFile(filepath.Join(tmpDir, "cpu.max"), []byte("50000 100000"), 0644)

					_ = os.WriteFile(procFile, []byte("0::/\n"), 0644)
				},
				afterFunc: func() {
					cgroupBasePath = origCgroupBasePath
					procCgroupPath = origProcCgroupPath
				},
				checkFunc: func(w want, metrics *CgroupMetrics, err error) error {
					if err != nil {
						return errors.Errorf("unexpected error: %v", err)
					}
					if metrics == nil {
						return errors.New("metrics should not be nil")
					}
					if metrics.Mode != CGV2 {
						return errors.Errorf("expected CGV2 mode, got %v", metrics.Mode)
					}
					if metrics.MemUsageBytes != w.metrics.MemUsageBytes {
						return errors.Errorf("mem usage: got %d, want %d", metrics.MemUsageBytes, w.metrics.MemUsageBytes)
					}
					return nil
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			metrics, err := readCgroupV2Metrics()
			if err := test.checkFunc(test.want, metrics, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func Test_readCgroupV1Metrics(t *testing.T) {
// 	type want struct {
// 		wantMetrics *CgroupMetrics
// 		err         error
// 	}
// 	type test struct {
// 		name       string
// 		want       want
// 		checkFunc  func(want, *CgroupMetrics, error) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, gotMetrics *CgroupMetrics, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotMetrics, w.wantMetrics) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotMetrics, w.wantMetrics)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			gotMetrics, err := readCgroupV1Metrics()
// 			if err := checkFunc(test.want, gotMetrics, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_calculateCPUUsageCores(t *testing.T) {
// 	type args struct {
// 		m1        *CgroupMetrics
// 		m2        *CgroupMetrics
// 		deltaTime time.Duration
// 	}
// 	type want struct {
// 		wantCalculatedStats CgroupStats
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, CgroupStats) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotCalculatedStats CgroupStats) error {
// 		if !reflect.DeepEqual(gotCalculatedStats, w.wantCalculatedStats) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotCalculatedStats, w.wantCalculatedStats)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           m1:CgroupMetrics{},
// 		           m2:CgroupMetrics{},
// 		           deltaTime:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           m1:CgroupMetrics{},
// 		           m2:CgroupMetrics{},
// 		           deltaTime:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			gotCalculatedStats := calculateCPUUsageCores(test.args.m1, test.args.m2, test.args.deltaTime)
// 			if err := checkFunc(test.want, gotCalculatedStats); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
