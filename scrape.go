package main

import (
	"fmt"
	"os"
	"net/url"
	"strings"
	"encoding/json"
	"github.com/joho/godotenv"
	"github.com/gocolly/colly/v2"

	"github.com/jonathanhamel4/goscrape/db"
	"github.com/jonathanhamel4/goscrape/types"
)

func scrape() {
	articles := make([]*types.Movie, 0)
	genres := make([]types.Genre, 0)

	c := colly.NewCollector(
		colly.AllowedDomains("www.imdb.com"),
	)

	articleCollector := c.Clone()

	// collect the data
	articleCollector.OnHTML("body", func(e *colly.HTMLElement) {
		movieIndex := &types.Movie{
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
		VerifyError(err)

		queryString, err := url.ParseQuery(resource.RawQuery)
		VerifyError(err)

		if len(queryString["genres"]) > 0 {
			fmt.Println("Discovered genre: ", queryString["genres"][0])
			genres = append(genres, types.Genre(queryString["genres"][0]))
			c.Visit(e.Request.AbsoluteURL(link))
		}
	})

	// Navigate to the movie information
	c.OnHTML("h3.lister-item-header a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		title, err := e.DOM.Html()
		VerifyError(err)

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

	db_provider.InsertMovies(articles)
}


func main() {
	defer scrape()

	err := godotenv.Load()
	if err != nil {
	  panic("Error loading .env file")
	}
	
	mongoConn := os.Getenv("MONGO_CONNECTION_STRING")

	db_provider.ConnectDB(mongoConn)
}