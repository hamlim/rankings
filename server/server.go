package server

import (
	"fmt"
	"net"
	"strings"
)

type Context struct {
	headers map[string]string
	method  string
	path    string
	version string
	body    string
}

type Response struct {
	StatusCode int
	Headers    map[string]string
	Body       string
}

func Create(port int, handleRequest func(ctx Context) (Response, error)) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Server is listening on port 8080")

	// accept incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		// handle the connection in a separate goroutine
		go handleConnection(conn, handleRequest)
	}
}

func handleConnection(conn net.Conn, handleRequest func(ctx Context) (Response, error)){
	defer conn.Close()

	fmt.Println("New connection from:", conn.RemoteAddr())
	data := make([]byte, 1024)
	_, err := conn.Read(data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// split the data by new lines
	requestContext := strings.Split(string(data), "\r\n\r\n")
	headerString := requestContext[0]
	body := requestContext[1]

	headers := strings.Split(headerString, "\r\n"	)

	meta := strings.Split(headers[0], " ")
	method := meta[0]
	path := meta[1]
	version := meta[2]

	// take an array of "key: value" pairs and convert it into a map
	headerMap := make(map[string]string)
	for _, header := range headers[1:] {
		headerParts := strings.Split(header, ": ")
		headerMap[headerParts[0]] = headerParts[1]
	}

	ctx := Context{
		headers: headerMap,
		method:  method,
		path:    path,
		version: version,
		body:    body,
	}

	response, err := handleRequest(ctx)

	if (err != nil) {
		conn.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
		conn.Write([]byte("Content-Type: text/plain\r\n"))
		conn.Write([]byte("\r\n"))
		conn.Write([]byte(err.Error()))
		return
	}

	conn.Write([]byte(fmt.Sprintf("HTTP/1.1 %d OK\r\n", response.StatusCode)))
	for key, value := range response.Headers {
		conn.Write([]byte(fmt.Sprintf("%s: %s\r\n", key, value)))
	}
	conn.Write([]byte("\r\n"))
	conn.Write([]byte(response.Body))
}

