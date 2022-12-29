package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"

	"github.com/seetohjinwei/gophercises/14-recover-middleware/rm"
)

const dev = true

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/", sourceCodeHandler)
	mux.HandleFunc("/panic/", panicDemo)
	mux.HandleFunc("/panic-after/", panicAfterDemo)
	mux.HandleFunc("/", hello)
	log.Fatal(http.ListenAndServe(":3000", recoverMiddleware(mux, dev)))
}

func sourceCodeHandler(w http.ResponseWriter, r *http.Request) {
	path := r.FormValue("path")
	f, err := os.Open(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	line := r.FormValue("line")
	lineInt, err := strconv.Atoi(line)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rm.Lex(w, f, lineInt, lineInt)

	// b := bytes.NewBuffer(nil)
	// _, err = io.Copy(b, f)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// quick.Highlight(w, b.String(), "go", "html", "monokai")
}

func recoverMiddleware(h http.Handler, dev bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
				stack := debug.Stack()

				if !dev {
					http.Error(w, "Something went wrong :/", http.StatusInternalServerError)
				} else {
					w.WriteHeader(http.StatusInternalServerError)
					rm.WriteStack(w, stack)
					fmt.Fprintf(w, "<h1>panic: %v</h1><pre>%s</pre>", r, stack)
				}
			}
		}()

		h.ServeHTTP(w, r)

		/*
			nw := &responseWriter{ResponseWriter: w}
			h.ServeHTTP(nw, r)
			nw.flush()
		*/
	}
}

/*

// Wrapping the interface to avoid writing to the original ResponseWriter until it is confirmed that there is no error!
// ...but this has some disadvantages, namely: other functions might rely on checking the type to validate certain
// behaviour. However, by wrapping it, we lose this ability and might mislead other code.
//
// e.g. the original ResponseWriter does not implement the Hijacker, but now, it seems like it does.
// e.g. the wrapped Flush implementation cannot return an error to indicate that the interface is not implemented
type responseWriter struct {
	http.ResponseWriter
	writes [][]byte
	status int
}

func (w *responseWriter) Write(bytes []byte) (int, error) {
	w.writes = append(w.writes, bytes)

	return len(bytes), nil
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
}

var ErrNotHijacker = errors.New("the ResponseWriter does not implement a Hijacker interface")

func (w *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, ErrNotHijacker
	}

	return hijacker.Hijack()
}

func (w *responseWriter) Flush() {
	flusher, ok := w.ResponseWriter.(http.Flusher)
	if !ok {
		return
	}

	flusher.Flush()
}

func (w *responseWriter) flush() error {
	if w.status != 0 {
		w.ResponseWriter.WriteHeader(w.status)
	}

	for _, write := range w.writes {
		_, err := w.ResponseWriter.Write(write)
		if err != nil {
			return err
		}
	}

	return nil
}

*/

func panicDemo(w http.ResponseWriter, r *http.Request) {
	funcThatPanics()
}

func panicAfterDemo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello!</h1>")
	funcThatPanics()
}

func funcThatPanics() {
	panic("Oh no!")
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>Hello!</h1>")
}
