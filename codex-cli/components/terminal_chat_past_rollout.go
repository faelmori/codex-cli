package components

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TerminalChatPastRolloutData struct {
	Session TerminalChatSession
	Items   []ResponseItem
}

func TerminalChatPastRollout(c *gin.Context) {
	data := TerminalChatPastRolloutData{
		Session: getSession(),
		Items:   getResponseItems(),
	}

	tmpl, err := template.ParseFiles("templates/terminal_chat_past_rollout.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading template: %v", err)
		return
	}

	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, "Error rendering template: %v", err)
	}
}

func getSession() TerminalChatSession {
	// Implement the logic to get the TerminalChatSession
	return TerminalChatSession{}
}

func getResponseItems() []ResponseItem {
	// Implement the logic to get the ResponseItems
	return []ResponseItem{}
}
