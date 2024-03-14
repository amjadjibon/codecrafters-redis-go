package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

var db = make(map[string]string)

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
	case "set":
		if len(requestParts) < 7 {
			return []byte("-ERR invalid request\r\n")
		}

		db[requestParts[4]] = requestParts[6]
		return []byte("+OK\r\n")
	case "get":
		if len(requestParts) < 5 {
			return []byte("-ERR invalid request\r\n")
		}

		value, ok := db[requestParts[4]]
		if !ok {
			return []byte("$-1\r\n")
		}

		return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(value), value))
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
