package main

import (
	dt "connections/connectionTimer"
	beta "connections/pgx_conn"
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
	var pq_rpns, pgx_rpns, pgx_multi_rpns float64
	var pq_times, pgx_times, pgx_multi_times []float64
	var pgx_multi_data_set [][][]interface{}
	rand.Seed(time.Now().UnixNano())

	m_batches := rand.Intn(9999) + 1 // 10000

	pq_conn := alpha.MakeConnection("swarmtest", "ingestion_test")
	pgx_conn := beta.MakeConnection("swarmtest", "ingestion_test", 1)
	n_connections := 8
	pgx_multi_conn := beta.MakeConnection("swarmtest", "ingestion_test", n_connections)

	for m := 0; m < m_batches; m++ {
		if m == 0 {
			pq_times = make([]float64, m_batches)
			pgx_times = make([]float64, m_batches)
			pgx_multi_times = make([]float64, m_batches)
		}
		n_rows := rand.Intn(99999) + 1 // 100000
		pq_data_set := gen.NewDataSet(n_rows)
		pgx_data_set := gen.NewPgxDataSet(n_rows)
		pgx_multi_data_set = build_multi_data_set(n_rows, n_connections)

		// pq_conn - refactor to function.
		rpns := float64(n_rows) / dt.TimeIngestion(pq_data_set, pq_conn)
		pq_times[m] = rpns * math.Pow10(9)
		pq_rpns += rpns
		if pq_conn.BatchIngested() {
			fmt.Printf("pq_conn Batch Ingested: %v of %v\n", m, m_batches)
			fmt.Printf("rows: %v, rows/s: %v\n", n_rows, rpns*math.Pow10(9))
		} else {
			fmt.Printf("pq_conn Batch Rejected: %v of %v\n    !!!", m, m_batches)
		}
		// pq_conn.DropAllRows()
		fmt.Println()

		// pgx_conn
		rpns, err := pgx_conn.CopyData(pgx_data_set)
		rpns = (float64(n_rows) / rpns)
		pgx_times[m] = rpns * math.Pow10(9)
		pgx_rpns += rpns
		if pgx_conn.BatchIngested() {
			fmt.Printf("pgx_conn Batch Ingested: %v of %v\n", m, m_batches)
			fmt.Printf("rows: %v, rows/s: %v\n", n_rows, rpns*math.Pow10(9))
		} else if err != nil {
			fmt.Printf("pgx_conn Batch Rejected by CopyFrom: %v of %v      !!!\n%v", m, m_batches, err)
		} else {
			fmt.Printf("pgx_conn Batch Rejected: %v of %v\n       !!!", m, m_batches)
		}
		// pq_conn.DropAllRows()
		fmt.Println()

		// pgx multi conn
		rpns_sum := 0.0

		for i, set := range pgx_multi_data_set {
			rpns, err := pgx_multi_conn.CopyData(set)
			if err != nil {
				fmt.Printf(
					"pgx_multi_conn Batch#%v Rejected by CopyFrom: %v of %v      !!!\n%v",
					i, m, m_batches, err)
			}
			rpns_sum += rpns
		}
		rpns = (float64(n_rows) / rpns_sum)
		pgx_multi_times[m] = rpns * math.Pow10(9)
		pgx_multi_rpns += rpns
		// should be using ary of bools for this check.
		if pgx_multi_conn.BatchIngested() {
			fmt.Printf("pgx_conn Batch Ingested: %v of %v\n", m, m_batches)
			fmt.Printf("rows: %v, rows/s: %v\n", n_rows, rpns*math.Pow10(9))
		} else if err != nil {
			fmt.Printf("pgx_multi_conn Batch Rejected by CopyFrom: %v of %v      !!!\n%v", m, m_batches, err)
		} else {
			fmt.Printf("pgx_multi_conn Batch Rejected: %v of %v\n       !!!", m, m_batches)
		}
		// pq_conn.DropAllRows()
		fmt.Println()
	}
	hist.BuildHistogram("pq_conn.png", pq_times)
	fmt.Println("Drawn pq_conn histogram")
	hist.BuildHistogram("pgx_conn.png", pgx_times)
	fmt.Println("Drawn pgx_conn histogram")

	// pgx multi conn n connections in pool
	var pgx_multiconn_times [][]float64
	for m := 0; m < m_batches; m++ {
		if m == 0 {
			pgx_multiconn_times = make([][]float64, 8)
			pgx_multi_times = make([]float64, m_batches)
		}

		n_rows := rand.Intn(99999) + 1 // 100000
		n_connections := make([]int, 8)
		for i := range n_connections {
			n_connections[i] = 2 * (i + 1)
		}
		var rpns float64
		for x, n := range n_connections {
			pgx_multi_data_set = build_multi_data_set(n_rows, n)
			pgx_multi_conn := beta.MakeConnection("swarmtest", "ingestion_test", n)
			rpns_sum := 0.0
			for i, set := range pgx_multi_data_set {
				rpns, err := pgx_multi_conn.CopyData(set)
				if err != nil {
					fmt.Printf(
						"pgx_multi_conn Batch#%v Rejected by CopyFrom: %v of %v      !!!\n%v",
						i, m, m_batches, err)
				}
				rpns_sum += rpns
			}
			rpns = (float64(n_rows) / rpns_sum)
			pgx_multi_times[m] = rpns * math.Pow10(9)
			pgx_multi_rpns += rpns
			// pq_conn.DropAllRows()
			pgx_multiconn_times[x] = pgx_multi_times
		}
		for i, data := range pgx_multiconn_times {
			histfilename := fmt.Sprintf("pgx_multi_conn_%v.png", 2*(i+1))
			hist.BuildHistogram(histfilename, data)
			fmt.Printf("Drawn pgx_multi_conn histogram for batch: %v\n", m)
		}
	}

	pq_rpns /= float64(m_batches)
	mu_rps := pq_rpns * math.Pow10(9)
	fmt.Printf("Mean pq_conn rows/s: %v\n", mu_rps)
	pgx_rpns /= float64(m_batches)
	mu_rps = pgx_rpns * math.Pow10(9)
	fmt.Printf("Mean pgx_multi_conn rows/s: %v\n", mu_rps)
	pgx_multi_rpns /= float64(m_batches)
	mu_rps = pgx_multi_rpns * math.Pow10(9)
	fmt.Printf("Mean pgx_conn rows/s: %v\n", mu_rps)
	fmt.Println()
}

func build_multi_data_set(n_rows, n_connections int) [][][]interface{} {
	var _n_rows int
	pgx_multi_data_set := make([][][]interface{}, n_connections)
	for i := 0; i < n_connections; i++ {
		if i == 0 {
			_n_rows = n_rows/n_connections + int(math.Mod(float64(n_rows), float64(n_connections)))
		} else {
			_n_rows = n_rows / n_connections
		}
		pgx_multi_data_set[i] = gen.NewPgxDataSet(_n_rows)
	}
	return pgx_multi_data_set
}
