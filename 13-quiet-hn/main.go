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
	"sync/atomic"
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

func getItem(client *hn.Client, ids []int, i int, next_i *atomic.Int32, wg *sync.WaitGroup, ch chan<- item) {
	defer wg.Done()

	id := ids[i]

	hnItem, err := client.GetItem(id)
	if err != nil {
		return
	}
	item := parseHNItem(hnItem, i)

	if isStoryLink(item) {
		ch <- item
	} else {
		i2 := next_i.Add(1)
		wg.Add(1)
		getItem(client, ids, int(i2), next_i, wg, ch)
	}
}

func getItems(client *hn.Client, ids []int, numStories int) []item {
	ch := make(chan item)
	var wg sync.WaitGroup
	next_i := atomic.Int32{}
	next_i.Store(int32(numStories - 1))

	var stories []item

	for i := 0; i < numStories; i++ {
		wg.Add(1)
		go getItem(client, ids, i, &next_i, &wg, ch)
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

	return stories
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
