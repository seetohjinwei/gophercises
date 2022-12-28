package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/seetohjinwei/gophercises/13-quiet-hn/hn"
	"golang.org/x/exp/slices"
)

const concurrencyMultiplier float32 = 1.25

func main() {
	// parse flags
	var port, numStories int
	flag.IntVar(&port, "port", 3000, "the port to start the web server on")
	flag.IntVar(&numStories, "num_stories", 30, "the number of top stories to display")
	flag.Parse()

	tpl := template.Must(template.ParseFiles("./index.gohtml"))

	http.HandleFunc("/", handler(numStories, tpl))

	// Start the server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handler(numStories int, tpl *template.Template) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		var client hn.Client
		ids, err := client.TopItems()
		if err != nil {
			http.Error(w, "Failed to load top stories", http.StatusInternalServerError)
			return
		}
		stories := getItems(&client, ids, numStories)
		data := templateData{
			Stories: stories,
			Time:    time.Now().Sub(start),
		}
		err = tpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	})
}

func getItems(client *hn.Client, ids []int, numStories int) []item {
	ch := make(chan item)
	var wg sync.WaitGroup

	var stories []item

	numStoriesToGet := int(float32(numStories) * concurrencyMultiplier)

	for i := 0; i < numStoriesToGet; i++ {
		id := ids[i]

		wg.Add(1)
		go func(i, id int) {
			defer wg.Done()

			hnItem, err := client.GetItem(id)
			if err != nil {
				return
			}
			item := parseHNItem(hnItem, i)

			if isStoryLink(item) {
				ch <- item
			}
		}(i, id)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for item := range ch {
		stories = append(stories, item)
	}

	// sort the stories
	slices.SortFunc(stories, func(a, b item) bool {
		return a.Index < b.Index
	})

	// return only `numStories` or less
	if len(stories) <= numStories {
		return stories
	}

	return stories[:numStories]
}

func isStoryLink(item item) bool {
	return item.Type == "story" && item.URL != ""
}

func parseHNItem(hnItem hn.Item, index int) item {
	ret := item{Item: hnItem, Index: index}
	url, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(url.Hostname(), "www.")
	}
	return ret
}

// item is the same as the hn.Item, but adds the Host field
type item struct {
	hn.Item
	Host  string
	Index int
}

type templateData struct {
	Stories []item
	Time    time.Duration
}
