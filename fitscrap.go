package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gocolly/colly"
)

type pageInfo struct {
	Headings map[string]int
	Excerpt  map[string]int
}

func main() {
	http.HandleFunc("/", crawl)
	http.ListenAndServe(":3000", nil)
}

func crawl(w http.ResponseWriter, r *http.Request) {
	scrappedPage := &pageInfo{Headings: make(map[string]int), Excerpt: make(map[string]int)}

	url := "https://hackernoon.com/"

	collector := colly.NewCollector(
		colly.AllowedDomains("hackernoon.com", "medium.com"),
	)

	collector.OnResponse(func(r *colly.Response) {
		log.Println("response received")
	})

	collector.OnError(func(r *colly.Response, err error) {
		log.Println("error:", err)
	})

	collector.OnHTML(".js-trackedPost a", func(e *colly.HTMLElement) {
		heading := e.ChildText("h3")
		excerpt := e.ChildText("div.u-contentSansThin")

		if heading != "" && excerpt != "" {
			scrappedPage.Headings[heading]++
			scrappedPage.Excerpt[excerpt]++

			jsonOutput, err := json.Marshal(scrappedPage)

			if err != nil {
				fmt.Println(err)
				return
			}

			// fmt.Println(string(jsonOutput))
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonOutput)
		}
	})

	collector.Visit(url)

}
