package generator

import (
	"math/big"
	"math/rand"
	"time"
)

type DataFields struct {
	src   big.Int
	dst   big.Int
	flags big.Int
}

var maxBigInt = big.NewInt(9223372036854775807)
var source = rand.NewSource(time.Now().Unix())
var rng = rand.New(source)

func randomBigInt() big.Int {
	r := big.NewInt(0)
	return *r.Rand(rng, maxBigInt)
}

func newDataFields() (dfs DataFields) {
	dfs.src = randomBigInt()
	dfs.dst = randomBigInt()
	dfs.flags = randomBigInt()
	return
}

func NewDataSet(n int) []DataFields {
	slice := make([]DataFields, n)
	for i := range slice {
		slice[i] = newDataFields()
	}
	return slice
}
