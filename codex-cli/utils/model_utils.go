package utils

import (
	"strings"
)

type ResponseItem struct {
	Type    string
	Role    string
	Content []ResponseContent
	Name    string
	Arguments string
	Output  string
}

type ResponseContent struct {
	Type     string
	Text     string
	Filename string
	Refusal  string
}

func CalculateContextPercentRemaining(items []ResponseItem, model string) int {
	maxTokens := getMaxTokensForModel(model)
	usedTokens := ApproximateTokensUsed(items)
	return int(float64(maxTokens-usedTokens) / float64(maxTokens) * 100)
}

func getMaxTokensForModel(model string) int {
	switch strings.ToLower(model) {
	case "gpt-3.5-turbo":
		return 4096
	case "gpt-4":
		return 8192
	default:
		return 2048
	}
}

func UniqueById(items []ResponseItem) []ResponseItem {
	seen := make(map[string]bool)
	var result []ResponseItem
	for _, item := range items {
		if !seen[item.Name] {
			seen[item.Name] = true
			result = append(result, item)
		}
	}
	return result
}
