package qsortparallel

import "sync"
import "sync/atomic"
import "runtime"

// count how many threads
// are ongoing to avoid scheduling overheads
// with more goroutines than cores
var concurrent_thread_count int32

var _thread_amount int32 = int32(runtime.NumCPU())


func Swap(array []int, i,j int) {
	temp := array[j]
	array[j] = array[i]
	array[i] = temp
}

// for quicksort
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

// InsertionSort: self explanatory but with start and end
// indices to be used with quicksort
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
// start,end is start index, end index
func QsortParallel(array []int, start, end int) {
	// invalid indices
	if (start > end || start < 0 || end < 0) {
		return
	}
	atomic.AddInt32(&concurrent_thread_count,1)
	var wg sync.WaitGroup
	
	wg.Add(1)
	_qsortParallelImpl(array, start, end, &wg)
	wg.Wait();

	atomic.AddInt32(&concurrent_thread_count,-1)
}

// QsortParallel: the implementation
func _qsortParallelImpl(array []int, start, end int, wg* sync.WaitGroup) {
	defer wg.Done()

	
	if (start > end) {
		return
	}

	if (end - start < 100) {
		InsertionSort(array,start,end)
		return
	}

	pivot := (start + end) / 2

	pivot_real_pos := Partition(array, pivot, start, end)

	wg.Add(2)
	if (atomic.LoadInt32(&concurrent_thread_count) <= _thread_amount) {

		atomic.AddInt32(&concurrent_thread_count,1)

		go _qsortParallelImpl(array, start, pivot_real_pos - 1, wg)

		atomic.AddInt32(&concurrent_thread_count,-1)
	} else {
		_qsortParallelImpl(array, start, pivot_real_pos - 1, wg)
	}

	if (atomic.LoadInt32(&concurrent_thread_count) <= _thread_amount) {

		atomic.AddInt32(&concurrent_thread_count,1)

		go _qsortParallelImpl(array, pivot_real_pos + 1, end, wg)

		atomic.AddInt32(&concurrent_thread_count,-1)
	} else {
		_qsortParallelImpl(array, pivot_real_pos + 1, end,  wg)
	}
}

// Merge: merge two slices of ints
// for mergesort
func Merge(a []int, b []int) []int {
	sizea := len(a);
	sizeb := len(b);

	new_arr := make([]int,sizea + sizeb)

	var a_i,b_i, n_i int;

	a_i = 0
	b_i = 0

	for ; n_i < sizea + sizeb; n_i++ {	
		if (a_i < sizea && b_i < sizeb && a[a_i] <= b[b_i]) || (b_i == sizeb) {
			new_arr[n_i] = a[a_i]
			a_i++
		} else {
			new_arr[n_i] = b[b_i]
			b_i++
		}

	}

	return new_arr
}


func _parallelMergeSort(a []int, ret chan []int, wg* sync.WaitGroup) {

	defer wg.Done()

	sizea := len(a)

	if (sizea < 64) {
		InsertionSort(a,0,sizea - 1)
		ret <- a
		return
	}

	mid := sizea / 2

	chan_left := make(chan []int,1)
	chan_right := make(chan []int,1)

	wg.Add(2)
	if (atomic.LoadInt32(&concurrent_thread_count) <= _thread_amount) {

		atomic.AddInt32(&concurrent_thread_count,1)

		go _parallelMergeSort(a[0:mid], chan_left, wg)

		atomic.AddInt32(&concurrent_thread_count,-1)
	} else {
		_parallelMergeSort(a[0:mid], chan_left, wg)
	}

	if (atomic.LoadInt32(&concurrent_thread_count) <= _thread_amount) {

		atomic.AddInt32(&concurrent_thread_count,1)

		go _parallelMergeSort(a[mid:], chan_right, wg)

		atomic.AddInt32(&concurrent_thread_count,-1)
	} else {
		_parallelMergeSort(a[mid:], chan_right, wg)
	}

	
	
	left_arr := <-chan_left
	right_arr :=  <-chan_right

	ret <- Merge(left_arr, right_arr)
}

// ParallelMergeSort: a parallel version of mergesort
// which assigns work to goroutines if threads are available
func ParallelMergeSort(a []int) []int {

	atomic.AddInt32(&concurrent_thread_count,1)
	var wg sync.WaitGroup

	res_chan := make(chan []int,1)
	
	wg.Add(1)
	_parallelMergeSort(a, res_chan,&wg)
	wg.Wait();


	atomic.AddInt32(&concurrent_thread_count,-1)

	return <-res_chan
}






