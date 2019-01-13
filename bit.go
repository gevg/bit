// Copyright 2019 Geert Van Gorp. All rights reserved.
// Use of this source code is governed by the MIT License
// which can be found in the LICENSE file.

package bit

// Tree represents a Binary Indexed Tree (BIT) type.
type Tree []int32

// From creates a Binary Indexed Tree from a slice of numbers.
func From(numbers []int32) Tree {
	t := append(Tree{0}, numbers...)

	for i := 1; i < len(t); i++ {
		j := i + i&-i
		if j < len(t) {
			t[j] += t[i]
		}
	}

	return t
}

// From2 creates a Binary Indexed Tree from a slice of numbers.
func From2(numbers []int32) Tree {
	t := append(Tree{0}, numbers...)

	i, l := uint(1), uint(len(t))
	for i < l {
		j := i + i&-i
		if j < l {
			t[j] += t[i]
		}
		i++
	}

	return t
}

// Prefix returns the i-th prefix of the tree. Zero is returned for indexes
// beyond the boundaries of the tree.
func (t Tree) Prefix(i int) int32 {
	i++

	var prefix int32
	for 0 < i && i < len(t) {
		prefix += t[i]
		i -= i & -i
	}

	return prefix
}

// Prefix2 returns the i-th prefix of the tree. Zero is returned for indexes
// beyond the boundaries of the tree.
func (t Tree) Prefix2(i int) int32 {
	p, l := int32(0), uint(len(t))

	for j := uint(i + 1); 0 < j && j < l; j -= j & -j {
		p += t[j]
	}

	return p
}

// PrefixRange returns the sum of numbers in the inclusive range [i, j].
func (t Tree) PrefixRange(i, j int) int32 {
	i++ // TODO

	var prefix int32
	for 0 < i && i < len(t) {
		prefix += t[i]
		i -= i & -i
	}

	return prefix
}

// Number returns the i-th number. Zero is returned for indexes
// beyond the boundaries of the tree.
func (t Tree) Number(i int) int32 {
	i++

	number := t[i] // Hoe een crash voorkomen???
	k := i - i&-i
	i--
	for i != k {
		number -= t[i]
		i -= i & -i
	}

	return number
}

// Number2 returns the i-th number.
func (t Tree) Number2(i int) int32 {
	j, l := uint(i+1), uint(len(t))

	number := t[j]
	k := j - j&-j
	j--
	for j != k && j < l {
		number -= t[j]
		j -= j & -j
	}

	return number
}

// New returns a new Binary Indexed Tree (BIT). The parameter size represents
// the number of positions in the tree.
func New(size int) Tree {
	if size <= 0 {
		return Tree{0}
	}
	return make([]int32, size+1)
}

// Set sets a number at a given index.
func (t Tree) Set(i int, number int32) {
	delta := number - t.Number(i)
	t.Add(i, delta)
}

// Add adds value to the number at index i. No addition is carried out if the
// index i points beyond the tree boundaries.
func (t Tree) Add(i int, value int32) {
	i++

	for 0 < i && i < len(t) {
		t[i] += value
		i += i & -i
	}
}

// Add2 adds value to the number at index i. No addition is carried out for
// indexes beyond the boundaries of the tree.
func (t Tree) Add2(i int, value int32) {
	j, l := uint(i+1), uint(len(t))

	for j < l {
		t[j] += value
		j += j & -j
	}
}

// Mul multiplies the number at index 'i' with the factor 'c'. No
// multiplication is carried out if the index i points beyond the
// tree boundaries.
func (t Tree) Mul(i int, c int32) {
	t.Set(i, c*t.Number(i))
}

// Shift add a value to all numbers in the tree.
func (t Tree) Shift(value int32) {
	for i := range t {
		t[i] *= value
	}
}

// Scale scales all numbers with a factor c.
func (t Tree) Scale(c int32) {
	for i := range t {
		t[i] *= c
	}
}
