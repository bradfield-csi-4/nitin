package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"net"
	"strconv"
	"syscall"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))
var seqNum uint16 = 0

const maxHops = 30

func main() {
	socket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_ICMP)
	check(err)

	err = syscall.SetsockoptTimeval(socket, syscall.SOL_SOCKET, syscall.SO_RCVTIMEO, &syscall.Timeval{Sec: 2})
	check(err)

	destAddr := getDestAddr(err, "google.com")

	for ttl := 1; ttl < maxHops; ttl++ {
		id := randUInt16()
		req := makeIcmpEchoRequest(id)

		err := syscall.SetsockoptInt(socket, syscall.IPPROTO_IP, syscall.IP_TTL, ttl)
		check(err)

		start := time.Now()
		err = syscall.Sendto(socket, req, 0, destAddr)
		check(err)

		resp := make([]byte, 100)

		_, _, err = syscall.Recvfrom(socket, resp, 0)
		roundTripTime := time.Since(start)

		if err != nil && err != syscall.EAGAIN {
			panic(err)
		}

		// ICMP response isn't received before timeout
		if err == syscall.EAGAIN {
			fmt.Printf("\t%v  * \n", seqNum)
			continue
		}

		icmp := parseIcmp(resp)

		if icmp.packetType == TimeExceeded || icmp.packetType == IcmpReply {
			if icmp.packetType == IcmpReply && id != icmp.identifier {
				fmt.Println("received unexpected ICMP Reply")
				continue
			}

			fmt.Printf("\t%v  %v (%v) - %v \n", seqNum, icmp.host, icmp.sourceAddr, roundTripTime)

			if icmp.packetType == IcmpReply {
				break
			}
		}
	}

	err = syscall.Close(socket)
	check(err)
}

func getDestAddr(err error, input string) *syscall.SockaddrInet4 {
	ips, err := net.LookupIP("google.com")
	check(err)

	addr := toByteArray(ips[0][12:])

	destAddr := &syscall.SockaddrInet4{
		Port: 80,
		Addr: addr,
	}

	fmt.Printf("traceroute to %v (%v), %v hops max\n", input, ips[0], maxHops)
	return destAddr
}

func toByteArray(byteSlice []byte) [4]byte {
	var addr [4]byte
	for i, b := range byteSlice {
		addr[i] = b
	}
	return addr
}

func check(err error) {
	if err != nil {
		if err != syscall.EINTR {
			panic(err)
		}
	}
}

func makeIcmpEchoRequest(identifier uint16) []byte {
	req := make([]byte, 8)

	// Set Type to IPv4, ICMP
	req[0] = uint8(8)

	// Set random integer as Identifier
	binary.BigEndian.PutUint16(req[4:6], identifier)

	// Set incremented Sequence Number
	seqNum++
	binary.BigEndian.PutUint16(req[6:8], seqNum)

	// Set Checksum
	binary.BigEndian.PutUint16(req[2:4], calcChecksum(req))

	return req
}

func randUInt16() uint16 {
	return uint16(r.Intn(int(math.Pow(2, 16))))
}

func calcChecksum(packet []byte) uint16 {
	var sum uint16
	for i := 0; i < len(packet); i += 2 {
		sum += binary.BigEndian.Uint16(packet[i : i+2])
	}

	return onesComp(sum)
}

func onesComp(num uint16) uint16 {
	return ((1 << 16) - 1) ^ num
}

type icmp struct {
	sourceAddr     string
	host           string
	packetType     int
	code           uint8
	checksum       uint16
	identifier     uint16
	sequenceNumber uint16
}

func parseIcmp(datagram []byte) icmp {
	ipHeaderLen := 4 * (datagram[0] & 0x0F)
	icmpPayload := datagram[ipHeaderLen:]

	sourceIP := getIPStr(toByteArray(datagram[12:16]))

	var host string
	nameParts, err := net.LookupAddr(sourceIP)
	if err == nil {
		host = nameParts[0]
	} else {
		host = sourceIP
	}

	return icmp{
		sourceAddr:     sourceIP,
		host:           host,
		packetType:     int(icmpPayload[0]),
		code:           icmpPayload[1],
		checksum:       binary.BigEndian.Uint16(icmpPayload[2:4]),
		identifier:     binary.BigEndian.Uint16(icmpPayload[4:6]),
		sequenceNumber: binary.BigEndian.Uint16(icmpPayload[6:8]),
	}
}

const (
	IcmpReply    int = 0
	IcmpRequest      = 8
	TimeExceeded     = 11
)

func getIPStr(ip [4]byte) string {
	var result string
	for i := 0; i < len(ip); i++ {
		result += strconv.Itoa(int(ip[i]))
		if i < len(ip)-1 {
			result += "."
		}
	}
	return result
}
