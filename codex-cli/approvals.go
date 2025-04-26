package main

import (
	"path/filepath"
	"strings"
	"regexp"
	"codex-cli/utils/agent"
)

type SafetyAssessment struct {
	ApplyPatch *ApplyPatchCommand
	Type       string
	Reason     string
	Group      string
	RunInSandbox bool
}

type ApplyPatchCommand struct {
	Patch string
}

type ApprovalPolicy string

const (
	Suggest   ApprovalPolicy = "suggest"
	AutoEdit  ApprovalPolicy = "auto-edit"
	FullAuto  ApprovalPolicy = "full-auto"
)

func CanAutoApprove(command []string, workdir string, policy ApprovalPolicy, writableRoots []string, env map[string]string) SafetyAssessment {
	if command[0] == "apply_patch" {
		if len(command) == 2 && command[1] != "" {
			return canAutoApproveApplyPatch(command[1], workdir, writableRoots, policy)
		}
		return SafetyAssessment{
			Type:   "reject",
			Reason: "Invalid apply_patch command",
		}
	}

	isSafe := isSafeCommand(command)
	if isSafe != nil {
		return SafetyAssessment{
			Type:         "auto-approve",
			Reason:       isSafe.Reason,
			Group:        isSafe.Group,
			RunInSandbox: false,
		}
	}

	if command[0] == "bash" && command[1] == "-lc" && command[2] != "" && len(command) == 3 {
		applyPatchArg := tryParseApplyPatch(command[2])
		if applyPatchArg != "" {
			return canAutoApproveApplyPatch(applyPatchArg, workdir, writableRoots, policy)
		}

		bashCmd := parseShellCommand(command[2], env)
		if bashCmd != nil {
			shellSafe := isEntireShellExpressionSafe(bashCmd)
			if shellSafe != nil {
				return SafetyAssessment{
					Type:         "auto-approve",
					Reason:       shellSafe.Reason,
					Group:        shellSafe.Group,
					RunInSandbox: false,
				}
			}
		}
	}

	if policy == FullAuto {
		return SafetyAssessment{
			Type:         "auto-approve",
			Reason:       "Full auto mode",
			Group:        "Running commands",
			RunInSandbox: true,
		}
	}

	return SafetyAssessment{Type: "ask-user"}
}

func canAutoApproveApplyPatch(applyPatchArg, workdir string, writableRoots []string, policy ApprovalPolicy) SafetyAssessment {
	if policy == Suggest {
		return SafetyAssessment{
			Type:       "ask-user",
			ApplyPatch: &ApplyPatchCommand{Patch: applyPatchArg},
		}
	}

	if isWritePatchConstrainedToWritablePaths(applyPatchArg, workdir, writableRoots) {
		return SafetyAssessment{
			Type:         "auto-approve",
			Reason:       "apply_patch command is constrained to writable paths",
			Group:        "Editing",
			RunInSandbox: false,
			ApplyPatch:   &ApplyPatchCommand{Patch: applyPatchArg},
		}
	}

	if policy == FullAuto {
		return SafetyAssessment{
			Type:         "auto-approve",
			Reason:       "Full auto mode",
			Group:        "Editing",
			RunInSandbox: true,
			ApplyPatch:   &ApplyPatchCommand{Patch: applyPatchArg},
		}
	}

	return SafetyAssessment{
		Type:       "ask-user",
		ApplyPatch: &ApplyPatchCommand{Patch: applyPatchArg},
	}
}

func isWritePatchConstrainedToWritablePaths(applyPatchArg, workdir string, writableRoots []string) bool {
	return allPathsConstrainedToWritablePaths(agent.IdentifyFilesNeeded(applyPatchArg), workdir, writableRoots) &&
		allPathsConstrainedToWritablePaths(agent.IdentifyFilesAdded(applyPatchArg), workdir, writableRoots)
}

func allPathsConstrainedToWritablePaths(candidatePaths []string, workdir string, writableRoots []string) bool {
	for _, candidatePath := range candidatePaths {
		if !isPathConstrainedToWritablePaths(candidatePath, workdir, writableRoots) {
			return false
		}
	}
	return true
}

func isPathConstrainedToWritablePaths(candidatePath, workdir string, writableRoots []string) bool {
	candidateAbsolutePath := resolvePathAgainstWorkdir(candidatePath, workdir)
	for _, writablePath := range writableRoots {
		if pathContains(writablePath, candidateAbsolutePath) {
			return true
		}
	}
	return false
}

func resolvePathAgainstWorkdir(candidatePath, workdir string) string {
	if filepath.IsAbs(candidatePath) {
		return candidatePath
	}
	if workdir != "" {
		return filepath.Join(workdir, candidatePath)
	}
	return filepath.Join(".", candidatePath)
}

func pathContains(parent, child string) bool {
	relative := filepath.Rel(parent, child)
	return relative != "" && !strings.HasPrefix(relative, "..") && !filepath.IsAbs(relative)
}

func tryParseApplyPatch(bashArg string) string {
	const prefix = "apply_patch"
	if !strings.HasPrefix(bashArg, prefix) {
		return ""
	}

	heredoc := strings.TrimPrefix(bashArg, prefix)
	heredocMatch := heredocRegex.FindStringSubmatch(heredoc)
	if len(heredocMatch) > 2 {
		return strings.TrimSpace(heredocMatch[2])
	}
	return strings.TrimSpace(heredoc)
}

var heredocRegex = regexp.MustCompile(`^\s*<<\s*['"]?(\w+)['"]?\n([\s\S]*?)\n\1`)

type SafeCommandReason struct {
	Reason string
	Group  string
}

func isSafeCommand(command []string) *SafeCommandReason {
	switch command[0] {
	case "cd":
		return &SafeCommandReason{Reason: "Change directory", Group: "Navigating"}
	case "ls":
		return &SafeCommandReason{Reason: "List directory", Group: "Searching"}
	case "pwd":
		return &SafeCommandReason{Reason: "Print working directory", Group: "Navigating"}
	case "true":
		return &SafeCommandReason{Reason: "No-op (true)", Group: "Utility"}
	case "echo":
		return &SafeCommandReason{Reason: "Echo string", Group: "Printing"}
	case "cat":
		return &SafeCommandReason{Reason: "View file contents", Group: "Reading files"}
	case "rg":
		return &SafeCommandReason{Reason: "Ripgrep search", Group: "Searching"}
	case "find":
		if !containsUnsafeFindOptions(command) {
			return &SafeCommandReason{Reason: "Find files or directories", Group: "Searching"}
		}
	case "grep":
		return &SafeCommandReason{Reason: "Text search (grep)", Group: "Searching"}
	case "head":
		return &SafeCommandReason{Reason: "Show file head", Group: "Reading files"}
	case "tail":
		return &SafeCommandReason{Reason: "Show file tail", Group: "Reading files"}
	case "wc":
		return &SafeCommandReason{Reason: "Word count", Group: "Reading files"}
	case "which":
		return &SafeCommandReason{Reason: "Locate command", Group: "Searching"}
	case "git":
		return isSafeGitCommand(command)
	case "cargo":
		if command[1] == "check" {
			return &SafeCommandReason{Reason: "Cargo check", Group: "Running command"}
		}
	case "sed":
		if isValidSedNArg(command) {
			return &SafeCommandReason{Reason: "Sed print subset", Group: "Reading files"}
		}
	}
	return nil
}

func containsUnsafeFindOptions(command []string) bool {
	unsafeOptions := []string{"-exec", "-execdir", "-ok", "-okdir", "-delete", "-fls", "-fprint", "-fprint0", "-fprintf"}
	for _, arg := range command {
		for _, unsafeOption := range unsafeOptions {
			if arg == unsafeOption {
				return true
			}
		}
	}
	return false
}

func isSafeGitCommand(command []string) *SafeCommandReason {
	switch command[1] {
	case "status":
		return &SafeCommandReason{Reason: "Git status", Group: "Versioning"}
	case "branch":
		return &SafeCommandReason{Reason: "List Git branches", Group: "Versioning"}
	case "log":
		return &SafeCommandReason{Reason: "Git log", Group: "Using git"}
	case "diff":
		return &SafeCommandReason{Reason: "Git diff", Group: "Using git"}
	case "show":
		return &SafeCommandReason{Reason: "Git show", Group: "Using git"}
	}
	return nil
}

func isValidSedNArg(command []string) bool {
	return len(command) == 4 && command[1] == "-n" && regexp.MustCompile(`^(\d+,)?\d+p$`).MatchString(command[2])
}

func isEntireShellExpressionSafe(parts []string) *SafeCommandReason {
	if len(parts) == 0 {
		return nil
	}

	var currentSegment []string
	var firstReason *SafeCommandReason

	flushSegment := func() bool {
		if len(currentSegment) == 0 {
			return true
		}
		assessment := isSafeCommand(currentSegment)
		if assessment == nil {
			return false
		}
		if firstReason == nil {
			firstReason = assessment
		}
		currentSegment = nil
		return true
	}

	for _, part := range parts {
		if part == "(" || part == ")" || part == "{" || part == "}" {
			return nil
		}
		if isShellOperator(part) {
			if !flushSegment() {
				return nil
			}
		} else {
			currentSegment = append(currentSegment, part)
		}
	}

	if !flushSegment() {
		return nil
	}

	return firstReason
}

func isShellOperator(part string) bool {
	safeShellOperators := []string{"&&", "||", "|", ";"}
	for _, op := range safeShellOperators {
		if part == op {
			return true
		}
	}
	return false
}

func parseShellCommand(command string, env map[string]string) []string {
	// Implement shell command parsing logic here
	return nil
}
