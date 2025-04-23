package utils

import (
	"os"
	"path/filepath"
	"strings"
)

func ShortenPath(p string, maxLength int) string {
	home, _ := os.UserHomeDir()
	displayPath := strings.Replace(p, home, "~", 1)
	if len(displayPath) <= maxLength {
		return displayPath
	}

	parts := strings.Split(displayPath, string(filepath.Separator))
	for i := range parts {
		candidate := filepath.Join(append([]string{"~", "..."}, parts[i:]...)...)
		if len(candidate) <= maxLength {
			return candidate
		}
	}

	return displayPath[len(displayPath)-maxLength:]
}

func ShortCwd(maxLength int) string {
	cwd, _ := os.Getwd()
	return ShortenPath(cwd, maxLength)
}
