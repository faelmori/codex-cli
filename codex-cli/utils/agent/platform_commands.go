package agent

import (
	"os"
	"runtime"
	"strings"
)

// CommandMap maps Unix commands to their Windows equivalents
var CommandMap = map[string]string{
	"ls":   "dir",
	"grep": "findstr",
	"cat":  "type",
	"rm":   "del",
	"cp":   "copy",
	"mv":   "move",
	"touch": "echo.>",
	"mkdir": "md",
}

// OptionMap maps common Unix command options to their Windows equivalents
var OptionMap = map[string]map[string]string{
	"ls": {
		"-l": "/p",
		"-a": "/a",
		"-R": "/s",
	},
	"grep": {
		"-i": "/i",
		"-r": "/s",
	},
}

// AdaptCommandForPlatform adapts a command for the current platform
func AdaptCommandForPlatform(command []string) []string {
	if runtime.GOOS != "windows" {
		return command
	}

	if len(command) == 0 {
		return command
	}

	cmd := command[0]
	if _, ok := CommandMap[cmd]; !ok {
		return command
	}

	adaptedCommand := make([]string, len(command))
	copy(adaptedCommand, command)
	adaptedCommand[0] = CommandMap[cmd]

	if options, ok := OptionMap[cmd]; ok {
		for i := 1; i < len(adaptedCommand); i++ {
			if option, ok := options[adaptedCommand[i]]; ok {
				adaptedCommand[i] = option
			}
		}
	}

	return adaptedCommand
}

// IsPathConstrainedToWritablePaths checks if a path is constrained to writable paths
func IsPathConstrainedToWritablePaths(candidatePath, workdir string, writableRoots []string) bool {
	candidateAbsolutePath := ResolvePathAgainstWorkdir(candidatePath, workdir)
	for _, writablePath := range writableRoots {
		if PathContains(writablePath, candidateAbsolutePath) {
			return true
		}
	}
	return false
}

// ResolvePathAgainstWorkdir resolves a path against the workdir
func ResolvePathAgainstWorkdir(candidatePath, workdir string) string {
	if filepath.IsAbs(candidatePath) {
		return candidatePath
	}
	if workdir != "" {
		return filepath.Join(workdir, candidatePath)
	}
	return filepath.Join(".", candidatePath)
}

// PathContains checks if a path contains another path
func PathContains(parent, child string) bool {
	relative := filepath.Rel(parent, child)
	return relative != "" && !strings.HasPrefix(relative, "..") && !filepath.IsAbs(relative)
}
