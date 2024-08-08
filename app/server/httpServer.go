package server

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/model"
)

type HttpServer struct {
	ListenPort uint16
	endpoints  []endpoint
}

type HttpRequest struct {
	Request    model.HttpRequest
	PathParams map[string]string
}

type EndpointAction func(req HttpRequest) (*model.HttpResponse, error)

type endpoint struct {
	pathPattern *regexp.Regexp
	method      model.HttpVerb
	action      EndpointAction
}

func MakeServer(port uint16) *HttpServer {
	return &HttpServer{
		ListenPort: port,
		endpoints:  make([]endpoint, 0),
	}
}

func (s *HttpServer) Route(method model.HttpVerb, path string, action EndpointAction) *HttpServer {
	// Transform path into regex pattern
	pathPattern := "^" + path + "$"
	pathPattern = strings.Replace(pathPattern, "{", "(?P<", -1)
	pathPattern = strings.Replace(pathPattern, "}", ">[^/]+)", -1)
	regex := regexp.MustCompile(pathPattern)

	endpoint := endpoint{
		pathPattern: regex,
		method:      method,
		action:      action,
	}

	s.endpoints = append(s.endpoints, endpoint)
	return s
}

func (s *HttpServer) Listen(onReady func(), onError func(err error)) {
	l, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(int(s.ListenPort)))
	if err != nil {
		onError(err)
		return
	}
	defer l.Close()
	onReady()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}
		go handleClient(conn, s)
	}
}

func handleClient(conn net.Conn, server *HttpServer) {
	defer conn.Close()

	err := handleConnection(conn, server)
	if err != nil {
		fmt.Printf("Error handling connection: %v\n", err)
	}
}

func handleConnection(conn net.Conn, server *HttpServer) error {
	request, err := model.ReadHttpRequest(conn)
	var response *model.HttpResponse

	if err != nil {
		response = model.MakePlainTextResponse(model.BAD_REQUEST, err.Error())
	} else {
		matchedEndpoints := make([]endpoint, 0)

		for _, endpoint := range server.endpoints {
			if endpoint.pathPattern.MatchString(request.Path) {
				matchedEndpoints = append(matchedEndpoints, endpoint)
			}
		}

		if len(matchedEndpoints) > 0 {
			var matchedEndpoint *endpoint
			var pathParams map[string]string

			for _, endpoint := range matchedEndpoints {
				if endpoint.method == request.Method {
					matchedEndpoint = &endpoint
					matches := endpoint.pathPattern.FindStringSubmatch(request.Path)
					pathParams = extractPathParams(endpoint.pathPattern, matches)
					break
				}
			}

			if matchedEndpoint != nil {
				if response, err = matchedEndpoint.action(HttpRequest{
					Request:    *request,
					PathParams: pathParams,
				}); err != nil {
					response = model.MakePlainTextResponse(model.INTERNAL_SERVER_ERROR, err.Error())
				}
			} else {
				response = model.MakeResponse(model.METHOD_NOT_ALLOWED, []byte{})
			}
		} else {
			response = model.MakeResponse(model.NOT_FOUND, []byte{})
		}
	}

	if request != nil {
		response.CompressBody(*request)
	}
	if _, err = conn.Write(response.WriteResponse()); err != nil {
		return fmt.Errorf("error writing response: %w", err)
	}

	return nil
}

func extractPathParams(pattern *regexp.Regexp, matches []string) map[string]string {
	params := make(map[string]string)
	names := pattern.SubexpNames()
	for i, match := range matches {
		if i > 0 && i <= len(names) {
			params[names[i]] = match
		}
	}
	return params
}
