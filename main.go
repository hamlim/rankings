package main

import (
	"github.com/hamlim/rankings/server"
)

func main() {
	server.Create(8080, func(ctx server.Context) (server.Response, error) { // Update the function signature
		return server.Response{
			StatusCode: 200,
			Headers:    map[string]string{"Content-Type": "text/plain"},
			Body:       "Hello, World!",
		}, nil
	})
}
