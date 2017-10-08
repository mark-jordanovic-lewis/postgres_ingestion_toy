package generator

import (
	"math/rand"
	"time"
)

type DataFields struct {
	Src   int64
	Dst   int64
	Flags int64
}

var maxBigInt = int64(9223372036854775807)
var source = rand.NewSource(time.Now().Unix())
var rng = rand.New(source)

// this does not cover -9223372036854775808, only -9223372036854775807
func randomBigInt() int64 {
	if rng.Float32() > 0.5 {
		return -1 * rng.Int63n(maxBigInt)
	}
	return rng.Int63n(maxBigInt)
}

func NewDataFields() (dfs DataFields) {
	dfs.Src = randomBigInt()
	dfs.Dst = randomBigInt()
	dfs.Flags = randomBigInt()
	return
}

func NewDataSet(n int) []DataFields {
	slice := make([]DataFields, n)
	for i := range slice {
		slice[i] = NewDataFields()
	}
	return slice
}

func NewPgxDataSet(n int) [][]interface{} {
	row := make([]interface{}, 3)
	pgx_data := make([][]interface{}, n)
	for i := 0; i < n; i++ {
		row[0] = randomBigInt()
		row[1] = randomBigInt()
		row[2] = randomBigInt()
		pgx_data[i] = row
	}
	return pgx_data
}
