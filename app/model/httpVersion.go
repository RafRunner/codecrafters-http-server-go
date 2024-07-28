package model

// HttpVersion represents the version of the HTTP protocol.
type HttpVersion string

const (
	HTTP10 HttpVersion = "HTTP/1.0"
	HTTP11 HttpVersion = "HTTP/1.1"
	HTTP20 HttpVersion = "HTTP/2.0"
)

func (v HttpVersion) IsValid() bool {
	switch v {
	case HTTP10, HTTP11, HTTP20:
		return true
	default:
		return false
	}
}
