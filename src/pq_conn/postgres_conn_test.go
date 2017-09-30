package pq_conn

import "testing"

func TestMakeConnection(t *testing.T) {
	t.Log("Checking connection to postgres database")
	conn := MakeConnection()
	if conn == nil {
		t.Errorf("Expected connection to postgres but got %v", conn)
	}
}
