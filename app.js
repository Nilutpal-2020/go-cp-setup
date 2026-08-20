// Algorithm snippets library data for Tab 2
const ALGO_SNIPPETS = [
  {
    title: "⚡ Fast I/O Scanner & Writer",
    category: "I/O",
    code: `type FastScanner struct {
    r *bufio.Reader
}
func NewFastScanner(r io.Reader) *FastScanner {
    return &FastScanner{r: bufio.NewReaderSize(r, 1<<20)}
}
func (fs *FastScanner) NextInt() int {
    b, _ := fs.r.ReadByte()
    for b <= ' ' || b > '~' { b, _ = fs.r.ReadByte() }
    sign := 1
    if b == '-' { sign = -1; b, _ = fs.r.ReadByte() }
    var res int
    for b >= '0' && b <= '9' {
        res = res*10 + int(b-'0')
        b, _ = fs.r.ReadByte()
    }
    return res * sign
}`
  },
  {
    title: "🔗 Disjoint Set Union (DSU / Union-Find)",
    category: "Data Structures",
    code: `type DSU struct {
    parent, size []int
    count        int
}
func NewDSU(n int) *DSU {
    p, s := make([]int, n), make([]int, n)
    for i := range p { p[i] = i; s[i] = 1 }
    return &DSU{parent: p, size: s, count: n}
}
func (d *DSU) Find(i int) int {
    if d.parent[i] == i { return i }
    d.parent[i] = d.Find(d.parent[i])
    return d.parent[i]
}
func (d *DSU) Union(i, j int) bool {
    rI, rJ := d.Find(i), d.Find(j)
    if rI == rJ { return false }
    if d.size[rI] < d.size[rJ] { rI, rJ = rJ, rI }
    d.parent[rJ] = rI
    d.size[rI] += d.size[rJ]
    d.count--
    return true
}`
  },
  {
    title: "🌲 Fenwick Tree (Binary Indexed Tree)",
    category: "Data Structures",
    code: `type FenwickTree struct {
    n    int
    tree []int64
}
func NewFenwickTree(n int) *FenwickTree {
    return &FenwickTree{n: n, tree: make([]int64, n+1)}
}
func (ft *FenwickTree) Add(i int, val int64) {
    for ; i <= ft.n; i += i & -i { ft.tree[i] += val }
}
func (ft *FenwickTree) Query(i int) int64 {
    var sum int64
    for ; i > 0; i -= i & -i { sum += ft.tree[i] }
    return sum
}
func (ft *FenwickTree) QueryRange(l, r int) int64 {
    if l > r { return 0 }
    return ft.Query(r) - ft.Query(l-1)
}`
  },
  {
    title: "🔢 Modular Arithmetic & Combinatorics (nCr)",
    category: "Math",
    code: `const MOD int64 = 1_000_000_007

func modPow(base, exp, mod int64) int64 {
    res := int64(1)
    base %= mod
    for exp > 0 {
        if exp%2 == 1 { res = (res * base) % mod }
        base = (base * base) % mod
        exp /= 2
    }
    return res
}
func modInverse(n, mod int64) int64 {
    return modPow(n, mod-2, mod)
}`
  },
  {
    title: "🔍 Sieve of Eratosthenes & SPF",
    category: "Math",
    code: `func SieveSPF(n int) (primes []int, spf []int) {
    spf = make([]int, n+1)
    for i := 2; i <= n; i++ {
        if spf[i] == 0 {
            spf[i] = i
            primes = append(primes, i)
        }
        for _, p := range primes {
            if p > spf[i] || i*p > n { break }
            spf[i*p] = p
        }
    }
    return primes, spf
}`
  },
  {
    title: "🧭 Dijkstra Shortest Path",
    category: "Graphs",
    code: `type Edge struct { to, weight int }
type Item struct { node int; dist int64 }
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].dist < pq[j].dist }
func (pq PriorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x any)        { *pq = append(*pq, x.(*Item)) }
func (pq *PriorityQueue) Pop() any {
    old := *pq; n := len(old); item := old[n-1]; *pq = old[0 : n-1]; return item
}`
  }
];

// Fallback data if fetched on local file protocol without server
const FALLBACK_DATA = [
  {
    platform: "leetcode",
    problem_name: "Two Sum",
    path: "archive/leetcode/2026-08-20_Two_Sum",
    date: "2026-08-20",
    tags: ["Arrays", "Hash Map"],
    notion_url: "https://app.notion.com/p/19d7320760cc80d79115e483af59b450"
  },
  {
    platform: "codeforces",
    problem_name: "starter_problem",
    path: "archive/codeforces/2026-08-20_starter_problem",
    date: "2026-08-20",
    tags: ["starter", "fast-io"],
    notion_url: ""
  }
];

let allProblems = [];
let currentFilter = "all";
let searchQuery = "";

document.addEventListener("DOMContentLoaded", () => {
  initTabs();
  initFilters();
  initSearch();
  initSnippets();
  initModal();
  loadProblemData();
});

async function loadProblemData() {
  try {
    const res = await fetch("archive/index.json");
    if (res.ok) {
      allProblems = await res.json();
    } else {
      allProblems = FALLBACK_DATA;
    }
  } catch (err) {
    console.warn("Using fallback local data:", err);
    allProblems = FALLBACK_DATA;
  }
  updateStats();
  renderProblems();
}

function updateStats() {
  const total = allProblems.length;
  const leetcode = allProblems.filter(p => p.platform?.toLowerCase() === "leetcode").length;
  const cf = allProblems.filter(p => p.platform?.toLowerCase() !== "leetcode").length;

  const topicsSet = new Set();
  allProblems.forEach(p => {
    if (Array.isArray(p.tags)) {
      p.tags.forEach(t => topicsSet.add(t.toLowerCase()));
    }
  });

  document.getElementById("stat-total-count").textContent = total;
  document.getElementById("stat-leetcode-count").textContent = leetcode;
  document.getElementById("stat-cf-count").textContent = cf;
  document.getElementById("stat-topics-count").textContent = topicsSet.size;
}

function renderProblems() {
  const container = document.getElementById("problems-container");
  container.innerHTML = "";

  const filtered = allProblems.filter(p => {
    const matchesPlatform = currentFilter === "all" || p.platform?.toLowerCase() === currentFilter.toLowerCase();
    
    const query = searchQuery.toLowerCase();
    const matchesQuery = !query || 
      p.problem_name?.toLowerCase().includes(query) ||
      p.platform?.toLowerCase().includes(query) ||
      p.date?.includes(query) ||
      (Array.isArray(p.tags) && p.tags.some(t => t.toLowerCase().includes(query)));

    return matchesPlatform && matchesQuery;
  });

  if (filtered.length === 0) {
    container.innerHTML = `
      <div class="empty-state">
        <h3>No problems found</h3>
        <p>Try clearing your search filters or start a new problem with <code>cpnew &lt;name&gt;</code>.</p>
      </div>
    `;
    return;
  }

  filtered.forEach(p => {
    const card = document.createElement("div");
    card.className = "problem-card";

    const platform = p.platform ? p.platform.toLowerCase() : "misc";
    const badgeClass = `badge-${platform}`;

    const tagsHtml = (p.tags || []).map(t => `<span class="tag">${escapeHtml(t)}</span>`).join("");

    const notionBtn = p.notion_url ? `
      <a href="${p.notion_url}" target="_blank" rel="noopener" class="action-btn notion-btn" title="Open in Notion">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"></path><polyline points="15 3 21 3 21 9"></polyline><line x1="10" y1="14" x2="21" y2="3"></line></svg>
        Notion
      </a>
    ` : "";

    const gitHubLink = `https://github.com/Nilutpal-2020/go-cp-setup/tree/main/${p.path || ''}`;

    card.innerHTML = `
      <div>
        <div class="card-header">
          <h3 class="problem-title">${escapeHtml(p.problem_name || "Problem")}</h3>
          <span class="platform-badge ${badgeClass}">${escapeHtml(p.platform || "misc")}</span>
        </div>
        <div class="tag-list">
          ${tagsHtml || '<span class="tag">Solution</span>'}
        </div>
      </div>
      <div class="card-footer">
        <span class="date-text">📅 ${p.date || 'Recent'}</span>
        <div class="card-actions">
          ${notionBtn}
          <a href="${gitHubLink}" target="_blank" rel="noopener" class="action-btn" title="View in GitHub">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor"><path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0024 12c0-6.63-5.37-12-12-12z"/></svg>
            Code
          </a>
        </div>
      </div>
    `;

    container.appendChild(card);
  });
}

function initTabs() {
  const tabBtns = document.querySelectorAll(".tab-btn");
  tabBtns.forEach(btn => {
    btn.addEventListener("click", () => {
      tabBtns.forEach(b => b.classList.remove("active"));
      btn.classList.add("active");

      const tabTarget = btn.getAttribute("data-tab");
      document.querySelectorAll(".tab-content").forEach(tc => {
        tc.style.display = tc.id === `tab-${tabTarget}` ? "block" : "none";
      });
    });
  });
}

function initFilters() {
  const filterBtns = document.querySelectorAll(".filter-btn");
  filterBtns.forEach(btn => {
    btn.addEventListener("click", () => {
      filterBtns.forEach(b => b.classList.remove("active"));
      btn.classList.add("active");
      currentFilter = btn.getAttribute("data-filter");
      renderProblems();
    });
  });
}

function initSearch() {
  const searchInput = document.getElementById("search-input");
  searchInput.addEventListener("input", (e) => {
    searchQuery = e.target.value;
    renderProblems();
  });
}

function initSnippets() {
  const container = document.getElementById("snippets-grid");
  container.innerHTML = "";

  ALGO_SNIPPETS.forEach(snippet => {
    const card = document.createElement("div");
    card.className = "snippet-card";
    card.innerHTML = `
      <div class="snippet-header">
        <h3>${escapeHtml(snippet.title)}</h3>
        <button class="copy-btn" onclick="copySnippetCode(this)">Copy</button>
      </div>
      <pre class="code-block"><code>${escapeHtml(snippet.code)}</code></pre>
    `;
    container.appendChild(card);
  });
}

window.copySnippetCode = function(button) {
  const codeBlock = button.closest(".snippet-card").querySelector("code");
  navigator.clipboard.writeText(codeBlock.textContent).then(() => {
    const orig = button.textContent;
    button.textContent = "Copied!";
    button.style.color = "var(--accent-green)";
    setTimeout(() => {
      button.textContent = orig;
      button.style.color = "";
    }, 1800);
  });
};

function initModal() {
  const modal = document.getElementById("code-modal");
  const closeBtn = document.getElementById("close-modal-btn");
  closeBtn.addEventListener("click", () => modal.classList.remove("active"));
  modal.addEventListener("click", (e) => {
    if (e.target === modal) modal.classList.remove("active");
  });
}

function escapeHtml(str) {
  if (!str) return "";
  return str.replace(/[&<>'"]/g, 
    tag => ({
      '&': '&amp;',
      '<': '&lt;',
      '>': '&gt;',
      "'": '&#39;',
      '"': '&quot;'
    }[tag] || tag)
  );
}
