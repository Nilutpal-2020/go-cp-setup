package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[1;31m"
	colorGreen  = "\033[1;32m"
	colorYellow = "\033[1;33m"
	colorBlue   = "\033[1;34m"
	colorPurple = "\033[1;35m"
	colorCyan   = "\033[1;36m"
	colorGray   = "\033[0;90m"
	colorBold   = "\033[1m"
)

const defaultNotionDatabaseID = "19d7320760cc80d79115e483af59b450"

type Config struct {
	NotionToken      string `json:"notion_token"`
	NotionDatabaseID string `json:"notion_database_id"`
}

type ArchiveEntry struct {
	Platform     string    `json:"platform"`
	ProblemName  string    `json:"problem_name"`
	Path         string    `json:"path"`
	Date         string    `json:"date"`
	Tags         []string  `json:"tags"`
	NotionPageID string    `json:"notion_page_id,omitempty"`
	NotionURL    string    `json:"notion_url,omitempty"`
	Timestamp    time.Time `json:"timestamp"`
}

type ProblemMetadata struct {
	ProblemName  string    `json:"problem_name"`
	CreatedAt    string    `json:"created_at"`
	Platform     string    `json:"platform,omitempty"`
	ProblemURL   string    `json:"problem_url,omitempty"`
	Difficulty   string    `json:"difficulty,omitempty"`
	NotionPageID string    `json:"notion_page_id,omitempty"`
	NotionURL    string    `json:"notion_url,omitempty"`
	Tags         []string  `json:"tags,omitempty"`
}

func getWorkspaceRoot() string {
	if env := os.Getenv("CP_WORKSPACE"); env != "" {
		return env
	}
	cwd, err := os.Getwd()
	if err == nil {
		curr := cwd
		for {
			if _, err := os.Stat(filepath.Join(curr, "go.mod")); err == nil {
				return curr
			}
			parent := filepath.Dir(curr)
			if parent == curr {
				break
			}
			curr = parent
		}
	}
	return "/Users/nilutpal/Documents/coding"
}

func getConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	cfgDir := filepath.Join(home, ".config", "cptool")
	_ = os.MkdirAll(cfgDir, 0755)
	return filepath.Join(cfgDir, "config.json")
}

func loadConfig() Config {
	cfg := Config{
		NotionDatabaseID: defaultNotionDatabaseID,
	}

	data, err := os.ReadFile(getConfigPath())
	if err == nil {
		_ = json.Unmarshal(data, &cfg)
	}

	if envToken := os.Getenv("NOTION_API_KEY"); envToken != "" {
		cfg.NotionToken = envToken
	} else if envToken := os.Getenv("NOTION_TOKEN"); envToken != "" {
		cfg.NotionToken = envToken
	}

	if envDb := os.Getenv("NOTION_DATABASE_ID"); envDb != "" {
		cfg.NotionDatabaseID = envDb
	}

	if cfg.NotionDatabaseID == "" {
		cfg.NotionDatabaseID = defaultNotionDatabaseID
	}

	return cfg
}

func saveConfig(cfg Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(getConfigPath(), data, 0644)
}

func getGitRepoURL() string {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	cmd.Dir = getWorkspaceRoot()
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	url := strings.TrimSpace(string(out))
	url = strings.TrimSuffix(url, ".git")
	if strings.HasPrefix(url, "git@github.com:") {
		url = strings.Replace(url, "git@github.com:", "https://github.com/", 1)
	}
	return url
}

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "new", "init":
		cmdNew(args)
	case "test", "run", "t":
		cmdTest(args)
	case "add-test", "add":
		cmdAddTest(args)
	case "backup", "archive", "save":
		cmdBackup(args)
	case "list", "ls":
		cmdList(args)
	case "search", "find":
		cmdSearch(args)
	case "stress":
		cmdStress(args)
	case "config":
		cmdConfig(args)
	case "notion-status", "notion":
		cmdNotionStatus(args)
	case "notion-sync", "sync":
		cmdNotionSync(args)
	case "help", "--help", "-h":
		printHelp()
	default:
		fmt.Printf("%sUnknown command: %s%s\n\n", colorRed, cmd, colorReset)
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Printf(`%s╔══════════════════════════════════════════════════════════╗
║              ⚡ Go Competitive Coding Tool ⚡             ║
╚══════════════════════════════════════════════════════════╝%s

%sUsage:%s
  cptool <command> [arguments]

%sCore Commands:%s
  %snew <name> [platform] [diff] [url/tags]%s
                             Reset & initialize active problem (auto Notion sync with default values)
  %stest%s (or %srun%s)            Compile & test solution against sample test cases
  %sadd-test%s                  Add a new test case to the active problem
  %sbackup <platform> <name> [tags]%s
                             Archive active problem and update Notion status to Solved
  %slist%s                      List all archived problems
  %ssearch <keyword>%s          Search through archived problems by name/tag
  %sstress [iterations]%s       Stress test active/main.go against brute.go with gen.go

%sNotion Integration Commands:%s
  %sconfig --notion-token <KEY>%s
                             Set your Notion API token
  %sconfig --notion-db <ID>%s   Set your Notion Database ID (default: 19d7320760cc80d79115e483af59b450)
  %snotion-status%s             Test connection and display Notion database fields
  %snotion-sync%s               Manually sync active problem to Notion

%sExamples:%s
  cptool new "Two Sum" leetcode easy "https://leetcode.com/problems/two-sum/" "Arrays,Hash Map"
  cptool new 1000A_Watermelon codeforces
  cptool test
  cptool backup leetcode "Two Sum" "Arrays,Hash Map"
`,
		colorCyan, colorReset,
		colorBold, colorReset,
		colorBold, colorReset,
		colorGreen, colorReset,
		colorGreen, colorReset, colorGreen, colorReset,
		colorGreen, colorReset,
		colorGreen, colorReset,
		colorGreen, colorReset,
		colorGreen, colorReset,
		colorGreen, colorReset,
		colorBold, colorReset,
		colorGreen, colorReset,
		colorGreen, colorReset,
		colorGreen, colorReset,
		colorGreen, colorReset,
		colorBold, colorReset,
	)
}

func parseNewArgs(args []string) (name, platform, difficulty, problemURL string, tags []string) {
	difficulty = "Medium"
	platform = "leetcode"

	if len(args) == 0 {
		name = "problem"
		return
	}

	name = args[0]

	for _, arg := range args[1:] {
		lower := strings.ToLower(arg)
		if strings.HasPrefix(arg, "http://") || strings.HasPrefix(arg, "https://") {
			problemURL = arg
		} else if lower == "easy" || lower == "medium" || lower == "hard" {
			difficulty = strings.Title(lower)
		} else if lower == "leetcode" || lower == "codeforces" || lower == "atcoder" || lower == "cses" || lower == "hackerrank" || lower == "misc" {
			platform = lower
		} else {
			// Check if comma-separated tags
			for _, t := range strings.Split(arg, ",") {
				t = strings.TrimSpace(t)
				if t != "" {
					tags = append(tags, t)
				}
			}
		}
	}

	return
}

func cmdNew(args []string) {
	ws := getWorkspaceRoot()
	problemName, platform, difficulty, problemURL, tags := parseNewArgs(args)

	activeDir := filepath.Join(ws, "active")
	testsDir := filepath.Join(activeDir, "tests")

	_ = os.RemoveAll(activeDir)
	_ = os.MkdirAll(testsDir, 0755)

	templatePath := filepath.Join(ws, "templates", "solution.go")
	templateData, err := os.ReadFile(templatePath)
	if err != nil {
		fmt.Printf("%sError reading template: %v%s\n", colorRed, err, colorReset)
		return
	}

	mainPath := filepath.Join(activeDir, "main.go")
	if err := os.WriteFile(mainPath, templateData, 0644); err != nil {
		fmt.Printf("%sError writing main.go: %v%s\n", colorRed, err, colorReset)
		return
	}

	// Create sample test files
	in1Path := filepath.Join(testsDir, "in1.txt")
	out1Path := filepath.Join(testsDir, "out1.txt")
	_ = os.WriteFile(in1Path, []byte("4\n"), 0644)
	_ = os.WriteFile(out1Path, []byte("4\n"), 0644)

	meta := ProblemMetadata{
		ProblemName: problemName,
		CreatedAt:   time.Now().Format(time.RFC3339),
		Platform:    platform,
		Difficulty:  difficulty,
		ProblemURL:  problemURL,
		Tags:        tags,
	}

	fmt.Printf("%s✔ Initialized new problem:%s %s%s%s\n", colorGreen, colorReset, colorBold, problemName, colorReset)
	fmt.Printf("  • Platform:   %s\n", platform)
	fmt.Printf("  • Difficulty: %s\n", difficulty)
	if problemURL != "" {
		fmt.Printf("  • Link:       %s\n", problemURL)
	}
	if len(tags) > 0 {
		fmt.Printf("  • Topics:     %s\n", strings.Join(tags, ", "))
	}
	fmt.Printf("  • Solution:   %s\n", mainPath)
	fmt.Printf("  • Test input: %s\n", in1Path)
	fmt.Printf("  • Test output:%s\n", out1Path)

	// Sync to Notion with all default values populated
	cfg := loadConfig()
	if cfg.NotionToken != "" {
		fmt.Printf("%s⏳ Syncing problem to Notion database with defaults...%s\n", colorGray, colorReset)
		pageID, pageURL, err := createNotionEntryWithDefaults(cfg, meta, "In Progress")
		if err != nil {
			fmt.Printf("%s⚠ Notion sync notice: %v%s\n", colorYellow, err, colorReset)
		} else {
			meta.NotionPageID = pageID
			meta.NotionURL = pageURL
			fmt.Printf("%s✔ Created Notion database entry:%s %s\n", colorGreen, colorReset, pageURL)
		}
	} else {
		fmt.Printf("%s💡 Notion Auto-Sync: Set your token with '%scptool config --notion-token <TOKEN>%s' to auto-create Notion entries.%s\n",
			colorGray, colorCyan, colorGray, colorReset)
	}

	metaData, _ := json.MarshalIndent(meta, "", "  ")
	_ = os.WriteFile(filepath.Join(activeDir, "problem.json"), metaData, 0644)

	fmt.Printf("\n%sRun '%scptool test%s' or '%scprun%s' to compile and test.%s\n", colorGray, colorCyan, colorGray, colorCyan, colorGray, colorReset)
}

func cmdTest(args []string) {
	ws := getWorkspaceRoot()
	activeDir := filepath.Join(ws, "active")
	mainPath := filepath.Join(activeDir, "main.go")

	if _, err := os.Stat(mainPath); os.IsNotExist(err) {
		fmt.Printf("%sError: active/main.go not found. Run 'cptool new <name>' first.%s\n", colorRed, colorReset)
		return
	}

	tmpBinary := filepath.Join(os.TempDir(), fmt.Sprintf("cp_sol_%d", time.Now().UnixNano()))
	defer os.Remove(tmpBinary)

	fmt.Printf("%s⏳ Compiling solution...%s\n", colorGray, colorReset)
	compileStart := time.Now()
	compileCmd := exec.Command("go", "build", "-o", tmpBinary, mainPath)
	var compileStderr bytes.Buffer
	compileCmd.Stderr = &compileStderr

	if err := compileCmd.Run(); err != nil {
		fmt.Printf("%s✖ Compilation Failed:%s\n\n%s\n", colorRed, colorReset, compileStderr.String())
		return
	}
	compileDuration := time.Since(compileStart)
	fmt.Printf("%s✔ Compiled successfully%s %s(%dms)%s\n\n", colorGreen, colorReset, colorGray, compileDuration.Milliseconds(), colorReset)

	testsDir := filepath.Join(activeDir, "tests")
	files, err := os.ReadDir(testsDir)
	if err != nil {
		fmt.Printf("%sNo tests directory found at %s%s\n", colorYellow, testsDir, colorReset)
		return
	}

	var testIndices []int
	re := regexp.MustCompile(`^in(\d+)\.txt$`)
	for _, f := range files {
		matches := re.FindStringSubmatch(f.Name())
		if len(matches) == 2 {
			idx, _ := strconv.Atoi(matches[1])
			testIndices = append(testIndices, idx)
		}
	}
	sort.Ints(testIndices)

	if len(testIndices) == 0 {
		fmt.Printf("%sNo test cases found in %s (e.g. in1.txt)%s\n", colorYellow, testsDir, colorReset)
		return
	}

	passedCount := 0
	totalCount := len(testIndices)

	for _, idx := range testIndices {
		inFile := filepath.Join(testsDir, fmt.Sprintf("in%d.txt", idx))
		outFile := filepath.Join(testsDir, fmt.Sprintf("out%d.txt", idx))

		inData, err := os.ReadFile(inFile)
		if err != nil {
			fmt.Printf("%sError reading %s: %v%s\n", colorRed, inFile, err, colorReset)
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		runCmd := exec.CommandContext(ctx, tmpBinary)
		runCmd.Stdin = bytes.NewReader(inData)
		var stdoutBuf, stderrBuf bytes.Buffer
		runCmd.Stdout = &stdoutBuf
		runCmd.Stderr = &stderrBuf

		runStart := time.Now()
		runErr := runCmd.Run()
		runDuration := time.Since(runStart)
		cancel()

		if ctx.Err() == context.DeadlineExceeded {
			fmt.Printf("%s[TLE] Test #%d%s (Time Limit Exceeded > 3.0s)\n", colorRed, idx, colorReset)
			printInputSnippet(inData)
			continue
		}

		if runErr != nil {
			fmt.Printf("%s[RTE] Test #%d%s (Runtime Error: %v)\n", colorRed, idx, colorReset, runErr)
			if stderrBuf.Len() > 0 {
				fmt.Printf("%sStderr:%s\n%s\n", colorRed, colorReset, stderrBuf.String())
			}
			printInputSnippet(inData)
			continue
		}

		hasExpected := false
		var expectedData []byte
		if exp, err := os.ReadFile(outFile); err == nil {
			hasExpected = true
			expectedData = exp
		}

		actualStr := strings.TrimSpace(stdoutBuf.String())
		if !hasExpected {
			fmt.Printf("%s[RUN] Test #%d%s (%s%dms%s)\n", colorBlue, idx, colorReset, colorGray, runDuration.Milliseconds(), colorReset)
			fmt.Printf("Output:\n%s\n", actualStr)
			if stderrBuf.Len() > 0 {
				fmt.Printf("%sDebug output:%s\n%s\n", colorGray, colorReset, stderrBuf.String())
			}
			passedCount++
			continue
		}

		expectedStr := strings.TrimSpace(string(expectedData))
		normActual := normalizeWhitespace(actualStr)
		normExpected := normalizeWhitespace(expectedStr)

		if normActual == normExpected {
			passedCount++
			fmt.Printf("%s[PASS] Test #%d%s %s(%dms)%s\n", colorGreen, idx, colorReset, colorGray, runDuration.Milliseconds(), colorReset)
			if stderrBuf.Len() > 0 {
				fmt.Printf("%sDebug output:%s\n%s\n", colorGray, colorReset, stderrBuf.String())
			}
		} else {
			fmt.Printf("%s[FAIL] Test #%d%s %s(%dms)%s\n", colorRed, idx, colorReset, colorGray, runDuration.Milliseconds(), colorReset)
			printInputSnippet(inData)
			fmt.Printf("  %sExpected:%s\n%s\n", colorGreen, colorReset, indent(expectedStr, "    "))
			fmt.Printf("  %sGot:%s\n%s\n", colorRed, colorReset, indent(actualStr, "    "))
			if stderrBuf.Len() > 0 {
				fmt.Printf("  %sDebug:%s\n%s\n", colorGray, colorReset, indent(stderrBuf.String(), "    "))
			}
		}
	}

	fmt.Println()
	if passedCount == totalCount {
		fmt.Printf("%s🎉 ALL %d TEST(S) PASSED! 🎉%s\n", colorGreen, totalCount, colorReset)
	} else {
		fmt.Printf("%s⚠ RESULT: %d/%d test(s) passed.%s\n", colorYellow, passedCount, totalCount, colorReset)
	}
}

func printInputSnippet(inData []byte) {
	s := strings.TrimSpace(string(inData))
	lines := strings.Split(s, "\n")
	if len(lines) > 5 {
		lines = append(lines[:5], "... (truncated)")
	}
	fmt.Printf("  %sInput:%s\n%s\n", colorGray, colorReset, indent(strings.Join(lines, "\n"), "    "))
}

func normalizeWhitespace(s string) string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	var cleaned []string
	for _, l := range lines {
		cleaned = append(cleaned, strings.TrimRight(l, " \t\r"))
	}
	return strings.Join(cleaned, "\n")
}

func indent(s, prefix string) string {
	lines := strings.Split(s, "\n")
	for i, l := range lines {
		lines[i] = prefix + l
	}
	return strings.Join(lines, "\n")
}

func cmdAddTest(args []string) {
	ws := getWorkspaceRoot()
	testsDir := filepath.Join(ws, "active", "tests")
	_ = os.MkdirAll(testsDir, 0755)

	files, _ := os.ReadDir(testsDir)
	maxIdx := 0
	re := regexp.MustCompile(`^in(\d+)\.txt$`)
	for _, f := range files {
		matches := re.FindStringSubmatch(f.Name())
		if len(matches) == 2 {
			idx, _ := strconv.Atoi(matches[1])
			if idx > maxIdx {
				maxIdx = idx
			}
		}
	}
	nextIdx := maxIdx + 1

	inPath := filepath.Join(testsDir, fmt.Sprintf("in%d.txt", nextIdx))
	outPath := filepath.Join(testsDir, fmt.Sprintf("out%d.txt", nextIdx))

	if len(args) >= 1 {
		_ = os.WriteFile(inPath, []byte(args[0]+"\n"), 0644)
		if len(args) >= 2 {
			_ = os.WriteFile(outPath, []byte(args[1]+"\n"), 0644)
		} else {
			_ = os.WriteFile(outPath, []byte(""), 0644)
		}
		fmt.Printf("%s✔ Created Test #%d with provided arguments.%s\n", colorGreen, nextIdx, colorReset)
		return
	}

	_ = os.WriteFile(inPath, []byte(""), 0644)
	_ = os.WriteFile(outPath, []byte(""), 0644)
	fmt.Printf("%s✔ Created Test #%d:%s\n", colorGreen, nextIdx, colorReset)
	fmt.Printf("  • Input file:  %s\n", inPath)
	fmt.Printf("  • Output file: %s\n", outPath)
}

func cmdBackup(args []string) {
	if len(args) < 2 {
		fmt.Printf("%sUsage: cptool backup <platform> <problem_name> [tags...]%s\n", colorRed, colorReset)
		fmt.Printf("Example: cptool backup leetcode \"Two Sum\" \"Arrays,Hash Map\"\n")
		return
	}

	platform := strings.ToLower(args[0])
	problemName := args[1]
	var tags []string
	if len(args) > 2 {
		rawTags := strings.Join(args[2:], ",")
		for _, t := range strings.Split(rawTags, ",") {
			t = strings.TrimSpace(t)
			if t != "" {
				tags = append(tags, t)
			}
		}
	}

	ws := getWorkspaceRoot()
	activeDir := filepath.Join(ws, "active")
	mainPath := filepath.Join(activeDir, "main.go")

	if _, err := os.Stat(mainPath); os.IsNotExist(err) {
		fmt.Printf("%sError: active/main.go not found to backup.%s\n", colorRed, colorReset)
		return
	}

	var meta ProblemMetadata
	if metaData, err := os.ReadFile(filepath.Join(activeDir, "problem.json")); err == nil {
		_ = json.Unmarshal(metaData, &meta)
	}
	if len(tags) > 0 {
		meta.Tags = tags
	}

	dateFolder := time.Now().Format("2006-01-02")
	folderSafeName := sanitizeFilename(problemName)
	destDir := filepath.Join(ws, "archive", platform, fmt.Sprintf("%s_%s", dateFolder, folderSafeName))
	_ = os.MkdirAll(destDir, 0755)

	// Copy active files to archive
	copyDir(activeDir, destDir)

	// Generate GitHub Solution Link if repository exists
	var solutionURL string
	gitURL := getGitRepoURL()
	relDest, _ := filepath.Rel(ws, destDir)
	if gitURL != "" {
		solutionURL = fmt.Sprintf("%s/tree/main/%s", gitURL, relDest)
	}

	// Update Notion status
	cfg := loadConfig()
	if cfg.NotionToken != "" {
		if meta.NotionPageID != "" {
			fmt.Printf("%s⏳ Updating Notion entry to Solved with tags and links...%s\n", colorGray, colorReset)
			if err := updateNotionEntryWithDefaults(cfg, meta.NotionPageID, "Solved", meta.Tags, solutionURL); err != nil {
				fmt.Printf("%s⚠ Notion update notice: %v%s\n", colorYellow, err, colorReset)
			} else {
				fmt.Printf("%s✔ Updated Notion page to Solved!%s\n", colorGreen, colorReset)
			}
		} else {
			fmt.Printf("%s⏳ Creating Notion archive entry...%s\n", colorGray, colorReset)
			pageID, pageURL, err := createNotionEntryWithDefaults(cfg, meta, "Solved")
			if err == nil {
				meta.NotionPageID = pageID
				meta.NotionURL = pageURL
				fmt.Printf("%s✔ Created Notion database entry: %s%s\n", colorGreen, pageURL, colorReset)
			}
		}
	}

	// Add entry to archive/index.json
	indexPath := filepath.Join(ws, "archive", "index.json")
	var entries []ArchiveEntry
	if data, err := os.ReadFile(indexPath); err == nil {
		_ = json.Unmarshal(data, &entries)
	}

	newEntry := ArchiveEntry{
		Platform:     platform,
		ProblemName:  problemName,
		Path:         relDest,
		Date:         dateFolder,
		Tags:         meta.Tags,
		NotionPageID: meta.NotionPageID,
		NotionURL:    meta.NotionURL,
		Timestamp:    time.Now(),
	}
	entries = append([]ArchiveEntry{newEntry}, entries...)

	indexData, _ := json.MarshalIndent(entries, "", "  ")
	_ = os.WriteFile(indexPath, indexData, 0644)

	fmt.Printf("%s✔ Problem successfully archived!%s\n", colorGreen, colorReset)
	fmt.Printf("  • Platform: %s%s%s\n", colorBold, platform, colorReset)
	fmt.Printf("  • Problem:  %s%s%s\n", colorBold, problemName, colorReset)
	if len(meta.Tags) > 0 {
		fmt.Printf("  • Topics:   %s%s%s\n", colorCyan, strings.Join(meta.Tags, ", "), colorReset)
	}
	if solutionURL != "" {
		fmt.Printf("  • Solution: %s\n", solutionURL)
	}
	fmt.Printf("  • Location: %s\n", destDir)
}

func sanitizeFilename(name string) string {
	name = strings.ReplaceAll(name, " ", "_")
	reg := regexp.MustCompile(`[^a-zA-Z0-9_\-]+`)
	return reg.ReplaceAllString(name, "")
}

func cmdList(args []string) {
	ws := getWorkspaceRoot()
	indexPath := filepath.Join(ws, "archive", "index.json")

	data, err := os.ReadFile(indexPath)
	if err != nil {
		fmt.Printf("%sNo archived problems found yet.%s\n", colorYellow, colorReset)
		return
	}

	var entries []ArchiveEntry
	if err := json.Unmarshal(data, &entries); err != nil || len(entries) == 0 {
		fmt.Printf("%sNo archived problems found.%s\n", colorYellow, colorReset)
		return
	}

	fmt.Printf("%s📦 Archived Problems (%d total):%s\n\n", colorBold, len(entries), colorReset)
	fmt.Printf("%-12s %-25s %-12s %-25s\n", "PLATFORM", "PROBLEM", "DATE", "TAGS")
	fmt.Printf("%s\n", strings.Repeat("─", 80))

	for _, e := range entries {
		tagStr := strings.Join(e.Tags, ", ")
		if tagStr == "" {
			tagStr = "-"
		}
		fmt.Printf("%-12s %-25s %-12s %-25s\n", e.Platform, e.ProblemName, e.Date, tagStr)
	}
	fmt.Println()
}

func cmdSearch(args []string) {
	if len(args) == 0 {
		fmt.Printf("%sUsage: cptool search <keyword>%s\n", colorRed, colorReset)
		return
	}

	query := strings.ToLower(args[0])
	ws := getWorkspaceRoot()
	indexPath := filepath.Join(ws, "archive", "index.json")

	data, err := os.ReadFile(indexPath)
	if err != nil {
		fmt.Printf("%sNo archived problems found.%s\n", colorYellow, colorReset)
		return
	}

	var entries []ArchiveEntry
	_ = json.Unmarshal(data, &entries)

	var matched []ArchiveEntry
	for _, e := range entries {
		match := strings.Contains(strings.ToLower(e.ProblemName), query) ||
			strings.Contains(strings.ToLower(e.Platform), query) ||
			strings.Contains(strings.ToLower(e.Date), query)

		for _, t := range e.Tags {
			if strings.Contains(strings.ToLower(t), query) {
				match = true
				break
			}
		}

		if match {
			matched = append(matched, e)
		}
	}

	if len(matched) == 0 {
		fmt.Printf("%sNo problems matching '%s'%s\n", colorYellow, query, colorReset)
		return
	}

	fmt.Printf("%s🔍 Found %d matching problem(s) for '%s':%s\n\n", colorGreen, len(matched), query, colorReset)
	for _, e := range matched {
		fmt.Printf("  • %s[%s]%s %s%s%s (%s)\n", colorCyan, e.Platform, colorReset, colorBold, e.ProblemName, colorReset, e.Date)
		if len(e.Tags) > 0 {
			fmt.Printf("    Tags:   %s%s%s\n", colorYellow, strings.Join(e.Tags, ", "), colorReset)
		}
		if e.NotionURL != "" {
			fmt.Printf("    Notion: %s\n", e.NotionURL)
		}
		fmt.Printf("    Path:   %s\n\n", filepath.Join(ws, e.Path))
	}
}

func cmdStress(args []string) {
	ws := getWorkspaceRoot()
	activeDir := filepath.Join(ws, "active")

	mainFile := filepath.Join(activeDir, "main.go")
	bruteFile := filepath.Join(activeDir, "brute.go")
	genFile := filepath.Join(activeDir, "gen.go")

	if _, err := os.Stat(bruteFile); os.IsNotExist(err) {
		fmt.Printf("%sError: active/brute.go not found for stress testing.%s\n", colorRed, colorReset)
		fmt.Printf("Create active/brute.go (brute force solution) and active/gen.go (test generator).\n")
		return
	}
	if _, err := os.Stat(genFile); os.IsNotExist(err) {
		fmt.Printf("%sError: active/gen.go not found for stress testing.%s\n", colorRed, colorReset)
		return
	}

	maxIters := 100
	if len(args) > 0 {
		if val, err := strconv.Atoi(args[0]); err == nil && val > 0 {
			maxIters = val
		}
	}

	tmpDir := os.TempDir()
	binMain := filepath.Join(tmpDir, "stress_main")
	binBrute := filepath.Join(tmpDir, "stress_brute")
	binGen := filepath.Join(tmpDir, "stress_gen")

	defer os.Remove(binMain)
	defer os.Remove(binBrute)
	defer os.Remove(binGen)

	fmt.Printf("%s⏳ Compiling solutions and generator...%s\n", colorGray, colorReset)
	if err := exec.Command("go", "build", "-o", binMain, mainFile).Run(); err != nil {
		fmt.Printf("%sError compiling main.go%s\n", colorRed, colorReset)
		return
	}
	if err := exec.Command("go", "build", "-o", binBrute, bruteFile).Run(); err != nil {
		fmt.Printf("%sError compiling brute.go%s\n", colorRed, colorReset)
		return
	}
	if err := exec.Command("go", "build", "-o", binGen, genFile).Run(); err != nil {
		fmt.Printf("%sError compiling gen.go%s\n", colorRed, colorReset)
		return
	}

	fmt.Printf("%s🚀 Starting stress testing (%d iterations)...%s\n\n", colorBold, maxIters, colorReset)

	for i := 1; i <= maxIters; i++ {
		genOut, err := exec.Command(binGen).Output()
		if err != nil {
			fmt.Printf("%sGenerator failed at iteration %d%s\n", colorRed, i, colorReset)
			return
		}

		cmdM := exec.Command(binMain)
		cmdM.Stdin = bytes.NewReader(genOut)
		outMain, errM := cmdM.Output()

		cmdB := exec.Command(binBrute)
		cmdB.Stdin = bytes.NewReader(genOut)
		outBrute, errB := cmdB.Output()

		if errM != nil || errB != nil {
			fmt.Printf("%s[CRASH] Test #%d failed!%s\n", colorRed, i, colorReset)
			fmt.Printf("Input:\n%s\n", string(genOut))
			return
		}

		sMain := normalizeWhitespace(string(outMain))
		sBrute := normalizeWhitespace(string(outBrute))

		if sMain != sBrute {
			fmt.Printf("%s✖ COUNTEREXAMPLE FOUND on test #%d!%s\n\n", colorRed, i, colorReset)
			fmt.Printf("%sInput:%s\n%s\n\n", colorYellow, colorReset, string(genOut))
			fmt.Printf("%sMain (Optimal) Output:%s\n%s\n\n", colorRed, colorReset, sMain)
			fmt.Printf("%sBrute Force Output:%s\n%s\n\n", colorGreen, colorReset, sBrute)

			_ = os.WriteFile(filepath.Join(activeDir, "tests", "stress_in.txt"), genOut, 0644)
			_ = os.WriteFile(filepath.Join(activeDir, "tests", "stress_expected.txt"), outBrute, 0644)
			fmt.Printf("%sSaved counterexample to active/tests/stress_in.txt%s\n", colorCyan, colorReset)
			return
		}

		if i%10 == 0 || i == maxIters {
			fmt.Printf("  • Passed %d/%d test cases\n", i, maxIters)
		}
	}

	fmt.Printf("\n%s🎉 STRESS TEST PASSED: No discrepancies found in %d iterations!%s\n", colorGreen, maxIters, colorReset)
}

func cmdConfig(args []string) {
	cfg := loadConfig()

	if len(args) == 0 {
		fmt.Printf("%s⚙ Current cptool Configuration:%s\n", colorBold, colorReset)
		tokenMasked := "(not set)"
		if cfg.NotionToken != "" {
			if len(cfg.NotionToken) > 10 {
				tokenMasked = cfg.NotionToken[:6] + "..." + cfg.NotionToken[len(cfg.NotionToken)-4:]
			} else {
				tokenMasked = "***"
			}
		}
		fmt.Printf("  • Notion Database ID: %s%s%s\n", colorCyan, cfg.NotionDatabaseID, colorReset)
		fmt.Printf("  • Notion API Token:   %s%s%s\n", colorCyan, tokenMasked, colorReset)
		fmt.Printf("  • Config File:        %s\n\n", getConfigPath())
		fmt.Printf("Usage to update:\n")
		fmt.Printf("  cptool config --notion-token <TOKEN>\n")
		fmt.Printf("  cptool config --notion-db <DATABASE_ID>\n")
		return
	}

	for i := 0; i < len(args); i++ {
		if (args[i] == "--notion-token" || args[i] == "-token" || args[i] == "-t") && i+1 < len(args) {
			cfg.NotionToken = args[i+1]
			i++
		} else if (args[i] == "--notion-db" || args[i] == "-db") && i+1 < len(args) {
			cfg.NotionDatabaseID = args[i+1]
			i++
		}
	}

	if err := saveConfig(cfg); err != nil {
		fmt.Printf("%sError saving config: %v%s\n", colorRed, err, colorReset)
		return
	}

	fmt.Printf("%s✔ Configuration saved successfully to %s%s\n", colorGreen, getConfigPath(), colorReset)
}

func cmdNotionStatus(args []string) {
	cfg := loadConfig()
	if cfg.NotionToken == "" {
		fmt.Printf("%sNo Notion API Token configured.%s\n", colorYellow, colorReset)
		fmt.Printf("Set one with: %scptool config --notion-token <YOUR_TOKEN>%s\n", colorCyan, colorReset)
		return
	}

	fmt.Printf("%s🔍 Connecting to Notion Database ID: %s...%s\n", colorGray, cfg.NotionDatabaseID, colorReset)
	dbInfo, err := getNotionDatabase(cfg)
	if err != nil {
		fmt.Printf("%s✖ Failed to connect to Notion: %v%s\n", colorRed, err, colorReset)
		return
	}

	fmt.Printf("%s✔ Successfully connected to Notion Database!%s\n\n", colorGreen, colorReset)
	if title, ok := dbInfo["title"].([]interface{}); ok && len(title) > 0 {
		if textObj, ok := title[0].(map[string]interface{}); ok {
			if plain, ok := textObj["plain_text"].(string); ok {
				fmt.Printf("  • Title: %s%s%s\n", colorBold, plain, colorReset)
			}
		}
	}
	if url, ok := dbInfo["url"].(string); ok {
		fmt.Printf("  • URL:   %s\n", url)
	}

	if props, ok := dbInfo["properties"].(map[string]interface{}); ok {
		fmt.Printf("\n%sDetected Database Properties (%d):%s\n", colorBold, len(props), colorReset)
		for name, p := range props {
			if pMap, ok := p.(map[string]interface{}); ok {
				pType, _ := pMap["type"].(string)
				fmt.Printf("  • %s%-26s%s (type: %s)\n", colorCyan, name, colorReset, pType)
			}
		}
	}
	fmt.Println()
}

func cmdNotionSync(args []string) {
	ws := getWorkspaceRoot()
	activeDir := filepath.Join(ws, "active")

	var meta ProblemMetadata
	data, err := os.ReadFile(filepath.Join(activeDir, "problem.json"))
	if err != nil {
		fmt.Printf("%sNo active problem found.%s\n", colorRed, colorReset)
		return
	}
	_ = json.Unmarshal(data, &meta)

	cfg := loadConfig()
	if cfg.NotionToken == "" {
		fmt.Printf("%sNo Notion token configured. Set with 'cptool config --notion-token <TOKEN>'%s\n", colorRed, colorReset)
		return
	}

	fmt.Printf("%s⏳ Syncing %s to Notion with full default properties...%s\n", colorGray, meta.ProblemName, colorReset)
	pageID, pageURL, err := createNotionEntryWithDefaults(cfg, meta, "In Progress")
	if err != nil {
		fmt.Printf("%s✖ Notion sync failed: %v%s\n", colorRed, err, colorReset)
		return
	}

	meta.NotionPageID = pageID
	meta.NotionURL = pageURL
	updatedData, _ := json.MarshalIndent(meta, "", "  ")
	_ = os.WriteFile(filepath.Join(activeDir, "problem.json"), updatedData, 0644)

	fmt.Printf("%s✔ Successfully synced to Notion!%s\n", colorGreen, colorReset)
	fmt.Printf("  • Page URL: %s\n", pageURL)
}

func copyDir(src, dst string) {
	entries, _ := os.ReadDir(src)
	for _, entry := range entries {
		s := filepath.Join(src, entry.Name())
		d := filepath.Join(dst, entry.Name())
		if entry.IsDir() {
			_ = os.MkdirAll(d, 0755)
			copyDir(s, d)
		} else {
			data, err := os.ReadFile(s)
			if err == nil {
				_ = os.WriteFile(d, data, 0644)
			}
		}
	}
}

// ==========================================
// Notion API Client Implementation
// ==========================================

func formatNotionUUID(id string) string {
	clean := strings.ReplaceAll(id, "-", "")
	if len(clean) == 32 {
		return fmt.Sprintf("%s-%s-%s-%s-%s", clean[0:8], clean[8:12], clean[12:16], clean[16:20], clean[20:32])
	}
	return id
}

func getNotionDatabase(cfg Config) (map[string]interface{}, error) {
	dbID := formatNotionUUID(cfg.NotionDatabaseID)
	req, err := http.NewRequest("GET", "https://api.notion.com/v1/databases/"+dbID, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.NotionToken)
	req.Header.Set("Notion-Version", "2022-06-28")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// normalizeTopic matches tag against known Notion multi-select options
func normalizeTopics(inputTags []string, availableOptions []string) []string {
	var res []string
	for _, t := range inputTags {
		tTrim := strings.TrimSpace(t)
		if tTrim == "" {
			continue
		}
		matched := false
		for _, opt := range availableOptions {
			if strings.EqualFold(opt, tTrim) {
				res = append(res, opt)
				matched = true
				break
			}
		}
		if !matched {
			// Title case fallback
			res = append(res, strings.Title(tTrim))
		}
	}
	return res
}

func createNotionEntryWithDefaults(cfg Config, meta ProblemMetadata, status string) (string, string, error) {
	dbID := formatNotionUUID(cfg.NotionDatabaseID)
	dbInfo, err := getNotionDatabase(cfg)
	if err != nil {
		return "", "", err
	}

	rawProps, _ := dbInfo["properties"].(map[string]interface{})
	props := make(map[string]interface{})

	// 1. Problem Title (title)
	for name, v := range rawProps {
		vMap, _ := v.(map[string]interface{})
		pType, _ := vMap["type"].(string)

		if pType == "title" {
			props[name] = map[string]interface{}{
				"title": []map[string]interface{}{
					{
						"text": map[string]string{
							"content": meta.ProblemName,
						},
					},
				},
			}
		}

		// 2. Current Status (select / status) -> 'In Progress' / 'Solved'
		if (strings.EqualFold(name, "Current Status") || strings.EqualFold(name, "Status")) && (pType == "select" || pType == "status") {
			statusVal := "In Progress"
			if status != "" {
				statusVal = status
			}
			props[name] = map[string]interface{}{
				pType: map[string]string{"name": statusVal},
			}
		}

		// 3. Date (date) -> Today
		if (strings.EqualFold(name, "Date") || strings.EqualFold(name, "Created")) && pType == "date" {
			props[name] = map[string]interface{}{
				"date": map[string]string{"start": time.Now().Format("2006-01-02")},
			}
		}

		// 4. Difficulty Level (select) -> Easy / Medium / Hard
		if strings.EqualFold(name, "Difficulty Level") && pType == "select" {
			diffVal := "Medium"
			if meta.Difficulty != "" {
				diffVal = strings.Title(strings.ToLower(meta.Difficulty))
			}
			props[name] = map[string]interface{}{
				"select": map[string]string{"name": diffVal},
			}
		}

		// 5. Submission Attempts (number) -> default 1
		if (strings.EqualFold(name, "Submission Attempts") || strings.EqualFold(name, "Attempts")) && pType == "number" {
			props[name] = map[string]interface{}{
				"number": 1,
			}
		}

		// 6. Review Needed? (checkbox) -> default false
		if strings.EqualFold(name, "Review Needed?") && pType == "checkbox" {
			props[name] = map[string]interface{}{
				"checkbox": false,
			}
		}

		// 7. LeetCode Contest Problem? (checkbox)
		if strings.EqualFold(name, "LeetCode Contest Problem?") && pType == "checkbox" {
			isContest := strings.Contains(strings.ToLower(meta.Platform), "contest")
			props[name] = map[string]interface{}{
				"checkbox": isContest,
			}
		}

		// 8. Problem Link (url)
		if strings.EqualFold(name, "Problem Link") && pType == "url" && meta.ProblemURL != "" {
			props[name] = map[string]interface{}{
				"url": meta.ProblemURL,
			}
		}

		// 9. Topics (multi_select)
		if (strings.EqualFold(name, "Topics") || strings.EqualFold(name, "Tags")) && pType == "multi_select" {
			var availableOptions []string
			if selObj, ok := vMap["multi_select"].(map[string]interface{}); ok {
				if opts, ok := selObj["options"].([]interface{}); ok {
					for _, opt := range opts {
						if optMap, ok := opt.(map[string]interface{}); ok {
							if optName, ok := optMap["name"].(string); ok {
								availableOptions = append(availableOptions, optName)
							}
						}
					}
				}
			}

			normTags := normalizeTopics(meta.Tags, availableOptions)
			var multiSelect []map[string]string
			for _, t := range normTags {
				multiSelect = append(multiSelect, map[string]string{"name": t})
			}
			if len(multiSelect) > 0 {
				props[name] = map[string]interface{}{
					"multi_select": multiSelect,
				}
			}
		}
	}

	reqBody := map[string]interface{}{
		"parent": map[string]string{
			"database_id": dbID,
		},
		"properties": props,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", "", err
	}

	req, err := http.NewRequest("POST", "https://api.notion.com/v1/pages", bytes.NewReader(jsonData))
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.NotionToken)
	req.Header.Set("Notion-Version", "2022-06-28")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var resObj map[string]interface{}
	_ = json.Unmarshal(body, &resObj)

	pageID, _ := resObj["id"].(string)
	pageURL, _ := resObj["url"].(string)

	return pageID, pageURL, nil
}

func updateNotionEntryWithDefaults(cfg Config, pageID, status string, tags []string, solutionURL string) error {
	pageID = formatNotionUUID(pageID)

	reqGet, err := http.NewRequest("GET", "https://api.notion.com/v1/pages/"+pageID, nil)
	if err != nil {
		return err
	}
	reqGet.Header.Set("Authorization", "Bearer "+cfg.NotionToken)
	reqGet.Header.Set("Notion-Version", "2022-06-28")

	client := &http.Client{Timeout: 10 * time.Second}
	respGet, err := client.Do(reqGet)
	if err != nil {
		return err
	}
	defer respGet.Body.Close()

	bodyGet, _ := io.ReadAll(respGet.Body)
	var pageObj map[string]interface{}
	_ = json.Unmarshal(bodyGet, &pageObj)

	rawProps, _ := pageObj["properties"].(map[string]interface{})
	props := make(map[string]interface{})

	for name, v := range rawProps {
		vMap, _ := v.(map[string]interface{})
		pType, _ := vMap["type"].(string)

		if (strings.EqualFold(name, "Current Status") || strings.EqualFold(name, "Status")) && (pType == "select" || pType == "status") {
			props[name] = map[string]interface{}{
				pType: map[string]string{"name": status},
			}
		}

		if (strings.EqualFold(name, "Topics") || strings.EqualFold(name, "Tags")) && pType == "multi_select" && len(tags) > 0 {
			var multiSelect []map[string]string
			for _, t := range tags {
				multiSelect = append(multiSelect, map[string]string{"name": strings.Title(strings.TrimSpace(t))})
			}
			props[name] = map[string]interface{}{
				"multi_select": multiSelect,
			}
		}

		if (strings.EqualFold(name, "Solution Link") || strings.EqualFold(name, "Solution")) && pType == "url" && solutionURL != "" {
			props[name] = map[string]interface{}{
				"url": solutionURL,
			}
		}
	}

	if len(props) == 0 {
		return nil
	}

	reqBody := map[string]interface{}{
		"properties": props,
	}
	jsonData, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("PATCH", "https://api.notion.com/v1/pages/"+pageID, bytes.NewReader(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.NotionToken)
	req.Header.Set("Notion-Version", "2022-06-28")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
