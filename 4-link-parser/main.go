package main

import (
	"flag"
	"fmt"
	"io"

	"github.com/seetohjinwei/gophercises/4-link-parser/files"
	"github.com/seetohjinwei/gophercises/4-link-parser/models"
	"github.com/seetohjinwei/gophercises/4-link-parser/url"
)

const defaultWebsite = "https://jinwei.dev"

func main() {
	var file string
	var website string
	flag.StringVar(&file, "file", "", "file path for html page")
	flag.StringVar(&website, "url", defaultWebsite, "url path for html page")
	flag.Parse()

	var f io.Reader

	if file != "" {
		f = files.Open(file)
	} else {
		f = url.Open(website)
	}

	links := models.GetLinks(f)

	for _, link := range links {
		fmt.Printf("%q: %q\n", link.Href, link.Text)
	}
}
