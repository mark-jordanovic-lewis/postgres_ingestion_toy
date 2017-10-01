package main

import (
	gen "generator"
	alpha "pq_conn"
)

func main() {
	data_set := gen.NewDataSet(100000)
	pq_conn := alpha.MakeConnection("swarm_benchmark")
	pq_conn.OpenTransaction()
	if pq_conn.ConnectionOpen {
		pq_conn.IngestData(data_set)
	}
}
