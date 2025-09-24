package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":42069")
	if err != nil {
		log.Fatal("error", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal("error", err)
	}

	defer conn.Close()

	bufferedReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := bufferedReader.ReadString('\n')
		if err != nil {
			log.Fatal("error", err)
		}
		_, err = conn.Write([]byte("Message recieved " + line))
		if err != nil {
			log.Fatal("error", err)
		}
	}
}
