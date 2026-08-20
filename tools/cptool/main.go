package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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

type ArchiveEntry struct {
	Platform    string    `json:"platform"`
	ProblemName string    `json:"problem_name"`
	Path        string    `json:"path"`
	Date        string    `json:"date"`
	Tags        []string  `json:"tags"`
	Timestamp   time.Time `json:"timestamp"`
}

func getWorkspaceRoot() string {
	// Try environment variable or traverse up from current dir or binary
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
	// Fallback to default
	return "/Users/nilutpal/Documents/coding"
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

%sCommands:%s
  %snew [name]%s                 Reset and setup active problem from template
  %stest%s (or %srun%s)            Compile & test solution against sample test cases
  %sadd-test%s                  Add a new test case to the active problem
  %sbackup <platform> <name> [tags]%s
                             Archive active problem to archive/<platform>/<name>
  %slist%s                      List all archived problems
  %ssearch <keyword>%s          Search through archived problems by name/tag
  %sstress [iterations]%s       Stress test active/main.go against brute.go with gen.go
  %shelp%s                      Show this help screen

%sExamples:%s
  cptool new 1000A_Watermelon
  cptool test
  cptool add-test
  cptool backup codeforces 1000A_Watermelon "math,brute force"
  cptool search math
  cptool stress 100
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
		colorGreen, colorReset,
		colorBold, colorReset,
	)
}

func cmdNew(args []string) {
	ws := getWorkspaceRoot()
	problemName := "problem"
	if len(args) > 0 {
		problemName = args[0]
	}

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

	// Write active/main.go
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

	// Save problem metadata
	info := map[string]interface{}{
		"problem_name": problemName,
		"created_at":   time.Now().Format(time.RFC3339),
	}
	infoData, _ := json.MarshalIndent(info, "", "  ")
	_ = os.WriteFile(filepath.Join(activeDir, "problem.json"), infoData, 0644)

	fmt.Printf("%s✔ Initialized new problem:%s %s%s%s\n", colorGreen, colorReset, colorBold, problemName, colorReset)
	fmt.Printf("  • Solution: %s\n", mainPath)
	fmt.Printf("  • Test input: %s\n", in1Path)
	fmt.Printf("  • Test output: %s\n", out1Path)
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

	// Compile
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

	// Discover test cases
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
		// Normalize line endings
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

	// Find next test index
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

	// Create empty template test files
	_ = os.WriteFile(inPath, []byte(""), 0644)
	_ = os.WriteFile(outPath, []byte(""), 0644)
	fmt.Printf("%s✔ Created Test #%d:%s\n", colorGreen, nextIdx, colorReset)
	fmt.Printf("  • Input file:  %s\n", inPath)
	fmt.Printf("  • Output file: %s\n", outPath)
}

func cmdBackup(args []string) {
	if len(args) < 2 {
		fmt.Printf("%sUsage: cptool backup <platform> <problem_name> [tags...]%s\n", colorRed, colorReset)
		fmt.Printf("Example: cptool backup codeforces 1000A_Watermelon \"math,greedy\"\n")
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

	dateFolder := time.Now().Format("2006-01-02")
	destDir := filepath.Join(ws, "archive", platform, fmt.Sprintf("%s_%s", dateFolder, problemName))
	_ = os.MkdirAll(destDir, 0755)

	// Copy all files from active to archive
	copyDir(activeDir, destDir)

	// Add entry to archive/index.json
	indexPath := filepath.Join(ws, "archive", "index.json")
	var entries []ArchiveEntry
	if data, err := os.ReadFile(indexPath); err == nil {
		_ = json.Unmarshal(data, &entries)
	}

	relPath, _ := filepath.Rel(ws, destDir)
	newEntry := ArchiveEntry{
		Platform:    platform,
		ProblemName: problemName,
		Path:        relPath,
		Date:        dateFolder,
		Tags:        tags,
		Timestamp:   time.Now(),
	}
	entries = append([]ArchiveEntry{newEntry}, entries...)

	indexData, _ := json.MarshalIndent(entries, "", "  ")
	_ = os.WriteFile(indexPath, indexData, 0644)

	fmt.Printf("%s✔ Problem successfully archived!%s\n", colorGreen, colorReset)
	fmt.Printf("  • Platform: %s%s%s\n", colorBold, platform, colorReset)
	fmt.Printf("  • Problem:  %s%s%s\n", colorBold, problemName, colorReset)
	if len(tags) > 0 {
		fmt.Printf("  • Tags:     %s%s%s\n", colorCyan, strings.Join(tags, ", "), colorReset)
	}
	fmt.Printf("  • Location: %s\n", destDir)
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
			fmt.Printf("    Tags: %s%s%s\n", colorYellow, strings.Join(e.Tags, ", "), colorReset)
		}
		fmt.Printf("    Path: %s\n\n", filepath.Join(ws, e.Path))
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
		// Generate random test
		genOut, err := exec.Command(binGen).Output()
		if err != nil {
			fmt.Printf("%sGenerator failed at iteration %d%s\n", colorRed, i, colorReset)
			return
		}

		// Run Main
		cmdM := exec.Command(binMain)
		cmdM.Stdin = bytes.NewReader(genOut)
		outMain, errM := cmdM.Output()

		// Run Brute
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

			// Save counterexample into active/tests/
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

