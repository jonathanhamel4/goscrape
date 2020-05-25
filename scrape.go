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
	Country string `selector:"#titleDetails > .txt-block:first-of-type a"`
	Language string `selector:"#titleDetails > .txt-block:nth-of-type(2) a"`
	Rating string `selector:"span[itemprop='ratingValue']"`
	RatingCount string `selector:"span[itemprop='ratingCount']"`
	Summary string `selector:"div.summary_text"`
	StoryLine string `selector:"#titleStoryLine div > p > span"`
	Runtime string `selector:"#titleDetails > .txt-block:nth-of-type(8) > time"`
	Genres []string `selector:".title_wrapper a[href]:not([title])"`
	ImdbUrl string
}

func verifyErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	articles := make([]*movie, 0)
	genres := make([]string, 0)

	c := colly.NewCollector(
		colly.AllowedDomains("www.imdb.com"),
	)

	articleCollector := c.Clone()

	// collect the data
	articleCollector.OnHTML("body", func(e *colly.HTMLElement) {
		movieIndex := &movie{
			ImdbUrl: e.Request.URL.String(),
		}
		e.Unmarshal(movieIndex)
		movieIndex.ImdbUrl = e.Request.URL.String()

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
			genres = append(genres, queryString["genres"][0])
			c.Visit(e.Request.AbsoluteURL(link))
		}
	})

	// Navigate to the movie information
	c.OnHTML("h3.lister-item-header a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		title, err := e.DOM.Html()
		verifyErr(err)

		if title != "" && strings.Contains(link, "title/tt") && len(articles) == 0 {
			fmt.Println("Discovered title: ", title)
			articleCollector.Visit(e.Request.AbsoluteURL(link))
		}
	})

	c.Visit("https://www.imdb.com/feature/genre/")

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	// Dump json to the standard output
	enc.Encode(articles)
	enc.Encode(genres)
}