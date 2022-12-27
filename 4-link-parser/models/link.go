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

func GetLinks(f io.Reader) []Link {
	z := html.NewTokenizer(f)

	links := []Link{}
	anchors := []string{}
	sb := strings.Builder{}

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

			anchors = append(anchors, href)
		case html.EndTagToken:
			if !tokeniser.IsAnchor(z) {
				continue
			}
			current := anchors[len(anchors)-1]
			anchors = anchors[:len(anchors)-1]
			if len(anchors) == 0 {
				ss := sb.String()
				sb.Reset()
				link := Link{
					Href: current,
					Text: ss,
				}
				links = append(links, link)
			}
		case html.TextToken:
			if len(anchors) == 0 {
				continue
			}
			b := z.Text()
			trimmed := bytes.TrimSpace(b)
			sb.Write(trimmed)
		}
	}

	return links
}
