package main

import (
	"flag"
	"github.com/apourchet/pingme/lib"
	"log"
)

var (
	port *string
)

func init() {
	port = flag.String("p", ping.DEFAULT_PORT, "Port to listen on")
	flag.Parse()
}

func main() {
	s := ping.NewServer(*port)
	log.Println("Listening on port " + *port)
	s.Serve()
}
