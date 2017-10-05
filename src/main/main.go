package main

import (
	"fmt"
	gen "generator"
	"math/rand"
	alpha "pq_conn"
	"time"
)

// TODO: add interface so can pass all XXXConnection types to copy rate calc

var source = rand.NewSource(time.Now().UnixNano())
var rng = rand.New(source)

func main() {
	var row_per_s, mu_rps float64
	m_batches := rand.Intn(10000)

	pq_conn := alpha.MakeConnection("swarmtest", "ingestion_test")

	for m := 0; m < m_batches; m++ {
		n_rows := rand.Intn(100000)
		data_set := gen.NewDataSet(n_rows)
		dt := float64(
			generateHistogramIO(data_set, pq_conn))
		mu_rps += float64(
			n_rows) / dt_ns
		pq_conn.DropAllRows()
	}
	mu_rps /= float64(m_batches)

	fmt.Printf("mean rows/s: %v\n", mu_rps)
}
