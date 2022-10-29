package main

import (
	"fmt"
	"syscall"
)

const PORT int = 1234

func main() {
	// Create TCP socket
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		panic(err)
	}

	// Bind socket to any available port
	err = syscall.Bind(fd, &syscall.SockaddrInet4{Port: PORT, Addr: [4]byte{0, 0, 0, 0}})
	if err != nil {
		panic(err)
	}

	// Set max connection limit and start listening on socket
	err = syscall.Listen(fd, 3)
	if err != nil {
		panic(err)
	}
	fmt.Printf("TCP Server listening for connections on port: %v\n", PORT)

	for {
		// Accepts connection request and creates new socket for it
		nfd, sa, err := syscall.Accept(fd)
		if err != nil {
			panic(err)
		}

		fmt.Printf("TCP Server connected to remote client: %v\n", sa)

		for {
			// Read message from socket
			msg := make([]byte, 10)
			size, _, err := syscall.Recvfrom(nfd, msg, 0)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Received message of size: %v\n", size)
			if size == 0 {
				break
			}

			// Write message back to socket
			err = syscall.Sendto(nfd, msg, 0, sa)
			if err != nil {
				panic(err)
			}
		}

		fmt.Printf("TCP Server disconnected from remote client: %v\n", sa)

	}
}
