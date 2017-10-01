package pq_conn

import (
	"generator"
	"testing"
)

func TestMakeConnection(t *testing.T) {
	t.Log("Checking connection to postgres database")
	conn := MakeConnection()
	if conn == nil {
		t.Errorf("Expected connection to postgres but got %v", conn)
	}
}

func TestCheckConnectionState(t *testing.T) {
	conn := MakeConnection()
	conn.CheckConnectionState()
	if conn.ConnectionOpen {
		t.Errorf("Not expecting to be connected to DB on init")
	}
	if conn.ConnectionBroken {
		t.Errorf("Connection to DB broken on init")
	}
}

func TestSetupTransaction(t *testing.T) {
	t.Log("Opening a transaction on the DB")
	conn := MakeConnection()
	conn.OpenTransaction()
	if conn.Txn == nil {
		t.Errorf("Expected transaction to be opened but got %v", conn.Txn)
	}
	if !conn.ConnectionOpen {
		t.Errorf("Transaction not connected.")
	}
}

func TestCopyData(t *testing.T) {
	var exit error
	var scanned SwarmRow
	t.Log("Copying in a row to the table")
	data := []generator.DataFields{generator.NewDataFields()}
	conn := MakeConnection()
	conn.OpenTransaction()
	if conn.ConnectionOpen {
		if !conn.IngestData(data) {
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
			}
		}
	} else {
		t.Errorf(
			"Transaction connected: %v\n conn.Txn: %T : %v",
			conn.ConnectionOpen, conn.Txn, conn.Txn)
	}
}
