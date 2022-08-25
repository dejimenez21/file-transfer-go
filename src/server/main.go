package main

import (
	"flag"
	"log"
)

func main() {
	log.SetFlags(0)

	port := flag.Int("port", 8888, "The port where the server will listen.")
	flag.Parse()

	s := newServer()
	s.StartServer(*port)
}
