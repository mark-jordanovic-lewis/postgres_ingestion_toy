package histogram

import (
	"fmt"
	gen "generator"
	"math"
	conn "connections/pg_conn"
	"time"
)

type histogrammable interface {
	OpenTransaction()
	IngestData([]DataFields)
}

func GenerateHistogramIO(data_set []gen.DataFields, conn *histogrammable rows_per_s float64 {
	conn.OpenTransaction()
	// change to for, concurrentify
	if conn.ConnectionOpen {
		fmt.Println("Connection Opened")
		rows_per_s = timeIngestion(conn, data_set)
	}
	return
}



func timeIngestion(conn *conn.PqConnection, data_set []gen.DataFields) {
  // split data and concurrentify
  rps := eatData(conn, data_set)

  if conn.BatchIngested {
    fmt.Printf("Batch Ingested: %v of %v\n", m, m_batches)
    fmt.Printf("rows: %v, rows/s: %v\n", n_rows, rps)
  } else {
    fmt.Printf(("Batch Rejected: %v of %v\n", m, m_batches) !!!)
  }
  return rps
}

func eatData(conn *conn.PqConnection, data_set []gen.DataFields) {
	start := time.Now()
	conn.IngestData(data_set)
	end := time.Now()
 return (end.UnixNano() - start.UnixNano())

}