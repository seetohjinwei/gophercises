package url

import (
	"io"
	"log"
	"net/http"
)

func Open(url string) io.Reader {
	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("error opening url %q: %v", url, err)
	}

	return res.Body
}
