package pq_conn

import (
	gen "generator"
	alpha "pq_conn"
	"testing"
)

func BenchmarkIngestData(b *testing.B) {
	data_set := gen.NewDataSet(100000)
	pq_conn := alpha.MakeConnection("swarmtest", "ingestion_test")
	pq_conn.OpenTransaction()
	// change to for, concurrentify
	if pq_conn.ConnectionOpen {
		b.StartTimer()
		pq_conn.IngestData(data_set)
		b.StopTimer()
		pq_conn.DropAllRows()
	}
}
