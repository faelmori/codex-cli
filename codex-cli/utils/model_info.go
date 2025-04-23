package utils

import (
	"fmt"
	"strings"
)

// ModelInfo represents information about a model.
type ModelInfo struct {
	Name        string
	Description string
	Version     string
}

// GetModelInfo returns information about a model based on its name.
func GetModelInfo(modelName string) (*ModelInfo, error) {
	models := map[string]ModelInfo{
		"o4-mini": {
			Name:        "o4-mini",
			Description: "A small, efficient model for quick completions.",
			Version:     "1.0.0",
		},
		"openai": {
			Name:        "openai",
			Description: "OpenAI's powerful model for advanced completions.",
			Version:     "2.1.0",
		},
	}

	model, exists := models[strings.ToLower(modelName)]
	if !exists {
		return nil, fmt.Errorf("model %s not found", modelName)
	}

	return &model, nil
}
