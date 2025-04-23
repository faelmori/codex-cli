package components

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TerminalChatToolCallCommandData struct {
	CommandForDisplay string
	Explanation       string
}

func TerminalChatToolCallCommand(c *gin.Context) {
	data := TerminalChatToolCallCommandData{
		CommandForDisplay: c.Query("commandForDisplay"),
		Explanation:       c.Query("explanation"),
	}

	tmpl, err := template.ParseFiles("templates/terminal_chat_tool_call_command.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading template: %v", err)
		return
	}

	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, "Error rendering template: %v", err)
	}
}
