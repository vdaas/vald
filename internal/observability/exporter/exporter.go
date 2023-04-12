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
package exporter

import "context"

type Exporter interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

var DefaultMillisecondsHistogramDistribution = []float64{
	0.01,
	0.05,
	0.1,
	0.3,
	0.6,
	0.8,
	1,
	2,
	3,
	4,
	5,
	6,
	8,
	10,
	13,
	16,
	20,
	25,
	30,
	40,
	50,
	65,
	80,
	100,
	130,
	160,
	200,
	250,
	300,
	400,
	500,
	650,
	800,
	1000,
	2000,
	5000,
	10000,
	20000,
	50000,
	100000,
}
