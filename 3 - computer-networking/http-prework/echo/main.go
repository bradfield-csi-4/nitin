package main

import (
	"fmt"
	"syscall"
)

const PORT int = 1234

func main() {
	// Create TCP socket
	listeningSocket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		panic(err)
	}

	// Bind socket to any available port
	err = syscall.Bind(listeningSocket, &syscall.SockaddrInet4{Port: PORT, Addr: [4]byte{0, 0, 0, 0}})
	if err != nil {
		panic(err)
	}

	// Set max connection limit and start listening on socket
	err = syscall.Listen(listeningSocket, 3)
	if err != nil {
		panic(err)
	}
	fmt.Printf("TCP Server listening for connections on port: %v\n", PORT)

	for {
		// Accepts connection request and creates new socket for it
		socket, addr, err := syscall.Accept(listeningSocket)
		if err != nil {
			panic(err)
		}

		fmt.Printf("TCP Server connected to remote client: %v\n", addr)

		buf := make([]byte, 10)

		for {
			// Read message from socket
			numBytes, _, err := syscall.Recvfrom(socket, buf, 0)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Received message of size: %v\n", numBytes)
			if numBytes == 0 {
				break
			}

			// Write message back to socket
			err = syscall.Sendto(socket, buf, 0, addr)
			if err != nil {
				panic(err)
			}
		}
		err = syscall.Close(socket)
		if err != nil {
			panic(err)
		}
		fmt.Printf("TCP Server disconnected from remote client: %v\n", addr)

	}
}
