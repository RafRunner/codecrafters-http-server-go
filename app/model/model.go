package model

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// HeaderVal represents a single header's key and value.
type HeaderVal struct {
	OriginalKey string
	Value       string
}

func MakeHeader(key, val string) *HeaderVal {
	return &HeaderVal{
		OriginalKey: key,
		Value:       val,
	}
}

func ReadHeaderLine(line string) (*HeaderVal, error) {
	parts := strings.SplitN(strings.TrimSpace(line), ":", 2)
	if len(parts) != 2 {
		return nil, errors.New("header line should have two parts separated by ':'")
	}

	return &HeaderVal{
		OriginalKey: strings.TrimSpace(parts[0]),
		Value:       strings.TrimSpace(parts[1]),
	}, nil
}

// HttpVerb represents an HTTP verb/method.
type HttpVerb string

const (
	GET     HttpVerb = "GET"
	POST    HttpVerb = "POST"
	PUT     HttpVerb = "PUT"
	DELETE  HttpVerb = "DELETE"
	PATCH   HttpVerb = "PATCH"
	OPTIONS HttpVerb = "OPTIONS"
	HEAD    HttpVerb = "HEAD"
)

// HttpVersion represents the version of the HTTP protocol.
type HttpVersion string

const (
	HTTP10 HttpVersion = "HTTP/1.0"
	HTTP11 HttpVersion = "HTTP/1.1"
	HTTP20 HttpVersion = "HTTP/2.0"
)

// HttpRequest represents an HTTP request.
type HttpRequest struct {
	Verb    HttpVerb
	Path    string
	Version HttpVersion
	Headers map[string][]HeaderVal
	Body    []byte
}

func ReadHttpRequest(conn net.Conn) (*HttpRequest, error) {
	reader := bufio.NewReader(conn)

	requestLine, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading request line: %w", err)
	}

	parts := strings.SplitN(strings.TrimSpace(requestLine), " ", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("expected request line to have 3 parts")
	}

	verb := HttpVerb(parts[0])
	switch verb {
	case GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD:
		// valid verb
	default:
		return nil, fmt.Errorf("unknown HTTP verb: %s", parts[0])
	}

	path := parts[1]

	version := HttpVersion(parts[2])
	switch version {
	case HTTP10, HTTP11, HTTP20:
		// valid version
	default:
		return nil, fmt.Errorf("unknown HTTP version: %s", parts[2])
	}

	headers := make(map[string][]HeaderVal)
	for {
		headerLine, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("error reading header line: %w", err)
		}

		headerTrimmed := strings.TrimSpace(headerLine)
		if headerTrimmed == "" {
			break
		}

		header, err := ReadHeaderLine(headerLine)
		if err != nil {
			return nil, err
		}

		lowerKey := strings.ToLower(header.OriginalKey)
		headers[lowerKey] = append(headers[lowerKey], *header)
	}

	var body []byte
	contentLength := headers["content-length"]
	if len(contentLength) > 0 {
		toRead, err := strconv.Atoi(contentLength[0].Value)
		if err != nil {
			return nil, fmt.Errorf("content-length contains invalid number")
		}
		body = make([]byte, toRead)

		_, err = reader.Read(body)
		if err != nil {
			return nil, fmt.Errorf("error reading body: %w", err)
		}
	} else {
		body = []byte{}
	}

	return &HttpRequest{
		Verb:    verb,
		Path:    path,
		Version: version,
		Headers: headers,
		Body:    body,
	}, nil
}

// HttpStatus represents an HTTP status code.
type HttpStatus int

const (
	OK                    HttpStatus = iota + 200
	CREATED               HttpStatus = 201
	NO_CONTENT            HttpStatus = 204
	BAD_REQUEST           HttpStatus = 400
	NOT_FOUND             HttpStatus = 404
	METHOD_NOT_ALLOWED    HttpStatus = 405
	INTERNAL_SERVER_ERROR HttpStatus = 500
)

func (s *HttpStatus) GetReasonPhrase() string {
	switch *s {
	case OK:
		return "OK"
	case CREATED:
		return "Created"
	case NO_CONTENT:
		return "No Content"
	case BAD_REQUEST:
		return "Bad Request"
	case NOT_FOUND:
		return "Not Found"
	case METHOD_NOT_ALLOWED:
		return "Method Not Allowed"
	case INTERNAL_SERVER_ERROR:
		return "Internal Server Error"
	default:
		return "Unknown Status"
	}
}

// HttpResponse represents an HTTP response.
type HttpResponse struct {
	Version HttpVersion
	Status  HttpStatus
	Headers map[string][]HeaderVal
	Body    []byte
}

func MakeResponse(status HttpStatus, body []byte) *HttpResponse {
	headers := make(map[string][]HeaderVal)
	response := &HttpResponse{
		Version: HTTP11,
		Status:  status,
		Headers: headers,
		Body:    body,
	}

	if len(body) > 0 {
		response.AddHeader("Content-Length", strconv.Itoa(len(body)))
	}

	return response
}

func MakePlainTextResponse(status HttpStatus, body string) *HttpResponse {
	response := MakeResponse(status, []byte(body))
	response.AddHeader("Content-Type", "text/plain")

	return response
}

func MakeFileResponse(file []byte) *HttpResponse {
	response := MakeResponse(OK, file)
	response.AddHeader("Content-Type", "application/octet-stream")

	return response
}

// AddHeader adds a header to the HttpResponse.
func (r *HttpResponse) AddHeader(key, val string) {
	r.Headers[key] = append(r.Headers[key], *MakeHeader(key, val))
}

// WriteResponse generates the HTTP response as a byte array.
func (r *HttpResponse) WriteResponse() []byte {
	responseLine := fmt.Sprintf("%s %d %s\r\n", r.Version, r.Status, r.Status.GetReasonPhrase())

	// Add headers to the response
	for _, header := range r.Headers {
		for _, val := range header {
			responseLine += fmt.Sprintf("%s: %s\r\n", val.OriginalKey, val.Value)
		}
	}
	responseLine += "\r\n"

	// Append the body to the response
	return append([]byte(responseLine), r.Body...)
}

func (r *HttpResponse) CompressBody(request HttpRequest) {
	accepted := request.Headers["accept-encoding"]

	if len(accepted) == 0 {
		return
	}
	supportedAlgs := make([]string, 0)

	for _, alg := range accepted {
		supportedAlgs = append(supportedAlgs, alg.Value)
	}

	if contains(supportedAlgs, "gzip") {
		r.AddHeader("Content-Encoding", "gzip")
	}
}

func contains[T comparable](slice []T, element T) bool {
	for _, it := range slice {
		if it == element {
			return true
		}
	}
	return false
}
