package pq_conn

import "testing"

func TestMakeConnection(t *testing.T) {
	t.Log("Checking connection to postgres database")
	conn := MakeConnection()
	if conn == nil {
		t.Errorf("Expected connection to postgres but got %v", conn)
	}
}

func TestPingDB(t *testing.T) {
	conn := MakeConnection()
	if !conn.PingDB() {
		t.Errorf("Connection not possible")
	}
}

func TestSetupTransaction(t *testing.T) {
	t.Log("Opening a transaction on the DB")
	conn := MakeConnection()
	txn := conn.OpenTransaction()
	if txn == nil {
		t.Errorf("Expected transaction to be opened but got %v", txn)
	}
}
