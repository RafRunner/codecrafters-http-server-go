package model

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

func (v HttpVerb) IsValid() bool {
	switch v {
	case GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD:
		return true
	}
	return false
}
