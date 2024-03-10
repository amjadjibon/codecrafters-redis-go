package main

import (
	"log"
	"net"
)

func main() {
	var addr = "0.0.0.0:6379"

	log.Println("start listening on", addr)

	listen, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		log.Fatalln(err)
	}

	_, err = listen.Accept()
	if err != nil {
		log.Fatalln(err)
	}
}
