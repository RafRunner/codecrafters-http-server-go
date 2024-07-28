package model

import (
	"errors"
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
