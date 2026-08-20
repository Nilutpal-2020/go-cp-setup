# Competitive Programming in Go - Cheatsheet & Best Practices

## 1. Fast I/O Rules for Go
- **Never use `fmt.Scan` / `fmt.Print`** for inputs/outputs exceeding $10^4$ operations. It is notoriously slow in Go CP.
- Always use `bufio.Reader` and `bufio.Writer` (or the included `FastScanner` in [templates/solution.go](file:///Users/nilutpal/Documents/coding/templates/solution.go)).
- Always remember `defer out.Flush()` when using `bufio.NewWriter`.
- For fast string concatenation, use `strings.Builder` or `[]byte` instead of `+`.

---

## 2. Standard Go CP Snippets Reference

### Sorting
```go
import "sort"

// Primitives
sort.Ints(arr)
sort.Float64s(arr)
sort.Strings(arr)

// Custom slice sort
sort.Slice(arr, func(i, j int) bool {
    return arr[i] < arr[j] // ascending
})

// Sort in descending order
sort.Slice(arr, func(i, j int) bool {
    return arr[i] > arr[j]
})
```

### Binary Search (`sort.Search`)
`sort.Search(n, f)` returns the smallest index `i` in `[0, n)` where `f(i) == true`.
```go
// Lower Bound (first index with arr[i] >= target)
idx := sort.Search(len(arr), func(i int) bool {
    return arr[i] >= target
})

// Upper Bound (first index with arr[i] > target)
idx := sort.Search(len(arr), func(i int) bool {
    return arr[i] > target
})
```

### Priority Queue (Min-Heap / Max-Heap)
```go
import "container/heap"

type IntHeap []int
func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] } // '<' for MinHeap, '>' for MaxHeap
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x any)        { *h = append(*h, x.(int)) }
func (h *IntHeap) Pop() any {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[0 : n-1]
    return x
}

// Usage:
h := &IntHeap{}
heap.Init(h)
heap.Push(h, 5)
minVal := heap.Pop(h).(int)
```

---

## 3. Data Structures & Algorithms available in `templates/snippets.go`
- **Disjoint Set Union (DSU / Union-Find)**: Path compression + Union by rank/size.
- **Fenwick Tree (BIT)**: $O(\log N)$ point updates and prefix range queries.
- **Segment Tree**: $O(\log N)$ point update & range sum/min/max.
- **Combinatorics ($nCr \pmod p$)**: Factorial and inverse factorial precomputation with Fermat's Little Theorem.
- **Sieve of Eratosthenes & Smallest Prime Factor (SPF)**: $O(N)$ sieve, $O(\log X)$ prime factorization.
- **Dijkstra Shortest Path**: Graph shortest path with priority queue.

---

## 4. Common Go Pitfalls in CP
1. **Integer Overflow**: `int` is 64-bit on 64-bit architectures, but when multiplying two large integers ($10^9 \times 10^9 = 10^{18}$), always be mindful of intermediate overflows. Cast to `int64`.
2. **Deep Recursion (Stack Overflow)**: Go goroutine stacks grow dynamically up to 1GB, so recursion limit is rarely an issue compared to Python/C++.
3. **Map Performance**: Go `map[K]V` is fast, but if memory/time is tight, pre-allocate with `make(map[K]V, expectedSize)`.
4. **Slices as Arrays**: Re-slicing does not copy underlying array: `arr[:k]` shares memory.
