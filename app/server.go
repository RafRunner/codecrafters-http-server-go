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
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Connection accepted")

	request := make([]byte, 1024)

	_, err = conn.Read(request)
	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		os.Exit(1)
	}

	fmt.Println(string(request))

	var response []byte
	if strings.HasPrefix(string(request), "GET / ") {
		response = []byte("HTTP/1.1 200 OK\r\n\r\n")
	} else {
		response = []byte("HTTP/1.1 404 Not Found\r\n\r\n")
	}

	_, err = conn.Write(response)
	if err != nil {
		fmt.Println("Error writing response: ", err.Error())
		os.Exit(1)
	}
}
