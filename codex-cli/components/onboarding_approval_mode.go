package components

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OnboardingApprovalModeData struct {
	CurrentMode string
	Modes       []string
}

func OnboardingApprovalMode(c *gin.Context) {
	data := OnboardingApprovalModeData{
		CurrentMode: c.Query("currentMode"),
		Modes:       getApprovalModes(),
	}

	tmpl, err := template.ParseFiles("templates/onboarding_approval_mode.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading template: %v", err)
		return
	}

	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, "Error rendering template: %v", err)
	}
}

func getApprovalModes() []string {
	return []string{"Auto-approve file reads, but ask me for edits and commands", "Auto-approve file reads and edits, but ask me for commands", "Auto-approve file reads, edits, and running commands network-disabled"}
}
