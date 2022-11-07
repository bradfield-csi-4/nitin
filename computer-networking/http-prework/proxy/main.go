package main

import (
	"fmt"
	"strings"
	"syscall"
)

const ListeningPort int = 80
const DestinationPort int = 9000
const BufferSize int = 100
const CachePath = "website/"

var cache = make(map[string][]byte)

func main() {
	serverSocket := startTCPServer()
	destinationSocket := connectToDestinationServer()

	for {
		// Accepts connection and creates new socket
		connectionSocket, clientAddress, err := syscall.Accept(serverSocket)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Connected to remote client: %v\n", clientAddress)

		// Reads request from client
		req := readFrom(connectionSocket)
		if len(req) == 0 {
			continue
		}
		path := strings.Split(strings.Split(string(req), "\r\n")[0], " ")[1]

		// Checks if response is available in cache
		var resp []byte
		if cachedResp, exists := cache[path]; exists {
			// Use cached response
			resp = cachedResp
		} else {
			// Forward request to destination server
			err = syscall.Sendto(destinationSocket, req, 0, destServerAddr())
			if err != nil {
				panic(err)
			}
			resp = readFrom(destinationSocket)
		}

		// Return response to client
		err = syscall.Sendto(connectionSocket, resp, 0, clientAddress)
		if err != nil {
			panic(err)
		}

		// Cache response
		if strings.Contains(path, CachePath) {
			cache[path] = resp
		}

		_ = syscall.Close(connectionSocket)
		fmt.Printf("Disconnected from remote client: %v\n", clientAddress)
	}
	//_ = syscall.Close(serverSocket)
	//_ = syscall.Close(destinationSocket)
	//fmt.Printf("Stopping server")
}

func readFrom(socket int) []byte {
	var req []byte
	buf := make([]byte, BufferSize)
	for {
		// Read message from socket
		numBytes, _, err := syscall.Recvfrom(socket, buf, 0)
		if err != nil {
			panic(err)
		}
		req = append(req, buf[:numBytes]...)
		// Reached end of the response
		if numBytes < BufferSize {
			break
		}
	}
	return req
}

func localhost() [4]byte {
	return [4]byte{0, 0, 0, 0}
}

func connectToDestinationServer() int {
	// Create TCP socket (for forwarding)
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		panic(err)
	}

	// Bind forwarding socket to any available port
	err = syscall.Bind(fd, &syscall.SockaddrInet4{Port: 0, Addr: localhost()})
	if err != nil {
		panic(err)
	}

	err = syscall.Connect(fd, destServerAddr())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Connected to destination server on port: %v\n", DestinationPort)
	return fd
}

func destServerAddr() *syscall.SockaddrInet4 {
	return &syscall.SockaddrInet4{Port: DestinationPort, Addr: localhost()}
}

func startTCPServer() int {
	// Create TCP socket (for listening)
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		panic(err)
	}

	// Bind listening socket
	err = syscall.Bind(fd, &syscall.SockaddrInet4{Port: ListeningPort, Addr: localhost()})
	if err != nil {
		panic(err)
	}

	// Set max connection limit and start listening on socket
	err = syscall.Listen(fd, 3)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Server listening for connections on port: %v\n", ListeningPort)
	return fd
}
