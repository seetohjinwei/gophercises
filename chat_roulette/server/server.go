package server

import (
	"io"
	"log"
	"net"
)

type Server struct {
	value string
}

func (s Server) Run() {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	var prevConn net.Conn = nil

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		if prevConn != nil {
			go io.Copy(prevConn, conn)
			go io.Copy(conn, prevConn)

			prevConn = nil
		} else {
			prevConn = conn
		}
	}
}
