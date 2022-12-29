package rm

import (
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
)

type StackFile struct {
	Path      string
	Name      string
	Line      int
	Character int
}

var lineRegex = regexp.MustCompile(`(?m)^\s+(.*):(\d+) \+0x([\da-z]+)$`)

func ParseStackTrace(stack []byte) []StackFile {
	matches := lineRegex.FindAllSubmatch(stack, -1)

	files := make([]StackFile, len(matches))
	for i, match := range matches {
		path := string(match[1])

		file := StackFile{
			Path:      path,
			Name:      filepath.Base(path),
			Line:      mustByteToInt(match[2]),
			Character: mustByteToInt(match[3]),
		}
		files[i] = file
	}

	return files
}

func mustByteToInt(bytes []byte) int {
	result, _ := strconv.Atoi(string(bytes))
	return result
}

func WriteStack(w http.ResponseWriter, stack []byte) {
	files := ParseStackTrace(stack)

	fmt.Fprintf(w, `<ul>`)

	for _, f := range files {
		fmt.Fprintf(w, `<li><a href="/debug?path=%s&line=%d">%s</a></li>`, f.Path, f.Line, f.Name)
	}

	fmt.Fprintf(w, `</ul>`)
}
