package files

import (
	"io"
	"log"
	"os"
)

func Open(path string) io.Reader {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("error opening %q: %v", path, err)
	}

	return f
}
