package main

import (
	"fmt"
	gen "generator"
	"math"
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
		pq_conn.OpenTransaction()
		// change to for, concurrentify
		if pq_conn.ConnectionOpen {
			fmt.Println("Connection Opened")
			start := time.Now()
			pq_conn.IngestData(data_set)
			end := time.Now()

			if pq_conn.BatchIngested {
				fmt.Println("Batch Ingested")
				dt_ns := (end.UnixNano() - start.UnixNano())
				dt_s := float64(dt_ns) / math.Pow10(9)
				row_per_s = float64(n_rows) / dt_s
				fmt.Printf("rows/s: %v\n", row_per_s)
				mu_rps += row_per_s
			}
		}
		// pq_conn.DropAllRows()
	}
	mu_rps /= float64(m_batches)
	fmt.Printf("mean rows/s: %v\n", mu_rps)
}
