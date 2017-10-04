package calculator

import (
	"math"
	"sort"
	"time"
)

// calculate the PqBatchCopy rate
func CalculatePqBatchCopyRate(timestamps []time.Time) float64 {
	var total int64

	unixstamps := make([]int64, len(timestamps))
	for i, ts := range timestamps {
		unixstamps[i] = ts.UnixNano()
	}
	sort.Sort(Int64Slice(unixstamps))

	t0 := unixstamps[0]
	for _, t1 := range unixstamps[1:] {
		total += t0 - t1
		t0 = t1
	}
	return (float64(len(timestamps)) / float64(total)) * math.Pow10(9) // ns -> s
}

// Needed for interfacing int64 with sort package.
// the Less in this is for reversing, switch it for a real sort
type Int64Slice []int64

func (slice Int64Slice) Len() int           { return len(slice) }
func (slice Int64Slice) Swap(i, j int)      { slice[i], slice[j] = slice[j], slice[i] }
func (slice Int64Slice) Less(i, j int) bool { return slice[i] > slice[j] }
