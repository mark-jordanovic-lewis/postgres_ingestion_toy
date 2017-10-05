package main

import (
	dt "connections/connectionTimer"
	alpha "connections/pq_conn"
	"fmt"
	gen "generator"
	hist "histogram"
	"math"
	"math/rand"
	"time"
)

// TODO: add interface so can pass all XXXConnection types to copy rate calc

func main() {
	var mu_rpns float64
	var m_times []float64
	rand.Seed(time.Now().UnixNano())

	m_batches := rand.Intn(9999) + 1 // 10000

	pq_conn := alpha.MakeConnection("swarmtest", "ingestion_test")

	for m := 0; m < m_batches; m++ {
		if m == 0 {
			m_times = make([]float64, m_batches)
		}
		n_rows := rand.Intn(999) + 1 // 100000
		data_set := gen.NewDataSet(n_rows)
		rpns := float64(n_rows) / dt.TimeIngestion(data_set, pq_conn)
		m_times[m] = rpns * math.Pow10(9)
		mu_rpns += rpns

		if pq_conn.BatchIngested() {
			fmt.Printf("Batch Ingested: %v of %v\n", m, m_batches)
			fmt.Printf("rows: %v, rows/s: %v\n", n_rows, rpns*math.Pow10(9))
		} else {
			fmt.Printf("Batch Rejected: %v of %v\n    !!!", m, m_batches)
		}
		pq_conn.DropAllRows()
		fmt.Println()
	}
	hist.BuildHistogram("first.png", m_times)
	fmt.Println("Drawing histogram")

	mu_rpns /= float64(m_batches)
	mu_rps := mu_rpns * math.Pow10(9)
	fmt.Println()
	fmt.Printf("Mean rows/s: %v\n", mu_rps)
}
