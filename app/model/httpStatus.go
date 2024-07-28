package model

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
