package pgx_conn

import (
	"generator"
	"testing"
	"time"
)

func TestMakeConnection(t *testing.T) {
	t.Log("Checking connection to postgres database")
	conn := MakeConnection("swarmtest", "ingestion_test")
	if conn == nil {
		t.Errorf("Expected connection to postgres but got %v", conn)
	}
}

func TestCheckConnectionState(t *testing.T) {
	conn := MakeConnection("swarmtest", "ingestion_test")
	conn.CheckConnectionState()
	if !conn.ConnectionOpen() {
		t.Errorf("Expecting to be connected to DB on init")
	}
	if conn.ConnectionBroken() {
		t.Errorf("Connection to DB broken on init")
	}
}

func TestCopyData(t *testing.T) {
	var exit error
	var scanned SwarmRow
	t.Log("Copying in a row to the table")
	data := generator.NewPgxDataSet(1)
	conn := MakeConnection("swarmtest", "ingestion_test")

	conn.CopyData(data)
	if !conn.BatchIngested() {
		t.Errorf("Could not copy data into DB")
	} else {
		query_response, err := conn.Conn.Query(`SELECT * FROM ingestion_test`)
		if err != nil {
			t.Errorf("Error in querying DB: %v", err.Error())
		}
		for query_response.Next() {
			if exit = query_response.Scan(
				&scanned.Ts,
				&scanned.Src,
				&scanned.Dst,
				&scanned.Flags); exit != nil {
				t.Errorf("Error in reading rows: %v", exit.Error())
			}
			if incorrectData(scanned, data[0]) {
				t.Errorf(
					"data out does not match data in: %v, %v, %v != %v, %v, %v",
					*scanned.Src, *scanned.Dst, *scanned.Flags,
					data[0][0], data[0][1], data[0][2])
			}
		}
	}
	// have to do a clean up of the test db
	conn.DropAllRows()
}

func incorrectData(a SwarmRow, b []interface{}) bool {
	return !(*a.Src == b[0] && *a.Dst == b[1] && *a.Flags == b[2])
}

// SwarmRow : row return struct
type SwarmRow struct {
	Ts    *time.Time
	Src   *int64
	Dst   *int64
	Flags *int64
}
