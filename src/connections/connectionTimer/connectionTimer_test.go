package connectionTimer

import (
	beta "connections/pgx_conn"
	alpha "connections/pq_conn"
	gen "generator"
	"testing"
)

// PQ Ingestion
func TestTimeIngestion(t *testing.T) {
	pq_conn := alpha.MakeConnection("swarmtest", "ingestion_test")
	data_set := gen.NewDataSet(100)
	dt := float64(TimeIngestion(data_set, pq_conn))
	if dt < 0 {
		t.Errorf("No good dt value: %v", dt)
	}
	pq_conn.DropAllRows()
}

// Pgx Ingestion
func TestTimeCopy(t *testing.T) {
	pgx_conn := beta.MakeConnection("swarmtest", "ingestion_test")
	data_set := gen.NewPgxDataSet(100)
	dt := TimeCopy(data_set, pgx_conn)
	if dt < 0 {
		t.Errorf("No good dt value: %v", dt)
	}
	pgx_conn.DropAllRows()
}
