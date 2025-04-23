package components

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HelpOverlayData struct {
	Commands []Command
	Shortcuts []Shortcut
}

type Command struct {
	Name        string
	Description string
}

type Shortcut struct {
	Key         string
	Description string
}

func HelpOverlay(c *gin.Context) {
	data := HelpOverlayData{
		Commands: []Command{
			{Name: "/help", Description: "show this help overlay"},
			{Name: "/model", Description: "switch the LLM model in-session"},
			{Name: "/approval", Description: "switch auto-approval mode"},
			{Name: "/history", Description: "show command & file history for this session"},
			{Name: "/clear", Description: "clear screen & context"},
			{Name: "/clearhistory", Description: "clear command history"},
			{Name: "/bug", Description: "file a bug report with session log"},
			{Name: "/diff", Description: "view working tree git diff"},
			{Name: "/compact", Description: "condense context into a summary"},
		},
		Shortcuts: []Shortcut{
			{Key: "Enter", Description: "send message"},
			{Key: "Ctrl+J", Description: "insert newline"},
			{Key: "Up/Down", Description: "scroll prompt history"},
			{Key: "Esc (✕2)", Description: "interrupt current action"},
			{Key: "Ctrl+C", Description: "quit Codex"},
		},
	}

	tmpl, err := template.ParseFiles("templates/help_overlay.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading template: %v", err)
		return
	}

	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, "Error rendering template: %v", err)
	}
}
