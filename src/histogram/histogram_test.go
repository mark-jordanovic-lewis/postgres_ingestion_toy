package histogram

import (
	"math/rand"
	"os"
	"testing"
)

func TestBuildHistogram(t *testing.T) {
	rand.Seed(int64(0))
	data := make([]float64, 500)
	for i := range data {
		data[i] = rand.Float64()
	}
	BuildHistogram("test_hist.png", data)
	if _, err := os.Open("test_hist.png"); err != nil {
		t.Errorf("No histogram file made")
	}
}
