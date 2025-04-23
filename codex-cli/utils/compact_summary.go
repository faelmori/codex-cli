package utils

import (
	"strings"

	"github.com/openai/openai-go"
)

type AppConfig struct {
	Provider string
}

type ResponseItem struct {
	Type    string
	Role    string
	Content []ResponseContent
}

type ResponseContent struct {
	Type     string
	Text     string
	Filename string
	Refusal  string
}

func getBaseUrl(provider string) string {
	// Implement the logic to get the base URL for the provider
	return ""
}

func getApiKey(provider string) string {
	// Implement the logic to get the API key for the provider
	return ""
}

func generateCompactSummary(items []ResponseItem, model string, flexMode bool, config AppConfig) (string, error) {
	oai := openai.NewClient(getApiKey(config.Provider), getBaseUrl(config.Provider))

	var conversationText strings.Builder
	for _, item := range items {
		if item.Type == "message" && (item.Role == "user" || item.Role == "assistant") {
			for _, content := range item.Content {
				if content.Type == "text" {
					conversationText.WriteString(item.Role + ": " + content.Text + "\n")
				}
			}
		}
	}

	response, err := oai.ChatCompletions.Create(openai.ChatCompletionRequest{
		Model: model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "system",
				Content: "You are an expert coding assistant. Your goal is to generate a concise, structured summary of the conversation below that captures all essential information needed to continue development after context replacement. Include tasks performed, code areas modified or reviewed, key decisions or assumptions, test results or errors, and outstanding tasks or next steps.",
			},
			{
				Role:    "user",
				Content: "Here is the conversation so far:\n" + conversationText.String() + "\n\nPlease summarize this conversation, covering:\n1. Tasks performed and outcomes\n2. Code files, modules, or functions modified or examined\n3. Important decisions or assumptions made\n4. Errors encountered and test or build results\n5. Remaining tasks, open questions, or next steps\nProvide the summary in a clear, concise format.",
			},
		},
		ServiceTier: func() string {
			if flexMode {
				return "flex"
			}
			return ""
		}(),
	})
	if err != nil {
		return "", err
	}

	return response.Choices[0].Message.Content, nil
}
