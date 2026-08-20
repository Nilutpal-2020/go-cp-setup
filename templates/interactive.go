package main

import (
	"bufio"
	"fmt"
	"os"
)

// ==============================================
// Interactive Problem Template (Auto-Flushing)
// ==============================================

var scanner = bufio.NewScanner(os.Stdin)

func nextString() string {
	if scanner.Scan() {
		return scanner.Text()
	}
	return ""
}

func nextInt() int {
	scanner.Scan()
	var val int
	fmt.Sscanf(scanner.Text(), "%d", &val)
	return val
}

// query sends a query and immediately flushes stdout
func query(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
	os.Stdout.Sync()
}

// answer submits the final answer and terminates
func answer(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
	os.Stdout.Sync()
}

func solveInteractive() {
	// Example interactive binary search:
	// low, high := 1, 1000000
	// for low <= high {
	//     mid := (low + high) / 2
	//     query("? %d", mid)
	//     res := nextString() // e.g. ">=", "<"
	//     if res == ">=" {
	//         low = mid + 1
	//     } else {
	//         high = mid - 1
	//     }
	// }
	// answer("! %d", high)
}

func main() {
	scanner.Split(bufio.ScanWords)
	solveInteractive()
}
