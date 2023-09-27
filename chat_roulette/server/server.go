package server

import (
	"io"
	"log"
	"net"
)

type Server struct {
	partner chan io.ReadWriteCloser
}

func New() Server {
	return Server{
		partner: make(chan io.ReadWriteCloser),
	}
}

func (s Server) cp(w io.Writer, r io.Reader, errc chan<- error) {
	_, err := io.Copy(w, r)
	errc <- err
}

func (s Server) chat(a, b io.ReadWriteCloser) {
	errc := make(chan error, 1)
	go s.cp(a, b, errc)
	go s.cp(b, a, errc)

	if err := <-errc; err != nil {
		log.Fatal(err)
	}

	a.Close()
	b.Close()
}

func (s Server) match(c io.ReadWriteCloser) {
	select {
	case s.partner <- c:
	case p := <-s.partner:
		s.chat(p, c)
	}
}

func (s Server) Run() {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go s.match(conn)
	}
}
