package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	port := 4221

	l, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(port))
	if err != nil {
		fmt.Println("Failed to bind to port ", port)
		os.Exit(1)
	}
	fmt.Println("Server listening on port ", port)

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		return
	}
	defer conn.Close()

	fmt.Println("Connection accepted")

	request := make([]byte, 1024)

	_, err = conn.Read(request)
	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		return
	}

	req := strings.Split(string(request), "\r\n")
	parts := strings.Split(req[0], " ")

	if len(parts) != 3 {
		fmt.Println("Malformed request. Expected 3 parts on first line: verb, path and protocol version")
		return
	}

	verb, path := parts[0], parts[1]

	var response string
	if verb == "GET" && path == "/" {
		response = "HTTP/1.1 200 OK\r\n\r\n"
	} else if verb == "GET" && strings.HasPrefix(path, "/echo/") {
		arg := path[6:]
		response = "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + strconv.Itoa(len(arg)) + "\r\n\r\n" + arg
	} else {
		response = "HTTP/1.1 404 Not Found\r\n\r\n"
	}

	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing response: ", err.Error())
		return
	}
}
