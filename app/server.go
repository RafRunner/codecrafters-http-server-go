package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/model"
)

func main() {
	port := 4221
	directory := flag.String("directory", "", "Directory where to look for files")
	flag.Parse()

	if *directory != "" {
		stat, err := os.Stat(*directory)
		if err != nil {
			fmt.Printf("Error accessing directory provided: %v\n", err)
			os.Exit(1)
		}
		if !stat.IsDir() {
			fmt.Println("Directory provided is not a directory")
			os.Exit(1)
		}
	}

	l, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(port))
	if err != nil {
		fmt.Printf("Failed to bind to port %d: %v\n", port, err)
		os.Exit(1)
	}
	defer l.Close()
	fmt.Printf("Server listening on port %d\n", port)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}
		go handleClient(conn, *directory)
	}
}

func handleClient(conn net.Conn, fileDir string) {
	defer conn.Close()
	fmt.Println("Connection accepted")

	err := handleConnection(conn, fileDir)
	if err != nil {
		fmt.Printf("Error handling connection: %v\n", err)
	}
}

func handleConnection(conn net.Conn, fileDir string) error {
	request, err := model.ReadHttpRequest(conn)
	var response *model.HttpResponse

	if err != nil {
		response = model.MakePlainTextResponse(model.BAD_REQUEST, err.Error())
	} else {
		verb, path := request.Verb, request.Path

		if verb == model.GET && path == "/" {
			response = model.MakeResponse(model.OK, []byte{})
		} else if verb == model.GET && strings.HasPrefix(path, "/echo/") {
			echo := path[6:]
			response = model.MakePlainTextResponse(model.OK, echo)
		} else if verb == model.GET && path == "/user-agent" {
			userAgent := request.Headers["user-agent"]

			if len(userAgent) > 0 {
				response = model.MakePlainTextResponse(model.OK, userAgent[0].Value)
			} else {
				response = model.MakePlainTextResponse(model.BAD_REQUEST, "User-Agent header not found")
			}
		} else if verb == model.GET && strings.HasPrefix(path, "/files/") && fileDir != "" {
			fileName := path[7:]
			filePath := fileDir + fileName

			file, err := os.Stat(filePath)
			if err != nil || file.IsDir() {
				response = model.MakeResponse(model.NOT_FOUND, []byte{})
			} else {
				fileBytes, err := os.ReadFile(filePath)
				if err != nil {
					response = model.MakePlainTextResponse(model.INTERNAL_SERVER_ERROR, err.Error())
				} else {
					response = model.MakeFileResponse(fileBytes)
				}
			}
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
