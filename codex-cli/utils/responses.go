package utils

import (
	"fmt"
	"strings"
	"time"
)

type ResponseCreateInput struct {
	Model                string
	Input                []ResponseInputItem
	PreviousResponseID   string
	Instructions         string
	Tools                []Tool
	Temperature          float64
	TopP                 float64
	ToolChoice           string
	User                 string
	Metadata             map[string]string
	ParallelToolCalls    bool
	Truncation           string
}

type ResponseOutput struct {
	ID                   string
	Object               string
	CreatedAt            int64
	Status               string
	Error                *ResponseError
	IncompleteDetails    *IncompleteDetails
	Instructions         string
	MaxOutputTokens      int
	Model                string
	Output               []ResponseItem
	ParallelToolCalls    bool
	PreviousResponseID   string
	Reasoning            *Reasoning
	Store                bool
	Temperature          float64
	Text                 TextFormat
	ToolChoice           string
	Tools                []Tool
	TopP                 float64
	Truncation           string
	Usage                *Usage
	User                 string
	Metadata             map[string]string
}

type ResponseError struct {
	Code    string
	Message string
}

type IncompleteDetails struct {
	Reason string
}

type ResponseItem struct {
	Type    string
	ID      string
	Status  string
	Role    string
	Content []ResponseContent
}

type ResponseContent struct {
	Type     string
	Text     string
	Filename string
	Refusal  string
}

type Reasoning struct {
	Effort  string
	Summary string
}

type TextFormat struct {
	Type string
}

type Tool struct {
	Type        string
	Name        string
	Description string
	Parameters  map[string]interface{}
}

type Usage struct {
	InputTokens         int
	InputTokensDetails  TokenDetails
	OutputTokens        int
	OutputTokensDetails TokenDetails
	TotalTokens         int
}

type TokenDetails struct {
	CachedTokens int
}

type ResponseEvent struct {
	Type       string
	Response   *ResponseOutput
	OutputItem *ResponseItem
	ContentPart *ResponseContent
	Delta      string
	Arguments  string
	Error      *ResponseError
}

type ResponseInputItem struct {
	Type    string
	Role    string
	Content []ResponseContent
}

type ToolCallData struct {
	ID        string
	Name      string
	Arguments string
}

type ResponseContentOutput struct {
	Type       string
	CallID     string
	Name       string
	Arguments  string
	Text       string
	Annotations []interface{}
}

var conversationHistories = make(map[string]ConversationHistory)

type ConversationHistory struct {
	PreviousResponseID string
	Messages           []ChatCompletionMessage
}

type ChatCompletionMessage struct {
	Role    string
	Content string
}

func generateID(prefix string) string {
	return fmt.Sprintf("%s_%s", prefix, time.Now().Format("20060102150405"))
}

func convertInputItemToMessage(item ResponseInputItem) ChatCompletionMessage {
	content := ""
	for _, c := range item.Content {
		if c.Type == "input_text" {
			content += c.Text
		}
	}
	return ChatCompletionMessage{
		Role:    item.Role,
		Content: content,
	}
}

func getFullMessages(input ResponseCreateInput) []ChatCompletionMessage {
	var baseHistory []ChatCompletionMessage
	if input.PreviousResponseID != "" {
		prev, exists := conversationHistories[input.PreviousResponseID]
		if !exists {
			panic(fmt.Sprintf("Previous response not found: %s", input.PreviousResponseID))
		}
		baseHistory = prev.Messages
	}

	var newInputMessages []ChatCompletionMessage
	for _, item := range input.Input {
		newInputMessages = append(newInputMessages, convertInputItemToMessage(item))
	}

	messages := append(baseHistory, newInputMessages...)
	if input.Instructions != "" && len(messages) > 0 && messages[0].Role != "system" && messages[0].Role != "developer" {
		return append([]ChatCompletionMessage{{Role: "system", Content: input.Instructions}}, messages...)
	}
	return messages
}

func convertTools(tools []Tool) []Tool {
	var convertedTools []Tool
	for _, tool := range tools {
		if tool.Type == "function" {
			convertedTools = append(convertedTools, tool)
		}
	}
	return convertedTools
}

func responsesCreateViaChatCompletions(input ResponseCreateInput) (ResponseOutput, error) {
	fullMessages := getFullMessages(input)
	chatTools := convertTools(input.Tools)

	responseID := generateID("resp")
	outputItemID := generateID("msg")
	var outputContent []ResponseContentOutput

	for _, message := range fullMessages {
		if message.Role == "assistant" {
			outputContent = append(outputContent, ResponseContentOutput{
				Type: "output_text",
				Text: message.Content,
			})
		}
	}

	responseOutput := ResponseOutput{
		ID:                 responseID,
		Object:             "response",
		CreatedAt:          time.Now().Unix(),
		Status:             "completed",
		Model:              input.Model,
		Output:             []ResponseItem{{Type: "message", ID: outputItemID, Status: "completed", Role: "assistant", Content: []ResponseContent{{Type: "output_text", Text: strings.Join(outputContent, "\n")}}}},
		ParallelToolCalls:  input.ParallelToolCalls,
		PreviousResponseID: input.PreviousResponseID,
		Temperature:        input.Temperature,
		Text:               TextFormat{Type: "text"},
		ToolChoice:         input.ToolChoice,
		Tools:              input.Tools,
		TopP:               input.TopP,
		Truncation:         input.Truncation,
		Usage:              &Usage{InputTokens: len(fullMessages), OutputTokens: len(outputContent), TotalTokens: len(fullMessages) + len(outputContent)},
		User:               input.User,
		Metadata:           input.Metadata,
	}

	conversationHistories[responseID] = ConversationHistory{
		PreviousResponseID: input.PreviousResponseID,
		Messages:           fullMessages,
	}

	return responseOutput, nil
}
