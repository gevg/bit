# Binary Indexed Tree

<p align="center">
  <a href="https://github.com/gevg/bit/actions"><img src="https://github.com/gevg/bit/workflows/ci/badge.svg" alt="Build Status" /></a>
  <a href="https://pkg.go.dev/github.com/gevg/bit"><img src="https://pkg.go.dev/badge/github.com/gevg/bit.svg" alt="Go Reference"></a>
  <a href="https://goreportcard.com/report/github.com/gevg/bit"><img src="https://goreportcard.com/badge/github.com/gevg/bit?style=flat-square" alt="Go Report Card" /></a>
  <a href="https://codecov.io/gh/gevg/bit"><img src="https://codecov.io/gh/gevg/bit/branch/master/graph/badge.svg" alt="codecov" /></a>
</p>

A Binary Indexed Tree (BIT), also known as Fenwick tree, is a data structure that can efficiently update elements and calculate the sum of a range of consecutive elements. It was originally proposed by Boris Ryabko [[1]] in 1989 and Peter Fenwick [[2]] in 1994.

If we query the sum of a range of consecutive numbers in an array, it takes O(n) time on average, by adding up the numbers. Alternatively, we can utilize a second array to save the prefix sums, hence the query can be done in O(1) time, but updating the elements becomes an O(n) operation. Therefore, a Fenwick tree is preferred when elements mutate frequently.

A Fenwick tree is a simple data structure to solve this problem. Construction takes O(n) time, but a BIT balances both the sum query and element update operations, performing no worse than O(log n) time.

## Features

This library provides a pure `go` implementation with an extensive API. Updates and queries are supported on both numbers and prefix sums. They can be performed on a single element, a range of elements or the whole tree. There is no need to keep hold of the numbers array once the BIT is constructed, as the numbers can be calculated from the BIT. If required, the Fenwick tree can re-use the numbers array, avoiding memory allocations during BIT construction.

The construction of a Fenwick tree is most efficient when the number of elements is known at construction, using `From(numbers)`. When only the number of elements is known at construction, the tree can be built using `New(n)` and numbers can be added through `SetNumber()` and/or `AddNumber()`. When the number of elements is unknown, the tree can be constructed with `New()` and numbers can be appended to the tree with `Append()`.

This repository uses zero-based indexing. While this slightly complicates the implementation, it makes the Fenwick tree more natural to use.

The implementation is fast. The code is mostly allocation-free, loops are free from bounds checks, and the algorithms are optimized without adding too much complexity.

## Installation

Install `bit` with `go get`:

```bash
$ go get -u github.com/gevg/bit
```

## Usage

Construct the tree, initializing it with a number array, and iterate the tree.

```go
// initialize a number array, construct the BIT tree and iterate the prefix sums
numbers := []int32{1, 3, 5, 7, 9, 11, 13, 15, 17, 19}
tree := From(numbers)

for i := range tree {
    fmt.Printf("tree(%d) = %d\n", i, tree.Sum(i))
}
```

Alternatively, you can construct a tree of a given length and add the numbers.
`AddNumber()` increases the value at a given index, while `SetNumber()` overwrites the existing value.

```go
// Construct a tree of length 10, add/set the elements.
tree := New(10)
tree.AddNumber(0, 1); tree.SetNumber(1, 3); tree.SetNumber(2, 5)
tree.RangeAdd(3, []int32{7, 9, 11}), tree.RangeSet(6, []int32{13, 15, 17, 19})

// Calculate the range sum for the indices 2 and 3.
sums := tree.RangeSum(2, 4)
```

Numbers or a range of numbers in an tree can be multiplied by a factor and the impact on the prefix sums can be viewed.

```go
// Scale the the numbers in the tree with a factor 5. Subsequently, scale the range [2, 8) with a factor 6. Additionally, multiply the number at index 5 with factor -2. We can also multiply each number in a range with a different factor. RangeMul multiplies a range of 3 numbers with different factors, starting at index 2. Subsequently, the prefix sums are shown.
tree.Scale(5)
tree.RangeScale(2, 8, 6)
tree.MulNumber(5, -2)
tree.RangeMul(2, []int32{2, 4, 6})
sums := make([]int32, 6)
n := tree.Sums(sums)
fmt.Printf("%d prefix sums: %v\n", n, sums)
```

Similar operations exist for addition. First, we shift all numbers by 5. Subsequently, shift the range [2, 8) with a value 6, and add -2 to the number at index 5.

```go
// Shift all numbers by 5. Subsequently, shift the numbers in the range [2, 8) with the value 6, and add the value -2 to the number at index 5. Add to the numbers of the range of length 3, starting at index 2, the values 2, 4, and 6 respectively. Copy the range of prefix sums in the sums slice and show it.
tree.Shift(5)
tree.RangeShift(2, 8, 6)
tree.AddNumber(5, -2)
tree.RangeAdd(2, []int32{2, 4, 6})
sums := make([]int32, 6)
n := tree.Sums(2, sums)
fmt.Printf("%d prefix sums: %v\n", n, sums)
```

We can search for a prefix sum smaller than or equal to a given value. The algorithm only works for monotonically increasing prefix sums, as this is the prevalent use case. Other cases can be added when requested (through the issue tracker).

```go
// Search the largest prefix sum smaller than or equal to 6 and print out its index and value.
fmt.Printf("(i, sum) = (%d, %d)\n", tree.SearchSum(6))
```

## To Do

Implementation of the following features:

- [ ] Add examples to documentation.
- [ ] Introduction of parameterized types as soon as they become available in the `go` language.
- [ ] 2D Fenwick tree?
- [ ] Cache-related performance improvements for large arrays at the cost of the zero allocation BIT construction?
- [ ] An alternative implementation of some features in assemby using AVX2 SIMD?

## License

`bit` is available under the [BSD 3-Clause License](LICENSE).

## References

[1]: https://boris.ryabko.net/ryabko1992.pdf
(1) [Boris Ryabko (1992). "A fast on-line adaptive code". IEEE Transactions on Information Theory. 28 (1): 1400–1404.](https://boris.ryabko.net/ryabko1992.pdf)

[2]: https://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.14.8917&rep=rep1&type=pdf
(2) [Fenwick, Peter M. (1994). "A New Data Structure for Cumulative Frequency Tables". Software: Practice and Experience. 24 (3): 327–36.](https://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.14.8917&rep=rep1&type=pdf)

[3]: https://arxiv.org/abs/1904.12370
(3) [Marchini, Stefano and Vigna, Sebastiano. (2020). "Compact Fenwick trees for dynamic ranking and selection". Software: Practice and Experience. 50 (3): 10.1002/spe.2791.](https://arxiv.org/abs/1904.12370)
