// Copyright 2019 Geert Van Gorp. All rights reserved.
// Use of this source code is governed by the MIT License
// which can be found in the LICENSE file.

package bit

import (
	"math/rand"
	"testing"
)

var testcases = []struct {
	numbers  []int32
	prefixes []int32
	tree     []int32
}{
	{
		[]int32{},
		[]int32{},
		[]int32{0},
	},
	{
		[]int32{-10},
		[]int32{-10},
		[]int32{0, -10},
	},
	{
		[]int32{1, 0, 2, 1, 1, 3, 0, 4, 2, 5, 2, 2, 3, 1, 0, 2},
		[]int32{1, 1, 3, 4, 5, 8, 8, 12, 14, 19, 21, 23, 26, 27, 27, 29},
		[]int32{0, 1, 1, 2, 4, 1, 4, 0, 12, 2, 7, 2, 11, 3, 4, 0, 29},
	},
}

func TestFrom(t *testing.T) {
	for i, tc := range testcases {
		tree := From(tc.numbers)
		for j := range tc.tree {
			if tree[j] != tc.tree[j] {
				t.Errorf("Testcase: %d, index: %d, got: %d != want: %d\n", i, j, tree[j], tc.tree[j])
			}
		}
	}
}

func TestFrom2(t *testing.T) {
	for i, tc := range testcases {
		tree := From2(tc.numbers)
		for j := range tc.tree {
			if tree[j] != tc.tree[j] {
				t.Errorf("Testcase: %d, index: %d, got: %d != want: %d\n", i, j, tree[j], tc.tree[j])
			}
		}
	}
}

func TestNewAndSet(t *testing.T) {
	for i, tc := range testcases {
		tree := New(len(tc.numbers))

		for j, v := range tc.numbers {
			tree.Set(j, v)
		}

		for j := range tc.tree {
			if tree[j] != tc.tree[j] {
				t.Errorf("Testcase: %d, index: %d, got: %d != want: %d\n", i, j, tree[j], tc.tree[j])
			}
		}
	}
}

func TestPrefix(t *testing.T) {
	for i, tc := range testcases {
		tree := From(tc.numbers)
		for j := range tc.prefixes {
			if got, want := tree.Prefix(j), tc.prefixes[j]; got != want {
				t.Errorf("Testcase: %d, index: %d, got: %d != want: %d\n", i, j, got, want)
			}
		}
	}
}

func TestPrefix2(t *testing.T) {
	for i, tc := range testcases {
		tree := From(tc.numbers)
		for j := range tc.prefixes {
			if got, want := tree.Prefix2(j), tc.prefixes[j]; got != want {
				t.Errorf("Testcase: %d, index: %d, got: %d != want: %d\n", i, j, got, want)
			}
		}
	}
}

func TestNumber(t *testing.T) {
	for i, tc := range testcases {
		tree := From(tc.numbers)
		for j := range tc.numbers {
			if got, want := tree.Number(j), tc.numbers[j]; got != want {
				t.Errorf("Testcase: %d, index: %d, got: %d != want: %d\n", i, j, got, want)
			}
		}
	}
}

func TestNumber2(t *testing.T) {
	for i, tc := range testcases {
		tree := From(tc.numbers)
		for j := range tc.numbers {
			if got, want := tree.Number2(j), tc.numbers[j]; got != want {
				t.Errorf("Testcase: %d, index: %d, got: %d != want: %d\n", i, j, got, want)
			}
		}
	}
}

func TestSet(t *testing.T) {
	for i, tc := range testcases {
		tree := From(tc.numbers)

		for j := range tc.numbers {
			tree.Set(j, -5)
		}

		for j := range tc.numbers {
			if got, want := tree.Number(j), int32(-5); got != want {
				t.Errorf("Testcase: %d, index: %d, got: %d != want: %d\n", i, j, got, want)
			}
		}
	}
}

func TestAdd(t *testing.T) {
	for i, tc := range testcases {
		tree := From(tc.numbers)

		for j := range tc.numbers {
			tree.Add(j, tree.Number(j))
		}

		for j := range tc.numbers {
			if got, want := tree.Number(j), 2*tc.numbers[j]; got != want {
				t.Errorf("Testcase: %d, index: %d, got: %d != want: %d\n", i, j, got, want)
			}
		}
	}
}

func TestAdd2(t *testing.T) {
	for i, tc := range testcases {
		tree := From(tc.numbers)

		for j := range tc.numbers {
			tree.Add2(j, tree.Number(j))
		}

		for j := range tc.numbers {
			if got, want := tree.Number2(j), 2*tc.numbers[j]; got != want {
				t.Errorf("Testcase: %d, index: %d, got: %d != want: %d\n", i, j, got, want)
			}
		}
	}
}

func TestMul(t *testing.T) {
	for i, tc := range testcases {
		tree := From(tc.numbers)

		for j := range tc.numbers {
			tree.Mul(j, -5)
		}

		for j := range tc.numbers {
			if got, want := tree.Number2(j), -5*tc.numbers[j]; got != want {
				t.Errorf("Testcase: %d, index: %d, got: %d != want: %d\n", i, j, got, want)
			}
		}
	}
}
func TestScale(t *testing.T) {
	for i, tc := range testcases {
		tree := From(tc.numbers)

		tree.Scale(5)

		for j := range tc.numbers {
			if got, want := tree.Number2(j), 5*tc.numbers[j]; got != want {
				t.Errorf("Testcase: %d, index: %d, got: %d != want: %d\n", i, j, got, want)
			}
		}
	}
}

func BenchmarkTree(b *testing.B) {
	in := make(Tree, 10000)

	rand.Seed(18)
	for i := range in {
		in[i] = rand.Int31n(100)
	}

	b.Run("Cum", func(b *testing.B) { // 10.4 us/op  40960 B/op  1 alloc/op
		var tree Tree
		for i := 0; i < b.N; i++ {
			tree = make(Tree, 10000)
			var tot int32
			for i, number := range in {
				tot += number
				tree[i] = tot
			}
		}
		_ = tree
	})

	b.Run("From", func(b *testing.B) { // 16.5 us/op  40964 B/op  2 allocs/op
		var tree Tree
		for i := 0; i < b.N; i++ {
			tree = From(in)
		}
		_ = tree
	})

	b.Run("From2", func(b *testing.B) { // 16.2 us/op  40964 B/op  2 allocs/op
		var tree Tree
		for i := 0; i < b.N; i++ {
			tree = From2(in)
		}
		_ = tree
	})

	b.Run("NewAndAdd", func(b *testing.B) { // 121 us/op  40960 B/op  1 alloc/op
		for i := 0; i < b.N; i++ {
			tree := New(len(in))
			for i, a := range in {
				tree.Add(i, a)
			}
		}
	})

	b.Run("NewAndSet", func(b *testing.B) { // 194 us/op  40960 B/op  1 alloc/op
		for i := 0; i < b.N; i++ {
			tree := New(len(in))
			for i, a := range in {
				tree.Set(i, a)
			}
		}
	})

	b.Run("Prefix", func(b *testing.B) { // 100 us/op  0 B/op  0 alloc/op
		tree := From(in)
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for j := range tree {
				tree.Prefix(j)
			}
		}
	})

	b.Run("Prefix2", func(b *testing.B) { // 105 us/op  0 B/op  0 alloc/op
		tree := From(in)
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for j := range tree {
				tree.Prefix2(j)
			}
		}
	})

	b.Run("Number", func(b *testing.B) { // 53.4 us/op  0 B/op  0 allocs/op
		tree := From(in)
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for j := range in {
				tree.Number(j)
			}
		}
	})

	b.Run("Number2", func(b *testing.B) { // 53.2 us/op  0 B/op  0 allocs/op
		tree := From(in)
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for j := range in {
				tree.Number2(j)
			}
		}
	})

	b.Run("Set", func(b *testing.B) { // 183 us/op  0 B/op  0 allocs/op
		tree := From(in)
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for j := range in {
				tree.Set(j, -5)
			}
		}
	})

	b.Run("Add", func(b *testing.B) { // 115 us/op  0 B/op  0 allocs/op
		tree := From(in)
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for j := range in {
				tree.Add(j, -5)
			}
		}
	})

	b.Run("Add2", func(b *testing.B) { // 108 us/op  0 B/op  0 allocs/op
		tree := From(in)
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for j := range in {
				tree.Add2(j, -5)
			}
		}
	})

	b.Run("Mul", func(b *testing.B) { // 262 us/op  0 B/op  0 allocs/op
		tree := From(in)
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for j := range in {
				tree.Mul(j, -5)
			}
		}
	})

	b.Run("FromAndScale", func(b *testing.B) { // 23.0 us/op  40964 B/op  2 allocs/op
		for i := 0; i < b.N; i++ {
			tree := From(in)
			tree.Scale(5)
		}
	})
}
