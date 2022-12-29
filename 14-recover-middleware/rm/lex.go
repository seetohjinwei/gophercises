package rm

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

var lexer = lexers.Get("go")
var style = styles.Get("monokai")

func Lex(w http.ResponseWriter, r io.Reader, lineStart, lineEnd int) {
	contents, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatalf("error lexing: %v", err)
	}
	iterator, err := lexer.Tokenise(nil, string(contents))
	if err != nil {
		log.Fatalf("error lexing: %v", err)
	}

	w.Header().Set("Content-Type", "text/html")

	lines := [][2]int{{lineStart, lineEnd}}
	formatter := html.New(html.TabWidth(2), html.WithLineNumbers(true), html.LineNumbersInTable(true), html.HighlightLines(lines))
	err = formatter.Format(w, style, iterator)
	if err != nil {
		log.Fatalf("error lexing: %v", err)
	}
}
