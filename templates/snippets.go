package main

import (
	"container/heap"
)

// ==========================================
// 1. Disjoint Set Union (DSU / Union-Find)
// ==========================================

type DSU struct {
	parent []int
	size   []int
	count  int // Number of connected components
}

func NewDSU(n int) *DSU {
	parent := make([]int, n)
	size := make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
		size[i] = 1
	}
	return &DSU{parent: parent, size: size, count: n}
}

func (d *DSU) Find(i int) int {
	if d.parent[i] == i {
		return i
	}
	d.parent[i] = d.Find(d.parent[i]) // Path compression
	return d.parent[i]
}

func (d *DSU) Union(i, j int) bool {
	rootI, rootJ := d.Find(i), d.Find(j)
	if rootI == rootJ {
		return false
	}
	// Union by size
	if d.size[rootI] < d.size[rootJ] {
		rootI, rootJ = rootJ, rootI
	}
	d.parent[rootJ] = rootI
	d.size[rootI] += d.size[rootJ]
	d.count--
	return true
}

func (d *DSU) Same(i, j int) bool {
	return d.Find(i) == d.Find(j)
}

func (d *DSU) ComponentSize(i int) int {
	return d.size[d.Find(i)]
}

// ==========================================
// 2. Fenwick Tree (Binary Indexed Tree - BIT)
// ==========================================

type FenwickTree struct {
	n    int
	tree []int64
}

func NewFenwickTree(n int) *FenwickTree {
	return &FenwickTree{n: n, tree: make([]int64, n+1)}
}

// Add val to 1-indexed position i
func (ft *FenwickTree) Add(i int, val int64) {
	for ; i <= ft.n; i += i & -i {
		ft.tree[i] += val
	}
}

// Query prefix sum [1..i] (1-indexed)
func (ft *FenwickTree) Query(i int) int64 {
	var sum int64
	for ; i > 0; i -= i & -i {
		sum += ft.tree[i]
	}
	return sum
}

// QueryRange returns sum in [l, r] (1-indexed)
func (ft *FenwickTree) QueryRange(l, r int) int64 {
	if l > r {
		return 0
	}
	return ft.Query(r) - ft.Query(l-1)
}

// ==========================================
// 3. Segment Tree (Point Update, Range Sum)
// ==========================================

type SegmentTree struct {
	n    int
	tree []int64
}

func NewSegmentTree(arr []int64) *SegmentTree {
	n := len(arr)
	st := &SegmentTree{n: n, tree: make([]int64, 4*n)}
	st.build(arr, 1, 0, n-1)
	return st
}

func (st *SegmentTree) build(arr []int64, node, start, end int) {
	if start == end {
		st.tree[node] = arr[start]
		return
	}
	mid := (start + end) / 2
	st.build(arr, 2*node, start, mid)
	st.build(arr, 2*node+1, mid+1, end)
	st.tree[node] = st.tree[2*node] + st.tree[2*node+1]
}

// Update point idx to val (0-indexed)
func (st *SegmentTree) Update(idx int, val int64) {
	st.update(1, 0, st.n-1, idx, val)
}

func (st *SegmentTree) update(node, start, end, idx int, val int64) {
	if start == end {
		st.tree[node] = val
		return
	}
	mid := (start + end) / 2
	if idx <= mid {
		st.update(2*node, start, mid, idx, val)
	} else {
		st.update(2*node+1, mid+1, end, idx, val)
	}
	st.tree[node] = st.tree[2*node] + st.tree[2*node+1]
}

// QueryRange returns sum in [l, r] (0-indexed)
func (st *SegmentTree) QueryRange(l, r int) int64 {
	return st.query(1, 0, st.n-1, l, r)
}

func (st *SegmentTree) query(node, start, end, l, r int) int64 {
	if r < start || end < l {
		return 0
	}
	if l <= start && end <= r {
		return st.tree[node]
	}
	mid := (start + end) / 2
	return st.query(2*node, start, mid, l, r) + st.query(2*node+1, mid+1, end, l, r)
}

// ==========================================
// 4. Modular Arithmetic & Combinatorics
// ==========================================

const MOD int64 = 1_000_000_007 // or 998_244_353

func modPow(base, exp, mod int64) int64 {
	res := int64(1)
	base %= mod
	for exp > 0 {
		if exp%2 == 1 {
			res = (res * base) % mod
		}
		base = (base * base) % mod
		exp /= 2
	}
	return res
}

func modInverse(n, mod int64) int64 {
	return modPow(n, mod-2, mod) // Fermat's Little Theorem (mod must be prime)
}

type Combinatorics struct {
	fact    []int64
	invFact []int64
	mod     int64
}

func NewCombinatorics(maxN int, mod int64) *Combinatorics {
	fact := make([]int64, maxN+1)
	invFact := make([]int64, maxN+1)
	fact[0] = 1
	invFact[0] = 1
	for i := 1; i <= maxN; i++ {
		fact[i] = (fact[i-1] * int64(i)) % mod
	}
	invFact[maxN] = modInverse(fact[maxN], mod)
	for i := maxN - 1; i >= 1; i-- {
		invFact[i] = (invFact[i+1] * int64(i+1)) % mod
	}
	return &Combinatorics{fact: fact, invFact: invFact, mod: mod}
}

// nCr computes n choose r % mod
func (c *Combinatorics) nCr(n, r int) int64 {
	if r < 0 || r > n {
		return 0
	}
	return c.fact[n] * c.invFact[r] % c.mod * c.invFact[n-r] % c.mod
}

// ==========================================
// 5. Prime Sieve (Eratosthenes & SPF)
// ==========================================

// SieveSPF computes Smallest Prime Factor for O(log N) factorizations
func SieveSPF(n int) (primes []int, spf []int) {
	spf = make([]int, n+1)
	for i := 2; i <= n; i++ {
		if spf[i] == 0 {
			spf[i] = i
			primes = append(primes, i)
		}
		for _, p := range primes {
			if p > spf[i] || i*p > n {
				break
			}
			spf[i*p] = p
		}
	}
	return primes, spf
}

// Factorize returns prime factorization of x in O(log x) using precomputed SPF
func Factorize(x int, spf []int) map[int]int {
	factors := make(map[int]int)
	for x > 1 {
		p := spf[x]
		factors[p]++
		x /= p
	}
	return factors
}

// ==========================================
// 6. Dijkstra Shortest Path
// ==========================================

type Edge struct {
	to, weight int
}

type Item struct {
	node int
	dist int64
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].dist < pq[j].dist }
func (pq PriorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*Item))
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

func Dijkstra(n int, adj [][]Edge, start int) []int64 {
	const INF = int64(1) << 60
	dist := make([]int64, n)
	for i := range dist {
		dist[i] = INF
	}
	dist[start] = 0

	pq := &PriorityQueue{}
	heap.Init(pq)
	heap.Push(pq, &Item{node: start, dist: 0})

	for pq.Len() > 0 {
		curr := heap.Pop(pq).(*Item)
		u, d := curr.node, curr.dist
		if d > dist[u] {
			continue
		}
		for _, e := range adj[u] {
			if dist[u]+int64(e.weight) < dist[e.to] {
				dist[e.to] = dist[u] + int64(e.weight)
				heap.Push(pq, &Item{node: e.to, dist: dist[e.to]})
			}
		}
	}
	return dist
}
