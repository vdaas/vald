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

// Package crud provides end-to-end tests using ann-benchmarks datasets.
//
// Every other file in this package carries the `e2e` build tag (see
// crud_test.go's TestMain, which requires a -config flag pointing at a
// live cluster configuration) and is therefore excluded from a default
// (non "-tags e2e") build - by design, this package's tests never run
// without both the build tag and a real cluster config. That is fine for
// `go vet $(ROOTDIR)/...`/`make lint`, which target the repo root via the
// `./...` wildcard and silently skip a directory that matches zero
// packages under the current build tags, but tooling that type-checks this
// package directory directly (e.g. `golangci-lint run ./tests/v2/e2e/crud/`
// without `-tags e2e`) instead hard-fails with "build constraints exclude
// all Go files". This file carries no build tag - contributing no test
// functions and no other symbols - solely so such tooling still finds at
// least one buildable file in the package.
package crud
