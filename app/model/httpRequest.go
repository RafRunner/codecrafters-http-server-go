package model

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
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
	if !verb.IsValid() {
		return nil, fmt.Errorf("unknown HTTP verb: %s", parts[0])
	}

	path := parts[1]

	version := HttpVersion(parts[2])
	if !version.IsValid() {
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
