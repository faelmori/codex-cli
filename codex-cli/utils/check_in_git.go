package utils

import (
	"os/exec"
	"strings"
)

// CheckInGit checks if the current directory is inside a Git repository.
func CheckInGit() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) == "true"
}
