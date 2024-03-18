package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
)

type Entry struct {
	value      string
	expiration int64
}

var db = make(map[string]Entry)

func parseRequest(request string) []byte {
	request = strings.TrimPrefix(request, "\r\n")
	requestParts := strings.Split(request, "\r\n")

	if len(requestParts) < 3 {
		return []byte("-ERR invalid request\r\n")
	}

	switch strings.ToLower(requestParts[2]) {
	case "ping":
		return ping()
	case "echo":
		return echo(requestParts)
	case "set":
		return set(requestParts)
	case "get":
		return get(requestParts)
	case "info":
		return info(requestParts)
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
	var port = flag.String("port", "6379", "port to listen on")
	flag.Parse()

	var addr = fmt.Sprintf(":%s", *port)

	log.Println("start listening on", addr)

	listen, err := net.Listen("tcp", addr)
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
