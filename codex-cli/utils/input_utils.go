package utils

import (
	"encoding/base64"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/h2non/filetype"
)

type ResponseInputItem struct {
	Role    string
	Content []ResponseContent
	Type    string
}

type ResponseContent struct {
	Type     string
	Text     string
	Filename string
	Detail   string
	ImageURL string
}

func CreateInputItem(text string, images []string) (ResponseInputItem, error) {
	inputItem := ResponseInputItem{
		Role:    "user",
		Content: []ResponseContent{{Type: "input_text", Text: text}},
		Type:    "message",
	}

	for _, filePath := range images {
		binary, err := ioutil.ReadFile(filePath)
		if err != nil {
			inputItem.Content = append(inputItem.Content, ResponseContent{
				Type: "input_text",
				Text: "[missing image: " + filepath.Base(filePath) + "]",
			})
			continue
		}

		kind, err := filetype.Match(binary)
		if err != nil {
			inputItem.Content = append(inputItem.Content, ResponseContent{
				Type: "input_text",
				Text: "[invalid image: " + filepath.Base(filePath) + "]",
			})
			continue
		}

		encoded := base64.StdEncoding.EncodeToString(binary)
		mime := kind.MIME.Value
		inputItem.Content = append(inputItem.Content, ResponseContent{
			Type:     "input_image",
			Detail:   "auto",
			ImageURL: "data:" + mime + ";base64," + encoded,
		})
	}

	return inputItem, nil
}
