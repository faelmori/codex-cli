package utils

import (
	"os/exec"
)

// GetGitDiff returns the current Git diff for the working directory.
// If the current working directory is not inside a Git repository,
// isGitRepo will be false and diff will be an empty string.
func GetGitDiff() (isGitRepo bool, diff string) {
	// Check if we are inside a git repository.
	if err := exec.Command("git", "rev-parse", "--is-inside-work-tree").Run(); err != nil {
		return false, ""
	}

	// Retrieve the diff including color codes.
	output, err := exec.Command("git", "diff", "--color").Output()
	if err != nil {
		return true, ""
	}

	return true, string(output)
}
