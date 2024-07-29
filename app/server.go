package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/codecrafters-io/http-server-starter-go/app/model"
	"github.com/codecrafters-io/http-server-starter-go/app/server"
)

func main() {
	var port uint16 = 4221
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

	svr := server.MakeServer(port)

	svr.Route(model.GET, "/", func(req server.HttpRequest) (*model.HttpResponse, error) {
		return model.MakeResponse(model.OK, []byte{}), nil
	})

	svr.Route(model.GET, "/echo/{echo}", func(req server.HttpRequest) (*model.HttpResponse, error) {
		return model.MakePlainTextResponse(model.OK, req.PathParams["echo"]), nil
	})

	svr.Route(model.GET, "/user-agent", func(req server.HttpRequest) (*model.HttpResponse, error) {
		userAgent := req.Request.Headers["user-agent"]

		if len(userAgent) > 0 {
			return model.MakePlainTextResponse(model.OK, userAgent[0].Value), nil
		} else {
			return model.MakePlainTextResponse(model.BAD_REQUEST, "User-Agent header not found"), nil
		}
	})

	svr.Route(model.GET, "/files/{fileName}", func(req server.HttpRequest) (*model.HttpResponse, error) {
		filePath := *directory + req.PathParams["fileName"]

		file, err := os.Stat(filePath)
		if err != nil || file.IsDir() {
			return model.MakeResponse(model.NOT_FOUND, []byte{}), nil
		} else {
			fileBytes, err := os.ReadFile(filePath)
			if err != nil {
				return model.MakePlainTextResponse(model.INTERNAL_SERVER_ERROR, err.Error()), nil
			} else {
				return model.MakeFileResponse(fileBytes), nil
			}
		}
	})

	svr.Route(model.POST, "/files/{fileName}", func(req server.HttpRequest) (*model.HttpResponse, error) {
		filePath := *directory + req.PathParams["fileName"]

		fileBytes := req.Request.Body

		err := os.WriteFile(filePath, fileBytes, 0644)
		if err != nil {
			return model.MakePlainTextResponse(model.INTERNAL_SERVER_ERROR, err.Error()), nil
		} else {
			return model.MakeResponse(model.CREATED, []byte{}), nil
		}
	})

	svr.Listen(func() {
		fmt.Printf("Server listening on port %d, file dir is: '%s'\n", port, *directory)
	}, func(err error) {
		fmt.Printf("Failed to bind to port %d: %v\n", svr.ListenPort, err)
		os.Exit(1)
	})
}
