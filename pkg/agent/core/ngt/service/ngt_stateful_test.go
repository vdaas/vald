//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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
	"testing"
	"time"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/commands"
	"github.com/leanovate/gopter/gen"
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/pkg/agent/core/ngt/model"
)

const (
	dimension = 3
)

var (
	seed               int64
	minSuccessfulTests int

	vecCh = make(chan *vector, 10)

	cfg = &config.NGT{
		Dimension:              dimension,
		DistanceType:           "l2",
		ObjectType:             "float",
		EnableInMemoryMode:     true,
		AutoIndexDurationLimit: "96h",
		AutoIndexCheckDuration: "96h",
		AutoSaveIndexDuration:  "96h",
		AutoIndexLength:        10000000000,
	}

	uuidGen = gen.AlphaString().SuchThat(func(s string) bool {
		return len(s) > 0
	})

	f32sliceGen = gen.SliceOfN(dimension, gen.Float32())
)

func init() {
	testing.Init()

	flag.Int64Var(&seed, "pbt-seed", 0, "seed number used for PBT")
	flag.IntVar(&minSuccessfulTests, "pbt-min-successful-tests", 10, "minimum number of successful tests in PBT")
}

type vector struct {
	uuid   string
	vector []float32
}

type ngtSystem struct {
	ctx           context.Context
	ech           <-chan error
	ngt           *ngt
	insertedUUIDs map[string][]float32
}

type resultContainer struct {
	err     error
	results []model.Distance
	uuid    string
	vector  []float32
	exists  bool
}

type ngtState struct {
	iqUUIDs      map[string][]float32 // insert queue
	dqUUIDs      map[string]struct{}  // delete queue
	indexedUUIDs map[string][]float32
}

func (sy *ngtSystem) Insert(v *vector) {
	sy.insertedUUIDs[v.uuid] = v.vector
}

func (sy *ngtSystem) Delete(uuid string) {
	delete(sy.insertedUUIDs, uuid)
}

func (sy *ngtSystem) PickInsertedVector() *vector {
	for uuid, v := range sy.insertedUUIDs {
		return &vector{
			uuid:   uuid,
			vector: v,
		}
	}

	return nil
}

func (st *ngtState) Exists(uuid string) bool {
	if _, ok := st.indexedUUIDs[uuid]; ok {
		return true
	}
	if _, ok := st.iqUUIDs[uuid]; ok {
		return true
	}

	return false
}

func (st *ngtState) Insert(v *vector) {
	st.iqUUIDs[v.uuid] = v.vector
}

func (st *ngtState) Indexing() {
	for uuid, v := range st.iqUUIDs {
		st.indexedUUIDs[uuid] = v
		delete(st.iqUUIDs, uuid)
	}

	for uuid := range st.dqUUIDs {
		delete(st.dqUUIDs, uuid)
		delete(st.indexedUUIDs, uuid)
	}
}

func (st *ngtState) Delete(uuid string) {
	st.dqUUIDs[uuid] = struct{}{}
}

func uuidGenSample() string {
	for {
		uuid, ok := uuidGen.Sample()
		if ok {
			return uuid.(string)
		}
	}
}

func f32sliceGenSample() []float32 {
	for {
		v, ok := f32sliceGen.Sample()
		if ok {
			return v.([]float32)
		}
	}
}

var (
	insertCommand = &commands.ProtoCommand{
		Name: "Insert",
		RunFunc: func(
			systemUnderTest commands.SystemUnderTest,
		) commands.Result {
			sy := systemUnderTest.(*ngtSystem)
			ngt := sy.ngt

			v := &vector{
				uuid:   uuidGenSample(),
				vector: f32sliceGenSample(),
			}
			vecCh <- v
			sy.Insert(v)

			err := ngt.Insert(v.uuid, v.vector)
			return &resultContainer{
				uuid:   v.uuid,
				vector: v.vector,
				err:    err,
			}
		},
		NextStateFunc: func(state commands.State) commands.State {
			select {
			case v := <-vecCh:
				state.(*ngtState).Insert(v)
			default:
			}
			return state
		},
		PreConditionFunc: func(state commands.State) bool {
			return true
		},
		PostConditionFunc: func(
			state commands.State,
			result commands.Result,
		) *gopter.PropResult {
			st := state.(*ngtState)
			rc := result.(*resultContainer)

			if st.Exists(rc.uuid) {
				if rc.err != nil {
					return &gopter.PropResult{Status: gopter.PropTrue}
				}
			}

			if rc.err != nil {
				return &gopter.PropResult{
					Status: gopter.PropFalse,
					Error:  rc.err,
					Labels: []string{"error"},
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

			err := ngt.CreateAndSaveIndex(ctx, 10000)
			return &resultContainer{
				err: err,
			}
		},
		NextStateFunc: func(state commands.State) commands.State {
			state.(*ngtState).Indexing()
			return state
		},
		PreConditionFunc: func(state commands.State) bool {
			st := state.(*ngtState)
			return len(st.iqUUIDs) > 0
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
					Labels: []string{"error"},
				}
			}

			return &gopter.PropResult{Status: gopter.PropTrue}
		},
	}

	existsCommand = &commands.ProtoCommand{
		Name: "Exists",
		RunFunc: func(
			systemUnderTest commands.SystemUnderTest,
		) commands.Result {
			ngt := systemUnderTest.(*ngtSystem).ngt

			uuid := uuidGenSample()

			_, exists := ngt.Exists(uuid)
			return &resultContainer{
				uuid:   uuid,
				exists: exists,
			}
		},
		NextStateFunc: func(state commands.State) commands.State {
			return state
		},
		PreConditionFunc: func(state commands.State) bool {
			return true
		},
		PostConditionFunc: func(
			state commands.State,
			result commands.Result,
		) *gopter.PropResult {
			st := state.(*ngtState)
			rc := result.(*resultContainer)

			if st.Exists(rc.uuid) {
				if !rc.exists {
					return &gopter.PropResult{
						Status: gopter.PropFalse,
						Labels: []string{"uuid exists"},
					}
				}
			} else {
				if rc.exists {
					return &gopter.PropResult{
						Status: gopter.PropFalse,
						Labels: []string{"uuid does not exist"},
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
			ngt := systemUnderTest.(*ngtSystem).ngt

			v := f32sliceGenSample()

			res, err := ngt.Search(v, 10, 0.1, -1.0)
			return &resultContainer{
				err:     err,
				results: res,
			}
		},
		NextStateFunc: func(state commands.State) commands.State {
			return state
		},
		PreConditionFunc: func(state commands.State) bool {
			return true
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
					Labels: []string{"error"},
				}
			}

			return &gopter.PropResult{Status: gopter.PropTrue}
		},
	}

	searchByIDCommand = &commands.ProtoCommand{
		Name: "SearchByID",
		RunFunc: func(
			systemUnderTest commands.SystemUnderTest,
		) commands.Result {
			sy := systemUnderTest.(*ngtSystem)
			ngt := systemUnderTest.(*ngtSystem).ngt

			v := sy.PickInsertedVector()
			if v == nil {
				return &resultContainer{}
			}

			res, err := ngt.SearchByID(v.uuid, 10, 0.1, -1.0)
			return &resultContainer{
				err:     err,
				results: res,
			}
		},
		NextStateFunc: func(state commands.State) commands.State {
			return state
		},
		PreConditionFunc: func(state commands.State) bool {
			return true
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
					Labels: []string{"error"},
				}
			}

			return &gopter.PropResult{Status: gopter.PropTrue}
		},
	}

	updateCommand = &commands.ProtoCommand{
		Name: "Update",
		RunFunc: func(
			systemUnderTest commands.SystemUnderTest,
		) commands.Result {
			sy := systemUnderTest.(*ngtSystem)
			ngt := systemUnderTest.(*ngtSystem).ngt

			v := sy.PickInsertedVector()
			if v == nil {
				return &resultContainer{}
			}

			v.vector = f32sliceGenSample()

			vecCh <- v
			sy.Insert(v)

			err := ngt.Update(v.uuid, v.vector)
			return &resultContainer{
				err: err,
			}
		},
		NextStateFunc: func(state commands.State) commands.State {
			select {
			case v := <-vecCh:
				state.(*ngtState).Insert(v)
			default:
			}
			return state
		},
		PreConditionFunc: func(state commands.State) bool {
			return true
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
					Labels: []string{"error"},
				}
			}

			return &gopter.PropResult{Status: gopter.PropTrue}
		},
	}

	deleteCommand = &commands.ProtoCommand{
		Name: "Delete",
		RunFunc: func(
			systemUnderTest commands.SystemUnderTest,
		) commands.Result {
			sy := systemUnderTest.(*ngtSystem)
			ngt := systemUnderTest.(*ngtSystem).ngt

			v := sy.PickInsertedVector()
			if v == nil {
				return &resultContainer{}
			}

			vecCh <- v
			sy.Delete(v.uuid)

			err := ngt.Delete(v.uuid)
			return &resultContainer{
				err: err,
			}
		},
		NextStateFunc: func(state commands.State) commands.State {
			select {
			case v := <-vecCh:
				state.(*ngtState).Delete(v.uuid)
			default:
			}
			return state
		},
		PreConditionFunc: func(state commands.State) bool {
			return true
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
					Labels: []string{"error"},
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
			n, err := New(cfg.Bind())
			if err != nil {
				t.Fatalf("error: %s", err)
			}

			ctx := context.Background()
			ech := n.Start(ctx)

			time.Sleep(time.Second)

			st := initialState.(*ngtState)
			for uuid, v := range st.indexedUUIDs {
				err = n.Insert(uuid, v)
				if err != nil {
					t.Fatalf("error: %s", err)
				}
			}

			if len(st.indexedUUIDs) > 0 {
				time.Sleep(10 * time.Millisecond)
				err = n.CreateAndSaveIndex(ctx, 1000)
				if err != nil {
					t.Logf("initialState indexedUUIDs: %#v", st.indexedUUIDs)
					t.Fatalf("error: %s", err)
				}
			}

			for uuid, v := range st.iqUUIDs {
				err = n.Insert(uuid, v)
				if err != nil {
					t.Fatalf("error: %s", err)
				}
			}

			nn, ok := n.(*ngt)
			if !ok {
				t.Fatal("cannot convert NGT to ngt")
			}

			return &ngtSystem{
				ctx:           ctx,
				ech:           ech,
				ngt:           nn,
				insertedUUIDs: map[string][]float32{},
			}
		},
		DestroySystemUnderTestFunc: func(sys commands.SystemUnderTest) {
			s := sys.(*ngtSystem)
			s.ngt.Close(s.ctx)
			s.ctx.Done()
		},
		InitialStateGen: gen.Const(&ngtState{
			iqUUIDs:      map[string][]float32{},
			dqUUIDs:      map[string]struct{}{},
			indexedUUIDs: map[string][]float32{},
		}),
		GenCommandFunc: func(state commands.State) gopter.Gen {
			st := state.(*ngtState)

			cs := make([]interface{}, 0)
			cs = append(cs, insertCommand)

			if len(st.iqUUIDs) > 0 {
				cs = append(cs, createIndexCommand)
			}

			if len(st.indexedUUIDs) > 0 {
				cs = append(
					cs,
					existsCommand,
					searchCommand,
					searchByIDCommand,
				)
			}

			if len(st.iqUUIDs) > 0 || len(st.indexedUUIDs) > 0 {
				cs = append(
					cs,
					updateCommand,
					deleteCommand,
				)
			}

			return gen.OneConstOf(cs...)
		},
		InitialPreConditionFunc: func(state commands.State) bool {
			return true
		},
	}
}

func TestStatefulNGT(t *testing.T) {
	// initialize logger
	log.Init(log.WithLoggerType("nop"))

	parameters := gopter.DefaultTestParameters()
	if seed != 0 {
		parameters.SetSeed(seed)
	}
	parameters.MinSuccessfulTests = minSuccessfulTests

	properties := gopter.NewProperties(parameters)
	properties.Property("NGT", commands.Prop(rootCommands(t)))

	properties.TestingRun(t)
}
