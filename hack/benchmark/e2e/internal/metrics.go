//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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
package internal

func sum(x []float64) float64 {
	s := 0.0
	for _, a := range x {
		s += a
	}
	return s
}

func mean(x []float64) float64 {
	return sum(x) / float64(len(x))
}

func std(x []float64) float64 {
	x2 := make([]float64, len(x))
	for i, a := range x {
		x2[i] = a * a
	}
	m := mean(x)
	return mean(x2) - m*m
}

func Recall(datasetNeighbors, runNeighbors []string, k int) (r float64) {
	dn := map[string]struct{}{}
	for _, n := range datasetNeighbors[:k] {
		dn[n] = struct{}{}
	}
	for i := 0; i < k; i++ {
		if _, ok := dn[runNeighbors[i]]; ok {
			r++
		}
	}
	return r / float64(k)
}

func Recalls(datasetNeighbors, runNeighbors [][]string, k int) (recalls []float64) {
	recalls = make([]float64, len(runNeighbors))
	for i, d := range datasetNeighbors {
		recalls[i] = Recall(d, runNeighbors[i], k)
	}
	return recalls
}

func MeanStdRecalls(datasetNeighbors, runNeighbors [][]string, k int) (float64, float64, []float64) {
	recalls := Recalls(datasetNeighbors, runNeighbors, k)
	return mean(recalls), std(recalls), recalls
}
