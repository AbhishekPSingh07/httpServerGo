package main

import (
	"fmt"
	request "httpServerGo/internal"
	"log"
	"net"
)

func main() {

	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("error", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("error", err)
		}

		request, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal("error", err)
		}

		fmt.Printf("Request Line\n")
		fmt.Printf("- Method: %s\n", request.RequestLine.Method)
		fmt.Printf("- Target: %s\n", request.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", request.RequestLine.HttpVersion)
	}

}
