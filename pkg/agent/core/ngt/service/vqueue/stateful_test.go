//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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

package vqueue

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
	"github.com/vdaas/vald/internal/log"
)

type qSystem struct {
	ctx context.Context
	q   *vqueue
}

type resultContainer struct {
	err    error
	vector []float32
	exists bool
}

type qState struct {
	idA idState
	idB idState
	idC idState
}

func (st *qState) Reset() {
	st.idA = NOT_EXIST
	st.idB = NOT_EXIST
	st.idC = NOT_EXIST
}

type idState uint8

const (
	idA = "id-a"
	idB = "id-b"
	idC = "id-c"

	NOT_EXIST idState = iota
	IN_IV_ONLY
	IN_DV_ONLY
	IN_BOTH_IV_LATER
)

var (
	seed               int64
	minSuccessfulTests int
	maxDiscardRatio    float64
	applicationLog     bool

	vectors = map[string][]float32{
		idA: {0.1, 0.2, 0.3},
		idB: {0.2, 0.3, 0.1},
		idC: {0.3, 0.1, 0.2},
	}
)

func init() {
	testing.Init()

	flag.Int64Var(&seed, "pbt-seed", 0, "seed number used for PBT")
	flag.IntVar(&minSuccessfulTests, "pbt-min-successful-tests", 10, "minimum number of successful tests in PBT")
	flag.Float64Var(&maxDiscardRatio, "pbt-max-discard-ratio", 5.0, "maximum discard ratio of PBT")
	flag.BoolVar(&applicationLog, "pbt-enable-application-log", false, "enable application log on PBT")
}

var (
	pushIACommand = &commands.ProtoCommand{
		Name: "PushInsert-A",
		RunFunc: func(
			systemUnderTest commands.SystemUnderTest,
		) commands.Result {
			sy := systemUnderTest.(*qSystem)
			q := sy.q

			err := q.PushInsert(idA, vectors[idA], 0)
			return &resultContainer{
				err: err,
			}
		},
		NextStateFunc: func(state commands.State) commands.State {
			st := state.(*qState)

			switch st.idA {
			case NOT_EXIST, IN_IV_ONLY:
				st.idA = IN_IV_ONLY
			case IN_DV_ONLY, IN_BOTH_IV_LATER:
				st.idA = IN_BOTH_IV_LATER
			}

			return st
		},
		PreConditionFunc: func(state commands.State) bool {
			_, ok := state.(*qState)
			return ok
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
						"PushInsert-A",
						"error",
						rc.err.Error(),
					},
				}
			}

			return &gopter.PropResult{Status: gopter.PropTrue}
		},
	}

	pushIBCommand = &commands.ProtoCommand{
		Name: "PushInsert-B",
		RunFunc: func(
			systemUnderTest commands.SystemUnderTest,
		) commands.Result {
			sy := systemUnderTest.(*qSystem)
			q := sy.q

			err := q.PushInsert(idB, vectors[idB], 0)
			return &resultContainer{
				err: err,
			}
		},
		NextStateFunc: func(state commands.State) commands.State {
			st := state.(*qState)

			switch st.idB {
			case NOT_EXIST, IN_IV_ONLY:
				st.idB = IN_IV_ONLY
			case IN_DV_ONLY, IN_BOTH_IV_LATER:
				st.idB = IN_BOTH_IV_LATER
			}

			return st
		},
		PreConditionFunc: func(state commands.State) bool {
			_, ok := state.(*qState)
			return ok
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
						"PushInsert-B",
						"error",
						rc.err.Error(),
					},
				}
			}

			return &gopter.PropResult{Status: gopter.PropTrue}
		},
	}

	pushICCommand = &commands.ProtoCommand{
		Name: "PushInsert-C",
		RunFunc: func(
			systemUnderTest commands.SystemUnderTest,
		) commands.Result {
			sy := systemUnderTest.(*qSystem)
			q := sy.q

			err := q.PushInsert(idC, vectors[idC], 0)
			return &resultContainer{
				err: err,
			}
		},
		NextStateFunc: func(state commands.State) commands.State {
			st := state.(*qState)

			switch st.idC {
			case NOT_EXIST, IN_IV_ONLY:
				st.idC = IN_IV_ONLY
			case IN_DV_ONLY, IN_BOTH_IV_LATER:
				st.idC = IN_BOTH_IV_LATER
			}

			return st
		},
		PreConditionFunc: func(state commands.State) bool {
			_, ok := state.(*qState)

			return ok
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
						"PushInsert-C",
						"error",
						rc.err.Error(),
					},
				}
			}

			return &gopter.PropResult{Status: gopter.PropTrue}
		},
	}

	pushDACommand = &commands.ProtoCommand{
		Name: "PushDelete-A",
		RunFunc: func(
			systemUnderTest commands.SystemUnderTest,
		) commands.Result {
			sy := systemUnderTest.(*qSystem)
			q := sy.q

			err := q.PushDelete(idA, 0)
			return &resultContainer{
				err: err,
			}
		},
		NextStateFunc: func(state commands.State) commands.State {
			st := state.(*qState)
			st.idA = IN_DV_ONLY

			return st
		},
		PreConditionFunc: func(state commands.State) bool {
			_, ok := state.(*qState)
			return ok
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
						"PushDelete-A",
						"error",
						rc.err.Error(),
					},
				}
			}

			return &gopter.PropResult{Status: gopter.PropTrue}
		},
	}

	pushDBCommand = &commands.ProtoCommand{
		Name: "PushDelete-B",
		RunFunc: func(
			systemUnderTest commands.SystemUnderTest,
		) commands.Result {
			sy := systemUnderTest.(*qSystem)
			q := sy.q

			err := q.PushDelete(idB, 0)
			return &resultContainer{
				err: err,
			}
		},
		NextStateFunc: func(state commands.State) commands.State {
			st := state.(*qState)
			st.idB = IN_DV_ONLY

			return st
		},
		PreConditionFunc: func(state commands.State) bool {
			_, ok := state.(*qState)
			return ok
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
						"PushDelete-B",
						"error",
						rc.err.Error(),
					},
				}
			}

			return &gopter.PropResult{Status: gopter.PropTrue}
		},
	}

	pushDCCommand = &commands.ProtoCommand{
		Name: "PushDelete-C",
		RunFunc: func(
			systemUnderTest commands.SystemUnderTest,
		) commands.Result {
			sy := systemUnderTest.(*qSystem)
			q := sy.q

			err := q.PushDelete(idC, 0)
			return &resultContainer{
				err: err,
			}
		},
		NextStateFunc: func(state commands.State) commands.State {
			st := state.(*qState)
			st.idC = IN_DV_ONLY

			return st
		},
		PreConditionFunc: func(state commands.State) bool {
			_, ok := state.(*qState)
			return ok
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
						"PushDelete-C",
						"error",
						rc.err.Error(),
					},
				}
			}

			return &gopter.PropResult{Status: gopter.PropTrue}
		},
	}

	existsIACommand = &commands.ProtoCommand{
		Name: "IVExists-A",
		RunFunc: func(
			systemUnderTest commands.SystemUnderTest,
		) commands.Result {
			sy := systemUnderTest.(*qSystem)
			q := sy.q

			exists := q.IVExists(idA)
			return &resultContainer{
				exists: exists,
			}
		},
		NextStateFunc: func(state commands.State) commands.State {
			return state
		},
		PreConditionFunc: func(state commands.State) bool {
			_, ok := state.(*qState)
			return ok
		},
		PostConditionFunc: func(
			state commands.State,
			result commands.Result,
		) *gopter.PropResult {
			st := state.(*qState)
			rc := result.(*resultContainer)

			exists := func() bool {
				switch st.idA {
				case NOT_EXIST, IN_DV_ONLY:
					return false
				case IN_IV_ONLY, IN_BOTH_IV_LATER:
					return true
				}
				return false
			}()
			if rc.exists != exists {
				return &gopter.PropResult{
					Status: gopter.PropFalse,
					Labels: []string{
						"IVExists-A",
						"IVExists returns invalid result",
						fmt.Sprintf("IVExists: %t, Actually: %t", rc.exists, exists),
					},
				}
			}

			return &gopter.PropResult{Status: gopter.PropTrue}
		},
	}

	existsIBCommand = &commands.ProtoCommand{
		Name: "IVExists-B",
		RunFunc: func(
			systemUnderTest commands.SystemUnderTest,
		) commands.Result {
			sy := systemUnderTest.(*qSystem)
			q := sy.q

			exists := q.IVExists(idB)
			return &resultContainer{
				exists: exists,
			}
		},
		NextStateFunc: func(state commands.State) commands.State {
			return state
		},
		PreConditionFunc: func(state commands.State) bool {
			_, ok := state.(*qState)
			return ok
		},
		PostConditionFunc: func(
			state commands.State,
			result commands.Result,
		) *gopter.PropResult {
			st := state.(*qState)
			rc := result.(*resultContainer)

			exists := func() bool {
				switch st.idB {
				case NOT_EXIST, IN_DV_ONLY:
					return false
				case IN_IV_ONLY, IN_BOTH_IV_LATER:
					return true
				}
				return false
			}()
			if rc.exists != exists {
				return &gopter.PropResult{
					Status: gopter.PropFalse,
					Labels: []string{
						"IVExists-B",
						"IVExists returns invalid result",
						fmt.Sprintf("IVExists: %t, Actually: %t", rc.exists, exists),
					},
				}
			}

			return &gopter.PropResult{Status: gopter.PropTrue}
		},
	}

	existsICCommand = &commands.ProtoCommand{
		Name: "IVExists-C",
		RunFunc: func(
			systemUnderTest commands.SystemUnderTest,
		) commands.Result {
			sy := systemUnderTest.(*qSystem)
			q := sy.q

			exists := q.IVExists(idC)
			return &resultContainer{
				exists: exists,
			}
		},
		NextStateFunc: func(state commands.State) commands.State {
			return state
		},
		PreConditionFunc: func(state commands.State) bool {
			_, ok := state.(*qState)
			return ok
		},
		PostConditionFunc: func(
			state commands.State,
			result commands.Result,
		) *gopter.PropResult {
			st := state.(*qState)
			rc := result.(*resultContainer)

			exists := func() bool {
				switch st.idC {
				case NOT_EXIST, IN_DV_ONLY:
					return false
				case IN_IV_ONLY, IN_BOTH_IV_LATER:
					return true
				}
				return false
			}()
			if rc.exists != exists {
				return &gopter.PropResult{
					Status: gopter.PropFalse,
					Labels: []string{
						"IVExists-C",
						"IVExists returns invalid result",
						fmt.Sprintf("IVExists: %t, Actually: %t", rc.exists, exists),
					},
				}
			}

			return &gopter.PropResult{Status: gopter.PropTrue}
		},
	}

	existsDACommand = &commands.ProtoCommand{
		Name: "DVExists-A",
		RunFunc: func(
			systemUnderTest commands.SystemUnderTest,
		) commands.Result {
			sy := systemUnderTest.(*qSystem)
			q := sy.q

			exists := q.DVExists(idA)
			return &resultContainer{
				exists: exists,
			}
		},
		NextStateFunc: func(state commands.State) commands.State {
			return state
		},
		PreConditionFunc: func(state commands.State) bool {
			_, ok := state.(*qState)
			return ok
		},
		PostConditionFunc: func(
			state commands.State,
			result commands.Result,
		) *gopter.PropResult {
			st := state.(*qState)
			rc := result.(*resultContainer)

			exists := func() bool {
				switch st.idA {
				case NOT_EXIST, IN_IV_ONLY, IN_BOTH_IV_LATER:
					return false
				case IN_DV_ONLY:
					return true
				}
				return false
			}()
			if rc.exists != exists {
				return &gopter.PropResult{
					Status: gopter.PropFalse,
					Labels: []string{
						"DVExists-A",
						"DVExists returns invalid result",
						fmt.Sprintf("DVExists: %t, Actually: %t", rc.exists, exists),
					},
				}
			}

			return &gopter.PropResult{Status: gopter.PropTrue}
		},
	}

	existsDBCommand = &commands.ProtoCommand{
		Name: "DVExists-B",
		RunFunc: func(
			systemUnderTest commands.SystemUnderTest,
		) commands.Result {
			sy := systemUnderTest.(*qSystem)
			q := sy.q

			exists := q.DVExists(idB)
			return &resultContainer{
				exists: exists,
			}
		},
		NextStateFunc: func(state commands.State) commands.State {
			return state
		},
		PreConditionFunc: func(state commands.State) bool {
			_, ok := state.(*qState)
			return ok
		},
		PostConditionFunc: func(
			state commands.State,
			result commands.Result,
		) *gopter.PropResult {
			st := state.(*qState)
			rc := result.(*resultContainer)

			exists := func() bool {
				switch st.idB {
				case NOT_EXIST, IN_IV_ONLY, IN_BOTH_IV_LATER:
					return false
				case IN_DV_ONLY:
					return true
				}
				return false
			}()
			if rc.exists != exists {
				return &gopter.PropResult{
					Status: gopter.PropFalse,
					Labels: []string{
						"DVExists-B",
						"DVExists returns invalid result",
						fmt.Sprintf("DVExists: %t, Actually: %t", rc.exists, exists),
					},
				}
			}

			return &gopter.PropResult{Status: gopter.PropTrue}
		},
	}

	existsDCCommand = &commands.ProtoCommand{
		Name: "DVExists-C",
		RunFunc: func(
			systemUnderTest commands.SystemUnderTest,
		) commands.Result {
			sy := systemUnderTest.(*qSystem)
			q := sy.q

			exists := q.DVExists(idC)
			return &resultContainer{
				exists: exists,
			}
		},
		NextStateFunc: func(state commands.State) commands.State {
			return state
		},
		PreConditionFunc: func(state commands.State) bool {
			_, ok := state.(*qState)
			return ok
		},
		PostConditionFunc: func(
			state commands.State,
			result commands.Result,
		) *gopter.PropResult {
			st := state.(*qState)
			rc := result.(*resultContainer)

			exists := func() bool {
				switch st.idC {
				case NOT_EXIST, IN_IV_ONLY, IN_BOTH_IV_LATER:
					return false
				case IN_DV_ONLY:
					return true
				}
				return false
			}()
			if rc.exists != exists {
				return &gopter.PropResult{
					Status: gopter.PropFalse,
					Labels: []string{
						"DVExists-C",
						"DVExists returns invalid result",
						fmt.Sprintf("DVExists: %t, Actually: %t", rc.exists, exists),
					},
				}
			}

			return &gopter.PropResult{Status: gopter.PropTrue}
		},
	}

	getvectorACommand = &commands.ProtoCommand{
		Name: "GetVector-A",
		RunFunc: func(
			systemUnderTest commands.SystemUnderTest,
		) commands.Result {
			sy := systemUnderTest.(*qSystem)
			q := sy.q

			vec, exists := q.GetVector(idA)
			return &resultContainer{
				vector: vec,
				exists: exists,
			}
		},
		NextStateFunc: func(state commands.State) commands.State {
			return state
		},
		PreConditionFunc: func(state commands.State) bool {
			_, ok := state.(*qState)
			return ok
		},
		PostConditionFunc: func(
			state commands.State,
			result commands.Result,
		) *gopter.PropResult {
			st := state.(*qState)
			rc := result.(*resultContainer)

			exists := func() bool {
				switch st.idA {
				case NOT_EXIST, IN_DV_ONLY:
					return false
				case IN_IV_ONLY, IN_BOTH_IV_LATER:
					return true
				}
				return false
			}()
			if rc.exists != exists {
				return &gopter.PropResult{
					Status: gopter.PropFalse,
					Labels: []string{
						"GetObject-A",
						fmt.Sprintf("GetObject tells idA exists(%t), but actually (%t)", rc.exists, exists),
					},
				}
			}

			if exists && !reflect.DeepEqual(rc.vector, vectors[idA]) {
				return &gopter.PropResult{
					Status: gopter.PropFalse,
					Labels: []string{
						"GetVector-A",
						"GetVector returns invalid result",
						fmt.Sprintf("GetVector: %#v, Actually: %#v", rc.vector, vectors[idA]),
					},
				}
			}

			return &gopter.PropResult{Status: gopter.PropTrue}
		},
	}

	getvectorBCommand = &commands.ProtoCommand{
		Name: "GetVector-B",
		RunFunc: func(
			systemUnderTest commands.SystemUnderTest,
		) commands.Result {
			sy := systemUnderTest.(*qSystem)
			q := sy.q

			vec, exists := q.GetVector(idB)
			return &resultContainer{
				vector: vec,
				exists: exists,
			}
		},
		NextStateFunc: func(state commands.State) commands.State {
			return state
		},
		PreConditionFunc: func(state commands.State) bool {
			_, ok := state.(*qState)
			return ok
		},
		PostConditionFunc: func(
			state commands.State,
			result commands.Result,
		) *gopter.PropResult {
			st := state.(*qState)
			rc := result.(*resultContainer)

			exists := func() bool {
				switch st.idB {
				case NOT_EXIST, IN_DV_ONLY:
					return false
				case IN_IV_ONLY, IN_BOTH_IV_LATER:
					return true
				}
				return false
			}()
			if rc.exists != exists {
				return &gopter.PropResult{
					Status: gopter.PropFalse,
					Labels: []string{
						"GetObject-B",
						fmt.Sprintf("GetObject tells idB exists(%t), but actually (%t)", rc.exists, exists),
					},
				}
			}

			if exists && !reflect.DeepEqual(rc.vector, vectors[idB]) {
				return &gopter.PropResult{
					Status: gopter.PropFalse,
					Labels: []string{
						"GetVector-B",
						"GetVector returns invalid result",
						fmt.Sprintf("GetVector: %#v, Actually: %#v", rc.vector, vectors[idB]),
					},
				}
			}

			return &gopter.PropResult{Status: gopter.PropTrue}
		},
	}

	getvectorCCommand = &commands.ProtoCommand{
		Name: "GetVector-C",
		RunFunc: func(
			systemUnderTest commands.SystemUnderTest,
		) commands.Result {
			sy := systemUnderTest.(*qSystem)
			q := sy.q

			vec, exists := q.GetVector(idC)
			return &resultContainer{
				vector: vec,
				exists: exists,
			}
		},
		NextStateFunc: func(state commands.State) commands.State {
			return state
		},
		PreConditionFunc: func(state commands.State) bool {
			_, ok := state.(*qState)
			return ok
		},
		PostConditionFunc: func(
			state commands.State,
			result commands.Result,
		) *gopter.PropResult {
			st := state.(*qState)
			rc := result.(*resultContainer)

			exists := func() bool {
				switch st.idC {
				case NOT_EXIST, IN_DV_ONLY:
					return false
				case IN_IV_ONLY, IN_BOTH_IV_LATER:
					return true
				}
				return false
			}()
			if rc.exists != exists {
				return &gopter.PropResult{
					Status: gopter.PropFalse,
					Labels: []string{
						"GetObject-C",
						fmt.Sprintf("GetObject tells idC exists(%t), but actually (%t)", rc.exists, exists),
					},
				}
			}

			if exists && !reflect.DeepEqual(rc.vector, vectors[idC]) {
				return &gopter.PropResult{
					Status: gopter.PropFalse,
					Labels: []string{
						"GetVector-C",
						"GetVector returns invalid result",
						fmt.Sprintf("GetVector: %#v, Actually: %#v", rc.vector, vectors[idC]),
					},
				}
			}

			return &gopter.PropResult{Status: gopter.PropTrue}
		},
	}

	popCommand = &commands.ProtoCommand{
		Name: "RangePopDelete+RangePopInsert",
		RunFunc: func(
			systemUnderTest commands.SystemUnderTest,
		) commands.Result {
			sy := systemUnderTest.(*qSystem)
			q := sy.q

			now := time.Now().UnixNano()

			q.RangePopDelete(
				sy.ctx,
				now,
				func(uuid string) bool { return true },
			)

			q.RangePopInsert(
				sy.ctx,
				now,
				func(uuid string, vector []float32) bool { return true },
			)

			return &resultContainer{}
		},
		NextStateFunc: func(state commands.State) commands.State {
			st := state.(*qState)

			st.idA = NOT_EXIST
			st.idB = NOT_EXIST
			st.idC = NOT_EXIST

			return st
		},
		PreConditionFunc: func(state commands.State) bool {
			_, ok := state.(*qState)
			return ok
		},
		PostConditionFunc: func(
			state commands.State,
			result commands.Result,
		) *gopter.PropResult {
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
			q, err := New()
			if err != nil {
				t.Fatalf("error: %s", err)
			}

			ctx := context.Background()

			return &qSystem{
				ctx: ctx,
				q:   q.(*vqueue),
			}
		},
		DestroySystemUnderTestFunc: func(sys commands.SystemUnderTest) {
			sy := sys.(*qSystem)
			sy.ctx.Done()
		},
		InitialStateGen: gen.Const(&qState{
			idA: NOT_EXIST,
			idB: NOT_EXIST,
			idC: NOT_EXIST,
		}),
		GenCommandFunc: func(state commands.State) gopter.Gen {
			return gen.OneConstOf(
				pushIACommand,
				pushIBCommand,
				pushICCommand,
				pushDACommand,
				pushDBCommand,
				pushDCCommand,
				existsIACommand,
				existsIBCommand,
				existsICCommand,
				existsDACommand,
				existsDBCommand,
				existsDCCommand,
				getvectorACommand,
				getvectorBCommand,
				getvectorCCommand,
				popCommand,
			)
		},
		InitialPreConditionFunc: func(state commands.State) bool {
			st := state.(*qState)

			st.Reset()

			return st.idA == NOT_EXIST && st.idB == NOT_EXIST && st.idC == NOT_EXIST
		},
	}
}

func TestStatefulQueue(t *testing.T) {
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
	properties.Property("vqueue", commands.Prop(rootCommands(t)))

	properties.TestingRun(t)
}
