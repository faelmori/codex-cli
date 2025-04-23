package utils

import (
	"encoding/json"
	"strings"
)

type ResponseItem struct {
	Type      string
	Name      string
	Arguments string
}

func ExtractAppliedPatches(items []ResponseItem) string {
	var patches []string

	for _, item := range items {
		if item.Type != "function_call" {
			continue
		}

		if item.Name != "apply_patch" {
			continue
		}

		var args map[string]string
		if err := json.Unmarshal([]byte(item.Arguments), &args); err != nil {
			continue
		}

		if patch, ok := args["patch"]; ok && patch != "" {
			patches = append(patches, strings.TrimSpace(patch))
		}
	}

	return strings.Join(patches, "\n\n")
}
