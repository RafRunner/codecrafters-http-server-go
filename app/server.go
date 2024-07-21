package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/model"
)

func main() {
	port := 4221

	l, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(port))
	if err != nil {
		fmt.Printf("Failed to bind to port %d: %v", port, err)
		os.Exit(1)
	}
	defer l.Close()
	fmt.Printf("Server listening on port %d", port)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v", err)
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Connection accepted")

	err := handleConnection(conn)
	if err != nil {
		fmt.Printf("Error handling connection: %v\n", err)
	}
}

func handleConnection(conn net.Conn) error {
	request, err := model.ReadHttpRequest(conn)
	var response *model.HttpResponse

	if err != nil {
		response = model.MakePlainTextResponse(model.BAD_REQUEST, err.Error())
	} else {
		verb, path := request.Verb, request.Path

		if verb == model.GET && path == "/" {
			response = model.MakeResponse(model.OK, []byte{})
		} else if verb == model.GET && strings.HasPrefix(path, "/echo/") {
			arg := path[6:]
			response = model.MakePlainTextResponse(model.OK, arg)
		} else if verb == model.GET && path == "/user-agent" {
			userAgent := request.Headers["user-agent"]
			response = model.MakePlainTextResponse(model.OK, userAgent[0].Value)
		} else {
			response = model.MakeResponse(model.NOT_FOUND, []byte{})
		}
	}

	_, err = conn.Write(response.WriteResponse())
	if err != nil {
		return fmt.Errorf("error writing response: %w", err)
	}

	return nil
}
