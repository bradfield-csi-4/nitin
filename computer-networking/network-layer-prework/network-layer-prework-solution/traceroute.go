package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"syscall"
	"time"
)

const (
	SOL_IP = 0x00 // not sure why this isn't available under `syscall`
)

type IpAddr [4]byte

func (x *IpAddr) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", x[0], x[1], x[2], x[3])
}

func (x *IpAddr) Equals(y *IpAddr) bool {
	return x[0] == y[0] && x[1] == y[1] && x[2] == y[2] && x[3] == y[3]
}

func IcmpTtlExpired(data []byte) bool {
	ihl := (data[0] & 0x0f) << 2
	return data[ihl] == 11
}

func IcmpDestinationUnreachable(data []byte) bool {
	ihl := (data[0] & 0x0f) << 2
	return data[ihl] == 3 && data[ihl+1] == 3
}

// Trace a basic path to given host
// Usage e.g.: sudo go run traceroute.go -q 2 -w 1 google.com
func main() {
	empty := make([]byte, 24)
	data := make([]byte, 4096)
	nullIp := IpAddr{}

	// Parse command line
	if len(os.Args) < 2 {
		log.Fatal("Usage: sudo go run traceroute.go [flags] foo.com")
	}
	addr, e := net.LookupIP(os.Args[len(os.Args)-1])
	check(e)
	host := [4]byte{}
	host[0] = addr[0][12]
	host[1] = addr[0][13]
	host[2] = addr[0][14]
	host[3] = addr[0][15]

	var waitTime, nQueries, firstTtl, maxTtl, port int
	flag.IntVar(&firstTtl, "f", 1, "Set the initial time to live")
	flag.IntVar(&maxTtl, "m", 64, "Set the max time to live")
	flag.IntVar(&port, "p", 33434, "Set the base port number used in probes")
	flag.IntVar(&nQueries, "q", 3, "Set the number of probes per ttl")
	flag.IntVar(&waitTime, "w", 5, "Set the time to wait for a response to a probe, in seconds")
	flag.Parse()

	// Establish sockets for both sender (UDP) and receiver (raw)

	sender, e := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	check(e)

	receiver, e := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_ICMP)
	check(e)

	tv := syscall.Timeval{Sec: int64(waitTime), Usec: 0}
	// Time out if we don't receive replies
	e = syscall.SetsockoptTimeval(receiver, syscall.SOL_SOCKET, syscall.SO_RCVTIMEO, &tv)
	check(e)

	for ttl := firstTtl; ttl <= maxTtl; ttl += 1 {
		fmt.Printf("%2d  ", ttl)
		priorIp := nullIp
		for i := 0; i < nQueries; i += 1 {
			e = syscall.SetsockoptInt(sender, SOL_IP, syscall.IP_TTL, ttl)
			check(e)
			start := time.Now()

			e = syscall.Sendto(sender, empty, 0, &syscall.SockaddrInet4{Port: port, Addr: host})
			check(e)
			for {
				_, sa, e := syscall.Recvfrom(receiver, data, 0)
				// We want to handle "resource temporarily unavailable" ourselves
				// TODO panic on *other* errors
				if e != nil {
					if priorIp.Equals(&nullIp) {
						fmt.Print("* ")
					} else {
						fmt.Print("\n    *")
					}
					break
				}
				// Extract the IP address and name of responding router
				x, ok := sa.(*syscall.SockaddrInet4)
				if !ok {
					continue
				}
				respondingIp := IpAddr(x.Addr)
				name := respondingIp.String()
				nameParts, e := net.LookupAddr(name)
				if e == nil {
					name = nameParts[0]
				}
				// Expect response to be an ICMP of "TTL expired" or "destination unreachable" type
				if data[9] != 1 || !(IcmpTtlExpired(data) || IcmpDestinationUnreachable(data)) {
					continue
				}

				elapsed := time.Since(start)
				if !respondingIp.Equals(&priorIp) {
					if !priorIp.Equals(&nullIp) {
						fmt.Print("\n    ")
					}
					fmt.Printf("%s (%s)", name, respondingIp.String())
					priorIp = respondingIp
				}
				fmt.Printf("  %.3f ms", float32(elapsed)/1000000)
				break
			}
		}
		fmt.Printf("\n")
		port += 1

		if IcmpDestinationUnreachable(data) {
			break
		}

	}

	syscall.Close(sender)
	syscall.Close(receiver)

	fmt.Printf("Ok\n")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
