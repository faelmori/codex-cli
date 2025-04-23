package components

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TerminalHeaderData struct {
	TerminalRows       int
	Version            string
	PWD                string
	Model              string
	Provider           string
	ApprovalPolicy     string
	ColorsByPolicy     map[string]string
	SessionID          string
	InitialImagePaths  []string
	FlexModeEnabled    bool
}

func TerminalHeader(c *gin.Context) {
	data := TerminalHeaderData{
		TerminalRows:       getTerminalRows(),
		Version:            c.Query("version"),
		PWD:                c.Query("pwd"),
		Model:              c.Query("model"),
		Provider:           c.Query("provider"),
		ApprovalPolicy:     c.Query("approvalPolicy"),
		ColorsByPolicy:     getColorsByPolicy(),
		SessionID:          c.Query("sessionId"),
		InitialImagePaths:  c.QueryArray("initialImagePaths"),
		FlexModeEnabled:    c.Query("flexModeEnabled") == "true",
	}

	tmpl, err := template.ParseFiles("templates/terminal_header.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading template: %v", err)
		return
	}

	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, "Error rendering template: %v", err)
	}
}

func getTerminalRows() int {
	// Implement the logic to get the terminal rows
	return 24
}

func getColorsByPolicy() map[string]string {
	// Implement the logic to get the colors by policy
	return map[string]string{
		"suggest":   "gray",
		"auto-edit": "greenBright",
		"full-auto": "green",
	}
}
