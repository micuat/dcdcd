package main

import (
	"math/rand"
	"strconv"
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

type QuoteContainer struct {
	Start  int
	Quotes []QuoteDiv
	Next   int
}

type QuoteDiv struct {
	storage.Quote
	Id        int
	EmbedMore bool
	Next      int
	Hashtag   string
}

func moreQuotes(hashtag string, start int, showStep int) QuoteContainer {
	qs := []QuoteDiv{}
	quotes := storage.SearchQuotes(hashtag)
	// TODO: error when length is 0
	for i := start; i < start+showStep; i++ {
		quote := quotes[rand.Intn(len(quotes))]
		qs = append(qs, QuoteDiv{
			Quote:     quote,
			Id:        i,
			EmbedMore: start+showStep < 100 && i == start+showStep-3,
			Next:      start + showStep,
			Hashtag:   "",
		})
	}
	return QuoteContainer{
		Start:  start,
		Next:   start + showStep,
		Quotes: qs,
	}
}

func main() {
	// fmt.Printf("quotes: %v", quotes)

	e := echo.New()
	// e.Use(middleware.Logger())
	e.Static("/static", "static")

	// count := Count{Count: 0}

	e.Renderer = NewTemplate()

	showStep := 10

	e.GET("/", func(c echo.Context) error {
		start := 0
		return c.Render(200, "index", moreQuotes("", start, showStep))
	})

	e.GET("/get/quotes", func(c echo.Context) error {
		startStr := c.QueryParam("start")
		start, err := strconv.Atoi(startStr)
		if err != nil {
			start = 0
		}
		hashtag := c.QueryParam("hashtag")
		return c.Render(200, "quotes-pane", moreQuotes(hashtag, start, showStep))
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
		// page.Data.Contacts = append(page.Data.Contacts, contact)
		// c.Render(200, "form", newFormData())
		// return c.Render(200, "oob-contact", contact)
		return c.Render(200, "submitted", nil)
	})
	e.Logger.Fatal(e.Start(":8090"))
}
