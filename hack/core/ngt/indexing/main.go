//
// Copyright (C) 2019 kpango (Yusuke Kato)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package main

import (
	"flag"
	"strings"

	"github.com/vdaas/vald/hack/core/ngt/assets"
	"github.com/vdaas/vald/internal/core/ngt"
	"github.com/vdaas/vald/internal/log"
)

var (
	distanceType = flag.String("distance-type", "L2", "choice indexing distance type")
)

func main() {
	log.Init(log.DefaultGlg())

	flag.Parse()
	distanceTypeOption := ngt.L2
	switch strings.ToLower(*distanceType) {
	case "l1":
		distanceTypeOption = ngt.L1
	case "l2":
		distanceTypeOption = ngt.L2
	case "angle":
		distanceTypeOption = ngt.Angle
	case "hamming":
		distanceTypeOption = ngt.Hamming
	case "cos":
		distanceTypeOption = ngt.Cosine
	}
	indexName := flag.Arg(0)

	d, err := assets.LoadDataset(flag.Arg(1))
	if err != nil {
		log.Warn(err)
		return
	}
	defer log.Infof("[%s] Done.", indexName)
	opts := []ngt.Option{
		ngt.WithDistanceType(distanceTypeOption),
	}
	loader := func() ([][]float64, error) {
		vectors, err := d.LoadTrain()
		if err != nil {
			return nil, err
		}
		defer log.Infof("[%s] %d items", d.Path, len(vectors))
		return vectors, nil
	}
	if err := assets.CreateIndex(indexName, loader, opts...); err != nil {
		log.Warn(err)
		return
	}
}
