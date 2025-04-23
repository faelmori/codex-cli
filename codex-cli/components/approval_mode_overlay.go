package components

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApprovalModeOverlayData struct {
	CurrentMode string
	Modes       []string
}

func ApprovalModeOverlay(c *gin.Context) {
	data := ApprovalModeOverlayData{
		CurrentMode: c.Query("currentMode"),
		Modes:       getApprovalModes(),
	}

	tmpl, err := template.ParseFiles("templates/approval_mode_overlay.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading template: %v", err)
		return
	}

	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, "Error rendering template: %v", err)
	}
}

func getApprovalModes() []string {
	// Implement the logic to get the list of approval modes
	return []string{"suggest", "auto-edit", "full-auto"}
}
