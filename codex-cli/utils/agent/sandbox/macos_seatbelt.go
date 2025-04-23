package sandbox

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type MacOSSeatbeltSandbox struct{}

func (s *MacOSSeatbeltSandbox) Execute(command []string, workdir string, writableRoots []string, env map[string]string) (string, error) {
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Dir = workdir
	cmd.Env = os.Environ()

	for key, value := range env {
		cmd.Env = append(cmd.Env, key+"="+value)
	}

	// Apply macOS seatbelt sandbox profile
	seatbeltProfile := generateSeatbeltProfile(writableRoots)
	cmd.Env = append(cmd.Env, "SEATBELT_PROFILE="+seatbeltProfile)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func generateSeatbeltProfile(writableRoots []string) string {
	profile := "(version 1)\n(deny default)\n"

	for _, root := range writableRoots {
		profile += "(allow file-write* (subpath \"" + root + "\"))\n"
	}

	return profile
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
