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
		x2[i] = a*a
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
