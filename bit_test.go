// Copyright 2019 Geert Van Gorp. All rights reserved.
// Use of this source code is governed by the MIT License
// which can be found in the LICENSE file.

package bit

import (
	"math/rand"
	"testing"
)

var testcases = []struct {
	numbers []int32
	sums    []int32
	tree    []int32
}{
	{
		[]int32{},
		[]int32{},
		[]int32{},
	},
	{
		[]int32{-10},
		[]int32{-10},
		[]int32{-10},
	},
	{
		[]int32{1, 3, 2, 2, 1, 0, 1, 2, 2, 1, 2},
		[]int32{1, 4, 6, 8, 9, 9, 10, 12, 14, 15, 17},
		[]int32{1, 4, 2, 8, 1, 1, 1, 12, 2, 3, 2},
	},
	{
		[]int32{1, 0, 2, 1, 1, 3, 0, 4, 2, 5, 2, 2, 3, 1, 0, 2},
		[]int32{1, 1, 3, 4, 5, 8, 8, 12, 14, 19, 21, 23, 26, 27, 27, 29},
		[]int32{1, 1, 2, 4, 1, 4, 0, 12, 2, 7, 2, 11, 3, 4, 0, 29},
	},
}

// Tests with an empty tree
func TestNew(t *testing.T) {
	tree := New()
	l := Len(tree)
	tree.Reset()
	dst := New()
	n := Copy(dst, tree)
	tree = Append(tree)
	l = Len(tree)
	sum := tree.Sum(0)
	sum = tree.RangeSum(0, 0)
	n = tree.Sums(nil)
	sum = tree.Number(0)
	n = tree.RangeNumbers(0, nil)
	n = tree.Numbers(nil)
	tree.SetNumber(0, -5)
	tree.AddNumber(0, -5)
	tree.MulNumber(0, -5)
	tree.Shift(-5)
	tree.Scale(-5)
	tree.RangeShift(0, 0, -5)
	tree.RangeScale(0, 0, -5)
	tree.RangeAdd(0, nil)
	tree.RangeMul(0, nil)
	tree.RangeSet(0, nil)
	tree.SearchSum(5)
	_, _, _ = n, l, sum
}

func TestFrom(t *testing.T) {
	for i, tc := range testcases {
		tree := From(tc.numbers)

		for j := range tc.tree {
			if tree[j] != tc.tree[j] {
				t.Errorf(
					"Testcase: %d, index: %d, got: %d != want: %d\n",
					i, j, tree[j], tc.tree[j],
				)
			}
		}
	}
}

func TestCopy(t *testing.T) {
	numbers := []int32{2, 4, 7}
	src := From(numbers)

	dst := make(Tree, Len(src))
	Copy(dst, src)
	for i, v := range dst {
		if src[i] != v {
			t.Errorf("Index: %d, src: %d != dst: %d\n", i, src[i], v)
		}
	}

	dst = make(Tree, Len(src)-1)
	Copy(dst, src)
	for i, v := range dst {
		if src[i] != v {
			t.Errorf("Index: %d, src: %d != dst: %d\n", i, src[i], v)
		}
	}

	dst = make(Tree, Len(src)+1)
	Copy(dst, src)
	for i, v := range src {
		if v != dst[i] {
			t.Errorf("Index: %d, src: %d != dst: %d\n", i, v, dst[i])
		}
	}
	for i := Len(src); i < Len(dst); i++ {
		if dst.Number(i) != 0 {
			t.Errorf("Index: %d, dst: %d != 0\n", i, dst[i])
		}
	}

}

func TestNewAndSet(t *testing.T) {
	for i, tc := range testcases {
		tree := New(len(tc.numbers))

		for j, v := range tc.numbers {
			tree.SetNumber(j, v)
		}

		for j, want := range tc.tree {
			if tree[j] != want {
				t.Errorf(
					"Testcase: %d, index: %d, got: %d != want: %d\n",
					i, j, tree[j], want,
				)
			}
		}
	}
}

func TestAppend(t *testing.T) {
	for i, tc := range testcases {
		var tree Tree
		for l := 0; l < len(tc.numbers); l++ {
			tree = From(tc.numbers[:l])

			for _, v := range tc.numbers[l:] {
				tree = Append(tree, v)
			}

			for j, want := range tc.tree {
				if tree[j] != want {
					t.Errorf(
						"Testcase: %d, l: %d, idx: %d, got: %d != want: %d\n",
						i, l, j, tree[j], want)
				}
			}

			tree.Reset()
			tree = From(tc.numbers[:l])
			tree = Append(tree, tc.numbers[l:]...)

			for j, want := range tc.tree {
				if tree[j] != want {
					t.Errorf(
						"Testcase: %d, l: %d, idx: %d, got: %d != want: %d\n",
						i, l, j, tree[j], want,
					)
				}
			}

			tree.Reset()
		}
	}
}

func TestSum(t *testing.T) {
	for i, tc := range testcases {
		tree := From(tc.numbers)
		for j := range tc.sums {
			if got, want := tree.Sum(j), tc.sums[j]; got != want {
				t.Errorf(
					"Testcase: %d, index: %d, got: %d != want: %d\n",
					i, j, got, want,
				)
			}
		}
	}
}

func TestRangeSum(t *testing.T) {
	for i, tc := range testcases {
		tree := From(tc.numbers)
		for j := range tc.sums {
			if got, want := tree.RangeSum(j, j+1), tc.numbers[j]; got != want {
				t.Errorf(
					"Testcase: %d, index: %d, got: %d != want: %d\n",
					i, j, got, want,
				)
			}
		}
	}
}

func TestRangeNumbers(t *testing.T) {
	buf := make([]int32, 1)
	for i, tc := range testcases {
		tree := From(tc.numbers)
		for j := range tc.numbers {
			tree.RangeNumbers(j, buf)
			want := tc.numbers[j]
			if buf[0] != want {
				t.Errorf(
					"Testcase: %d, index: %d, got: %d != want: %d\n",
					i, j, buf[0], want,
				)
			}
		}
	}
}

func TestRangeAdd(t *testing.T) {
	for i, tc := range testcases {
		tree := From(tc.numbers)
		tree.RangeAdd(4, []int32{-5, -5, -5, -5, -5})

		for j, want := range tc.numbers {
			if 4 <= j && j < 9 {
				want -= 5
			}
			if got := tree.Number(j); got != want {
				t.Errorf(
					"Testcase: %d, index: %d, got: %d != want: %d\n",
					i, j, tree.Number(j), want,
				)
			}
		}
	}
}

func TestRangeMul(t *testing.T) {
	for i, tc := range testcases {
		tree := From(tc.numbers)
		tree.RangeMul(4, []int32{-5, -5, -5, -5, -5})

		for j, want := range tc.numbers {
			if 4 <= j && j < 9 {
				want *= -5
			}
			if got := tree.Number(j); got != want {
				t.Errorf(
					"Testcase: %d, index: %d, got: %d != want: %d\n",
					i, j, tree.Number(j), want,
				)
			}
		}
	}
}

func TestRangeShift(t *testing.T) {
	for i, tc := range testcases {
		tree := From(tc.numbers)
		tree.RangeShift(4, 9, 5)

		for j, want := range tc.numbers {
			if 4 <= j && j < 9 {
				want += 5
			}
			if got := tree.Number(j); got != want {
				t.Errorf(
					"Testcase: %d, index: %d, got: %d != want: %d\n",
					i, j, tree.Number(j), want,
				)
			}
		}
	}
}

func TestRangeScale(t *testing.T) {
	for i, tc := range testcases {
		tree := From(tc.numbers)
		tree.RangeScale(4, 9, 5)

		for j, want := range tc.numbers {
			if 4 <= j && j < 9 {
				want *= 5
			}
			if got := tree.Number(j); got != want {
				t.Errorf(
					"Testcase: %d, index: %d, got: %d != want: %d\n",
					i, j, tree.Number(j), want,
				)
			}
		}
	}
}

func TestNumber(t *testing.T) {
	for i, tc := range testcases {
		tree := From(tc.numbers)
		for j, want := range tc.numbers {
			if got := tree.Number(j); got != want {
				t.Errorf(
					"Testcase: %d, index: %d, number got: %d != want: %d\n",
					i, j, got, want,
				)
			}
		}
	}
}

func TestSetNumber(t *testing.T) {
	for i, tc := range testcases {
		tree := From(tc.numbers)

		for j := range tc.numbers {
			tree.SetNumber(j, -5)
		}

		for j := range tc.numbers {
			if got, want := tree.Number(j), int32(-5); got != want {
				t.Errorf(
					"Testcase: %d, index: %d, number got: %d != want: %d\n",
					i, j, got, want,
				)
			}
		}
	}
}

func TestRangeSet(t *testing.T) {
	for i, tc := range testcases {
		tree := From(tc.numbers)
		tree.RangeSet(4, []int32{-5, -5, -5, -5, -5})

		for j, want := range tc.numbers {
			if 4 <= j && j < 9 {
				want = -5
			}
			if got := tree.Number(j); got != want {
				t.Errorf(
					"Testcase: %d, index: %d, got: %d != want: %d\n",
					i, j, tree.Number(j), want,
				)
			}
		}
	}
}

func TestAddNumber(t *testing.T) {
	for i, tc := range testcases {
		tree := From(tc.numbers)

		for j := range tc.numbers {
			tree.AddNumber(j, tree.Number(j))
		}

		for j := range tc.numbers {
			if got, want := tree.Number(j), 2*tc.numbers[j]; got != want {
				t.Errorf(
					"Testcase: %d, index: %d, got: %d != want: %d\n",
					i, j, got, want,
				)
			}
		}
	}
}

func TestMulNumber(t *testing.T) {
	const x = 8

	for i, tc := range testcases {
		tree := From(tc.numbers)

		for j := range tc.numbers {
			tree.MulNumber(j, x)
		}

		for j := range tc.numbers {
			if got, want := tree.Number(j), x*tc.numbers[j]; got != want {
				t.Errorf(
					"Testcase: %d, index: %d, got: %d != want: %d\n",
					i, j, got, want,
				)
			}
		}
	}
}

func TestShift(t *testing.T) {
	const x = 8

	for i, tc := range testcases {
		tree := From(tc.numbers)

		tree.Shift(x)

		for j := range tc.numbers {
			if got, want := tree.Number(j), tc.numbers[j]+x; got != want {
				t.Errorf(
					"Testcase: %d, index: %d, got: %d != want: %d\n",
					i, j, got, want,
				)
			}
		}
	}
}

func TestScale(t *testing.T) {
	for i, tc := range testcases {
		tree := From(tc.numbers)

		tree.Scale(5)

		for j := range tc.numbers {
			if got, want := tree.Number(j), 5*tc.numbers[j]; got != want {
				t.Errorf(
					"Testcase: %d, index: %d, got: %d != want: %d\n",
					i, j, got, want,
				)
			}
		}
	}
}

func TestSearchSum(t *testing.T) {
	for _, tc := range testcases {
		tree := From(tc.numbers)
		for j, toSearch := range tc.sums {
			got, gsum := tree.SearchSum(toSearch)
			if gsum != toSearch {
				t.Errorf(
					"got: (%d, %d) == want: (%d, %d)?\n",
					got, gsum, j, toSearch,
				)
			}
		}
	}
}

func TestNumbers(t *testing.T) {
	for i, tc := range testcases {
		tree := From(tc.numbers)
		numbers := make([]int32, len(tc.numbers))
		_ = tree.Numbers(numbers)

		for j, num := range tc.numbers {
			if numbers[j] != num {
				t.Errorf(
					"Testcase: %d, index: %d, got: %d == want: %d\n",
					i, j, numbers[j], num,
				)
			}
		}
	}
}

func TestSums(t *testing.T) {
	for i, tc := range testcases {
		tree := From(tc.numbers)
		sums := make([]int32, len(tc.numbers))
		tree.Sums(sums)
		for j, sum := range tc.sums {
			if sums[j] != sum {
				t.Errorf(
					"Testcase: %d, index: %d, got: %d == want: %d\n",
					i, j, sums[j], sum,
				)
			}
		}
	}
}

func TestLargerSamples(t *testing.T) {
	const (
		n       = 10
		lo      = n / 3
		hi      = 2 * lo
		max     = 100
		maxDiv2 = max / 2
		shift   = -9
		factor  = -5
	)

	numbers := make([]int32, n)
	rand.Seed(18)
	for i := range numbers {
		numbers[i] = rand.Int31n(max) // - maxDiv2
	}

	tree, newtree := From(numbers), New(len(numbers))

	var sum int32
	rangeOne := make([]int32, 1)
	for i, num := range numbers {
		// Sum()
		sum += num
		if tsum := tree.Sum(i); tsum != sum {
			t.Errorf("index: %d, sum got: %d != want: %d\n", i, tsum, sum)
		}

		// RangeSum()
		if number := tree.RangeSum(i, i+1); number != num {
			t.Errorf("index: %d, rangeSum got: %d != want: %d\n", i, number, num)
		}

		// Number()
		if number := tree.Number(i); number != num {
			t.Errorf("index: %d, number got: %d != want: %d\n", i, number, num)
		}

		// RangeNumbers()
		if tree.RangeNumbers(i, rangeOne); rangeOne[0] != num {
			t.Errorf(
				"index: %d, rangeNumber got: %d != want: %d\n",
				i, rangeOne[0], num,
			)
		}

		// RangeSet()
		rangeOne = []int32{num}
		if tree.RangeSet(i, rangeOne); rangeOne[0] != num {
			t.Errorf(
				"index: %d, rangeNumber got: %d != want: %d\n",
				i, rangeOne[0], num,
			)
		}

		// RangeAdd()
		tree.RangeAdd(i, []int32{num})
		if tree.RangeAdd(i, []int32{-num}); tree.Number(i) != num {
			t.Errorf(
				"index: %d, rangeNumber got: %d != want: %d\n",
				i, rangeOne[0], num,
			)
		}

		// RangeMul()
		tree.RangeMul(i, []int32{2})
		if tree.RangeAdd(i, []int32{-num}); tree.Number(i) != num {
			t.Errorf(
				"index: %d, rangeNumber got: %d != want: %d\n",
				i, rangeOne[0], num,
			)
		}

		// SearchSum()
		if _, number := tree.SearchSum(sum); number != sum {
			t.Errorf(
				"index: %d, SearchSum got: %d != want: %d\n",
				i, number, sum,
			)
		}

		// AddNumber()
		newtree.AddNumber(i, num)
		if number := newtree.Number(i); number != num {
			t.Errorf(
				"index: %d, number got: %d != want: %d\n",
				i, number, num,
			)
		}

		// MulNumber()
		number := newtree.MulNumber(i, factor)
		if number != factor*num {
			t.Errorf(
				"index: %d, number got: %d != want: %d\n",
				i, number, factor*num,
			)
		}

		// SetNumber()
		newtree.SetNumber(i, num)
		if number := newtree.Number(i); number != num {
			t.Errorf("index: %d, number got: %d != want: %d\n", i, number, num)
		}
	}

	// Numbers()
	tree.Numbers(numbers)
	for i, want := range numbers {
		got := numbers[i]
		if got != want {
			t.Errorf("index: %d, number got: %d != want: %d\n", i, got, want)
		}
	}

	// Append()
	for l := 0; l < len(numbers); l++ {
		partTree := From(numbers[:l])
		fullTree := Append(partTree, numbers[l:]...)
		for i, want := range numbers {
			got := fullTree.Number(i)
			if got != want {
				t.Errorf(
					"length: %d, index: %d, number got: %d != want: %d\n",
					l, i, got, want,
				)
			}
		}
	}

	// Sums()
	mySums := make([]int32, Len(tree))
	tree.Sums(mySums)
	for i := range numbers {
		got := mySums[i]
		want := tree.Sum(i)
		if got != want {
			t.Errorf("index: %d, sum got: %d != want: %d\n", i, got, want)
		}
	}

	// Shift()
	newtree.Shift(shift)
	for i, want := range numbers {
		newtree.AddNumber(i, -shift)
		got := newtree.Number(i)
		if got != want {
			t.Errorf("index: %d, number got: %d != want: %d\n", i, got, want)
		}
	}

	// Scale()
	newtree.Scale(factor)
	for i, want := range numbers {
		got := newtree.Number(i)
		want *= factor
		if got != want {
			t.Errorf("index: %d, number got: %d != want: %d\n", i, got, want)
		}
	}

	// RangeShift()
	Copy(newtree, tree)
	newtree.RangeShift(lo, hi, shift)
	for i, want := range numbers {
		if lo <= i && i < hi {
			newtree.AddNumber(i, -shift)
		}
		got := newtree.Number(i)
		if got != want {
			t.Errorf("index: %d, number got: %d != want: %d\n", i, got, want)
		}
	}

	// RangeScale()
	Copy(newtree, tree)
	newtree.RangeScale(lo, hi, factor)
	for i, want := range numbers {
		if lo <= i && i < hi {
			want *= factor
		}
		got := newtree.Number(i)
		if got != want {
			t.Errorf("index: %d, number got: %d != want: %d\n", i, got, want)
		}
	}
}

func BenchmarkTree(b *testing.B) {
	const (
		n       = 10_000
		max     = 100
		maxDiv2 = max / 2
	)

	in1 := make([]int32, n)
	in2 := make([]int32, n)

	rand.Seed(18)
	for i := range in1 {
		in1[i] = rand.Int31n(max) //- maxDiv2
	}
	copy(in2, in1)

	b.Run("Reference", func(b *testing.B) {
		var tree Tree
		for i := 0; i < b.N; i++ {
			tree = make([]int32, n)
			var sum int32
			for j, number := range in1 {
				sum += number
				tree[j] = sum
			}
		}
		_ = tree
	})

	b.Run("From", func(b *testing.B) {
		var tree Tree
		for i := 0; i < b.N; i++ {
			tree = From(in1)
		}
		_ = tree
	})

	b.Run("From(re-use)", func(b *testing.B) {
		var tree Tree
		for i := 0; i < b.N; i++ {
			tree = From(in2, true)
		}
		_ = tree
	})

	b.Run("Copy", func(b *testing.B) {
		src := From(in1)
		dst := New(Len(src))
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			Copy(dst, src)
		}
	})

	b.Run("Copy+Rebuild", func(b *testing.B) {
		src := New()
		dst := New(len(in1))
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			Copy(dst, src)
		}
	})

	b.Run("Append-oneNumAtATime", func(b *testing.B) {
		tree := New()
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for _, num := range in1 {
				tree = Append(tree, num)
			}
			tree.Reset()
		}
	})

	b.Run("Append-AllNums", func(b *testing.B) {
		tree := New()
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree = Append(tree, in1...)
			tree.Reset()
		}
	})

	b.Run("Append-HalfNums", func(b *testing.B) {
		l := len(in1) / 2
		refTree := From(in1[:l])
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree := New(l)
			Copy(tree, refTree)
			tree = Append(tree, in1[l:]...)
		}
	})

	b.Run("Append-VarNums", func(b *testing.B) {
		refTree := From(in1)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for l := 0; l < len(in1); l++ {
				tree := New(l)
				Copy(tree, refTree[:l])
				tree = Append(tree, in1[l:]...)
			}
		}
	})

	b.Run("SetNumber", func(b *testing.B) {
		tree := New(len(in1))
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for j := range in1 {
				tree.SetNumber(j, -5)
			}
		}
	})

	b.Run("RangeSet", func(b *testing.B) {
		tree := From(in1)
		buf := make([]int32, 100)
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for j := range in1 {
				tree.RangeSet(j, buf)
			}
		}
	})

	b.Run("Sum", func(b *testing.B) {
		tree := From(in1)
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for j := range tree {
				tree.Sum(j)
			}
		}
	})

	b.Run("RangeSum", func(b *testing.B) {
		tree := From(in1)
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for j := range tree {
				tree.RangeSum(j, j+100)
			}
		}
	})

	b.Run("Sums", func(b *testing.B) {
		tree := From(in1)
		sums := make([]int32, Len(tree))
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree.Sums(sums)
		}
	})

	b.Run("Number", func(b *testing.B) {
		tree := From(in1)
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for j := range in1 {
				tree.Number(j)
			}
		}
	})

	b.Run("RangeNumbers", func(b *testing.B) {
		tree := From(in1)
		buf := make([]int32, 100)
		b.ReportAllocs()
		b.ResetTimer()

		var n int
		for i := 0; i < b.N; i++ {
			for j := range in1 {
				n = tree.RangeNumbers(j, buf)
			}
		}
		_ = n
	})

	b.Run("Numbers", func(b *testing.B) {
		tree := From(in1)
		numbers := make([]int32, Len(tree))
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree.Numbers(numbers)
		}
	})

	b.Run("AddNumber", func(b *testing.B) {
		tree := From(in1)
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for j := range in1 {
				tree.AddNumber(j, -5)
			}
		}
	})

	b.Run("RangeAdd", func(b *testing.B) {
		tree := From(in1)
		buf := make([]int32, 100)
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for j := range in1 {
				tree.RangeAdd(j, buf)
			}
		}
	})

	b.Run("RangeShift", func(b *testing.B) {
		tree := From(in1)
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for j := range in1 {
				tree.RangeShift(j, j+100, -5)
			}
		}
	})

	b.Run("Shift", func(b *testing.B) {
		tree := From(in1)
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree.Shift(1)
		}
	})

	b.Run("MulNumber", func(b *testing.B) {
		tree := From(in1)
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for j := range in1 {
				tree.MulNumber(j, -5)
			}
		}
	})

	b.Run("RangeMul", func(b *testing.B) {
		tree := From(in1)
		buf := make([]int32, 100)
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for j := range in1 {
				tree.RangeMul(j, buf)
			}
		}
	})

	b.Run("RangeScale", func(b *testing.B) {
		tree := From(in1)
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for j := range in1 {
				tree.RangeScale(j, j+100, -5)
			}
		}
	})

	b.Run("Scale", func(b *testing.B) {
		tree := From(in1)
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			tree.Scale(1)
		}
	})

	b.Run("SearchSum", func(b *testing.B) {
		tree := From(in1)
		b.ReportAllocs()
		b.ResetTimer()

		var idx int
		var v int32
		for i := 0; i < b.N; i++ {
			idx, v = tree.SearchSum(100)
		}
		_, _ = idx, v
	})
}
