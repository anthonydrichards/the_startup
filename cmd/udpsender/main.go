package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const updAddress = "localhost:42069"

func main() {
	addr, err := net.ResolveUDPAddr("udp", updAddress)
	if err != nil {
		log.Fatalf("could not resolve address: %s\n", err)
		panic(err)
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("could not dial UDP: %s\n", err)
	}
	defer conn.Close()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("could not read from stdin: %s\n", err)
			continue
		}
		_, err = conn.Write([]byte(line))
		if err != nil {
			log.Printf("could not send UDP packet: %s\n", err)
		}
	}
}
