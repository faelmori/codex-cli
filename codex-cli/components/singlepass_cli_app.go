package components

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SinglePassAppData struct {
	OriginalPrompt string
	Config         AppConfig
	RootPath       string
}

func SinglePassApp(c *gin.Context) {
	data := SinglePassAppData{
		OriginalPrompt: c.Query("originalPrompt"),
		Config:         getConfig(),
		RootPath:       c.Query("rootPath"),
	}

	tmpl, err := template.ParseFiles("templates/singlepass_cli_app.html")
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
