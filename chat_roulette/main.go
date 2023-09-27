package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/seetohjinwei/gophercises/chat_roulette/client"
	"github.com/seetohjinwei/gophercises/chat_roulette/server"
)

var modeFlag = flag.String("mode", "client", "client or server, defaults to client")

type Application interface {
	Run()
}

func main() {
	flag.Parse()
	fmt.Println("mode:", *modeFlag)

	var app Application

	switch *modeFlag {
	case "server":
		app = server.Server{}
	case "client":
		app = client.Client{}
	default:
		os.Exit(1)
	}

	app.Run()
}
