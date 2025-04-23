package singlepass

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type CodeDiff struct {
	FileName string
	Diff     string
}

func GenerateCodeDiff(filePath string) (CodeDiff, error) {
	cmd := exec.Command("git", "diff", "--", filePath)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return CodeDiff{}, fmt.Errorf("failed to generate code diff: %w", err)
	}

	return CodeDiff{
		FileName: filePath,
		Diff:     out.String(),
	}, nil
}

func ApplyCodeDiff(diff CodeDiff) error {
	cmd := exec.Command("git", "apply")
	cmd.Stdin = strings.NewReader(diff.Diff)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to apply code diff: %w", err)
	}
	return nil
}
