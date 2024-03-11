package main

import (
	"log"
	"net"
)

func handleConn(conn net.Conn) {
	defer conn.Close()

	log.Println("new connection from", conn.RemoteAddr())

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("received", n, "bytes:", string(buf[:n]))

	_, err = conn.Write([]byte("+PONG\r\n"))
	if err != nil {
		log.Println(err)
		return
	}
}

func main() {
	var addr = "0.0.0.0:6379"

	log.Println("start listening on", addr)

	listen, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		log.Fatalln(err)
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		handleConn(conn)
	}
}
