package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gocolly/colly"
)

type pageInfo struct {
	Heading map[string]int
	Excerpt map[string]int
}

type pages struct {
	articles []pageInfo
}

func (box *pages) AddItem(item MyBoxItem) []MyBoxItem {
	box.Items = append(box.Items, item)
	return box.Items
}

func main() {
	http.HandleFunc("/", crawl)
	http.ListenAndServe(":3000", nil)
}

func crawl(w http.ResponseWriter, r *http.Request) {

	scrappedPage := &pageInfo{Heading: make(map[string]int), Excerpt: make(map[string]int)}

	pages := make([]pageInfo, 20)

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
			scrappedPage.Heading[heading]++
			scrappedPage.Excerpt[excerpt]++

			pages = append({Heading: "ef", Excerpt: "fe"})

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
