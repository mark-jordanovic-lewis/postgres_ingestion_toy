package calculator

import (
	"math"
	"math/rand"
	"testing"
	"time"
)

var source = rand.NewSource(time.Now().Unix())
var rng = rand.New(source)

func TestCalculatePqBatchCopyRate(t *testing.T) {
	var mu float64

	periods := make([]float64, 100)
	for i := range periods {
		periods[i] = float64(rng.Int63n(5))
	}
	for _, t := range periods {
		mu += t
	}
	mu = (100.0 / mu) * math.Pow10(3) // ms -> s

	tss := make([]time.Time, 100)
	for i, _ := range tss {
		tss[i] = time.Now()
		time.Sleep(time.Duration(periods[i]) * time.Millisecond)
	}

	// it's hard to verify time based functions. Not the most precise but it'll do.
	// mean is always larger than pure mu (dependent on other procs being run on machine)
	if mean := CalculatePqBatchCopyRate(tss); mean-mu > math.Sqrt(mean) {
		t.Errorf("Unreasonably bad match of mean time taken: %v ! ~= %v", mean, mu)
	}

}
