package utils

import (
	"math"
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

func ApproximateTokensUsed(items []ResponseItem) int {
	charCount := 0

	for _, item := range items {
		switch item.Type {
		case "message":
			if item.Role != "user" && item.Role != "assistant" {
				continue
			}

			for _, c := range item.Content {
				if c.Type == "input_text" || c.Type == "output_text" {
					charCount += len(c.Text)
				} else if c.Type == "refusal" {
					charCount += len(c.Refusal)
				} else if c.Type == "input_file" {
					charCount += len(c.Filename)
				}
				// images and other content types are ignored (0 chars)
			}

		case "function_call":
			charCount += len(item.Name) + len(item.Arguments)

		case "function_call_output":
			charCount += len(item.Output)

		default:
			continue
		}
	}

	return int(math.Ceil(float64(charCount) / 4))
}
