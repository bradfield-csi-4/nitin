package main

import (
	"./dns"
	"fmt"
	"log"
	"os"
	"syscall"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatal("Usage: go run dns.go [domain] [type] (e.g `google.com A`)")
	}

	googlePublicDNS := syscall.SockaddrInet4{
		Port: 53,
		Addr: [4]byte{8, 8, 8, 8},
	}

	// Construct the query
	query := dns.NewQuery(os.Args[1], os.Args[2])

	// Create the socket
	fd, e := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	check(e)

	// Bind to any available port
	e = syscall.Bind(fd, &syscall.SockaddrInet4{Port: 0, Addr: [4]byte{0, 0, 0, 0}})
	check(e)

	// Sendto DNS server
	e = syscall.Sendto(fd, dns.Serialize(query), 0, &googlePublicDNS)
	check(e)

	// Ignore messages from other hosts!

	// Receive/print the response
	out := make([]byte, 4096)
	for {
		_, from, e := syscall.Recvfrom(fd, out, 0)
		check(e)

		// Expect ipv4
		fromip4, ok := from.(*syscall.SockaddrInet4)
		if !ok {
			continue
		}

		// ignore responses from other hosts!
		if fromip4.Addr != googlePublicDNS.Addr || fromip4.Port != googlePublicDNS.Port {
			continue
		}

		response := dns.Deserialize(out)

		// ignore responses to *other* queries
		if !dns.QueryResponseMatch(query, response) {
			continue
		}

		fmt.Println(";; Got answer:")
		fmt.Print(response)
		break
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
