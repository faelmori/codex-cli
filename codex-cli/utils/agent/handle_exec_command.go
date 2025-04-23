package agent

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type ExecInput struct {
	Cmd                   []string
	Workdir               string
	TimeoutInMillis       int
	AdditionalWritableRoots []string
}

type ExecResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

const DefaultTimeout = 10000 // 10 seconds

func HandleExecCommand(input ExecInput) ExecResult {
	cmd := exec.Command(input.Cmd[0], input.Cmd[1:]...)
	cmd.Dir = input.Workdir

	// Set environment variables
	cmd.Env = os.Environ()

	// Set additional writable roots
	for _, root := range input.AdditionalWritableRoots {
		cmd.Env = append(cmd.Env, "WRITABLE_ROOT="+root)
	}

	// Set timeout
	timeout := time.Duration(input.TimeoutInMillis) * time.Millisecond
	if timeout == 0 {
		timeout = DefaultTimeout * time.Millisecond
	}

	// Run the command with timeout
	done := make(chan error, 1)
	go func() {
		done <- cmd.Run()
	}()

	select {
	case <-time.After(timeout):
		cmd.Process.Kill()
		return ExecResult{
			Stdout:   "",
			Stderr:   "Command timed out",
			ExitCode: 1,
		}
	case err := <-done:
		stdout, _ := cmd.Output()
		stderr, _ := cmd.CombinedOutput()
		exitCode := 0
		if err != nil {
			exitCode = 1
		}
		return ExecResult{
			Stdout:   string(stdout),
			Stderr:   string(stderr),
			ExitCode: exitCode,
		}
	}
}

func ResolvePathAgainstWorkdir(candidatePath, workdir string) string {
	if filepath.IsAbs(candidatePath) {
		return candidatePath
	}
	if workdir != "" {
		return filepath.Join(workdir, candidatePath)
	}
	return filepath.Join(".", candidatePath)
}

func IsPathConstrainedToWritablePaths(candidatePath, workdir string, writableRoots []string) bool {
	candidateAbsolutePath := ResolvePathAgainstWorkdir(candidatePath, workdir)
	for _, writablePath := range writableRoots {
		if PathContains(writablePath, candidateAbsolutePath) {
			return true
		}
	}
	return false
}

func PathContains(parent, child string) bool {
	relative := filepath.Rel(parent, child)
	return relative != "" && !strings.HasPrefix(relative, "..") && !filepath.IsAbs(relative)
}
