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

package vald

import "time"

const (
	TimeFormat                                   = time.RFC3339Nano
	UncommittedAnnotationsKey                    = "vald.vdaas.org/uncommitted"
	UnsavedProcessedVQAnnotationsKey             = "vald.vdaas.org/unsaved-processed-vq"
	UnsavedCreateIndexExecutionNumAnnotationsKey = "vald.vdaas.org/unsaved-create-index-execution"
	LastTimeSaveIndexTimestampAnnotationsKey     = "vald.vdaas.org/last-time-save-index-timestamp"
	IndexCountAnnotationsKey                     = "vald.vdaas.org/index-count"
	LastTimeSnapshotTimestampAnnotationsKey      = "vald.vdaas.org/last-time-snapshot-timestamp"
)
