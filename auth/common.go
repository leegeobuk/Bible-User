package auth

import "errors"

const (
	origin = "http://localhost:3000"
)

var (
	headers = map[string]string{
		"Access-Control-Allow-Credentials": "true",
		"Access-Control-Allow-Headers":     "Content-Type",
		"Access-Control-Allow-Origin":      origin,
	}
	errAccountExist    = errors.New("error account already exists")
	errAccountNotExist = errors.New("error account doesn't exist")
)

func addHeaders(m map[string]string, h map[string]string) {
	for k, v := range h {
		m[k] = v
	}
}
