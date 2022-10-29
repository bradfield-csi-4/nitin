package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strings"
	"syscall"
)

// See RFC 1035 § 3.2.2 for a full list of types
var qtypes = map[string]int{
	"A":     1,
	"NS":    2,
	"CNAME": 5,
	"SOA":   6,
	"MX":    15,
	"TXT":   16,
}

func main() {
	if len(os.Args) != 3 {
		log.Fatal("Usage: go run dns.go [domain] [type] (e.g `google.com A`)")
	}

	hostname := os.Args[1]
	qtype := os.Args[2]

	// Initialize query
	query := make([]byte, 0)

	// Append identifier
	id := uint16(rand.Intn(0xffff))
	query = binary.BigEndian.AppendUint16(query, id)

	// Append flags
	var flags uint16
	rd := uint16(math.Pow(2, 15)) >> 7 // RD (Recursion Desired)
	flags = flags | rd
	query = binary.BigEndian.AppendUint16(query, flags)

	// Append counts
	query = binary.BigEndian.AppendUint16(query, 1) // QDCOUNT
	query = binary.BigEndian.AppendUint16(query, 0) // ANCOUNT
	query = binary.BigEndian.AppendUint16(query, 0) // NSCOUNT
	query = binary.BigEndian.AppendUint16(query, 0) // ARCOUNT

	// Append question section
	// Append qname
	for _, label := range strings.Split(hostname, ".") {
		query = append(query, byte(len(label)))
		query = append(query, []byte(label)...)
	}
	query = append(query, 0)

	// Append qtype & qclass
	query = binary.BigEndian.AppendUint16(query, uint16(qtypes[qtype])) // Type: A
	query = binary.BigEndian.AppendUint16(query, 1)                     // Class: IN

	// Create Socket
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err != nil {
		panic(err)
	}

	// Send Query
	err = syscall.Sendto(fd, query, 0, &syscall.SockaddrInet4{Port: 53, Addr: [4]byte{8, 8, 8, 8}})
	if err != nil {
		panic(err)
	}

	// Decode Response
	response := make([]byte, 4096)
	_, _, err = syscall.Recvfrom(fd, response, 0)
	if err != nil {
		panic(err)
	}
	//fmt.Println(toHexString(response))

	// Assert identifier match
	responseId := binary.BigEndian.Uint16(response[:2])
	if id != responseId {
		log.Fatalf("Query and Response identifier mismatch")
	}

	responseFlags := binary.BigEndian.Uint16(response[2:4])
	rCode := responseFlags & 0x000f
	if rCode != 0 {
		log.Fatalf("Received non-zero response code: %v", rCode)
	}

	qcount := binary.BigEndian.Uint16(response[4:6])
	if qcount != 1 {
		log.Fatalf("Received unexpected question count: %v", qcount)
	}

	ansCount := binary.BigEndian.Uint16(response[6:8])
	if ansCount <= 0 {
		log.Fatalf("Received non-positive answer count: %v", ansCount)
	}

	answerRecord := response[28:]

	rData := answerRecord[12:16]

	fmt.Printf("IP Address: %v\n", rData)

}

func toHexString(bytes []byte) string {
	var result string
	for i, bite := range hex.EncodeToString(bytes) {
		result += string(bite)
		if i%2 != 0 {
			result += " "
		}
		if (i+1)%32 == 0 {
			result += "\n"
		}
	}
	return result
}
