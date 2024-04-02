package storage

import (
	"encoding/json"
	"fmt"
	"os"
)

type Quote struct {
	Text string `json:"text"`
	Link string `json:"link"`
}

func NewQuote(text string, link string) Quote {
	return Quote{Text: text, Link: link}
}

func AddQuote(text string, link string) {
	quotes = append(quotes, NewQuote(text, link))
}

var quotes []Quote

func init() {
	rootPath, _ := os.Getwd()
	jsonData, err := os.ReadFile(rootPath + "/data.json")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	// fmt.Println(string(jsonData))

	err = json.Unmarshal(jsonData, &quotes)
	if err != nil {
		fmt.Printf("could not unmarshal json: %s\n", err)
		return
	}
	// fmt.Printf("json map: %v\n", quotes)
}

func GetQuotes() []Quote {
	return quotes
}
