#!/usr/bin/env zsh
# ==============================================================================
# Go Competitive Programming Terminal Aliases & Shortcuts
# Source this file in your ~/.zshrc via:
#   source /Users/nilutpal/Documents/coding/aliases.zsh
# ==============================================================================

export CP_WORKSPACE="/Users/nilutpal/Documents/coding"
export PATH="$CP_WORKSPACE/bin:$PATH"

# Quick Navigation
alias cpcd="cd $CP_WORKSPACE"
alias cpactive="cd $CP_WORKSPACE/active"
alias cparchive="cd $CP_WORKSPACE/archive"

# Solution Execution & Testing
alias cprun="$CP_WORKSPACE/bin/cptool test"
alias cptest="$CP_WORKSPACE/bin/cptool test"
alias cpdebug="DEBUG=1 $CP_WORKSPACE/bin/cptool test"

# Problem Management
alias cpnew="$CP_WORKSPACE/bin/cptool new"
alias cpadd="$CP_WORKSPACE/bin/cptool add-test"
alias cpbackup="$CP_WORKSPACE/bin/cptool backup"
alias cparch="$CP_WORKSPACE/bin/cptool backup"
alias cplist="$CP_WORKSPACE/bin/cptool list"
alias cpsearch="$CP_WORKSPACE/bin/cptool search"
alias cpstress="$CP_WORKSPACE/bin/cptool stress"

# Notion Database Sync
alias cpnotion="$CP_WORKSPACE/bin/cptool notion-status"
alias cpsync="$CP_WORKSPACE/bin/cptool notion-sync"
alias cpconfig="$CP_WORKSPACE/bin/cptool config"

# Open active problem in VS Code / Editor (if available)
cpopen() {
    if command -v code >/dev/null 2>&1; then
        code "$CP_WORKSPACE/active/main.go"
    else
        open "$CP_WORKSPACE/active/main.go"
    fi
}

# Help Command
cphelp() {
    cat << 'EOF'
⚡ Go Competitive Programming Shortcuts ⚡

Navigation:
  cpcd                 Jump to competitive programming workspace
  cpactive             Jump to active problem folder
  cparchive            Jump to archive folder
  cpopen               Open active/main.go in code editor

Problem Lifecycle:
  cpnew <name> [plat]  Create new problem (auto syncs to Notion database)
  cprun / cptest       Compile & run tests against in1.txt, in2.txt...
  cpdebug              Run tests with debug logs enabled (DEBUG=1)
  cpadd [in] [out]     Add new test case
  cpbackup <platform> <name> [tags]
                       Archive active problem & update Notion status to Solved
  cplist               List all archived problems
  cpsearch <keyword>   Search archived solutions
  cpstress [count]     Stress test active/main.go vs brute.go

Notion Integration:
  cpconfig             View or edit Notion configuration
  cpnotion             Test connection to Notion database
  cpsync               Manually sync active problem to Notion

Cheatsheet:
  view $CP_WORKSPACE/CHEATSHEET.md
EOF
}

# Auto installer helper
install_cp_aliases() {
    ZSHRC="$HOME/.zshrc"
    SOURCE_LINE="source $CP_WORKSPACE/aliases.zsh"
    if grep -Fxq "$SOURCE_LINE" "$ZSHRC" 2>/dev/null; then
        echo "✔ CP aliases already present in $ZSHRC"
    else
        echo "" >> "$ZSHRC"
        echo "# Go Competitive Programming Workspace Shortcuts" >> "$ZSHRC"
        echo "$SOURCE_LINE" >> "$ZSHRC"
        echo "✔ Added CP shortcuts to $ZSHRC. Please run 'source ~/.zshrc' or open a new terminal."
    fi
}
