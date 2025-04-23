package utils

import (
	"os"
	"path/filepath"
	"strings"
)

type FileSystemSuggestion struct {
	Path  string
	IsDir bool
}

func GetFileSystemSuggestions(basePath string, query string) ([]FileSystemSuggestion, error) {
	var suggestions []FileSystemSuggestion

	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.Contains(strings.ToLower(info.Name()), strings.ToLower(query)) {
			suggestions = append(suggestions, FileSystemSuggestion{
				Path:  path,
				IsDir: info.IsDir(),
			})
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return suggestions, nil
}
