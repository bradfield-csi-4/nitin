package main

import (
	"fmt"
	"syscall"
)

const ListenPort int = 1234
const ForwardPort int = 2345

func main() {
	// Create TCP socket (for listening)
	ld, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		panic(err)
	}

	// Bind listening socket
	err = syscall.Bind(ld, &syscall.SockaddrInet4{Port: ListenPort, Addr: [4]byte{0, 0, 0, 0}})
	if err != nil {
		panic(err)
	}

	// Set max connection limit and start listening on socket
	err = syscall.Listen(ld, 3)
	if err != nil {
		panic(err)
	}
	fmt.Printf("TCP Server listening for connections on port: %v\n", ListenPort)

	// Create TCP socket (for forwarding)
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		panic(err)
	}

	// Bind forwarding socket to any available port
	err = syscall.Bind(fd, &syscall.SockaddrInet4{Port: 0, Addr: [4]byte{0, 0, 0, 0}})
	if err != nil {
		panic(err)
	}

	destAddr := &syscall.SockaddrInet4{Port: ForwardPort, Addr: [4]byte{0, 0, 0, 0}}
	err = syscall.Connect(fd, destAddr)
	if err != nil {
		return
	}

	for {
		// Accepts connection request and creates new socket for it
		nfd, sa, err := syscall.Accept(ld)
		if err != nil {
			panic(err)
		}
		fmt.Printf("TCP Server connected to remote client: %v\n", sa)

		buf := make([]byte, 10)

		for {
			// Read message from socket
			size, _, err := syscall.Recvfrom(nfd, buf, 0)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Received message of size: %v\n", size)
			if size == 0 {
				break
			}

			// Send message back to forwarding socket
			err = syscall.Sendto(fd, buf, 0, destAddr)
			if err != nil {
				panic(err)
			}
		}

		fmt.Printf("TCP Server disconnected from remote client: %v\n", sa)

	}
}
