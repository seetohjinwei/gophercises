package tokeniser

import (
	"bytes"
	"errors"

	"golang.org/x/net/html"
)

var anchor []byte = []byte("a")

func IsAnchor(z *html.Tokenizer) bool {
	tn, _ := z.TagName()
	return bytes.Equal(tn, anchor)
}

var href []byte = []byte("href")
var ErrHrefNotFound = errors.New("href not found")

func GetHref(z *html.Tokenizer) (string, error) {
	hasMore := true
	for hasMore {
		key, value, more := z.TagAttr()
		if bytes.Equal(key, href) {
			return string(value), nil
		}
		hasMore = more
	}

	return "", ErrHrefNotFound
}
