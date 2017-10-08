package pq_conn

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
	if conn.ConnectionOpen() {
		t.Errorf("Not expecting to be connected to DB on init")
	}
	if conn.ConnectionBroken() {
		t.Errorf("Connection to DB broken on init")
	}
}

func TestSetupTransaction(t *testing.T) {
	t.Log("Opening a transaction on the DB")
	conn := MakeConnection("swarmtest", "ingestion_test")
	conn.OpenTransaction()
	if conn.Txn == nil {
		t.Errorf("Expected transaction to be opened but got %v", conn.Txn)
	}
	if !conn.ConnectionOpen() {
		t.Errorf("Transaction not connected.")
	}
}

func TestDropAllRows(t *testing.T) {
	data := generator.NewDataSet(5)
	conn := MakeConnection("swarmtest", "ingestion_test")
	conn.OpenTransaction()
	conn.IngestData(data)
	if !conn.DropAllRows() {
		t.Errorf("Reported that there was an error.")
	}
	rows, _ := conn.Conn.Query("SELECT * FROM ingestion_test")
	i := 0
	for rows.Next() {
		i++
	}
	if i != 0 {
		t.Error("There is still data in the table.")
	}
}

func TestCopyData(t *testing.T) {
	var exit error
	var scanned SwarmRow
	t.Log("Copying in a row to the table")
	data := []generator.DataFields{generator.NewDataFields()}
	conn := MakeConnection("swarmtest", "ingestion_test")
	conn.OpenTransaction()
	if conn.ConnectionOpen() {
		conn.IngestData(data)
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
						data[0].Src, data[0].Dst, data[0].Flags)
				}
			}
		}
	} else {
		t.Errorf(
			"Transaction connected: %v\n conn.Txn: %T : %v",
			conn.ConnectionOpen(), conn.Txn, conn.Txn)
	}
	// have to do a clean up of the test db
	conn.DropAllRows()
}

func TestSelectTimestamps(t *testing.T) {
	data := generator.NewDataSet(50)
	conn := MakeConnection("swarmtest", "ingestion_test")
	conn.OpenTransaction()
	conn.IngestData(data)
	tss := conn.SelectTimeStamps()
	l := len(tss)
	if l == 0 {
		t.Errorf("Expected to get some data but got none.")
	}
	if l != 50 {
		t.Errorf("Expected to get 50 data but got %v.", l)
	}
	conn.DropAllRows()
}

// Test Helpers

func incorrectData(a SwarmRow, b generator.DataFields) bool {
	return !(*a.Src == b.Src && *a.Dst == b.Dst && *a.Flags == b.Flags)
}

// SwarmRow : row return struct
type SwarmRow struct {
	Ts    *time.Time
	Src   *int64
	Dst   *int64
	Flags *int64
}
