# ⚡ Go Competitive Programming Workspace

A high-performance, all-in-one competitive programming environment in Go with automated multi-case testing, fast I/O templates, stress testing, problem archiving, terminal shortcuts, and **automatic Notion database synchronization**.

---

## 🚀 1-Minute Quick Start

### Step 1: Clone & Setup
```bash
git clone https://github.com/Nilutpal-2020/go-cp-setup.git ~/Documents/coding
cd ~/Documents/coding
```

### Step 2: Build the CLI & Install Aliases
```bash
go build -o bin/cptool ./tools/cptool
chmod +x aliases.zsh bin/cptool
zsh -c "source ./aliases.zsh && install_cp_aliases"
source ~/.zshrc
```

### Step 3: Start Coding!
```bash
cpnew "Two Sum" leetcode easy "https://leetcode.com/problems/two-sum/" "Arrays,Hash Map"
cprun
```

---

## 📖 How to Use the Tool (Step-by-Step)

### 1. Starting a New Problem (`cpnew`)
When you begin a problem on LeetCode, Codeforces, AtCoder, or CSES:

```bash
# Format: cpnew <name> [platform] [difficulty] [problem_url] [topics...]

# Examples:
cpnew "Two Sum" leetcode easy "https://leetcode.com/problems/two-sum/" "Arrays,Hash Map"
cpnew "1000A_Watermelon" codeforces easy
cpnew "Coin Combinations I" cses medium "https://cses.fi/problemset/task/1635" "Dynamic Programming"
```

**What happens automatically:**
1. Resets the `active/` workspace.
2. Loads the optimized **Fast I/O template** into `active/main.go`.
3. Creates sample test input/output files (`active/tests/in1.txt`, `active/tests/out1.txt`).
4. **Syncs to your Notion Database**: Creates a page with `Status: In Progress`, Date, Difficulty, Problem Link, and Topics.

---

### 2. Writing Your Solution (`active/main.go`)

Open `active/main.go` in your favorite editor (or type `cpopen`):

```go
package main

import (
    "bufio"
    "fmt"
    "os"
)

func solve(in *FastScanner, out *bufio.Writer) {
    // 1. Read input using FastScanner
    n := in.NextInt()
    target := in.NextInt()
    nums := in.NextIntSlice(n)

    // 2. Debug logs (only visible with 'cpdebug', ignored during contest submissions)
    debug("Solving for target =", target)

    // 3. Output answer with buffered writer
    fmt.Fprintln(out, ans)
}
```

#### FastScanner Cheat-Sheet:
- `in.NextInt()` : Read next integer
- `in.NextInt64()` : Read next 64-bit integer (`int64`)
- `in.Next()` : Read next string token
- `in.NextFloat64()` : Read next `float64`
- `in.NextIntSlice(n)` : Read a slice of `n` integers

---

### 3. Adding Test Cases (`cpadd`)

You can edit test files directly in `active/tests/in1.txt`, `out1.txt` or add new test cases directly from the terminal:

```bash
# Add a test case with input "4" and expected output "4"
cpadd "4" "4"

# Add an empty test case to fill in later
cpadd
```

---

### 4. Running & Testing Your Code (`cprun` / `cptest`)

To compile and run your solution against all test cases:

```bash
cprun
```

**Example Output:**
```text
⏳ Compiling solution...
✔ Compiled successfully (15ms)

[PASS] Test #1 (2ms)
[PASS] Test #2 (1ms)

🎉 ALL 2 TEST(S) PASSED! 🎉
```

#### Handling Test Outcomes:
- **`[PASS]`**: Output matched expected output exactly.
- **`[FAIL]`**: Output differed. Shows pretty side-by-side diff between **Expected** and **Got**.
- **`[TLE]`**: Execution exceeded timeout (> 3.0s).
- **`[RTE]`**: Program panicked or threw runtime error.
- **Debug Logs**: Run `cpdebug` to view all `debug(...)` logs printed to stderr.

---

### 5. Archiving & Backing Up (`cpbackup`)

Once your solution is accepted on the online judge, archive it to keep your record and search it later:

```bash
# Format: cpbackup <platform> <problem_name> [tags...]

# Example:
cpbackup leetcode "Two Sum" "Arrays,Hash Map"
```

**What happens automatically:**
1. Copies solution and test cases to `archive/<platform>/<date>_<name>/`.
2. Indexes the submission into `archive/index.json`.
3. **Updates Notion Page**: Sets `Status: Solved`, updates tags, and attaches a direct GitHub link to the archived solution!

---

### 6. Searching Past Solutions (`cpsearch` / `cplist`)

Find any previously solved problem or code snippet in seconds:

```bash
# Search by topic, keyword, or platform
cpsearch "hash map"
cpsearch "binary search"
cpsearch codeforces

# List all archived solutions in a table
cplist
```

---

### 7. Stress Testing Edge Cases (`cpstress`)

Catch edge cases and subtle bugs by stress testing your solution against a brute-force validator using random inputs:

1. Create `active/brute.go` (a simple, guaranteed-correct $O(N^2)$ solution).
2. Create `active/gen.go` (generates random inputs).
3. Run:
```bash
# Run 100 random test iterations
cpstress 100
```
If a counterexample is found, `cptool` immediately prints the failing input and saves it to `active/tests/stress_in.txt`.

---

## 📑 Notion Database Integration Setup

The workspace is configured to sync with your **DSA Tracker** Notion table:
[`19d7320760cc80d79115e483af59b450`](https://app.notion.com/p/19d7320760cc80d79115e483af59b450)

### 1-Time Setup (2 minutes):
1. Visit [Notion Integrations](https://www.notion.so/profile/integrations) and click **New integration** (name it "CP Tool").
2. Copy your **Internal Integration Secret** (`ntn_...` or `secret_...`).
3. Open your [Notion Database Page](https://app.notion.com/p/19d7320760cc80d79115e483af59b450), click `...` in top right ➔ **Connections** ➔ **Connect to** ➔ Select **CP Tool**.
4. Configure your token:
   ```bash
   cpconfig --notion-token secret_your_token_here
   ```
5. Test the connection:
   ```bash
   cpnotion
   ```

*(Note: If you work offline or without a token, the tool works 100% locally without interruption).*

---

## ⌨️ Complete Command & Shortcut Reference

| Shortcut | Full Command | Description |
| :--- | :--- | :--- |
| `cpcd` | `cd $CP_WORKSPACE` | Jump to competitive programming workspace |
| `cpactive` | `cd $CP_WORKSPACE/active` | Jump to active problem folder |
| `cparchive` | `cd $CP_WORKSPACE/archive` | Jump to archive folder |
| `cpopen` | `code active/main.go` | Open active solution in editor |
| `cpnew` | `cptool new <args>` | Start new problem (auto Notion sync) |
| `cprun` / `cptest` | `cptool test` | Build & run tests against all sample inputs |
| `cpdebug` | `DEBUG=1 cptool test` | Run tests with debug statements visible |
| `cpadd` | `cptool add-test [in] [out]` | Add a new test case |
| `cpbackup` | `cptool backup <plat> <name>` | Archive solution & update Notion to Solved |
| `cplist` | `cptool list` | List all archived solutions |
| `cpsearch` | `cptool search <query>` | Search archived solutions by keyword/tag |
| `cpstress` | `cptool stress [iters]` | Run automated stress testing |
| `cpnotion` | `cptool notion-status` | Inspect Notion database connection & fields |
| `cpsync` | `cptool notion-sync` | Manually sync active problem to Notion |
| `cpconfig` | `cptool config` | View/set Notion token and database ID |
| `cphelp` | `cphelp` | Show terminal helper cheat-sheet |

---

## 🧠 Algorithms & Templates Library

Pre-tested, copy-paste ready implementations are located in [templates/snippets.go](file:///Users/nilutpal/Documents/coding/templates/snippets.go) and documented in [CHEATSHEET.md](file:///Users/nilutpal/Documents/coding/CHEATSHEET.md):

- **Disjoint Set Union (DSU / Union-Find)**: Path compression & union by rank.
- **Fenwick Tree (BIT)**: $O(\log N)$ point updates & range prefix queries.
- **Segment Tree**: $O(\log N)$ point updates & range queries.
- **Modular Arithmetic & Combinatorics**: $nCr \pmod p$, Modular inverse, Modular exponentiation.
- **Sieve of Eratosthenes & SPF**: $O(N)$ sieve, $O(\log X)$ prime factorization.
- **Dijkstra Shortest Path**: Graph shortest path with `container/heap`.
- **Interactive Problems**: See [templates/interactive.go](file:///Users/nilutpal/Documents/coding/templates/interactive.go) for query auto-flushing.
