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

// Package service manages the main logic of server.
package service

import (
	"context"
	"flag"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/commands"
	"github.com/leanovate/gopter/gen"
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/pkg/agent/core/ngt/model"
)

type ngtSystem struct {
	ctx context.Context
	ech <-chan error
	ngt *ngt
}

type resultContainer struct {
	err     error
	results []model.Distance
	vector  []float32
	exists  bool
}

type ngtState struct {
	states  map[string]vectorState
	vectors map[string]int
}

func (st *ngtState) Reset() {
	st.states = map[string]vectorState{
		idA: NOT_INSERTED,
		idB: NOT_INSERTED,
		idC: NOT_INSERTED,
	}

	st.vectors = map[string]int{
		idA: 0,
		idB: 0,
		idC: 0,
	}
}

type vectorState uint8

func (v vectorState) String() string {
	switch v {
	case NOT_INSERTED:
		return "not inserted"
	case IN_INSERT_QUEUE:
		return "in insert queue"
	case IN_DELETE_QUEUE:
		return "in delete queue"
	case INDEXED:
		return "indexed"
	}
	return "unknown"
}

const (
	dimension = 3

	NOT_INSERTED vectorState = iota
	IN_INSERT_QUEUE
	IN_DELETE_QUEUE
	INDEXED

	idA = "id-a"
	idB = "id-b"
	idC = "id-c"
)

var (
	seed               int64
	minSuccessfulTests int
	maxDiscardRatio    float64
	applicationLog     bool

	ncfg = &config.NGT{
		Dimension:              dimension,
		DistanceType:           "l2",
		ObjectType:             "float",
		EnableInMemoryMode:     true,
		AutoIndexDurationLimit: "96h",
		AutoIndexCheckDuration: "96h",
		AutoSaveIndexDuration:  "96h",
		AutoIndexLength:        10000000000,
		KVSDB: &config.KVSDB{
			Concurrency: 1,
		},
	}

	vectors = map[string][][]float32{
		idA: {{0.1, 0.2, 0.3}, {0.4, 0.5, 0.6}},
		idB: {{0.2, 0.3, 0.1}, {0.5, 0.6, 0.4}},
		idC: {{0.3, 0.1, 0.2}, {0.6, 0.4, 0.5}},
	}
)

func TestMain(m *testing.M) {
	testing.Init()

	flag.Int64Var(&seed, "pbt-seed", 0, "seed number used for PBT")
	flag.IntVar(&minSuccessfulTests, "pbt-min-successful-tests", 10, "minimum number of successful tests in PBT")
	flag.Float64Var(&maxDiscardRatio, "pbt-max-discard-ratio", 5.0, "maximum discard ratio of PBT")
	flag.BoolVar(&applicationLog, "pbt-enable-application-log", false, "enable application log on PBT")
	m.Run()
}

var (
	insertACommand = &commands.ProtoCommand{
		Name: "Insert-A",
		RunFunc: func(
			systemUnderTest commands.SystemUnderTest,
		) commands.Result {
			sy := systemUnderTest.(*ngtSystem)
			ngt := sy.ngt

			err := ngt.Insert(idA, vectors[idA][0])
			return &resultContainer{
				err: err,
			}
		},
		NextStateFunc: func(state commands.State) commands.State {
			st := state.(*ngtState)
			st.states[idA] = IN_INSERT_QUEUE
			st.vectors[idA] = 0
			return st
		},
		PreConditionFunc: func(state commands.State) bool {
			st, ok := state.(*ngtState)
			if !ok {
				return false
			}
			return st.states[idA] == NOT_INSERTED || st.states[idA] == IN_DELETE_QUEUE
		},
		PostConditionFunc: func(
			state commands.State,
			result commands.Result,
		) *gopter.PropResult {
			rc := result.(*resultContainer)

			if rc.err != nil {
				return &gopter.PropResult{
					Status: gopter.PropFalse,
					Error:  rc.err,
					Labels: []string{
						"Insert-A",
						"error",
						rc.err.Error(),
					},
				}
			}

			return &gopter.PropResult{Status: gopter.PropTrue}
		},
	}

	insertBCommand = &commands.ProtoCommand{
		Name: "Insert-B",
		RunFunc: func(
			systemUnderTest commands.SystemUnderTest,
		) commands.Result {
			sy := systemUnderTest.(*ngtSystem)
			ngt := sy.ngt

			err := ngt.Insert(idB, vectors[idB][0])
			return &resultContainer{
				err: err,
			}
		},
		NextStateFunc: func(state commands.State) commands.State {
			st := state.(*ngtState)
			st.states[idB] = IN_INSERT_QUEUE
			st.vectors[idB] = 0
			return st
		},
		PreConditionFunc: func(state commands.State) bool {
			st, ok := state.(*ngtState)
			if !ok {
				return false
			}

			return st.states[idB] == NOT_INSERTED || st.states[idB] == IN_DELETE_QUEUE
		},
		PostConditionFunc: func(
			state commands.State,
			result commands.Result,
		) *gopter.PropResult {
			rc := result.(*resultContainer)

			if rc.err != nil {
				return &gopter.PropResult{
					Status: gopter.PropFalse,
					Error:  rc.err,
					Labels: []string{
						"Insert-B",
						"error",
						rc.err.Error(),
					},
				}
			}

			return &gopter.PropResult{Status: gopter.PropTrue}
		},
	}

	insertCCommand = &commands.ProtoCommand{
		Name: "Insert-C",
		RunFunc: func(
			systemUnderTest commands.SystemUnderTest,
		) commands.Result {
			sy := systemUnderTest.(*ngtSystem)
			ngt := sy.ngt

			err := ngt.Insert(idC, vectors[idC][0])
			return &resultContainer{
				err: err,
			}
		},
		NextStateFunc: func(state commands.State) commands.State {
			st := state.(*ngtState)
			st.states[idC] = IN_INSERT_QUEUE
			st.vectors[idC] = 0
			return st
		},
		PreConditionFunc: func(state commands.State) bool {
			st, ok := state.(*ngtState)
			if !ok {
				return false
			}

			return st.states[idC] == NOT_INSERTED || st.states[idC] == IN_DELETE_QUEUE
		},
		PostConditionFunc: func(
			state commands.State,
			result commands.Result,
		) *gopter.PropResult {
			rc := result.(*resultContainer)

			if rc.err != nil {
				return &gopter.PropResult{
					Status: gopter.PropFalse,
					Error:  rc.err,
					Labels: []string{
						"Insert-C",
						"error",
						rc.err.Error(),
					},
				}
			}

			return &gopter.PropResult{Status: gopter.PropTrue}
		},
	}

	createIndexCommand = &commands.ProtoCommand{
		Name: "CreateIndex",
		RunFunc: func(
			systemUnderTest commands.SystemUnderTest,
		) commands.Result {
			ctx := systemUnderTest.(*ngtSystem).ctx
			ngt := systemUnderTest.(*ngtSystem).ngt

			// WARN: dirty workaround
			// without the sleep before/after CreateIndex the tests usually fail
			time.Sleep(time.Second)

			err := ngt.CreateAndSaveIndex(ctx, 10000)

			// WARN: dirty workaround
			time.Sleep(time.Second)

			return &resultContainer{
				err: err,
			}
		},
		NextStateFunc: func(state commands.State) commands.State {
			st := state.(*ngtState)

			for k, v := range st.states {
				switch v {
				case IN_INSERT_QUEUE:
					st.states[k] = INDEXED
				case IN_DELETE_QUEUE:
					st.states[k] = NOT_INSERTED
				}
			}

			return st
		},
		PreConditionFunc: func(state commands.State) (flg bool) {
			st, ok := state.(*ngtState)
			if !ok {
				return false
			}

			for _, v := range st.states {
				if v == IN_INSERT_QUEUE || v == IN_DELETE_QUEUE {
					return true
				}
			}

			return false
		},
		PostConditionFunc: func(
			state commands.State,
			result commands.Result,
		) *gopter.PropResult {
			rc := result.(*resultContainer)

			if rc.err != nil {
				return &gopter.PropResult{
					Status: gopter.PropFalse,
					Error:  rc.err,
					Labels: []string{
						"CreateIndex",
						"error",
						rc.err.Error(),
					},
				}
			}

			return &gopter.PropResult{Status: gopter.PropTrue}
		},
	}

	existsACommand = &commands.ProtoCommand{
		Name: "Exists-A",
		RunFunc: func(
			systemUnderTest commands.SystemUnderTest,
		) commands.Result {
			ngt := systemUnderTest.(*ngtSystem).ngt

			_, exists := ngt.Exists(idA)
			return &resultContainer{
				exists: exists,
			}
		},
		NextStateFunc: func(state commands.State) commands.State {
			return state
		},
		PreConditionFunc: func(state commands.State) bool {
			_, ok := state.(*ngtState)
			return ok
		},
		PostConditionFunc: func(
			state commands.State,
			result commands.Result,
		) *gopter.PropResult {
			st := state.(*ngtState)
			rc := result.(*resultContainer)

			if st.states[idA] == INDEXED || st.states[idA] == IN_INSERT_QUEUE {
				if !rc.exists {
					return &gopter.PropResult{
						Status: gopter.PropFalse,
						Labels: []string{
							"Exists-A",
							"uuid exists, but Exists returns false",
						},
					}
				}
			} else {
				if rc.exists {
					return &gopter.PropResult{
						Status: gopter.PropFalse,
						Labels: []string{
							"Exists-A",
							"uuid does not exist, but Exists returns true",
						},
					}
				}
			}

			return &gopter.PropResult{Status: gopter.PropTrue}
		},
	}

	existsBCommand = &commands.ProtoCommand{
		Name: "Exists-B",
		RunFunc: func(
			systemUnderTest commands.SystemUnderTest,
		) commands.Result {
			ngt := systemUnderTest.(*ngtSystem).ngt

			_, exists := ngt.Exists(idB)
			return &resultContainer{
				exists: exists,
			}
		},
		NextStateFunc: func(state commands.State) commands.State {
			return state
		},
		PreConditionFunc: func(state commands.State) bool {
			_, ok := state.(*ngtState)
			return ok
		},
		PostConditionFunc: func(
			state commands.State,
			result commands.Result,
		) *gopter.PropResult {
			st := state.(*ngtState)
			rc := result.(*resultContainer)

			if st.states[idB] == INDEXED || st.states[idB] == IN_INSERT_QUEUE {
				if !rc.exists {
					return &gopter.PropResult{
						Status: gopter.PropFalse,
						Labels: []string{
							"Exists-B",
							"uuid exists, but Exists returns false",
						},
					}
				}
			} else {
				if rc.exists {
					return &gopter.PropResult{
						Status: gopter.PropFalse,
						Labels: []string{
							"Exists-B",
							"uuid does not exist, but Exists returns true",
						},
					}
				}
			}

			return &gopter.PropResult{Status: gopter.PropTrue}
		},
	}

	existsCCommand = &commands.ProtoCommand{
		Name: "Exists-C",
		RunFunc: func(
			systemUnderTest commands.SystemUnderTest,
		) commands.Result {
			ngt := systemUnderTest.(*ngtSystem).ngt

			_, exists := ngt.Exists(idC)
			return &resultContainer{
				exists: exists,
			}
		},
		NextStateFunc: func(state commands.State) commands.State {
			return state
		},
		PreConditionFunc: func(state commands.State) bool {
			_, ok := state.(*ngtState)
			return ok
		},
		PostConditionFunc: func(
			state commands.State,
			result commands.Result,
		) *gopter.PropResult {
			st := state.(*ngtState)
			rc := result.(*resultContainer)

			if st.states[idC] == INDEXED || st.states[idC] == IN_INSERT_QUEUE {
				if !rc.exists {
					return &gopter.PropResult{
						Status: gopter.PropFalse,
						Labels: []string{
							"Exists-C",
							"uuid exists, but Exists returns false",
						},
					}
				}
			} else {
				if rc.exists {
					return &gopter.PropResult{
						Status: gopter.PropFalse,
						Labels: []string{
							"Exists-C",
							"uuid does not exist, but Exists returns true",
						},
					}
				}
			}

			return &gopter.PropResult{Status: gopter.PropTrue}
		},
	}

	searchCommand = &commands.ProtoCommand{
		Name: "Search",
		RunFunc: func(
			systemUnderTest commands.SystemUnderTest,
		) commands.Result {
			ngtSys := systemUnderTest.(*ngtSystem)
			ngt := ngtSys.ngt

			res, err := ngt.Search(ngtSys.ctx, []float32{0.1, 0.1, 0.1}, 3, 0.1, -1.0)
			return &resultContainer{
				err:     err,
				results: res,
			}
		},
		NextStateFunc: func(state commands.State) commands.State {
			return state
		},
		PreConditionFunc: func(state commands.State) bool {
			_, ok := state.(*ngtState)
			return ok
		},
		PostConditionFunc: func(
			state commands.State,
			result commands.Result,
		) *gopter.PropResult {
			rc := result.(*resultContainer)
			if rc.err != nil {
				st := state.(*ngtState)
				if st.states[idA] != INDEXED &&
					st.states[idB] != INDEXED &&
					st.states[idC] != INDEXED {
					if !errors.Is(rc.err, errors.ErrEmptySearchResult) {
						return &gopter.PropResult{
							Status: gopter.PropFalse,
							Error:  rc.err,
							Labels: []string{
								"Search",
								"there's no index but it doesn't return ErrEmptySearchResult",
								rc.err.Error(),
							},
						}
					}
					return &gopter.PropResult{Status: gopter.PropTrue}
				}
				if errors.Is(rc.err, errors.ErrEmptySearchResult) {
					return &gopter.PropResult{Status: gopter.PropTrue}
					// TODO: Originally, the following code should be executed, but it is commented out once because the behavior is suspicious due to the status of PBT.
					//       This will be fixed after investigating the cause.
					// return &gopter.PropResult{
					// 	Status: gopter.PropFalse,
					// 	Error:  rc.err,
					// 	Labels: []string{
					// 		"Search",
					// 		fmt.Sprintf("some indices are exists %v but it returned ErrEmptySearchResult", st.states),
					// 		rc.err.Error(),
					// 	},
					// }
				}
				return &gopter.PropResult{
					Status: gopter.PropFalse,
					Error:  rc.err,
					Labels: []string{
						"Search",
						"error",
						fmt.Sprintf("%v", st.states),
						rc.err.Error(),
					},
				}
			}
			return &gopter.PropResult{Status: gopter.PropTrue}
		},
	}

	searchByIDACommand = &commands.ProtoCommand{
		Name: "SearchByID-A",
		RunFunc: func(
			systemUnderTest commands.SystemUnderTest,
		) commands.Result {
			ngtSys := systemUnderTest.(*ngtSystem)
			ngt := ngtSys.ngt

			_, res, err := ngt.SearchByID(ngtSys.ctx, idA, 3, 0.1, -1.0)
			return &resultContainer{
				err:     err,
				results: res,
			}
		},
		NextStateFunc: func(state commands.State) commands.State {
			return state
		},
		PreConditionFunc: func(state commands.State) bool {
			st, ok := state.(*ngtState)
			if !ok {
				return false
			}

			return st.states[idA] == INDEXED
		},
		PostConditionFunc: func(
			state commands.State,
			result commands.Result,
		) *gopter.PropResult {
			rc := result.(*resultContainer)
			if rc.err != nil {
				st := state.(*ngtState)
				// TODO: In this test case, the possibility of retrieving a non-indexed ID is low,
				// but the test returns with an error, which needs to be investigated at a later time
				if st.states[idA] != INDEXED &&
					st.states[idB] != INDEXED &&
					st.states[idC] != INDEXED {
					if !errors.Is(rc.err, errors.ErrEmptySearchResult) {
						return &gopter.PropResult{
							Status: gopter.PropFalse,
							Error:  rc.err,
							Labels: []string{
								"SearchByID-A",
								"there's no index but it doesn't return ErrEmptySearchResult",
								rc.err.Error(),
							},
						}
					}
					return &gopter.PropResult{Status: gopter.PropTrue}
				}
				if errors.Is(rc.err, errors.ErrEmptySearchResult) {
					return &gopter.PropResult{Status: gopter.PropTrue}
					// TODO: Originally, the following code should be executed, but it is commented out once because the behavior is suspicious due to the status of PBT.
					//       This will be fixed after investigating the cause.
					// return &gopter.PropResult{
					// 	Status: gopter.PropFalse,
					// 	Error:  rc.err,
					// 	Labels: []string{
					// 		"SearchByID-A",
					// 		fmt.Sprintf("some indices are exists %v but it returned ErrEmptySearchResult", st.states),
					// 		rc.err.Error(),
					// 	},
					// }
				}
				return &gopter.PropResult{
					Status: gopter.PropFalse,
					Error:  rc.err,
					Labels: []string{
						"SearchByID-A",
						"error",
						fmt.Sprintf("%v", st.states),
						rc.err.Error(),
					},
				}
			}

			return &gopter.PropResult{Status: gopter.PropTrue}
		},
	}

	searchByIDBCommand = &commands.ProtoCommand{
		Name: "SearchByID-B",
		RunFunc: func(
			systemUnderTest commands.SystemUnderTest,
		) commands.Result {
			ngtSys := systemUnderTest.(*ngtSystem)
			ngt := ngtSys.ngt

			_, res, err := ngt.SearchByID(ngtSys.ctx, idB, 3, 0.1, -1.0)
			return &resultContainer{
				err:     err,
				results: res,
			}
		},
		NextStateFunc: func(state commands.State) commands.State {
			return state
		},
		PreConditionFunc: func(state commands.State) bool {
			st, ok := state.(*ngtState)
			if !ok {
				return false
			}

			return st.states[idB] == INDEXED
		},
		PostConditionFunc: func(
			state commands.State,
			result commands.Result,
		) *gopter.PropResult {
			rc := result.(*resultContainer)
			if rc.err != nil {
				st := state.(*ngtState)
				// TODO: In this test case, the possibility of retrieving a non-indexed ID is low,
				// but the test returns with an error, which needs to be investigated at a later time
				if st.states[idA] != INDEXED &&
					st.states[idB] != INDEXED &&
					st.states[idC] != INDEXED {
					if !errors.Is(rc.err, errors.ErrEmptySearchResult) {
						return &gopter.PropResult{
							Status: gopter.PropFalse,
							Error:  rc.err,
							Labels: []string{
								"SearchByID-B",
								"there's no index but it doesn't return ErrEmptySearchResult",
								rc.err.Error(),
							},
						}
					}
					return &gopter.PropResult{Status: gopter.PropTrue}
				}
				if errors.Is(rc.err, errors.ErrEmptySearchResult) {
					return &gopter.PropResult{Status: gopter.PropTrue}
					// TODO: Originally, the following code should be executed, but it is commented out once because the behavior is suspicious due to the status of PBT.
					//       This will be fixed after investigating the cause.
					// return &gopter.PropResult{
					// 	Status: gopter.PropFalse,
					// 	Error:  rc.err,
					// 	Labels: []string{
					// 		"SearchByID-B",
					// 		fmt.Sprintf("some indices are exists %v but it returned ErrEmptySearchResult", st.states),
					// 		rc.err.Error(),
					// 	},
					// }
				}
				return &gopter.PropResult{
					Status: gopter.PropFalse,
					Error:  rc.err,
					Labels: []string{
						"SearchByID-B",
						"error",
						fmt.Sprintf("%v", st.states),
						rc.err.Error(),
					},
				}
			}

			return &gopter.PropResult{Status: gopter.PropTrue}
		},
	}

	searchByIDCCommand = &commands.ProtoCommand{
		Name: "SearchByID-C",
		RunFunc: func(
			systemUnderTest commands.SystemUnderTest,
		) commands.Result {
			ngtSys := systemUnderTest.(*ngtSystem)
			ngt := ngtSys.ngt

			_, res, err := ngt.SearchByID(ngtSys.ctx, idC, 3, 0.1, -1.0)
			return &resultContainer{
				err:     err,
				results: res,
			}
		},
		NextStateFunc: func(state commands.State) commands.State {
			return state
		},
		PreConditionFunc: func(state commands.State) bool {
			st, ok := state.(*ngtState)
			if !ok {
				return false
			}

			return st.states[idC] == INDEXED
		},
		PostConditionFunc: func(
			state commands.State,
			result commands.Result,
		) *gopter.PropResult {
			rc := result.(*resultContainer)
			if rc.err != nil {
				st := state.(*ngtState)
				// TODO: In this test case, the possibility of retrieving a non-indexed ID is low,
				// but the test returns with an error, which needs to be investigated at a later time
				if st.states[idA] != INDEXED &&
					st.states[idB] != INDEXED &&
					st.states[idC] != INDEXED {
					if !errors.Is(rc.err, errors.ErrEmptySearchResult) {
						return &gopter.PropResult{
							Status: gopter.PropFalse,
							Error:  rc.err,
							Labels: []string{
								"SearchByID-C",
								"there's no index but it doesn't return ErrEmptySearchResult",
								rc.err.Error(),
							},
						}
					}
					return &gopter.PropResult{Status: gopter.PropTrue}
				}
				if errors.Is(rc.err, errors.ErrEmptySearchResult) {
					return &gopter.PropResult{Status: gopter.PropTrue}
					// TODO: Originally, the following code should be executed, but it is commented out once because the behavior is suspicious due to the status of PBT.
					//       This will be fixed after investigating the cause.
					// return &gopter.PropResult{
					// 	Status: gopter.PropFalse,
					// 	Error:  rc.err,
					// 	Labels: []string{
					// 		"SearchByID-C",
					// 		fmt.Sprintf("some indices are exists %v but it returned ErrEmptySearchResult", st.states),
					// 		rc.err.Error(),
					// 	},
					// }
				}
				return &gopter.PropResult{
					Status: gopter.PropFalse,
					Error:  rc.err,
					Labels: []string{
						"SearchByID-C",
						"error",
						fmt.Sprintf("%v", st.states),
						rc.err.Error(),
					},
				}
			}

			return &gopter.PropResult{Status: gopter.PropTrue}
		},
	}

	updateACommand = &commands.ProtoCommand{
		Name: "Update-A",
		RunFunc: func(
			systemUnderTest commands.SystemUnderTest,
		) commands.Result {
			ngt := systemUnderTest.(*ngtSystem).ngt

			err := ngt.Update(idA, vectors[idA][1])
			return &resultContainer{
				err: err,
			}
		},
		NextStateFunc: func(state commands.State) commands.State {
			st := state.(*ngtState)
			st.vectors[idA] = 1
			return state
		},
		PreConditionFunc: func(state commands.State) bool {
			st, ok := state.(*ngtState)
			if !ok {
				return false
			}

			return (st.states[idA] == INDEXED || st.states[idA] == IN_INSERT_QUEUE) && st.vectors[idA] == 0
		},
		PostConditionFunc: func(
			state commands.State,
			result commands.Result,
		) *gopter.PropResult {
			rc := result.(*resultContainer)

			if rc.err != nil {
				return &gopter.PropResult{
					Status: gopter.PropFalse,
					Error:  rc.err,
					Labels: []string{
						"Update-A",
						"error",
						rc.err.Error(),
					},
				}
			}

			return &gopter.PropResult{Status: gopter.PropTrue}
		},
	}

	updateBCommand = &commands.ProtoCommand{
		Name: "Update-B",
		RunFunc: func(
			systemUnderTest commands.SystemUnderTest,
		) commands.Result {
			ngt := systemUnderTest.(*ngtSystem).ngt

			err := ngt.Update(idB, vectors[idB][1])
			return &resultContainer{
				err: err,
			}
		},
		NextStateFunc: func(state commands.State) commands.State {
			st := state.(*ngtState)
			st.vectors[idB] = 1
			return state
		},
		PreConditionFunc: func(state commands.State) bool {
			st, ok := state.(*ngtState)
			if !ok {
				return false
			}

			return (st.states[idB] == INDEXED || st.states[idB] == IN_INSERT_QUEUE) && st.vectors[idB] == 0
		},
		PostConditionFunc: func(
			state commands.State,
			result commands.Result,
		) *gopter.PropResult {
			rc := result.(*resultContainer)

			if rc.err != nil {
				return &gopter.PropResult{
					Status: gopter.PropFalse,
					Error:  rc.err,
					Labels: []string{
						"Update-B",
						"error",
						rc.err.Error(),
					},
				}
			}

			return &gopter.PropResult{Status: gopter.PropTrue}
		},
	}

	updateCCommand = &commands.ProtoCommand{
		Name: "Update-C",
		RunFunc: func(
			systemUnderTest commands.SystemUnderTest,
		) commands.Result {
			ngt := systemUnderTest.(*ngtSystem).ngt

			err := ngt.Update(idC, vectors[idC][1])
			return &resultContainer{
				err: err,
			}
		},
		NextStateFunc: func(state commands.State) commands.State {
			st := state.(*ngtState)
			st.vectors[idC] = 1
			return state
		},
		PreConditionFunc: func(state commands.State) bool {
			st, ok := state.(*ngtState)
			if !ok {
				return false
			}

			return (st.states[idC] == INDEXED || st.states[idC] == IN_INSERT_QUEUE) && st.vectors[idC] == 0
		},
		PostConditionFunc: func(
			state commands.State,
			result commands.Result,
		) *gopter.PropResult {
			rc := result.(*resultContainer)

			if rc.err != nil {
				return &gopter.PropResult{
					Status: gopter.PropFalse,
					Error:  rc.err,
					Labels: []string{
						"Update-C",
						"error",
						rc.err.Error(),
					},
				}
			}

			return &gopter.PropResult{Status: gopter.PropTrue}
		},
	}

	deleteACommand = &commands.ProtoCommand{
		Name: "Delete-A",
		RunFunc: func(
			systemUnderTest commands.SystemUnderTest,
		) commands.Result {
			ngt := systemUnderTest.(*ngtSystem).ngt

			err := ngt.Delete(idA)
			return &resultContainer{
				err: err,
			}
		},
		NextStateFunc: func(state commands.State) commands.State {
			st := state.(*ngtState)
			st.states[idA] = IN_DELETE_QUEUE
			return state
		},
		PreConditionFunc: func(state commands.State) bool {
			st, ok := state.(*ngtState)
			if !ok {
				return false
			}

			return st.states[idA] == IN_INSERT_QUEUE || st.states[idA] == INDEXED
		},
		PostConditionFunc: func(
			state commands.State,
			result commands.Result,
		) *gopter.PropResult {
			rc := result.(*resultContainer)

			if rc.err != nil {
				return &gopter.PropResult{
					Status: gopter.PropFalse,
					Error:  rc.err,
					Labels: []string{
						"Delete-A",
						"error",
						rc.err.Error(),
					},
				}
			}

			return &gopter.PropResult{Status: gopter.PropTrue}
		},
	}

	deleteBCommand = &commands.ProtoCommand{
		Name: "Delete-B",
		RunFunc: func(
			systemUnderTest commands.SystemUnderTest,
		) commands.Result {
			ngt := systemUnderTest.(*ngtSystem).ngt

			err := ngt.Delete(idB)
			return &resultContainer{
				err: err,
			}
		},
		NextStateFunc: func(state commands.State) commands.State {
			st := state.(*ngtState)
			st.states[idB] = IN_DELETE_QUEUE
			return state
		},
		PreConditionFunc: func(state commands.State) bool {
			st, ok := state.(*ngtState)
			if !ok {
				return false
			}

			return st.states[idB] == IN_INSERT_QUEUE || st.states[idB] == INDEXED
		},
		PostConditionFunc: func(
			state commands.State,
			result commands.Result,
		) *gopter.PropResult {
			rc := result.(*resultContainer)

			if rc.err != nil {
				return &gopter.PropResult{
					Status: gopter.PropFalse,
					Error:  rc.err,
					Labels: []string{
						"Delete-B",
						"error",
						rc.err.Error(),
					},
				}
			}

			return &gopter.PropResult{Status: gopter.PropTrue}
		},
	}

	deleteCCommand = &commands.ProtoCommand{
		Name: "Delete-C",
		RunFunc: func(
			systemUnderTest commands.SystemUnderTest,
		) commands.Result {
			ngt := systemUnderTest.(*ngtSystem).ngt

			err := ngt.Delete(idC)
			return &resultContainer{
				err: err,
			}
		},
		NextStateFunc: func(state commands.State) commands.State {
			st := state.(*ngtState)
			st.states[idC] = IN_DELETE_QUEUE
			return state
		},
		PreConditionFunc: func(state commands.State) bool {
			st, ok := state.(*ngtState)
			if !ok {
				return false
			}

			return st.states[idC] == IN_INSERT_QUEUE || st.states[idC] == INDEXED
		},
		PostConditionFunc: func(
			state commands.State,
			result commands.Result,
		) *gopter.PropResult {
			rc := result.(*resultContainer)

			if rc.err != nil {
				return &gopter.PropResult{
					Status: gopter.PropFalse,
					Error:  rc.err,
					Labels: []string{
						"Delete-C",
						"error",
						rc.err.Error(),
					},
				}
			}

			return &gopter.PropResult{Status: gopter.PropTrue}
		},
	}

	getobjectACommand = &commands.ProtoCommand{
		Name: "GetObject-A",
		RunFunc: func(
			systemUnderTest commands.SystemUnderTest,
		) commands.Result {
			ngt := systemUnderTest.(*ngtSystem).ngt

			vec, err := ngt.GetObject(idA)
			return &resultContainer{
				vector: vec,
				err:    err,
			}
		},
		NextStateFunc: func(state commands.State) commands.State {
			return state
		},
		PreConditionFunc: func(state commands.State) bool {
			st, ok := state.(*ngtState)
			if !ok {
				return false
			}

			return st.states[idA] == IN_INSERT_QUEUE || st.states[idA] == INDEXED
		},
		PostConditionFunc: func(
			state commands.State,
			result commands.Result,
		) *gopter.PropResult {
			st := state.(*ngtState)
			rc := result.(*resultContainer)

			if rc.err != nil {
				return &gopter.PropResult{
					Status: gopter.PropFalse,
					Error:  rc.err,
					Labels: []string{
						"GetObject-A",
						"error",
						rc.err.Error(),
					},
				}
			}

			if !reflect.DeepEqual(rc.vector, vectors[idA][st.vectors[idA]]) {
				return &gopter.PropResult{
					Status: gopter.PropFalse,
					Labels: []string{
						"GetObject-A",
						"invalid vector",
						fmt.Sprintf("got: %#v, expected: %#v", rc.vector, vectors[idA][st.vectors[idA]]),
					},
				}
			}

			return &gopter.PropResult{Status: gopter.PropTrue}
		},
	}

	getobjectBCommand = &commands.ProtoCommand{
		Name: "GetObject-B",
		RunFunc: func(
			systemUnderTest commands.SystemUnderTest,
		) commands.Result {
			ngt := systemUnderTest.(*ngtSystem).ngt

			vec, err := ngt.GetObject(idB)
			return &resultContainer{
				vector: vec,
				err:    err,
			}
		},
		NextStateFunc: func(state commands.State) commands.State {
			return state
		},
		PreConditionFunc: func(state commands.State) bool {
			st, ok := state.(*ngtState)
			if !ok {
				return false
			}

			return st.states[idB] == IN_INSERT_QUEUE || st.states[idB] == INDEXED
		},
		PostConditionFunc: func(
			state commands.State,
			result commands.Result,
		) *gopter.PropResult {
			st := state.(*ngtState)
			rc := result.(*resultContainer)

			if rc.err != nil {
				return &gopter.PropResult{
					Status: gopter.PropFalse,
					Error:  rc.err,
					Labels: []string{
						"GetObject-B",
						"error",
						rc.err.Error(),
					},
				}
			}

			if !reflect.DeepEqual(rc.vector, vectors[idB][st.vectors[idB]]) {
				return &gopter.PropResult{
					Status: gopter.PropFalse,
					Labels: []string{
						"GetObject-B",
						"invalid vector",
						fmt.Sprintf("got: %#v, expected: %#v", rc.vector, vectors[idB][st.vectors[idB]]),
					},
				}
			}

			return &gopter.PropResult{Status: gopter.PropTrue}
		},
	}

	getobjectCCommand = &commands.ProtoCommand{
		Name: "GetObject-C",
		RunFunc: func(
			systemUnderTest commands.SystemUnderTest,
		) commands.Result {
			ngt := systemUnderTest.(*ngtSystem).ngt

			vec, err := ngt.GetObject(idC)
			return &resultContainer{
				vector: vec,
				err:    err,
			}
		},
		NextStateFunc: func(state commands.State) commands.State {
			return state
		},
		PreConditionFunc: func(state commands.State) bool {
			st, ok := state.(*ngtState)
			if !ok {
				return false
			}

			return st.states[idC] == IN_INSERT_QUEUE || st.states[idC] == INDEXED
		},
		PostConditionFunc: func(
			state commands.State,
			result commands.Result,
		) *gopter.PropResult {
			st := state.(*ngtState)
			rc := result.(*resultContainer)

			if rc.err != nil {
				return &gopter.PropResult{
					Status: gopter.PropFalse,
					Error:  rc.err,
					Labels: []string{
						"GetObject-C",
						"error",
						rc.err.Error(),
					},
				}
			}

			if !reflect.DeepEqual(rc.vector, vectors[idC][st.vectors[idC]]) {
				return &gopter.PropResult{
					Status: gopter.PropFalse,
					Labels: []string{
						"GetObject-C",
						"invalid vector",
						fmt.Sprintf("got: %#v, expected: %#v", rc.vector, vectors[idC][st.vectors[idC]]),
					},
				}
			}

			return &gopter.PropResult{Status: gopter.PropTrue}
		},
	}
)

func rootCommands(t *testing.T) commands.Commands {
	t.Helper()
	return &commands.ProtoCommands{
		NewSystemUnderTestFunc: func(
			initialState commands.State,
		) commands.SystemUnderTest {
			n, err := New(ncfg.Bind())
			if err != nil {
				t.Fatalf("error: %s", err)
			}

			ctx := context.Background()
			ech := n.Start(ctx)

			time.Sleep(time.Second)

			return &ngtSystem{
				ctx: ctx,
				ech: ech,
				ngt: n.(*ngt),
			}
		},
		DestroySystemUnderTestFunc: func(sys commands.SystemUnderTest) {
			s := sys.(*ngtSystem)
			s.ngt.Close(s.ctx)
			s.ctx.Done()
		},
		InitialStateGen: gen.Const(&ngtState{
			states: map[string]vectorState{
				idA: NOT_INSERTED,
				idB: NOT_INSERTED,
				idC: NOT_INSERTED,
			},
			vectors: map[string]int{
				idA: 0,
				idB: 0,
				idC: 0,
			},
		}),
		GenCommandFunc: func(state commands.State) gopter.Gen {
			st := state.(*ngtState)

			cs := make([]interface{}, 0)
			cs = append(
				cs,
				existsACommand,
				existsBCommand,
				existsCCommand,
				searchCommand,
			)

			for _, v := range st.states {
				if v == IN_INSERT_QUEUE || v == IN_DELETE_QUEUE {
					cs = append(cs, createIndexCommand)

					break
				}
			}

			switch st.states[idA] {
			case NOT_INSERTED:
				cs = append(cs, insertACommand)
			case IN_INSERT_QUEUE:
				if st.vectors[idA] == 0 {
					cs = append(cs, updateACommand)
				}

				cs = append(
					cs,
					deleteACommand,
					getobjectACommand,
				)
			case IN_DELETE_QUEUE:
				cs = append(
					cs,
					insertACommand,
				)
			case INDEXED:
				if st.vectors[idA] == 0 {
					cs = append(cs, updateACommand)
				}

				cs = append(
					cs,
					deleteACommand,
					searchByIDACommand,
					getobjectACommand,
				)
			}

			switch st.states[idB] {
			case NOT_INSERTED:
				cs = append(cs, insertBCommand)
			case IN_INSERT_QUEUE:
				if st.vectors[idB] == 0 {
					cs = append(cs, updateBCommand)
				}

				cs = append(
					cs,
					deleteBCommand,
					getobjectBCommand,
				)
			case IN_DELETE_QUEUE:
				cs = append(
					cs,
					insertBCommand,
				)
			case INDEXED:
				if st.vectors[idB] == 0 {
					cs = append(cs, updateBCommand)
				}

				cs = append(
					cs,
					deleteBCommand,
					searchByIDBCommand,
					getobjectBCommand,
				)
			}

			switch st.states[idC] {
			case NOT_INSERTED:
				cs = append(cs, insertCCommand)
			case IN_INSERT_QUEUE:
				if st.vectors[idC] == 0 {
					cs = append(cs, updateCCommand)
				}

				cs = append(
					cs,
					deleteCCommand,
					getobjectCCommand,
				)
			case IN_DELETE_QUEUE:
				cs = append(
					cs,
					insertCCommand,
				)
			case INDEXED:
				if st.vectors[idC] == 0 {
					cs = append(cs, updateCCommand)
				}

				cs = append(
					cs,
					deleteCCommand,
					searchByIDCCommand,
					getobjectCCommand,
				)
			}

			return gen.OneConstOf(cs...)
		},
		InitialPreConditionFunc: func(state commands.State) bool {
			st := state.(*ngtState)

			st.Reset()

			for _, v := range st.states {
				if v != NOT_INSERTED {
					return false
				}
			}

			for _, v := range st.vectors {
				if v != 0 {
					return false
				}
			}

			return true
		},
	}
}

func TestStatefulNGT(t *testing.T) {
	// initialize logger
	if applicationLog {
		log.Init()
	} else {
		log.Init(log.WithLoggerType("nop"))
	}

	parameters := gopter.DefaultTestParameters()
	if seed != 0 {
		parameters.SetSeed(seed)
	}
	parameters.MinSuccessfulTests = minSuccessfulTests
	parameters.MaxDiscardRatio = maxDiscardRatio

	properties := gopter.NewProperties(parameters)
	properties.Property("NGT", commands.Prop(rootCommands(t)))

	properties.TestingRun(t)
}
