package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
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

		entry := Entry{value: requestParts[6]}
		if len(requestParts) > 8 {

			if requestParts[8] != "px" {
				return []byte("-ERR invalid request\r\n")
			}

			ttl, err := time.ParseDuration(requestParts[10] + "ms")
			if err != nil {
				return []byte("-ERR invalid request\r\n")
			}

			entry.expiration = time.Now().Add(ttl).UnixNano()
		}

		db[requestParts[4]] = entry
		return []byte("+OK\r\n")
	case "get":
		if len(requestParts) < 5 {
			return []byte("-ERR invalid request\r\n")
		}

		entry, ok := db[requestParts[4]]
		if !ok {
			return []byte("$-1\r\n")
		}

		if entry.expiration > 0 && entry.expiration < time.Now().UnixNano() {
			delete(db, requestParts[4])
			return []byte("$-1\r\n")
		}

		return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(entry.value), entry.value))
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
