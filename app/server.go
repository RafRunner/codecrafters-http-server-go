package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

func main() {
	port := 4221

	l, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(port))
	if err != nil {
		fmt.Printf("Failed to bind to port %d", port)
		os.Exit(1)
	}
	fmt.Printf("Server listening on port %d", port)

	_, err = l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	fmt.Print("Connection accepted")
}
