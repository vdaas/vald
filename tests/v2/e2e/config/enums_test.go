//go:build e2e

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

package config

import (
	"testing"

	"github.com/vdaas/vald/internal/test"
)

// TestOperationType_Bind is a regression test for OperationType.Bind(): every
// canonical YAML token used across tests/v2/e2e/assets/*.yaml (e.g.
// "create_index" in agent_crud.yaml) must bind back to its OperationType
// constant. This previously silently dropped OpCreateIndex/OpSaveIndex/
// OpCreateAndSaveIndex (the "create_index"/"save_index"/
// "create_and_save_index" tokens had no matching case), which made
// Execution.Bind fail agent_crud.yaml-style configs with "unsupported
// operation type" once e.Type came back as the empty OperationType.
func TestOperationType_Bind(t *testing.T) {
	if err := test.Run(t.Context(), t, func(t *testing.T, ot OperationType) (OperationType, error) {
		t.Helper()
		return ot.Bind()
	}, []test.Case[OperationType, OperationType]{
		{Name: "search", Args: OpSearch, Want: test.Result[OperationType]{Val: OpSearch}},
		{Name: "search_by_id", Args: OpSearchByID, Want: test.Result[OperationType]{Val: OpSearchByID}},
		{Name: "linear_search", Args: OpLinearSearch, Want: test.Result[OperationType]{Val: OpLinearSearch}},
		{Name: "linear_search_by_id", Args: OpLinearSearchByID, Want: test.Result[OperationType]{Val: OpLinearSearchByID}},
		{Name: "insert", Args: OpInsert, Want: test.Result[OperationType]{Val: OpInsert}},
		{Name: "update", Args: OpUpdate, Want: test.Result[OperationType]{Val: OpUpdate}},
		{Name: "upsert", Args: OpUpsert, Want: test.Result[OperationType]{Val: OpUpsert}},
		{Name: "remove", Args: OpRemove, Want: test.Result[OperationType]{Val: OpRemove}},
		{Name: "remove_by_timestamp", Args: OpRemoveByTimestamp, Want: test.Result[OperationType]{Val: OpRemoveByTimestamp}},
		{Name: "object", Args: OpObject, Want: test.Result[OperationType]{Val: OpObject}},
		{Name: "list_object", Args: OpListObject, Want: test.Result[OperationType]{Val: OpListObject}},
		{Name: "timestamp", Args: OpTimestamp, Want: test.Result[OperationType]{Val: OpTimestamp}},
		{Name: "exists", Args: OpExists, Want: test.Result[OperationType]{Val: OpExists}},
		{Name: "index_info", Args: OpIndexInfo, Want: test.Result[OperationType]{Val: OpIndexInfo}},
		{Name: "index_detail", Args: OpIndexDetail, Want: test.Result[OperationType]{Val: OpIndexDetail}},
		{Name: "index_statistics", Args: OpIndexStatistics, Want: test.Result[OperationType]{Val: OpIndexStatistics}},
		{Name: "index_statistics_detail", Args: OpIndexStatisticsDetail, Want: test.Result[OperationType]{Val: OpIndexStatisticsDetail}},
		{Name: "index_property", Args: OpIndexProperty, Want: test.Result[OperationType]{Val: OpIndexProperty}},
		{Name: "flush", Args: OpFlush, Want: test.Result[OperationType]{Val: OpFlush}},
		{
			// Regression case: this previously bound to the empty
			// OperationType ("") instead of OpCreateIndex.
			Name: "create_index",
			Args: OpCreateIndex,
			Want: test.Result[OperationType]{Val: OpCreateIndex},
		},
		{
			Name: "save_index",
			Args: OpSaveIndex,
			Want: test.Result[OperationType]{Val: OpSaveIndex},
		},
		{
			Name: "create_and_save_index",
			Args: OpCreateAndSaveIndex,
			Want: test.Result[OperationType]{Val: OpCreateAndSaveIndex},
		},
		{Name: "kubernetes", Args: OpKubernetes, Want: test.Result[OperationType]{Val: OpKubernetes}},
		{Name: "client", Args: OpClient, Want: test.Result[OperationType]{Val: OpClient}},
		{Name: "wait", Args: OpWait, Want: test.Result[OperationType]{Val: OpWait}},
		{
			// Empty OperationType is rejected up front with ErrInvalidConfig
			// rather than silently passing through.
			Name: "empty type is an error",
			Args: OperationType(""),
			Want: test.Result[OperationType]{Val: "", Err: nil},
			CheckFunc: func(t *testing.T, want test.Result[OperationType], got test.Result[OperationType]) error {
				t.Helper()
				if got.Err == nil {
					t.Error("expected an error for an empty OperationType, got nil")
				}
				if got.Val != "" {
					t.Errorf("expected empty OperationType on error, got %q", got.Val)
				}
				return nil
			},
		},
	}...); err != nil {
		t.Error(err)
	}
}
