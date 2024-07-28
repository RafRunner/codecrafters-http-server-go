package model

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"strconv"
	"strings"
)

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
	lowerKey := strings.ToLower(key)
	r.Headers[lowerKey] = append(r.Headers[lowerKey], *MakeHeader(key, val))
}

// SetHeader adds or alters a existing header to a value
func (r *HttpResponse) SetHeader(key, val string) {
	lowerKey := strings.ToLower(key)
	existing := r.Headers[lowerKey]

	if len(existing) == 0 {
		r.AddHeader(key, val)
	} else {
		existing[0].Value = val
	}
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
	if len(r.Body) == 0 {
		return
	}

	accepted := request.Headers["accept-encoding"]

	if len(accepted) == 0 {
		return
	}
	supportedAlgs := make([]string, 0)

	for _, acceptedHeader := range accepted {
		for _, alg := range strings.Split(acceptedHeader.Value, ",") {
			supportedAlgs = append(supportedAlgs, strings.TrimSpace(alg))
		}
	}

	if contains(supportedAlgs, "gzip") {
		// Compress the body using gzip
		compressedBody, err := compressGzip(r.Body)
		if err != nil {
			fmt.Printf("Failed to compress body: %v\n", err)
			return
		}

		r.Body = compressedBody
		r.SetHeader("Content-Length", strconv.Itoa(len(r.Body)))
		r.AddHeader("Content-Encoding", "gzip")
	}
}

func compressGzip(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)

	_, err := writer.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed to write data to gzip writer: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close gzip writer: %w", err)
	}

	return buf.Bytes(), nil
}

func contains[T comparable](slice []T, element T) bool {
	for _, it := range slice {
		if it == element {
			return true
		}
	}
	return false
}
