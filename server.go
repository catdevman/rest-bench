package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strings"
)

func main() {
	ln, err := net.Listen("tcp", ":3000") // Listen on port 8080
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept() // Accept incoming connections
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}
func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// Read the request line
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading request line:", err)
		return
	}

	// Check if it's a POST request to /myroute
	if strings.HasPrefix(requestLine, "POST /") {
		// Read headers until an empty line is encountered
		for {
			headerLine, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading headers:", err)
				return
			}
			headerLine = strings.TrimSpace(headerLine)
			if headerLine == "" {
				break // End of headers
			}
		}

		// Read the request body
		var body bytes.Buffer
		_, err = io.Copy(&body, reader)
		if err != nil {
			fmt.Println("Error reading request body:", err)
			return
		}

		// Unmarshal the JSON body
		var data map[string]any
		err = json.Unmarshal(body.Bytes(), &data)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			return
		}

		// Send a response
		response := "HTTP/1.1 200 OK\r\nContent-Type: application/json\r\n\r\n{\"status\": \"ok\"}"
		conn.Write([]byte(response))
	} else {
		// Handle other requests or send a 404 Not Found
		response := "HTTP/1.1 404 Not Found\r\nContent-Type: text/plain\r\n\r\nNot Found"
		conn.Write([]byte(response))
	}
}
