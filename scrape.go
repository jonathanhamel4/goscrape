package main

import (
	"fmt"
	"os"
	"net/url"
	"strings"
	"encoding/json"
	"github.com/gocolly/colly/v2"
)

type movie struct {
	Title string `selector:".title_wrapper h1"`
}

func verifyErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	articles := make([]*movie, 0)

	c := colly.NewCollector(
		colly.AllowedDomains("www.imdb.com"),
	)

	articleCollector := c.Clone()

	// collect the data
	articleCollector.OnHTML("body", func(e *colly.HTMLElement) {
		movieIndex := &movie{ }
		e.Unmarshal(movieIndex)
		articles = append(articles, movieIndex)
	})

	// find the root genre
	c.OnHTML(".ninja_image a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		resource, err := url.Parse(link)
		verifyErr(err)

		queryString, err := url.ParseQuery(resource.RawQuery)
		verifyErr(err)

		if len(queryString["genres"]) > 0 {
			fmt.Println("Discovered genre: ", queryString["genres"][0])
			c.Visit(e.Request.AbsoluteURL(link))
		}
	})

	// Navigate to the movie information
	c.OnHTML("h3.lister-item-header a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		title, err := e.DOM.Html()
		verifyErr(err)

		if title != "" && strings.Contains(link, "title/tt") {
			fmt.Println("Discovered title: ", title)
			articleCollector.Visit(e.Request.AbsoluteURL(link))
		}
	})

	c.Visit("https://www.imdb.com/search/title/?genres=comedy&explore=title_type,genres&pf_rd_m=A2FGELUUNOQJNL&pf_rd_p=3396781f-d87f-4fac-8694-c56ce6f490fe&pf_rd_r=D2GKRDFXVJGEER11W3FT&pf_rd_s=center-1&pf_rd_t=15051&pf_rd_i=genre&ref_=ft_gnr_pr1_i_1")

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	// Dump json to the standard output
	enc.Encode(articles)
}