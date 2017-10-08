package connectionTimer

import (
	pq "connections/pq_conn"
	"fmt"
	gen "generator"
	"time"
)

// PQ Connection Timer \\
// =================== \\
// make serial timer and concurrent timerource
func TimeIngestion(data_set []gen.DataFields, conn *pq.PqConnection) (rows_per_s float64) {
	conn.OpenTransaction()
	// change to for, concurrentify
	if conn.ConnectionOpen() {
		fmt.Println("Connection Opened")
		rows_per_s = eatData(conn, data_set)
	}

	return
}

func eatData(conn *pq.PqConnection, data_set []gen.DataFields) float64 {
	start := time.Now()
	conn.IngestData(data_set)
	end := time.Now()
	return float64((end.UnixNano() - start.UnixNano()))

}
