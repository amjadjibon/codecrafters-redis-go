package main

import (
	"fmt"
	"time"
)

func ping() []byte {
	return []byte("+PONG\r\n")
}

func echo(requestParts []string) []byte {
	if len(requestParts) < 5 {
		return []byte("-ERR invalid request\r\n")
	}

	return []byte("+" + requestParts[4] + "\r\n")
}

func set(requestParts []string) []byte {
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
}

func get(requestParts []string) []byte {
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

func info(requestParts []string) []byte {
	if len(requestParts) < 5 {
		return []byte("-ERR invalid request\r\n")
	}

	switch requestParts[4] {
	case "replication":
		return []byte("+role:master\r\n")
	default:
		return []byte("-ERR invalid request\r\n")
	}
}
