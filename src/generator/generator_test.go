package generator

import (
	"testing"
)

func TestDifferentBigIntsGenerated(t *testing.T) {
	t.Log("Generating two different random bigint")
	a := randomBigInt()
	b := randomBigInt()
	if a == b {
		t.Errorf("Expected two different random numbers but got %v == %v", a, b)
	}
}

func TestNewDataSet(t *testing.T) {
	t.Log("Generating Data Set")
	n := 30
	data := NewDataSet(n)
	ensureSliceDataFields(interface{}(data), t)
	if len := len(data); len != 30 {
		t.Errorf("Expected slice of length 30 but got %v", len)
	}
}

func TestMaxBatchSize(t *testing.T) {
	t.Log("Creating largest possible batch size for swarm64 code test")
	n := 100000
	data := NewDataSet(n)
	if len := len(data); len != 100000 {
		t.Errorf("Expected slice of length 100000 but got %v", len)
	}
}

func ensureSliceDataFields(i interface{}, t *testing.T) {
	switch iType := i.(type) {
	case []DataFields:
		t.Log("Correct type returned")
	default:
		t.Errorf("Expected []DataFields but got %v", iType)
	}
}
