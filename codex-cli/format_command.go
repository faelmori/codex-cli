package main

import (
	"strings"
)

// formatCommandForDisplay formats the args of an exec command for display as a single string.
func formatCommandForDisplay(command []string) string {
	if len(command) == 3 && command[0] == "bash" && command[1] == "-lc" {
		inner := command[2]
		if strings.HasPrefix(inner, "'") && strings.HasSuffix(inner, "'") {
			inner = inner[1 : len(inner)-1]
		}
		return inner
	}
	return strings.Join(command, " ")
}
