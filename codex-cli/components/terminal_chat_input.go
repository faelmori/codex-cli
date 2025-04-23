package components

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TerminalChatInputData struct {
	IsNew                bool
	Loading              bool
	ConfirmationPrompt   string
	Explanation          string
	ContextLeftPercent   int
	Active               bool
	ThinkingSeconds      int
	Items                []ResponseItem
}

func TerminalChatInput(c *gin.Context) {
	data := TerminalChatInputData{
		IsNew:                c.Query("isNew") == "true",
		Loading:              c.Query("loading") == "true",
		ConfirmationPrompt:   c.Query("confirmationPrompt"),
		Explanation:          c.Query("explanation"),
		ContextLeftPercent:   c.Query("contextLeftPercent"),
		Active:               c.Query("active") == "true",
		ThinkingSeconds:      c.Query("thinkingSeconds"),
		Items:                getResponseItems(),
	}

	tmpl, err := template.ParseFiles("templates/terminal_chat_input.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading template: %v", err)
		return
	}

	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, "Error rendering template: %v", err)
	}
}

func getResponseItems() []ResponseItem {
	// Implement the logic to get the ResponseItems
	return []ResponseItem{}
}
