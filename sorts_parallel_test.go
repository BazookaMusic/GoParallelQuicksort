package qsortparallel

import (
	"testing"
	"math/rand"
)

func TestQsortParallel(t *testing.T) {
	const N_ELEMS = 100000000
	var nums [N_ELEMS]int

	// repeat 10 tests
	for j := 0; j < 1; j++ {
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



func TestMergeSort(t *testing.T) {
	const N_ELEMS = 100000000
	a := make([]int, N_ELEMS);

	for i := 0; i < N_ELEMS; i++ {
		a[i] = rand.Int() % N_ELEMS
	}

	a = ParallelMergeSort(a)

	before := a[0]
		for i := 1; i < N_ELEMS; i++ {
			if (a[i] < before) {
				t.Errorf("Wrong order in the elements %d %d", a[i], before)
			}
			before = a[i]
	}
}



