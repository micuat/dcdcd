package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"

	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

type Count struct {
	Count int
}

func main() {
	jsonData, err := os.ReadFile("data.json")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonData))

	var data []map[string]interface{}
	err = json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		fmt.Printf("could not unmarshal json: %s\n", err)
		return
	}

	fmt.Printf("json map: %v\n", data)

	e := echo.New()
	// e.Use(middleware.Logger())

	count := Count{Count: 0}

	e.Renderer = newTemplate()

	e.GET("/", func(c echo.Context) error {
		quote := data[rand.Intn(len(data))]
		return c.Render(200, "index", quote)
	})

	e.POST("/count", func(c echo.Context) error {
		count.Count++
		quote := data[rand.Intn(len(data))]
		return c.Render(200, "count", quote)
	})

	e.Logger.Fatal(e.Start(":8090"))
}
