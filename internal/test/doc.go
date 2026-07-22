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

// Package test is the facade of Vald's layered testing framework. The
// implementation lives in three single-responsibility subpackages, and this
// package re-exports their APIs (generic type aliases plus thin forwarding
// functions) so call sites can use everything through one import:
//
//   - capability: the testing.TB capability layer — the Runner[X]
//     constraint (a testing.TB that spawns subtests of its own type), the
//     errors.As-style As[C] probe resolving benchmark-/test-only surfaces
//     through wrapper chains, the named helpers built on it (IsBenchmark,
//     Loop, Measured, ReportMetric, timer control, Parallel), and the
//     type-erased Node.
//   - table: the table-driven core — CaseFor/Result/DoFor and Run, generic
//     over Runner[X] so the same tables drive tests and benchmarks.
//   - testdata: fixture path resolution under internal/test/data, keeping
//     filesystem dependencies out of the framework layers.
//
// Typical use stays unchanged from the historical API: test.Run with
// test.Case rows in tests, test.BenchmarkCase in benchmarks, test.NewNode
// for non-generic orchestration trees, and test.GetTestdataPath for
// fixtures. Code needing a capability without a named helper probes for its
// own narrow interface via test.As.
package test
