package components

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TypeaheadItem struct {
	Label string
	Value string
}

type TypeaheadOverlayData struct {
	Title       string
	Description string
	Items       []TypeaheadItem
	CurrentValue string
	Limit       int
}

func TypeaheadOverlay(c *gin.Context) {
	data := TypeaheadOverlayData{
		Title:       c.Query("title"),
		Description: c.Query("description"),
		Items:       getTypeaheadItems(),
		CurrentValue: c.Query("currentValue"),
		Limit:       getLimit(c.Query("limit")),
	}

	tmpl, err := template.ParseFiles("templates/typeahead_overlay.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading template: %v", err)
		return
	}

	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, "Error rendering template: %v", err)
	}
}

func getTypeaheadItems() []TypeaheadItem {
	// Implement the logic to get the TypeaheadItems
	return []TypeaheadItem{}
}

func getLimit(limitStr string) int {
	// Implement the logic to parse and return the limit
	return 10
}
