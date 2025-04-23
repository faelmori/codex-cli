package components

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DiffOverlayData struct {
	DiffText string
}

func DiffOverlay(c *gin.Context) {
	data := DiffOverlayData{
		DiffText: getDiffText(),
	}

	tmpl, err := template.ParseFiles("templates/diff_overlay.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading template: %v", err)
		return
	}

	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, "Error rendering template: %v", err)
	}
}

func getDiffText() string {
	// Implement the logic to get the diff text
	return "(no changes)"
}
