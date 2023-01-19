package main

import (
	"./dns"
	"fmt"
	"log"
	"net"
	"os"
)

const (
	GooglePublicDNS = "8.8.8.8:53"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatal("Usage: go run dns.go [domain] [type] (e.g `google.com A`)")
	}

	// Construct the query
	query := dns.NewQuery(os.Args[1], os.Args[2])

	// Send the question
	conn, e := net.Dial("udp", GooglePublicDNS)
	check(e)
	_, e = conn.Write(dns.Serialize(query))
	check(e)

	// Receive/print the response
	out := make([]byte, 4096)
	for {
		conn.Read(out)
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
