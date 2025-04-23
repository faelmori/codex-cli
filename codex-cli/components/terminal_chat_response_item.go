package components

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseItem struct {
	Type    string
	Role    string
	Content []ContentItem
}

type ContentItem struct {
	Type     string
	Text     string
	Filename string
	Refusal  string
}

type TerminalChatResponseItemData struct {
	Item       ResponseItem
	FullStdout bool
}

func TerminalChatResponseItem(c *gin.Context) {
	data := TerminalChatResponseItemData{
		Item:       getResponseItem(),
		FullStdout: c.Query("fullStdout") == "true",
	}

	tmpl, err := template.ParseFiles("templates/terminal_chat_response_item.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading template: %v", err)
		return
	}

	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, "Error rendering template: %v", err)
	}
}

func getResponseItem() ResponseItem {
	// Implement the logic to get the ResponseItem
	return ResponseItem{}
}
