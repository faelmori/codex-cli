package singlepass

import (
	"fmt"
	"strings"
)

type ContextLimit struct {
	MaxTokens int
}

func (cl *ContextLimit) ApplyLimit(content string) (string, error) {
	tokens := strings.Fields(content)
	if len(tokens) > cl.MaxTokens {
		return "", fmt.Errorf("content exceeds the maximum token limit of %d", cl.MaxTokens)
	}
	return content, nil
}

func (cl *ContextLimit) TruncateContent(content string) string {
	tokens := strings.Fields(content)
	if len(tokens) > cl.MaxTokens {
		return strings.Join(tokens[:cl.MaxTokens], " ")
	}
	return content
}
