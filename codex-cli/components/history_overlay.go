package components

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HistoryOverlayData struct {
	Commands []string
	Files    []string
}

func HistoryOverlay(c *gin.Context) {
	data := HistoryOverlayData{
		Commands: getCommands(),
		Files:    getFiles(),
	}

	tmpl, err := template.ParseFiles("templates/history_overlay.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading template: %v", err)
		return
	}

	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, "Error rendering template: %v", err)
	}
}

func getCommands() []string {
	// Implement the logic to get the list of commands
	return []string{}
}

func getFiles() []string {
	// Implement the logic to get the list of files
	return []string{}
}
