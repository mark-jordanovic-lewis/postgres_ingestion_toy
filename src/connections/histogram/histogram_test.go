package histogram

import (
	alpha "connections/pq_conn"
	gen "generator"
	"os"
	"testing"
)

func TestGenerateHistogramIO(t *testing.T) {
	pq_conn := alpha.MakeConnection("swarmtest", "ingestion_test")
	data_set := gen.NewDataSet(100)
	dt := float64(generateHistogramIO(data_set, pq_conn))
	if dt < 0 {
		t.Error("No good dt value: %v", dt)
	}
	os.Open("pq_histgram")
	pq_conn.DropAllRows()
}
