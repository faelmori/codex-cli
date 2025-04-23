package components

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SelectInputData struct {
	Items          []Item
	IsFocused      bool
	InitialIndex   int
	Limit          int
	SelectedIndex  int
	RotateIndex    int
	Confirmation   string
	Explanation    string
}

type Item struct {
	Key   string
	Label string
	Value string
}

func SelectInput(c *gin.Context) {
	data := SelectInputData{
		Items:          getItems(),
		IsFocused:      c.Query("isFocused") == "true",
		InitialIndex:   getIntQuery(c, "initialIndex", 0),
		Limit:          getIntQuery(c, "limit", 10),
		SelectedIndex:  getIntQuery(c, "selectedIndex", 0),
		RotateIndex:    getIntQuery(c, "rotateIndex", 0),
		Confirmation:   c.Query("confirmation"),
		Explanation:    c.Query("explanation"),
	}

	tmpl, err := template.ParseFiles("templates/select_input.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading template: %v", err)
		return
	}

	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, "Error rendering template: %v", err)
	}
}

func getItems() []Item {
	// Implement the logic to get the Items
	return []Item{}
}

func getIntQuery(c *gin.Context, key string, defaultValue int) int {
	value, err := c.GetQuery(key)
	if !err {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}
