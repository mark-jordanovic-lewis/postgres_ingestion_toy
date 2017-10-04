package pq_conn

import (
	"fmt"
	"testing"
)

func BenchmarkIngestData(b *testing.B) {
	var row_rate float64
	data_set := gen.NewDataSet(100000)
	pq_conn := alpha.MakeConnection("swarmtest", "ingestion_test")
	pq_conn.OpenTransaction()
	// change to for, concurrentify
	if pq_conn.ConnectionOpen {
		fmt.Println("Connection Opened")
		pq_conn.IngestData(data_set)

		if pq_conn.BatchIngested {
			fmt.Println("Batch Ingested")
			timestamps := pq_conn.SelectTimeStamps()
			fmt.Println(timestamps)
			row_rate = calculator.CalculatePqBatchCopyRate(timestamps)
		}
		pq_conn.DropAllRows()
	}
}
