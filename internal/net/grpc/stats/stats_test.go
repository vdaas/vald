//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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
	"testing"

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

// func Test_getCPUStats(t *testing.T) {
// 	t.Parallel()
// 	type args struct {
// 		path string
// 	}
// 	type want struct {
// 		wantErr bool
// 		wantMin float64
// 		wantMax float64
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, float64, error) error
// 		beforeFunc func(*testing.T) string
// 		afterFunc  func(*testing.T, string)
// 		setupFunc  func(*testing.T) *server
// 	}
// 	defaultCheckFunc := func(w want, got float64, err error) error {
// 		if (err != nil) != w.wantErr {
// 			return errors.Errorf("wantErr = %v, error = %v", w.wantErr, err)
// 		}
// 		if !w.wantErr {
// 			if got < w.wantMin || got > w.wantMax {
// 				return errors.Errorf("CPU usage should be between %f and %f, got: %f", w.wantMin, w.wantMax, got)
// 			}
// 		}
// 		return nil
// 	}

// 	tests := []test{
// 		{
// 			name: "success to get CPU stats with valid /proc/stat format",
// 			want: want{
// 				wantErr: false,
// 				wantMin: 0.0,
// 				wantMax: 100.0,
// 			},
// 			beforeFunc: func(t *testing.T) string {
// 				tmpFile, err := os.CreateTemp("", "proc_stat_test_*.txt")
// 				if err != nil {
// 					t.Fatal(err)
// 				}

// 				content := `cpu  123456 0 78910 987654 0 0 0 0 0 0
// cpu0 61728 0 39455 493827 0 0 0 0 0 0
// intr 123456
// ctxt 789012
// `
// 				if _, err := tmpFile.WriteString(content); err != nil {
// 					t.Fatal(err)
// 				}
// 				tmpFile.Close()

// 				return tmpFile.Name()
// 			},
// 			afterFunc: func(t *testing.T, path string) {
// 				os.Remove(path)
// 			},
// 			checkFunc: defaultCheckFunc,
// 		},
// 		{
// 			name: "return error when file not found",
// 			args: args{
// 				path: "/nonexistent/file",
// 			},
// 			want: want{
// 				wantErr: true,
// 			},
// 			checkFunc: defaultCheckFunc,
// 		},
// 		{
// 			name: "return error when cpu line not found",
// 			want: want{
// 				wantErr: true,
// 			},
// 			beforeFunc: func(t *testing.T) string {
// 				tmpFile, err := os.CreateTemp("", "proc_stat_test_*.txt")
// 				if err != nil {
// 					t.Fatal(err)
// 				}

// 				content := `intr 123456
// ctxt 789012
// btime 1234567890
// `
// 				if _, err := tmpFile.WriteString(content); err != nil {
// 					t.Fatal(err)
// 				}
// 				tmpFile.Close()

// 				return tmpFile.Name()
// 			},
// 			afterFunc: func(t *testing.T, path string) {
// 				os.Remove(path)
// 			},
// 			checkFunc: defaultCheckFunc,
// 		},
// 		{
// 			name: "return error when cpu line format is invalid",
// 			want: want{
// 				wantErr: true,
// 			},
// 			beforeFunc: func(t *testing.T) string {
// 				tmpFile, err := os.CreateTemp("", "proc_stat_test_*.txt")
// 				if err != nil {
// 					t.Fatal(err)
// 				}

// 				content := `cpu  123 456
// `
// 				if _, err := tmpFile.WriteString(content); err != nil {
// 					t.Fatal(err)
// 				}
// 				tmpFile.Close()

// 				return tmpFile.Name()
// 			},
// 			afterFunc: func(t *testing.T, path string) {
// 				os.Remove(path)
// 			},
// 			checkFunc: defaultCheckFunc,
// 		},
// 		{
// 			name: "return error when cpu field parse fails",
// 			want: want{
// 				wantErr: true,
// 			},
// 			beforeFunc: func(t *testing.T) string {
// 				tmpFile, err := os.CreateTemp("", "proc_stat_test_*.txt")
// 				if err != nil {
// 					t.Fatal(err)
// 				}

// 				content := `cpu  invalid 0 78910 987654 0 0 0 0 0 0
// `
// 				if _, err := tmpFile.WriteString(content); err != nil {
// 					t.Fatal(err)
// 				}
// 				tmpFile.Close()

// 				return tmpFile.Name()
// 			},
// 			afterFunc: func(t *testing.T, path string) {
// 				os.Remove(path)
// 			},
// 			checkFunc: defaultCheckFunc,
// 		},
// 		{
// 			name: "return 0 usage when deltaTotal is 0",
// 			want: want{
// 				wantErr: false,
// 				wantMin: 0.0,
// 				wantMax: 0.0,
// 			},
// 			beforeFunc: func(t *testing.T) string {
// 				tmpFile, err := os.CreateTemp("", "proc_stat_test_*.txt")
// 				if err != nil {
// 					t.Fatal(err)
// 				}

// 				content := `cpu  123456 0 78910 987654 0 0 0 0 0 0
// `
// 				if _, err := tmpFile.WriteString(content); err != nil {
// 					t.Fatal(err)
// 				}
// 				tmpFile.Close()

// 				return tmpFile.Name()
// 			},
// 			afterFunc: func(t *testing.T, path string) {
// 				os.Remove(path)
// 			},
// 			checkFunc: defaultCheckFunc,
// 		},
// 		{
// 			name: "first call behavior - initialize stats",
// 			want: want{
// 				wantErr: false,
// 				wantMin: 0.0,
// 				wantMax: 100.0,
// 			},
// 			beforeFunc: func(t *testing.T) string {
// 				tmpFile, err := os.CreateTemp("", "proc_stat_test_*.txt")
// 				if err != nil {
// 					t.Fatal(err)
// 				}

// 				content := `cpu  100000 0 50000 800000 0 0 0 0 0 0
// cpu0 50000 0 25000 400000 0 0 0 0 0 0
// `
// 				if _, err := tmpFile.WriteString(content); err != nil {
// 					t.Fatal(err)
// 				}
// 				tmpFile.Close()

// 				return tmpFile.Name()
// 			},
// 			afterFunc: func(t *testing.T, path string) {
// 				os.Remove(path)
// 			},
// 			checkFunc: func(w want, got float64, err error) error {
// 				if err := defaultCheckFunc(w, got, err); err != nil {
// 					return err
// 				}
// 				return nil
// 			},
// 		},
// 		{
// 			name: "subsequent calls - calculate usage from previous stats",
// 			want: want{
// 				wantErr: false,
// 				wantMin: 0.0,
// 				wantMax: 100.0,
// 			},
// 			setupFunc: func(t *testing.T) *server {
// 				s := &server{}
// 				tmpFile1, err := os.CreateTemp("", "proc_stat_setup_*.txt")
// 				if err != nil {
// 					t.Fatal(err)
// 				}
// 				defer os.Remove(tmpFile1.Name())

// 				content1 := `cpu  100000 0 50000 800000 0 0 0 0 0 0
// `
// 				if _, err := tmpFile1.WriteString(content1); err != nil {
// 					t.Fatal(err)
// 				}
// 				tmpFile1.Close()

// 				_, err = s.getCPUStats(tmpFile1.Name())
// 				if err != nil {
// 					t.Fatal(err)
// 				}

// 				time.Sleep(10 * time.Millisecond)
// 				return s
// 			},
// 			beforeFunc: func(t *testing.T) string {
// 				tmpFile2, err := os.CreateTemp("", "proc_stat_test_*.txt")
// 				if err != nil {
// 					t.Fatal(err)
// 				}

// 				content2 := `cpu  100200 0 50100 800100 0 0 0 0 0 0
// `
// 				if _, err := tmpFile2.WriteString(content2); err != nil {
// 					t.Fatal(err)
// 				}
// 				tmpFile2.Close()

// 				return tmpFile2.Name()
// 			},
// 			afterFunc: func(t *testing.T, path string) {
// 				os.Remove(path)
// 			},
// 			checkFunc: func(w want, got float64, err error) error {
// 				if err := defaultCheckFunc(w, got, err); err != nil {
// 					return err
// 				}
// 				return nil
// 			},
// 		},
// 	}

// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			var testFilePath string
// 			var s *server

// 			if test.setupFunc != nil {
// 				s = test.setupFunc(tt)
// 			} else {
// 				s = &server{}
// 			}

// 			if test.beforeFunc != nil {
// 				testFilePath = test.beforeFunc(tt)
// 				test.args.path = testFilePath
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, testFilePath)
// 			}
// 			if test.checkFunc == nil {
// 				test.checkFunc = defaultCheckFunc
// 			}

// 			got, err := s.getCPUStats(test.args.path)
// 			if err := test.checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}

// 			if test.name == "first call behavior - initialize stats" {
// 				if s.statsLoadedCnt.Load() < 1 {
// 					tt.Errorf("statsLoadedCnt should be at least 1, got: %d", s.statsLoadedCnt.Load())
// 				}
// 			}
// 			if test.name == "subsequent calls - calculate usage from previous stats" {
// 				if s.statsLoadedCnt.Load() < 2 {
// 					tt.Errorf("statsLoadedCnt should be at least 2, got: %d", s.statsLoadedCnt.Load())
// 				}
// 			}
// 		})
// 	}
// }

// func Test_getMemoryStats(t *testing.T) {
// 	t.Parallel()
// 	type want struct {
// 		wantErr bool
// 		wantMin float64
// 		wantMax float64
// 	}
// 	type test struct {
// 		name       string
// 		want       want
// 		checkFunc  func(want, float64, error) error
// 		beforeFunc func(*testing.T) (string, string)
// 		afterFunc  func(*testing.T, string, string)
// 	}
// 	defaultCheckFunc := func(w want, got float64, err error) error {
// 		if (err != nil) != w.wantErr {
// 			return errors.Errorf("wantErr = %v, error = %v", w.wantErr, err)
// 		}
// 		if !w.wantErr {
// 			if got < w.wantMin || got > w.wantMax {
// 				return errors.Errorf("memory usage should be between %f and %f, got: %f", w.wantMin, w.wantMax, got)
// 			}
// 		}
// 		return nil
// 	}

// 	tests := []test{
// 		{
// 			name: "success to get memory stats with valid data",
// 			want: want{
// 				wantErr: false,
// 				wantMin: 0.0,
// 				wantMax: 100.0,
// 			},
// 			beforeFunc: func(t *testing.T) (string, string) {
// 				meminfoFile, err := os.CreateTemp("", "meminfo_test_*.txt")
// 				if err != nil {
// 					t.Fatal(err)
// 				}
// 				meminfoContent := `MemTotal:        8142848 kB
// MemFree:         1234567 kB
// MemAvailable:    5678901 kB
// Buffers:          123456 kB
// `
// 				if _, err := meminfoFile.WriteString(meminfoContent); err != nil {
// 					t.Fatal(err)
// 				}
// 				meminfoFile.Close()

// 				statusFile, err := os.CreateTemp("", "status_test_*.txt")
// 				if err != nil {
// 					t.Fatal(err)
// 				}
// 				statusContent := `Name:   test_process
// Pid:    12345
// VmRSS:   81428 kB
// VmSize: 164856 kB
// `
// 				if _, err := statusFile.WriteString(statusContent); err != nil {
// 					t.Fatal(err)
// 				}
// 				statusFile.Close()

// 				return meminfoFile.Name(), statusFile.Name()
// 			},
// 			afterFunc: func(t *testing.T, meminfoPath, statusPath string) {
// 				os.Remove(meminfoPath)
// 				os.Remove(statusPath)
// 			},
// 			checkFunc: defaultCheckFunc,
// 		},
// 		{
// 			name: "return error when meminfo file not found",
// 			want: want{
// 				wantErr: true,
// 			},
// 			beforeFunc: func(t *testing.T) (string, string) {
// 				statusFile, err := os.CreateTemp("", "status_test_*.txt")
// 				if err != nil {
// 					t.Fatal(err)
// 				}
// 				statusContent := `VmRSS:   81428 kB`
// 				if _, err := statusFile.WriteString(statusContent); err != nil {
// 					t.Fatal(err)
// 				}
// 				statusFile.Close()

// 				return "/nonexistent/meminfo", statusFile.Name()
// 			},
// 			afterFunc: func(t *testing.T, meminfoPath, statusPath string) {
// 				os.Remove(statusPath)
// 			},
// 			checkFunc: defaultCheckFunc,
// 		},
// 		{
// 			name: "return error when status file not found",
// 			want: want{
// 				wantErr: false,
// 				wantMin: 0.0,
// 				wantMax: 2.0,
// 			},
// 			beforeFunc: func(t *testing.T) (string, string) {
// 				meminfoFile, err := os.CreateTemp("", "meminfo_test_*.txt")
// 				if err != nil {
// 					t.Fatal(err)
// 				}
// 				meminfoContent := `MemTotal:        8142848 kB`
// 				if _, err := meminfoFile.WriteString(meminfoContent); err != nil {
// 					t.Fatal(err)
// 				}
// 				meminfoFile.Close()

// 				return meminfoFile.Name(), "/nonexistent/status"
// 			},
// 			afterFunc: func(t *testing.T, meminfoPath, statusPath string) {
// 				os.Remove(meminfoPath)
// 			},
// 			checkFunc: func(w want, got float64, err error) error {
// 				if (err != nil) != w.wantErr {
// 					return errors.Errorf("wantErr = %v, error = %v", w.wantErr, err)
// 				}
// 				if !w.wantErr && (got < w.wantMin || got > w.wantMax) {
// 					return errors.Errorf("memory usage should be between %f and %f, got: %f", w.wantMin, w.wantMax, got)
// 				}
// 				return nil
// 			},
// 		},
// 	}

// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			var meminfoPath, statusPath string

// 			if test.beforeFunc != nil {
// 				meminfoPath, statusPath = test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, meminfoPath, statusPath)
// 			}
// 			if test.checkFunc == nil {
// 				test.checkFunc = defaultCheckFunc
// 			}

// 			originalMemoryStatsPath := memoryStatsPath
// 			originalProcessMemoryStatsPath := processMemoryStatsPath
// 			defer func() {
// 				memoryStatsPath = originalMemoryStatsPath
// 				processMemoryStatsPath = originalProcessMemoryStatsPath
// 			}()

// 			memoryStatsPath = meminfoPath
// 			processMemoryStatsPath = statusPath

// 			got, err := getMemoryStats()
// 			if err := test.checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }

// func Test_getTotalMemory(t *testing.T) {
// 	t.Parallel()
// 	type want struct {
// 		wantErr bool
// 		want    uint64
// 	}
// 	type test struct {
// 		name       string
// 		want       want
// 		checkFunc  func(want, uint64, error) error
// 		beforeFunc func(*testing.T) string
// 		afterFunc  func(*testing.T, string)
// 	}
// 	defaultCheckFunc := func(w want, got uint64, err error) error {
// 		if (err != nil) != w.wantErr {
// 			return errors.Errorf("wantErr = %v, error = %v", w.wantErr, err)
// 		}
// 		if !w.wantErr && w.want != 0 && got != w.want {
// 			return errors.Errorf("got = %d, want = %d", got, w.want)
// 		}
// 		if !w.wantErr && got == 0 {
// 			return errors.New("total memory should not be zero")
// 		}
// 		return nil
// 	}

// 	tests := []test{
// 		{
// 			name: "success to get total memory",
// 			want: want{
// 				wantErr: false,
// 				want:    8338276352, // 8142848 * 1024
// 			},
// 			beforeFunc: func(t *testing.T) string {
// 				tmpFile, err := os.CreateTemp("", "meminfo_test_*.txt")
// 				if err != nil {
// 					t.Fatal(err)
// 				}

// 				content := `MemTotal:        8142848 kB
// MemFree:         1234567 kB
// MemAvailable:    5678901 kB
// Buffers:          123456 kB
// Cached:          1987654 kB
// `
// 				if _, err := tmpFile.WriteString(content); err != nil {
// 					t.Fatal(err)
// 				}
// 				tmpFile.Close()

// 				return tmpFile.Name()
// 			},
// 			afterFunc: func(t *testing.T, path string) {
// 				os.Remove(path)
// 			},
// 			checkFunc: defaultCheckFunc,
// 		},
// 		{
// 			name: "return error when file not found",
// 			want: want{
// 				wantErr: true,
// 			},
// 			beforeFunc: func(t *testing.T) string {
// 				return "/nonexistent/meminfo"
// 			},
// 			checkFunc: defaultCheckFunc,
// 		},
// 		{
// 			name: "return error when MemTotal not found",
// 			want: want{
// 				wantErr: true,
// 			},
// 			beforeFunc: func(t *testing.T) string {
// 				tmpFile, err := os.CreateTemp("", "meminfo_test_*.txt")
// 				if err != nil {
// 					t.Fatal(err)
// 				}

// 				content := `MemFree:         1234567 kB
// MemAvailable:    5678901 kB
// Buffers:          123456 kB
// `
// 				if _, err := tmpFile.WriteString(content); err != nil {
// 					t.Fatal(err)
// 				}
// 				tmpFile.Close()

// 				return tmpFile.Name()
// 			},
// 			afterFunc: func(t *testing.T, path string) {
// 				os.Remove(path)
// 			},
// 			checkFunc: defaultCheckFunc,
// 		},
// 		{
// 			name: "return error when MemTotal parse fails",
// 			want: want{
// 				wantErr: true,
// 			},
// 			beforeFunc: func(t *testing.T) string {
// 				tmpFile, err := os.CreateTemp("", "meminfo_test_*.txt")
// 				if err != nil {
// 					t.Fatal(err)
// 				}

// 				content := `MemTotal:        invalid kB
// MemFree:         1234567 kB
// `
// 				if _, err := tmpFile.WriteString(content); err != nil {
// 					t.Fatal(err)
// 				}
// 				tmpFile.Close()

// 				return tmpFile.Name()
// 			},
// 			afterFunc: func(t *testing.T, path string) {
// 				os.Remove(path)
// 			},
// 			checkFunc: defaultCheckFunc,
// 		},
// 	}

// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			var testFilePath string

// 			if test.beforeFunc != nil {
// 				testFilePath = test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, testFilePath)
// 			}
// 			if test.checkFunc == nil {
// 				test.checkFunc = defaultCheckFunc
// 			}

// 			got, err := getTotalMemory(testFilePath)
// 			if err := test.checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }

// func Test_getProcessMemory(t *testing.T) {
// 	t.Parallel()
// 	type want struct {
// 		wantErr bool
// 		want    uint64
// 	}
// 	type test struct {
// 		name       string
// 		want       want
// 		checkFunc  func(want, uint64, error) error
// 		beforeFunc func(*testing.T) string
// 		afterFunc  func(*testing.T, string)
// 	}
// 	defaultCheckFunc := func(w want, got uint64, err error) error {
// 		if (err != nil) != w.wantErr {
// 			return errors.Errorf("wantErr = %v, error = %v", w.wantErr, err)
// 		}
// 		if !w.wantErr && w.want != 0 && got != w.want {
// 			return errors.Errorf("got = %d, want = %d", got, w.want)
// 		}
// 		if !w.wantErr && got == 0 {
// 			return errors.New("process memory should not be zero")
// 		}
// 		return nil
// 	}

// 	tests := []test{
// 		{
// 			name: "success to get process memory",
// 			want: want{
// 				wantErr: false,
// 				want:    83382272, // 81428 * 1024
// 			},
// 			beforeFunc: func(t *testing.T) string {
// 				tmpFile, err := os.CreateTemp("", "status_test_*.txt")
// 				if err != nil {
// 					t.Fatal(err)
// 				}

// 				content := `Name:   test_process
// State:  S (sleeping)
// Tgid:   12345
// Pid:    12345
// VmRSS:   81428 kB
// VmSize: 164856 kB
// VmData:  12345 kB
// `
// 				if _, err := tmpFile.WriteString(content); err != nil {
// 					t.Fatal(err)
// 				}
// 				tmpFile.Close()

// 				return tmpFile.Name()
// 			},
// 			afterFunc: func(t *testing.T, path string) {
// 				os.Remove(path)
// 			},
// 			checkFunc: defaultCheckFunc,
// 		},
// 		{
// 			name: "return error when file not found",
// 			want: want{
// 				wantErr: true,
// 			},
// 			beforeFunc: func(t *testing.T) string {
// 				return "/nonexistent/status"
// 			},
// 			checkFunc: defaultCheckFunc,
// 		},
// 		{
// 			name: "return error when VmRSS not found",
// 			want: want{
// 				wantErr: true,
// 			},
// 			beforeFunc: func(t *testing.T) string {
// 				tmpFile, err := os.CreateTemp("", "status_test_*.txt")
// 				if err != nil {
// 					t.Fatal(err)
// 				}

// 				content := `Name:   test_process
// State:  S (sleeping)
// Tgid:   12345
// Pid:    12345
// VmSize: 164856 kB
// `
// 				if _, err := tmpFile.WriteString(content); err != nil {
// 					t.Fatal(err)
// 				}
// 				tmpFile.Close()

// 				return tmpFile.Name()
// 			},
// 			afterFunc: func(t *testing.T, path string) {
// 				os.Remove(path)
// 			},
// 			checkFunc: defaultCheckFunc,
// 		},
// 		{
// 			name: "return error when VmRSS parse fails",
// 			want: want{
// 				wantErr: true,
// 			},
// 			beforeFunc: func(t *testing.T) string {
// 				tmpFile, err := os.CreateTemp("", "status_test_*.txt")
// 				if err != nil {
// 					t.Fatal(err)
// 				}

// 				content := `Name:   test_process
// VmRSS:   invalid kB
// VmSize: 164856 kB
// `
// 				if _, err := tmpFile.WriteString(content); err != nil {
// 					t.Fatal(err)
// 				}
// 				tmpFile.Close()

// 				return tmpFile.Name()
// 			},
// 			afterFunc: func(t *testing.T, path string) {
// 				os.Remove(path)
// 			},
// 			checkFunc: defaultCheckFunc,
// 		},
// 	}

// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			var testFilePath string

// 			if test.beforeFunc != nil {
// 				testFilePath = test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, testFilePath)
// 			}
// 			if test.checkFunc == nil {
// 				test.checkFunc = defaultCheckFunc
// 			}

// 			got, err := getProcessMemory(testFilePath)
// 			if err := test.checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
