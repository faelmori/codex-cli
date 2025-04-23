package components

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TerminalChatData struct {
	Config                AppConfig
	Prompt                string
	ImagePaths            []string
	ApprovalPolicy        ApprovalPolicy
	AdditionalWritableRoots []string
	FullStdout            bool
}

func TerminalChat(c *gin.Context) {
	data := TerminalChatData{
		Config:                getConfig(),
		Prompt:                c.Query("prompt"),
		ImagePaths:            c.QueryArray("imagePaths"),
		ApprovalPolicy:        getApprovalPolicy(),
		AdditionalWritableRoots: getAdditionalWritableRoots(),
		FullStdout:            c.Query("fullStdout") == "true",
	}

	tmpl, err := template.ParseFiles("templates/terminal_chat.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading template: %v", err)
		return
	}

	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, "Error rendering template: %v", err)
	}
}

func getConfig() AppConfig {
	// Implement the logic to get the AppConfig
	return AppConfig{}
}

func getApprovalPolicy() ApprovalPolicy {
	// Implement the logic to get the ApprovalPolicy
	return ApprovalPolicy{}
}

func getAdditionalWritableRoots() []string {
	// Implement the logic to get the AdditionalWritableRoots
	return []string{}
}
