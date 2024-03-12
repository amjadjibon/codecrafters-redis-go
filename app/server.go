package main

import (
	"log"
	"net"
	"strings"
)

func parseRequest(request string) []byte {
	request = strings.TrimPrefix(request, "\r\n")
	requestParts := strings.Split(request, "\r\n")

	if len(requestParts) < 3 {
		return []byte("-ERR invalid request\r\n")
	}

	switch strings.ToLower(requestParts[2]) {
	case "ping":
		return []byte("+PONG\r\n")
	case "echo":
		if len(requestParts) < 5 {
			return []byte("-ERR invalid request\r\n")
		}

		return []byte("+" + requestParts[4] + "\r\n")
	}

	return []byte("-ERR unknown command\r\n")
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	log.Println("new connection from", conn.RemoteAddr())

	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("received", n, "bytes:", string(buf[:n]))

		_, err = conn.Write(parseRequest(string(buf[:n])))
		if err != nil {
			log.Println(err)
			return
		}
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

		go handleConn(conn)
	}
}
