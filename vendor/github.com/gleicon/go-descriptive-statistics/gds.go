// Author Gleicon Moraes (github.com/gleicon)

package descriptive_statistics

import (
	"fmt"
	"math"
	"sort"
	"strconv"
)

type Enum []float64

func NewEnumFromStringVector(sv []string) *Enum {
	l := len(sv)
	e := make(Enum, l)
	for count, value := range sv {
		fv, err := strconv.ParseFloat(value, 64)
		if err != nil {
			fmt.Printf("Error converting array to float array %s", err)
		}
		e[count] = fv
	}
	return &e
}

func (e Enum) Len() int {
	return len(e)
}
func (e Enum) Less(i, j int) bool {
	return e[i] < e[j]
}
func (e Enum) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e Enum) Sum() float64 {
	var sum float64
	for _, value := range e {
		sum += value
	}
	return sum
}

func (e Enum) Count() int {
	return len(e)
}

func (e Enum) Number() float64 {
	return float64(len(e))
}

func (e Enum) Mean() float64 {
	return e.Sum() / e.Number()
}

func (e Enum) Median() float64 {
	return e.Percentile(50.0)
}

func (e Enum) Variance() float64 {
	m := e.Mean()
	t := make(Enum, len(e))
	for i, value := range e {
		t[i] = math.Pow((m - value), 2)
	}
	r := t.Sum() / e.Number()
	return r
}

func (e Enum) StandardDeviation() float64 {
	return math.Sqrt(e.Variance())
}

func (e Enum) Percentile(p float64) float64 {
	s := make(Enum, len(e))
	copy(s, e)
	sort.Sort(s)
	if p == 100.0 {
		return s[len(s)-1]
	}

	rank := p / 100.0 * (e.Number() - 1)
	lrank := math.Floor(rank)
	d := rank - lrank
	lower := s[int(rank)]
	upper := s[int(rank)+1]
	r := lower + (upper-lower)*d
	return r
}
