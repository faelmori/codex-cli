package singlepass

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ContextFile struct {
	Path    string
	Content string
}

func LoadContextFiles(dir string) ([]ContextFile, error) {
	var contextFiles []ContextFile

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			contextFiles = append(contextFiles, ContextFile{
				Path:    path,
				Content: string(content),
			})
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to load context files: %w", err)
	}

	return contextFiles, nil
}

func SaveContextFiles(contextFiles []ContextFile, dir string) error {
	for _, file := range contextFiles {
		path := filepath.Join(dir, file.Path)
		err := os.MkdirAll(filepath.Dir(path), 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		err = os.WriteFile(path, []byte(file.Content), 0644)
		if err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}
	}

	return nil
}
