// Copyright 2019 Geert Van Gorp. All rights reserved.
// Use of this source code is governed by the MIT License
// which can be found in the LICENSE file.

package bit

import (
	"math/bits"
)

// Tree represents a Binary Indexed Tree (BIT) type.
type Tree []int32

// New creates a Binary Indexed Tree of n elements.
// If n is not provided, the tree length defaults to zero.
func New(n ...int) Tree {
	if len(n) == 0 || n[0] <= 0 {
		return Tree{}
	}
	return make([]int32, n[0])
}

// Len returns the number of elements in the tree.
func Len(t Tree) int {
	return len(t)
}

// From creates a Binary Indexed Tree from a slice of numbers.
// When the reUse option is set, the tree will use the numbers
// input slice as its backing store, avoiding new allocations.
// Default behavior for the tree is to allocate its own backing store.
func From(numbers []int32, reUse ...bool) Tree {
	var t Tree

	if len(reUse) == 0 || reUse[0] == false {
		t = make(Tree, len(numbers))
		copy(t, numbers)
	} else {
		t = numbers
	}

	for i := range t {
		if j := i | (i + 1); 0 <= j && j < len(t) {
			t[j] += t[i]
		}
	}

	return t
}

// Reset initializes the length of the tree to zero, but keeps the
// backing store. After Reset, the tree can be re-used with Append.
func (t *Tree) Reset() {
	*t = (*t)[:0]
}

// Copy does a deep copy of the src tree. If the dst tree is smaller
// than src, only part of the BIT is copied, up to the length of dst.
// Copy returns the number of elements copied.
func Copy(dst, src Tree) int {
	if len(dst) <= len(src) {
		return copy(dst, src)
	}

	n := copy(dst, src)

	// compute partial sums for dst[n:]
	for i := n; 0 <= i && i < len(dst); i++ {
		var num int32
		j, k := i, i&(i+1)
		for k < j && 0 < j && j < len(dst) {
			num += dst[j-1]
			j &= j - 1
		}
		dst[i] = num
	}
	return len(dst)
}

// Append adds numbers to the back of the tree.
func Append(t Tree, number ...int32) Tree {
	if len(number) == 1 {
		iNum := len(t)
		t = append(t, number[0])

		_ = t[iNum]
		i, j := iNum, iNum&(iNum+1)
		for j < i && 0 < i && i < len(t) {
			t[iNum] += t[i-1]
			i &= i - 1
		}
		return t
	}

	if len(number) > 1 {
		l := len(t)
		t = append(t, number...)

		var imin int
		if 0 < l {
			imin = 1<<(bits.Len(uint(l))-1) - 1
		}

		for i := imin; 0 <= i && i < len(t); i++ {
			j := i | (i + 1)
			if 0 <= j && l <= j && j < len(t) {
				t[j] += t[i]
			}
		}
		return t
	}

	return t
}

// Sum returns the prefix sum at index i of the tree. If i is larger than the
// largest index of the tree, the prefix sum of the largest index is returned.
func (t Tree) Sum(i int) int32 {
	if len(t) <= i {
		i = len(t) - 1
	}

	// compute prefix sum at index i by adding relevant partial sums.
	var sum int32
	for 0 <= i && i < len(t) {
		sum += t[i]
		i = i&(i+1) - 1
	}

	return sum
}

// RangeSum returns the prefix sum of the [lo, hi) range. In case of a partial
// overlap of the range with the tree, RangeSum will return the prefix sum
// of the intersection of the given interval with the interval of the tree.
func (t Tree) RangeSum(lo, hi int) int32 {
	if hi-lo < 0 {
		return 0
	}

	lo, hi = lo-1, hi-1
	lenhi := bits.LeadingZeros64(uint64(hi))
	lenlo := bits.LeadingZeros64(uint64(lo))

	var sum int32
	if lenhi != lenlo {
		// compute prefix sum at index hi by adding relevant partial sums
		for 0 <= hi && hi < len(t) {
			sum += t[hi]
			hi = hi&(hi+1) - 1
		}

		// compute prefix sum at index lo by adding relevant partial sums
		for 0 <= lo && lo < len(t) {
			sum -= t[lo]
			lo = lo&(lo+1) - 1
		}
		return sum
	}

	for {
		switch {
		case lo < hi && 0 <= hi && hi < len(t):
			sum += t[hi]
			hi = hi&(hi+1) - 1
		case hi < lo && 0 <= lo && lo < len(t):
			sum -= t[lo]
			lo = lo&(lo+1) - 1
		default:
			return sum
		}
	}
}

// Sums returns the prefix sums of the tree. If the length of the sums slice
// is too small, the Sums function fills the slice starting from index 0 and
// stops when the slice is full. Sums returns the number of elements in the
// sums slice.
func (t Tree) Sums(sums []int32) int {
	i := 0
	for i < len(sums) {
		var sum int32
		j := i
		// calculate sum[i] prefix sum by adding relevant partial sums
		for 0 <= j && j < len(t) {
			sum += t[j]
			j = j&(j+1) - 1
		}
		sums[i] = sum
		i++
	}

	return i
}

// Number returns the number at index i.
// If i is outside of the tree, 0 will be returned.
func (t Tree) Number(i int) int32 {
	if i < 0 || len(t) <= i {
		return 0
	}

	// calculate number by subtracting relevant partial sums from t[i]
	number := t[i]
	j := i & (i + 1)
	for j < i && 0 < i && i < len(t) {
		number -= t[i-1]
		i &= i - 1
	}

	return number
}

// RangeNumbers returns in the buf variable a slice of numbers, as defined by
// the given boundaries. The upper bound is not included. If the lo index is
// out of boundaries, zero will be returned.
func (t Tree) RangeNumbers(lo int, buf []int32) int {
	if lo < 0 || lo >= len(t) {
		return 0
	}

	i, j := 0, lo
	for i < len(buf) && j < len(t) {
		buf[i] = t.Number(j)
		i, j = i+1, j+1
	}

	return i
}

// Numbers returns the numbers in the tree. The caller provides the array to
// store the numbers. If the numbers slice is too short, only numbers up to
// the length of the slice will be returned.
func (t Tree) Numbers(numbers []int32) int {
	n := copy(numbers, t)

	i := len(t)&^1 - 1
	for 0 < i && i < n && i < len(numbers) {
		k := i & (i + 1)
		for j := i; k < j && 0 < j && j < len(numbers); j &= j - 1 {
			numbers[i] -= numbers[j-1]
		}
		i -= 2
	}

	return n
}

// SetNumber sets a number at a given index. If the
// index is outside of the tree, no change is made.
func (t Tree) SetNumber(i int, number int32) {
	if i < 0 || len(t) <= i {
		return
	}

	// calculate delta by subtracting relevant partial sums from number
	j := i
	number -= t[i]
	k := i & (i + 1)
	for k < i && 0 < i && i < len(t) {
		number += t[i-1]
		i &= i - 1
	}

	// add delta to relevant partial sums
	for 0 <= j && j < len(t) {
		t[j] += number
		j |= j + 1
	}
}

// AddNumber adds the given value to the tree at index i. If
// the index is outside of the tree boundaries, no value is added.
func (t Tree) AddNumber(i int, value int32) {
	// add value to relevant partial sums
	for 0 <= i && i < len(t) {
		t[i] += value
		i |= i + 1
	}
}

// MulNumber multiplies the number at index i with the given value. If the
// index is outside of the tree boundaries, no modifications are done.
func (t Tree) MulNumber(i int, value int32) int32 {
	if i < 0 || len(t) <= i {
		return 0
	}

	// calculate number by subtracting relevant partial sums
	j := i
	number := t[i]
	k := i & (i + 1)
	for k < i && 0 < i && i < len(t) {
		number -= t[i-1]
		i &= i - 1
	}

	// calculate delta that needs to be added
	delta := number * (value - 1)

	// add delta by adding to the relevant partial sums
	for 0 <= j && j < len(t) {
		t[j] += delta
		j |= j + 1
	}

	return number + delta
}

// Shift increases all numbers in the tree with the given value.
func (t Tree) Shift(value int32) {
	for i := range t {
		t[i] += value << bits.TrailingZeros64(uint64(i)+1)
	}
}

// Scale scales all numbers in the tree with the given factor.
func (t Tree) Scale(value int32) {
	for i := range t {
		t[i] *= value
	}
}

// RangeAdd adds a slice of numbers to the numbers in the tree
// at index i and subsequent indices.
func (t Tree) RangeAdd(i int, numbers []int32) {
	for j := 0; j < len(numbers) && i < len(t); i, j = i+1, j+1 {
		t.AddNumber(i, numbers[j])
	}
}

// RangeMul multiplies a slice of numbers to the numbers in the tree,
// starting at index i.
func (t Tree) RangeMul(i int, factors []int32) {
	for j := 0; j < len(factors) && i < len(t); i, j = i+1, j+1 {
		t.MulNumber(i, factors[j])
	}
}

// RangeSet sets a slice of numbers in the tree, starting at index i.
func (t Tree) RangeSet(i int, numbers []int32) {
	for j := 0; j < len(numbers) && i < len(t); i, j = i+1, j+1 {
		t.SetNumber(i, numbers[j])
	}
}

// RangeShift adds the given value to all numbers in the [lo, hi) index
// range of the tree. If lo and/or hi are outside the boundaries of the
// tree, the [lo, hi) range will be intersected with the tree range.
func (t Tree) RangeShift(lo, hi int, value int32) {
	cumVal := value
	if lo < 0 {
		lo = 0
	}

	for lo < hi && 0 <= lo && lo < len(t) {
		val := value << bits.TrailingZeros64(uint64(lo)+1)
		minVal := minabs(cumVal, val)

		t[lo] += minVal

		if j := lo | (lo + 1); hi <= j {
			for 0 <= j && j < len(t) {
				t[j] += minVal
				j |= j + 1
			}
		}

		cumVal += value
		lo++
	}
}

func minabs(a, b int32) int32 {
	if a <= b {
		if 0 <= a {
			return a
		}
		return b
	}

	if 0 <= b {
		return b
	}
	return a
}

// RangeScale scales all numbers in the [lo, hi) range
// of the tree with the given multiplier.
func (t Tree) RangeScale(lo, hi int, multiplier int32) {
	for i := lo; i < hi && 0 <= i && i < len(t); i++ {
		t.MulNumber(i, multiplier)
	}
}

// SearchSum returns the largest index and corresponding prefix sum that is
// smaller than or equal to the given value. In case the tree is empty, -1 is
// returned.This operation assumes the prefix sums to increase monotonically.
func (t Tree) SearchSum(value int32) (int, int32) {
	if len(t) == 0 {
		return -1, 0
	}

	lo, hi := 0, 1<<(bits.Len(uint(len(t)))-1)
	toSearch := value

	for hi != 0 {
		if m := lo + hi; 0 < m && m <= len(t) {
			if toSearch >= t[m-1] {
				lo += hi
				toSearch -= t[m-1]
			}
		}
		hi >>= 1
	}

	return lo - 1, value - toSearch
}
