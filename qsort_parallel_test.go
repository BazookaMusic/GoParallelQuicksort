package qsortparallel

import (
	"testing"
	"math/rand"
)

func TestQsortParallel(t *testing.T) {
	const N_ELEMS = 1000000
	var nums [N_ELEMS]int

	// repeat 10 tests
	for j := 0; j < 100; j++ {
		for i := 0; i < N_ELEMS; i++ {
			nums[i] = rand.Int() % N_ELEMS
		}
		
		QsortParallel(nums[:], 0, N_ELEMS - 1)

		before := nums[0]
		for i := 1; i < N_ELEMS; i++ {
			if (nums[i] < before) {
				t.Errorf("Wrong order in the elements")
			}
			before = nums[i]
		}
	}

	// wrong indexes should just do nothing
	QsortParallel(nums[:],-1,1)

}

func BenchmarkQsortParallel(b *testing.B) {
	const N_ELEMS = 1000000
	var nums [N_ELEMS]int

	for i := 0; i < N_ELEMS; i++ {

		for j:= 0 ; j < b.N; j++ {
			nums[i] = rand.Int() % b.N
		}
		
	}

	b.ResetTimer()

	QsortParallel(nums[:],0,N_ELEMS - 1)

}

