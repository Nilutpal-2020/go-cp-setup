# ⚡ Go Competitive Programming Setup

A fast, streamlined competitive programming environment in Go with automated multi-case testing, fast I/O templates, stress testing, archiving, terminal shortcuts, and **automated Notion database sync**.

---

## 🚀 Quick Start & Shortcuts

All commands are available directly from your terminal:

| Shortcut | Description |
| :--- | :--- |
| `cpcd` | Jump to the competitive programming workspace root |
| `cpnew <name> [platform]` | Scaffold a new problem & auto-create entry in your Notion database |
| `cprun` / `cptest` | Compile and run your solution against all test cases with colored diffs & time tracking |
| `cpdebug` | Run tests with debug mode enabled (`DEBUG=1`) |
| `cpadd [in] [out]` | Add a new custom test case to the active problem |
| `cpbackup <platform> <name> [tags]` | Backup solution to `archive/` & update Notion entry to **Solved** |
| `cplist` | List all archived problems with timestamps & tags |
| `cpsearch <keyword>` | Search archived solutions by name, platform, date, or tag |
| `cpstress [count]` | Stress test `active/main.go` against `active/brute.go` with `active/gen.go` |
| `cpnotion` | Test connection to your Notion database and display properties |
| `cpsync` | Manually sync current active problem to Notion |
| `cpconfig` | View or update Notion API token and Database ID |
| `cpopen` | Open `active/main.go` in your code editor |
| `cphelp` | Display list of shortcuts and usage examples |

---

## 📑 Automated Notion Database Integration

Your workspace is configured to sync with Notion Database:
**[`19d7320760cc80d79115e483af59b450`](https://app.notion.com/p/19d7320760cc80d79115e483af59b450)**

### 1-Time Setup (Authorize Notion)

1. Go to [Notion Integrations](https://www.notion.so/profile/integrations) and click **New integration** (e.g. name it "CP Tool").
2. Copy your **Internal Integration Secret** (starts with `ntn_` or `secret_`).
3. Open your [Notion Database Page](https://app.notion.com/p/19d7320760cc80d79115e483af59b450), click the top-right `...` menu -> **Connections** -> **Connect to** -> select your integration ("CP Tool").
4. Configure the token in your terminal:
   ```bash
   cptool config --notion-token secret_your_token_here
   ```
   *(or `export NOTION_API_KEY="secret_your_token_here"` in your shell)*
5. Test the connection:
   ```bash
   cpnotion
   ```

---

## 📁 Directory Structure

```text
/Users/nilutpal/Documents/coding/
├── go.mod                      # Go module (cp)
├── bin/
│   └── cptool                  # Pre-compiled high-performance CP CLI runner
├── active/                     # Current problem workspace
│   ├── main.go                 # Active solution file
│   ├── problem.json            # Problem metadata & Notion page link
│   └── tests/                  # Test input/output files
│       ├── in1.txt, out1.txt
│       └── in2.txt, out2.txt
├── archive/                    # Archived / backed up solutions
│   ├── index.json              # Searchable index of all submissions
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
cpnew 1000A_Watermelon codeforces
```
- Sets up `active/main.go` from Fast I/O template.
- Automatically creates an entry in your Notion database with `Status = "In Progress"`.

### 2. Add Test Cases & Run
```bash
cpadd "8" "YES"
cprun
```

### 3. Backup & Mark as Solved
```bash
cpbackup codeforces 1000A_Watermelon "math,brute force"
```
- Archives files to `archive/codeforces/YYYY-MM-DD_1000A_Watermelon`.
- Automatically updates the Notion database status to **`Solved`** with attached tags!

---

## 🧠 Cheatsheet & Library

Check [CHEATSHEET.md](file:///Users/nilutpal/Documents/coding/CHEATSHEET.md) and [templates/snippets.go](file:///Users/nilutpal/Documents/coding/templates/snippets.go) for pre-implemented algorithms.
