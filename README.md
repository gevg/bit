# Binary Indexed Tree

<p align="center">
  <a href="https://github.com/gevg/bit/actions"><img src="https://github.com/gevg/bit/workflows/Go/badge.svg?style=flat-square" alt="Build Status" /></a>
  <a href="https://pkg.go.dev/github.com/gevg/bit"><img src="https://img.shields.io/badge/doc-reference-007d9b?logo=go&style=flat-square" alt="go.dev" /></a>
  <a href="https://goreportcard.com/report/github.com/gevg/bit"><img src="https://goreportcard.com/badge/github.com/gevg/bit?style=flat-square" alt="Go Report Card" /></a>
  <a href="https://codecov.io/gh/valyala/fastjson"><img src="https://codecov.io/gh/valyala/fastjson/branch/master/graph/badge.svg" alt="codecov" /></a>
</p>

<p align="center">
  <a href="https://github.com/gevg/bit/actions"><img src="https://github.com/gevg/bit/workflows/Go/badge.svg?style=flat-square" alt="Build Status" /></a>
  <a href="https://pkg.go.dev/github.com/gevg/bit"><img src="https://img.shields.io/badge/doc-reference-007d9b?logo=go&style=flat-square" alt="go.dev" /></a>
  <a href="https://goreportcard.com/report/github.com/gevg/bit"><img src="https://goreportcard.com/badge/github.com/gevg/bit?style=flat-square" alt="Go Report Card" /></a>
  <a href="https://codecov.io/gh/gevg/bit"><img src="https://codecov.io/gh/gevg/bit/branch/master/graph/badge.svg" alt="codecov" /></a>
</p>

[![Build Status](https://github.com/gevg/bit/workflows/Go/badge.svg?style=flat-square)](https://github.com/gevg/bit/actions)
[![go.dev](https://img.shields.io/badge/doc-reference-007d9b?logo=go&style=flat-square)](https://pkg.go.dev/github.com/gevg/bit)
[![Go Report Card](https://goreportcard.com/badge/github.com/gevg/bit?style=flat-square)](https://goreportcard.com/report/github.com/gevg/bit)
[![codecov](https://codecov.io/gh/valyala/fastjson/branch/master/graph/badge.svg)](https://codecov.io/gh/valyala/fastjson)

A Binary Indexed Tree (BIT), also known as Fenwick tree, is a data structure that can efficiently update elements and calculate the sum of a range of consecutive numbers. Is was originally proposed by Boris Ryabko [[1]] in 1989 and Peter Fenwick [[2]] in 1994.

If we want to query the sum of a range of consecutive numbers in an array, it will take O(n) time on average, by adding the numbers up. Alternatively, we can add another array saving the prefix sum, hence the query can be done in O(1) time, but updating the elements becomes an O(n) operation. Therefore, a BIT is preferred when elements mutate frequently.

A Fenwick tree is a simple data structure to solve this problem. Construction takes O(n) time, but a BIT balances both the sum query and element update operations, performing no worse than O(log n) time.

## Features

This library provides a pure `go` implementation with an extensive API. Updates and queries are supported on both numbers and prefix sums. Updates and queries can be performed on a single element, a range of elements or the whole data structure. There is no need to keep hold of the number array, as the numbers can be calculated from the BIT. If required, the Fenwick tree can re-use the number array, avoiding memory allocations during the BIT construction.

It is assumed that the number of elements is known at construction time. Elements can be inserted at construction or when they become available. Adding the number elements at construction is much more efficient though.
(Should we make the Fenwick tree extensible???)

This repository uses zero-based indexing. While this slightly complicates the implementation, it makes the Fenwick tree more natural to use.

The implementation is fast. The code is allocation-free, loops are free from bounds checks, and the algorithms are optimized where possible.

## Installation

Install `bit` with `go get`:

```bash
$ go get -u github.com/gevg/bit
```

## Examples

Construct the tree, initializing it with the number array, and iterate the tree.

```go
// initialize a number array, construct the BIT tree and iterate the prefix sums
numbers := []int32{1, 3, 5, 7, 9}
tree := From(numbers)

for i := range numbers {
    fmt.Printf("tree(%d) = %d\n", i, tree.Sum(i))
}
```

Alternatively, construct a tree of length 5 and subsequently add the numbers to the tree.
`Add()` increases the value at a given index, while `Set()` overwrites the existing value.

```go
// construct a tree of length 5, add/set 5 elements
tree := New(5)
tree.Add(0, 1); tree.Add(1, 3); tree.Add(2, 5); tree.Set(3, 7); tree.Set(4, 9)

// calculate the range sum for the indices 2, 3 AND 4
sums := tree.RangeSum(2, 4)
```

```go
// construct a tree of zero length, append 5 elements
tree := New()
tree.Append(1); tree.Append(3); tree.Append([]int32{5, 7, 9}...)
```

```go
// set the 3rd number to 8, scale the number elements by a factor 5 and print them out.
tree.Set(2, 8)
tree.Scale(5)
fmt.Printf("Numbers: %v\n", tree.Numbers())
```

```go
// Search the largest prefix sum smaller than or equal to 6 and print out its index and value.
fmt.Printf("(i, sum) = (%d, %d)\n", tree.Search(6))
```

## Real Examples

url: https://pkg.go.dev/github.com/gevg/bit .

## To Do

Implementation of the following features:

- [ ] API extensions: Append, RangeSet, RangeAdd, RangeMul, RangeShift, RangeScale
- [ ] Implement examples
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
