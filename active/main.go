package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

// ==========================================
// Competitive Programming Fast I/O Template
// ==========================================

type FastScanner struct {
	r *bufio.Reader
}

func NewFastScanner(r io.Reader) *FastScanner {
	return &FastScanner{r: bufio.NewReaderSize(r, 1<<20)}
}

// Next reads next non-whitespace string token
func (fs *FastScanner) Next() string {
	b, err := fs.r.ReadByte()
	for err == nil && (b <= ' ' || b > '~') {
		b, err = fs.r.ReadByte()
	}
	if err != nil {
		return ""
	}
	var res []byte
	for err == nil && b > ' ' && b <= '~' {
		res = append(res, b)
		b, err = fs.r.ReadByte()
	}
	return string(res)
}

// NextInt reads next 64-bit integer
func (fs *FastScanner) NextInt() int {
	b, err := fs.r.ReadByte()
	for err == nil && (b <= ' ' || b > '~') {
		b, err = fs.r.ReadByte()
	}
	if err != nil {
		return 0
	}
	sign := 1
	if b == '-' {
		sign = -1
		b, _ = fs.r.ReadByte()
	}
	var res int
	for b >= '0' && b <= '9' {
		res = res*10 + int(b-'0')
		b, err = fs.r.ReadByte()
		if err != nil {
			break
		}
	}
	return res * sign
}

// NextInt64 reads next int64
func (fs *FastScanner) NextInt64() int64 {
	return int64(fs.NextInt())
}

// NextFloat64 reads next float64
func (fs *FastScanner) NextFloat64() float64 {
	s := fs.Next()
	if s == "" {
		return 0
	}
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

// NextIntSlice reads a slice of n integers
func (fs *FastScanner) NextIntSlice(n int) []int {
	res := make([]int, n)
	for i := 0; i < n; i++ {
		res[i] = fs.NextInt()
	}
	return res
}

// ==========================================
// Debug Helpers (Active only when DEBUG=1)
// ==========================================

var isDebug = os.Getenv("DEBUG") == "1"

func debug(args ...interface{}) {
	if isDebug {
		fmt.Fprintln(os.Stderr, args...)
	}
}

func debugf(format string, args ...interface{}) {
	if isDebug {
		fmt.Fprintf(os.Stderr, format, args...)
	}
}

// ==========================================
// Math Utilities
// ==========================================

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int) int {
	if a == 0 || b == 0 {
		return 0
	}
	return (a / gcd(a, b)) * b
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

// ==========================================
// Solution Logic
// ==========================================

func solve(in *FastScanner, out *bufio.Writer) {
	// Read problem input
	n := in.NextInt()
	if n == 0 {
		return
	}
	debug("Solving for n =", n)

	// Example output
	fmt.Fprintln(out, n)
}

func main() {
	in := NewFastScanner(os.Stdin)
	out := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer out.Flush()

	// Set multiTest to true for problems with T test cases
	multiTest := false

	testCases := 1
	if multiTest {
		testCases = in.NextInt()
	}

	for t := 1; t <= testCases; t++ {
		solve(in, out)
	}
}
