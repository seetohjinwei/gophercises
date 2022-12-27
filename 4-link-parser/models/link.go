package models

import (
	"bytes"
	"io"
	"log"
	"strings"

	"github.com/seetohjinwei/gophercises/4-link-parser/tokeniser"
	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

type node struct {
	href string
	sb   *strings.Builder
}

func GetLinks(f io.Reader) []Link {
	z := html.NewTokenizer(f)

	links := []Link{}
	anchors := []node{}

tokenise:
	for {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			break tokenise
		case html.StartTagToken:
			if !tokeniser.IsAnchor(z) {
				continue
			}

			href, err := tokeniser.GetHref(z)
			if err != nil {
				log.Fatalf("error: %v", err)
			}

			n := node{href: href, sb: &strings.Builder{}}
			anchors = append(anchors, n)
		case html.EndTagToken:
			if !tokeniser.IsAnchor(z) {
				continue
			}

			n := anchors[len(anchors)-1]       // top
			anchors = anchors[:len(anchors)-1] // pop

			ss := n.sb.String()
			link := Link{
				Href: n.href,
				Text: ss,
			}
			links = append(links, link)
		case html.TextToken:
			if len(anchors) == 0 {
				continue
			}
			b := z.Text()
			trimmed := bytes.TrimSpace(b)
			n := anchors[len(anchors)-1]
			n.sb.Write(trimmed)
		}
	}

	return links
}
