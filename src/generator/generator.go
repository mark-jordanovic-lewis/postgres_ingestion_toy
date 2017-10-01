package generator

import (
	"math/big"
	"math/rand"
	"time"
)

type DataFields struct {
	Src   big.Int
	Dst   big.Int
	Flags big.Int
}

var maxBigInt = big.NewInt(9223372036854775807)
var source = rand.NewSource(time.Now().Unix())
var rng = rand.New(source)

// this does not cover -9223372036854775808, only -9223372036854775807
func randomBigInt() big.Int {
	r := big.NewInt(0)
	r.Rand(rng, maxBigInt)
	if rng.Float32() > 0.5 {
		r.Neg(r)
	}
	return *r
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
