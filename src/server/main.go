package main

import "flag"

func main() {
	port := flag.Int("port", 8888, "The port where the server will listen.")
	flag.Parse()

	s := newServer()
	s.StartServer(*port)
}
