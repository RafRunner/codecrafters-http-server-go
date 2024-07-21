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

	request, err := model.ReadHttpRequest(conn)
	var response *model.HttpResponse

	if err != nil {
		response = model.MakeResponse(model.INTERNAL_SERVER_ERROR, []byte(err.Error()))
	} else {
		verb, path := request.Verb, request.Path

		if verb == "GET" && path == "/" {
			response = model.MakeResponse(model.OK, []byte{})
		} else if verb == "GET" && strings.HasPrefix(path, "/echo/") {
			arg := path[6:]
			response = model.MakeResponse(model.OK, []byte(arg))
			response.AddHeader("Content-Type", "text/plain")
		} else {
			response = model.MakeResponse(model.NOT_FOUND, []byte{})
		}
	}

	_, err = conn.Write(response.WriteResponse())
	if err != nil {
		fmt.Println("Error writing response: ", err.Error())
		return
	}
}
