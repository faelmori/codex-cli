package components

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ModelOverlayData struct {
	CurrentModel       string
	CurrentProvider    string
	HasLastResponse    bool
	AvailableModels    []string
	AvailableProviders []string
}

func ModelOverlay(c *gin.Context) {
	data := ModelOverlayData{
		CurrentModel:       c.Query("currentModel"),
		CurrentProvider:    c.Query("currentProvider"),
		HasLastResponse:    c.Query("hasLastResponse") == "true",
		AvailableModels:    getAvailableModels(),
		AvailableProviders: getAvailableProviders(),
	}

	tmpl, err := template.ParseFiles("templates/model_overlay.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading template: %v", err)
		return
	}

	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, "Error rendering template: %v", err)
	}
}

func getAvailableModels() []string {
	// Implement the logic to get the available models
	return []string{"model1", "model2", "model3"}
}

func getAvailableProviders() []string {
	// Implement the logic to get the available providers
	return []string{"provider1", "provider2", "provider3"}
}
