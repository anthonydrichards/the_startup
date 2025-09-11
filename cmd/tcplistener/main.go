package main

import (
	"fmt"
	"log"
	"net"
	"the_startup/internal/request"
)

const port = ":42069"

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("could not start server: %s\n", err)
	}
	defer listener.Close()
	fmt.Printf("Listening for TCP connections on %s\n", port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("could not accept connection: %s\n", err)
			continue
		}
		fmt.Printf("Accepted connection from %s\n", conn.RemoteAddr().String())
		request, err := request.RequestFromReader(conn)
		if err != nil {
			log.Printf("could not parse request: %s\n", err)
			continue
		}
		fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n", request.RequestLine.Method, request.RequestLine.RequestTarget, request.RequestLine.HttpVersion)
		fmt.Printf("Headers:\n")
		for k, v := range request.Headers {
			fmt.Printf("- %s: %s\n", k, v)
		}
		fmt.Printf("Closed connection from %s\n", conn.RemoteAddr().String())
	}

}
