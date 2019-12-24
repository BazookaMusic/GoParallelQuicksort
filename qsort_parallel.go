package qsortparallel

import "sync"
import "sync/atomic"
import "runtime"

var concurrent_thread_count int32

var THREAD_AMOUNT int32 = int32(runtime.NumCPU())

func Swap(array []int, i,j int) {
	temp := array[j]
	array[j] = array[i]
	array[i] = temp
}


func Partition(array []int, pivot,start,end int) int {
	Swap(array, pivot, end)

	if (start > end) {
		return -1
	}

	actual_pivot := array[end]

	i := start - 1

	for j := start; j < end; j++ {
		if array[j] < actual_pivot {
			i++
			Swap(array, i, j)
		}
	}

	i++

	Swap(array, end, i)
	return i
}

// InsertionSort: sort with insertion sort if array size
// is too small
func InsertionSort(array []int, start, end int) {
	var min_till_now int
	for i := start; i < end; i++ {

		min_index := i
		min_till_now = array[i]

		for j := i + 1; j <= end; j++ {
			if array[j] < min_till_now {
				min_till_now = array[j]
				min_index = j
			}

		}

		Swap(array,i, min_index)
	}
}

//QsortParallel: sort with qsort in parallel
// start is start,end is start index, end index
func QsortParallel(array []int, start, end int) {
	// invalid indices
	if (start > end || start < 0 || end < 0) {
		return
	}
	atomic.AddInt32(&concurrent_thread_count,1)
	var wg sync.WaitGroup
	
	wg.Add(1)
	_QsortParallelImpl(array, start, end, &wg)
	wg.Wait();
	atomic.AddInt32(&concurrent_thread_count,1)
}

// QsortParallel: the implementation
func _QsortParallelImpl(array []int, start, end int, wg* sync.WaitGroup) {
	defer wg.Done()

	
	if (start > end) {
		return
	}

	if (end - start < 64) {
		InsertionSort(array,start,end)
		return
	}

	pivot := (start + end) / 2

	pivot_real_pos := Partition(array, pivot, start, end)

	wg.Add(2)
	if (atomic.LoadInt32(&concurrent_thread_count) < THREAD_AMOUNT) {
		atomic.AddInt32(&concurrent_thread_count,1)
		go _QsortParallelImpl(array, start, pivot_real_pos - 1, wg)
		atomic.AddInt32(&concurrent_thread_count,-1)
	} else {
		_QsortParallelImpl(array, start, pivot_real_pos - 1, wg)
	}

	if (atomic.LoadInt32(&concurrent_thread_count) < THREAD_AMOUNT) {
		atomic.AddInt32(&concurrent_thread_count,1)
		go _QsortParallelImpl(array, pivot_real_pos + 1, end, wg)
		atomic.AddInt32(&concurrent_thread_count,-1)
	} else {
		_QsortParallelImpl(array, pivot_real_pos + 1, end,  wg)
	}
}

// func MergeSortParallelImpl(array []int, start, end int, wg* sync.WaitGroup) {
// 	defer wg.Done()

// 	atomic.AddInt32(&concurrent_thread_count,1)
	
// 	if (start > end) {
// 		return
// 	}

// 	if (end - start < 64) {
// 		InsertionSort(array,start,end)
// 		return
// 	}


// }
