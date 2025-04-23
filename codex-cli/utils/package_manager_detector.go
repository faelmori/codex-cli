package utils

import (
	"os/exec"
	"path/filepath"
	"strings"
)

type AgentName string

const (
	Npm  AgentName = "npm"
	Pnpm AgentName = "pnpm"
	Bun  AgentName = "bun"
)

func isInstalled(manager AgentName) bool {
	_, err := exec.LookPath(string(manager))
	return err == nil
}

func getGlobalBinDir(manager AgentName) string {
	if !isInstalled(manager) {
		return ""
	}

	switch manager {
	case Npm:
		output, err := exec.Command("npm", "prefix", "-g").Output()
		if err != nil {
			return ""
		}
		return filepath.Join(strings.TrimSpace(string(output)), "bin")
	case Pnpm:
		output, err := exec.Command("pnpm", "bin", "-g").Output()
		if err != nil {
			return ""
		}
		return strings.TrimSpace(string(output))
	case Bun:
		output, err := exec.Command("bun", "pm", "bin", "-g").Output()
		if err != nil {
			return ""
		}
		return strings.TrimSpace(string(output))
	default:
		return ""
	}
}

func DetectInstallerByPath() AgentName {
	invoked, err := filepath.Abs(os.Args[0])
	if err != nil {
		return ""
	}

	supportedManagers := []AgentName{Npm, Pnpm, Bun}

	for _, manager := range supportedManagers {
		binDir := getGlobalBinDir(manager)
		if binDir != "" && strings.HasPrefix(invoked, binDir) {
			return manager
		}
	}

	return ""
}
