package main

import (
	"fmt"
	"math/rand"

	"html/template"
	"io"

	"dcdcd.glitches.me/storage"
	"github.com/labstack/echo/v4"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

// type Count struct {
// 	Count int
// }

func main() {
	quotes := storage.GetQuotes()
	fmt.Printf("quotes: %v", quotes)

	e := echo.New()
	// e.Use(middleware.Logger())

	// count := Count{Count: 0}

	e.Renderer = NewTemplate()

	e.GET("/", func(c echo.Context) error {
		quote := quotes[rand.Intn(len(quotes))]
		return c.Render(200, "index", quote)
	})

	// e.POST("/count", func(c echo.Context) error {
	// 	count.Count++
	// 	quote := quotes[rand.Intn(len(quotes))]
	// 	return c.Render(200, "count", quote)
	// })

	e.Logger.Fatal(e.Start(":8090"))
}
