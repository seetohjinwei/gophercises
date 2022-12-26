/*
go build -o links . && ./links
*/

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/seetohjinwei/gophercises/2-url-short/url"
	yaml "github.com/seetohjinwei/gophercises/2-url-short/yamlparser"
)

func generateLinks(filePath string) []url.Url {
	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("error while reading yaml file: %v", err)
	}

	return yaml.Parse(contents)
}

func redirectHandler(link url.Url) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%q: redirecting to %q", link.Short, link.Full)
		http.Redirect(w, r, link.Full, http.StatusMovedPermanently)
	}
}

func handleLinks(mux *http.ServeMux, links []url.Url) {
	for _, link := range links {
		mux.HandleFunc("/"+link.Short, redirectHandler(link))
	}
}

const defaultYaml = "links.yaml"
const port = ":8085"

var filePath string

func getFlags() {
	flag.StringVar(&filePath, "file", defaultYaml, "file path for links, see 'links.yaml' for an example")
	flag.Parse()
}

func main() {
	getFlags()

	mux := http.NewServeMux()

	// Everything else routes here.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		content := "URL Shortener!"
		w.Write([]byte(content))
	})

	links := generateLinks(filePath)

	handleLinks(mux, links)

	log.Printf("Starting on port: %v", port)
	http.ListenAndServe(port, mux)
}
