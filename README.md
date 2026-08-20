# ⚡ Go Competitive Programming Setup

A fast, streamlined competitive programming environment in Go with automated multi-case testing, fast I/O templates, stress testing, archiving, and terminal shortcuts.

---

## 🚀 Quick Start & Shortcuts

All commands are available directly from your terminal:

| Shortcut | Description |
| :--- | :--- |
| `cpcd` | Jump to the competitive programming workspace root |
| `cpnew <name>` | Scaffold a new problem with fast I/O template and sample tests |
| `cprun` / `cptest` | Compile and run your solution against all test cases with colored diffs & time tracking |
| `cpdebug` | Run tests with debug mode enabled (`DEBUG=1`) |
| `cpadd [in] [out]` | Add a new custom test case to the active problem |
| `cpbackup <platform> <name> [tags]` | Backup/archive the solution to `archive/<platform>/<date>_<name>` |
| `cplist` | List all archived problems with timestamps & tags |
| `cpsearch <keyword>` | Search archived solutions by name, platform, date, or tag |
| `cpstress [count]` | Stress test `active/main.go` against `active/brute.go` with `active/gen.go` |
| `cpopen` | Open `active/main.go` in your code editor |
| `cphelp` | Display list of shortcuts and usage examples |

---

## 📁 Directory Structure

```text
/Users/nilutpal/Documents/coding/
├── go.mod                      # Go module (cp)
├── bin/
│   └── cptool                  # Pre-compiled high-performance CP CLI runner
├── active/                     # Current problem workspace
│   ├── main.go                 # Active solution file
│   ├── problem.json            # Active problem metadata
│   └── tests/                  # Test input/output files
│       ├── in1.txt, out1.txt
│       └── in2.txt, out2.txt
├── archive/                    # Archived / backed up solutions
│   ├── index.json              # Fast searchable index of all submissions
│   ├── codeforces/
│   ├── leetcode/
│   ├── atcoder/
│   └── misc/
├── templates/
│   ├── solution.go             # Standard Fast I/O starter template
│   ├── interactive.go          # Interactive problem starter template
│   └── snippets.go             # Ready-to-use algorithms (DSU, Fenwick, SegTree, Sieve, Dijkstra)
├── tools/
│   └── cptool/main.go          # CLI source code
├── aliases.zsh                 # Shell aliases (integrated with ~/.zshrc)
├── CHEATSHEET.md               # Go CP algorithms & syntax cheatsheet
└── README.md                   # Workspace guide
```

---

## 💡 Typical Problem Workflow

### 1. Start a New Problem
```bash
cpnew 1000A_Watermelon
```
This resets `active/` and loads the fast I/O template into `active/main.go`.

### 2. Add Test Cases
Paste sample inputs and outputs into `active/tests/in1.txt` and `active/tests/out1.txt`, or run:
```bash
cpadd "8" "YES"
```

### 3. Run and Verify
```bash
cprun
```
You'll see colorized results:
```text
✔ Compiled successfully (12ms)
[PASS] Test #1 (2ms)
[PASS] Test #2 (1ms)

🎉 ALL 2 TEST(S) PASSED! 🎉
```

### 4. Backup & Archive
Once accepted, archive your solution:
```bash
cpbackup codeforces 1000A_Watermelon "math,brute force"
```
The solution and its tests will be copied to `archive/codeforces/YYYY-MM-DD_1000A_Watermelon` and indexed in `archive/index.json`.

### 5. Search Past Solutions
```bash
cpsearch math
```

---

## ⚡ Fast I/O & Debugging in `main.go`

- **FastScanner**: Read tokens and numbers with `in.NextInt()`, `in.Next()`, `in.NextInt64()`, `in.NextFloat64()`, `in.NextIntSlice(n)`.
- **FastWriter**: Buffered output with `fmt.Fprintln(out, ...)` and deferred flush.
- **Debug Logs**: `debug("val =", x)` only prints when run with `cpdebug` (or `DEBUG=1`), ensuring no time penalty on online judges.

---

## 🧠 Cheatsheet & Library

Check [CHEATSHEET.md](file:///Users/nilutpal/Documents/coding/CHEATSHEET.md) and [templates/snippets.go](file:///Users/nilutpal/Documents/coding/templates/snippets.go) for:
- Disjoint Set Union (DSU)
- Binary Indexed Tree / Fenwick Tree
- Segment Tree (Point updates & range queries)
- Modular Inverse, Power, and $nCr \pmod p$ Combinatorics
- Sieve of Eratosthenes & Smallest Prime Factor (SPF)
- Dijkstra Shortest Path with `container/heap`
