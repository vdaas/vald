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

package cassandra

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"testing"

	"github.com/vdaas/vald/internal/db/nosql/cassandra"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/log/logger"
)

var (
	metaTable = "backup_vector"

	uuidColumn   = "uuid"
	vectorColumn = "vector"
	metaColumn   = "meta"
	ipsColumn    = "ips"

	metaColumnSlice = []string{uuidColumn, vectorColumn, metaColumn, ipsColumn}

	dropStmt = "DROP TABLE IF EXISTS vald.backup_vector;"

	schema = `
CREATE TABLE vald.backup_vector (
  uuid   text,
  vector blob,
  ips    list<text>,
  PRIMARY KEY (uuid)
);
`

	c cassandra.Cassandra
)

type ReadMetaVector struct {
	UUID   string   `json:"uuid"   yaml:"uuid"`
	Vector string   `json:"vector" yaml:"vector"`
	Meta   string   `json:"meta"   yaml:"meta"`
	IPs    []string `json:"ips"    yaml:"ips"`
}

type MetaVector struct {
	UUID   string   `json:"uuid"   db:"uuid"`
	Vector []byte   `json:"vector" db:"vector"`
	Meta   string   `json:"meta"   db:"meta"`
	IPs    []string `json:"ips"    db:"ips"`
}

func init() {
	log.Init(log.WithLoggerType(logger.NOP.String()))
	var err error
	c, err = cassandra.New(
		cassandra.WithHosts(
			"",
		),
		cassandra.WithCQLVersion("3.0.0"),
		cassandra.WithProtoVersion(0),
		cassandra.WithTimeout("3s"),
		cassandra.WithConnectTimeout("5s"),
		cassandra.WithPort(9042),
		cassandra.WithKeyspace("vald"),
		cassandra.WithNumConns(2),
		cassandra.WithConsistency("quorum"),
		cassandra.WithUsername("root"),
		cassandra.WithPassword(""),
		cassandra.WithRetryPolicyNumRetries(3),
		cassandra.WithRetryPolicyMinDuration("1s"),
		cassandra.WithRetryPolicyMaxDuration("5s"),
		cassandra.WithReconnectionPolicyMaxRetries(3),
		cassandra.WithReconnectionPolicyInitialInterval("1s"),
		cassandra.WithSocketKeepalive("0s"),
		cassandra.WithMaxPreparedStmts(1000),
		cassandra.WithMaxRoutingKeyInfo(1000),
		cassandra.WithPageSize(5000),
		cassandra.WithEnableHostVerification(false),
		cassandra.WithDefaultTimestamp(false),
		cassandra.WithReconnectInterval("1s"),
		cassandra.WithMaxWaitSchemaAgreement("1s"),
		cassandra.WithIgnorePeerAddr(false),
		cassandra.WithDisableInitialHostLookup(false),
		cassandra.WithDisableNodeStatusEvents(false),
		cassandra.WithDisableTopologyEvents(false),
		cassandra.WithDisableSkipMetadata(false),
		cassandra.WithDefaultIdempotence(false),
		cassandra.WithWriteCoalesceWaitTime("1s"),
	)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	err = c.Open(ctx)
	if err != nil {
		panic(err)
	}
}

func loadData() []MetaVector {
	f, err := os.Open("./testdata.json")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	var vs []*ReadMetaVector
	if err := json.NewDecoder(f).Decode(&vs); err != nil {
		panic(err)
	}

	vv := make([]MetaVector, 0, len(vs))
	for _, v := range vs {
		vv = append(vv, MetaVector{
			UUID:   v.UUID,
			Vector: []byte(v.Vector),
			Meta:   v.Meta,
			IPs:    v.IPs,
		})
	}

	return vv
}

func initTable(b *testing.B, metas []MetaVector) {
	if err := c.Query(dropStmt, []string{}).Exec(); err != nil {
		b.Fatal(err)
	}

	if err := c.Query(schema, []string{}).Exec(); err != nil {
		b.Fatal(err)
	}

	ib := cassandra.Insert(metaTable).Columns(metaColumnSlice...)
	bt := cassandra.Batch()
	entities := make(map[string]interface{}, len(metas)*4)
	for i, m := range metas {
		prefix := "p" + strconv.Itoa(i)
		bt = bt.AddWithPrefix(prefix, ib)
		entities[prefix+"."+uuidColumn] = m.UUID
		entities[prefix+"."+vectorColumn] = m.Vector
		entities[prefix+"."+metaColumn] = m.Meta
		entities[prefix+"."+ipsColumn] = m.IPs
	}
	if err := c.Query(bt.ToCql()).BindMap(entities).ExecRelease(); err != nil {
		b.Fatal(err)
	}
}

func BenchmarkGocqlxSelectBindMap(b *testing.B) {
	var val MetaVector

	metas := loadData()
	initTable(b, metas)

	keys := make([]map[string]interface{}, 0, len(metas))
	for _, m := range metas {
		keys = append(keys, map[string]interface{}{
			uuidColumn: m.UUID,
		})
	}

	b.StopTimer()
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		// query
		if err := c.Query(
			cassandra.Select(
				metaTable,
				metaColumnSlice,
				cassandra.Eq(uuidColumn),
			),
		).BindMap(keys[i%len(metas)]).GetRelease(&val); err != nil {
			b.Errorf("Error: %s", err)
		}

		// verify
		if val.UUID != keys[i%len(metas)][uuidColumn] {
			b.Errorf("Verify failed: %s != %s", val.UUID, keys[i%len(metas)][uuidColumn])
		}
	}
	b.StopTimer()
}

func BenchmarkGocqlxSelectBindStruct(b *testing.B) {
	var val MetaVector

	metas := loadData()
	initTable(b, metas)

	keys := make([]MetaVector, 0, len(metas))
	for _, m := range metas {
		keys = append(keys, MetaVector{
			UUID: m.UUID,
		})
	}

	b.StopTimer()
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		// query
		if err := c.Query(
			cassandra.Select(
				metaTable,
				metaColumnSlice,
				cassandra.Eq(uuidColumn),
			),
		).BindStruct(keys[i%len(metas)]).GetRelease(&val); err != nil {
			b.Errorf("Error: %s", err)
		}

		// verify
		if val.UUID != keys[i%len(metas)].UUID {
			b.Errorf("Verify failed: %s != %s", val.UUID, keys[i%len(metas)].UUID)
		}
	}
	b.StopTimer()
}
