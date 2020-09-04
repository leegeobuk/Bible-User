package auth

const (
	origin = "http://localhost:3000"
)

var (
	headers = map[string]string{
		"Access-Control-Allow-Credentials": "true",
		"Access-Control-Allow-Headers":     "Content-Type",
		"Access-Control-Allow-Origin":      origin,
	}
)
