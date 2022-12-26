package yaml

import (
	"log"

	"github.com/seetohjinwei/gophercises/2-url-short/url"
	"gopkg.in/yaml.v3"
)

type T struct {
	Shorts []string `yaml:"shorts"`
	Full   string   `yaml:"full"`
}

func Parse(data []byte) []url.Url {
	ts := []T{}
	err := yaml.Unmarshal(data, &ts)
	if err != nil {
		log.Fatalf("error while un-marshalling yaml file: %v", err)
	}

	links := []url.Url{}

	for _, t := range ts {
		for _, short := range t.Shorts {
			link := url.Url{Short: short, Full: t.Full}
			links = append(links, link)
		}
	}

	return links
}
