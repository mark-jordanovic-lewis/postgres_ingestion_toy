package generator

import (
	"math/big"
	"testing"
)

func TestRandomBigInt(t *testing.T) {
	t.Log("Generating a random bigint")
	r := randomBigInt()
	ensureBigInt(interface{}(r), t)
}

func TestDifferentBigIntsGenerated(t *testing.T) {
	t.Log("Generating two different random bigint")
	a := randomBigInt()
	b := randomBigInt()
	if (&a).Cmp(&b) == 0 {
		t.Errorf("Expected two different random numbers but got %v == %v", a, b)
	}
}

func TestNewDataField(t *testing.T) {
	t.Log("Generating a DataFields object")
	dfs := newDataFields()
	ensureBigInt(interface{}(dfs.src), t)
	ensureBigInt(interface{}(dfs.dst), t)
	ensureBigInt(interface{}(dfs.flags), t)
}

func TestNewDataSet(t *testing.T) {
	t.Log("Generating Data Set")
	n := 30
	data := NewDataSet(n)
	ensureSliceDateFields(interface{}(data), t)
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

func ensureBigInt(i interface{}, t *testing.T) {
	switch iType := i.(type) {
	case big.Int:
		t.Log("Correct type returned")
	default:
		t.Errorf("Expected big.Int but got %v", iType)
	}
}

func ensureSliceDateFields(i interface{}, t *testing.T) {
	switch iType := i.(type) {
	case []DataFields:
		t.Log("Correct type returned")
	default:
		t.Errorf("Expected []DataFields but got %v", iType)
	}
}
