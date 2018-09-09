package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gocolly/colly"
)

type pageInfo struct {
	StatusCode int
	Headings   map[string]int
	Excerpt    map[string]int
}

func main() {
	fmt.Println("Visiting")

	scrappedPage := &pageInfo{Headings: make(map[string]int), Excerpt: make(map[string]int)}

	url := "https://hackernoon.com/"

	collector := colly.NewCollector(
		colly.AllowedDomains("hackernoon.com", "medium.com"),
	)

	collector.OnResponse(func(r *colly.Response) {
		log.Println("response received", r.StatusCode)
		scrappedPage.StatusCode = r.StatusCode
	})

	collector.OnError(func(r *colly.Response, err error) {
		log.Println("error:", r.StatusCode, err)
		scrappedPage.StatusCode = r.StatusCode
	})

	collector.OnHTML(".js-trackedPost a", func(e *colly.HTMLElement) {
		heading := e.ChildText("h3")
		excerpt := e.ChildText("div.u-contentSansThin")

		if heading != "" && excerpt != "" {
			scrappedPage.Headings[heading]++
			scrappedPage.Excerpt[excerpt]++

			// fmt.Println("Extracted headline:", heading)
			// fmt.Println("Extracted exerpt:", excerpt)

			jsonOutput, err := json.Marshal(scrappedPage)

			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(string(jsonOutput))
		}
	})

	collector.Visit(url)
}
