package main

import "flag"

const EOT byte = 0x04

func main() {
	port := flag.Int("port", 8888, "The port where the server will listen.")
	flag.Parse()

	s := server{}
	s.startServer(*port)
}
