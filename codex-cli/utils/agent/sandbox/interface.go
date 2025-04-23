package sandbox

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Sandbox interface {
	Execute(command []string, workdir string, writableRoots []string, env map[string]string) (string, error)
}

type DefaultSandbox struct{}

func (s *DefaultSandbox) Execute(command []string, workdir string, writableRoots []string, env map[string]string) (string, error) {
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Dir = workdir
	cmd.Env = os.Environ()

	for key, value := range env {
		cmd.Env = append(cmd.Env, key+"="+value)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func IsPathConstrainedToWritablePaths(candidatePath, workdir string, writableRoots []string) bool {
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
