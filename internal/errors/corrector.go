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

// Package errors provides error types and function
package errors

// ErrIndexReplicaOne represents an error that nothing to correct when index replica is 1.
var ErrIndexReplicaOne = New("nothing to correct when index replica is 1")

// ErrNoAvailableAgentToInsert represents an error that no available agent to insert replica.
var ErrNoAvailableAgentToInsert = New("no available agent to insert replica")

// ErrFailedToCorrectReplicaNum represents an error that failed to correct replica number after correction process.
var ErrFailedToCorrectReplicaNum = New("failed to correct replica number after correction process")
