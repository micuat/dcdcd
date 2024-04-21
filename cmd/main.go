package main

import (
	"math/rand"
	"strings"

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
	// fmt.Printf("quotes: %v", quotes)

	e := echo.New()
	// e.Use(middleware.Logger())
	e.Static("/static", "static")

	// count := Count{Count: 0}

	e.Renderer = NewTemplate()

	e.GET("/", func(c echo.Context) error {
		quote := quotes[rand.Intn(len(quotes))]
		return c.Render(200, "index", quote)
	})

	e.POST("/newquote", func(c echo.Context) error {
		quote := c.FormValue("quote")
		url := c.FormValue("url")
		hashtags := c.FormValue("hashtags")

		// if page.Data.hasEmail(email) {
		// 	formData := newFormData()
		// 	formData.Values["name"] = name
		// 	formData.Values["email"] = email
		// 	formData.Errors["email"] = "Email already exists!"
		// 	return c.Render(422, "form", formData)
		// }

		// q := storage.NewQuote(quote, url)
		storage.AddQuote(quote, url, strings.Split(hashtags, ","))
		quotes = storage.GetQuotes()
		// page.Data.Contacts = append(page.Data.Contacts, contact)
		// c.Render(200, "form", newFormData())
		// return c.Render(200, "oob-contact", contact)
		return c.Render(200, "submitted", nil)
	})
	e.Logger.Fatal(e.Start(":8090"))
}
