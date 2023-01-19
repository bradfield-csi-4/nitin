package main

import (
	"bufio"
	"encoding/binary"
	"os"
	"syscall"
	"transport/util"
)

const (
	SourcePort       = 3333
	ProxyServerPort  = 6666
	SocketTimeoutSec = 2
)

func main() {
	// Create UDP Socket
	socket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
	util.Check(err)

	err = syscall.Bind(socket, &syscall.SockaddrInet4{Port: SourcePort, Addr: [4]byte{127, 0, 0, 1}})
	util.Check(err)

	proxyAddr := &syscall.SockaddrInet4{Port: ProxyServerPort, Addr: [4]byte{127, 0, 0, 1}}

	// Set Recv timeout
	err = syscall.SetsockoptTimeval(
		socket,
		syscall.SOL_SOCKET,
		syscall.SO_RCVTIMEO,
		&syscall.Timeval{Sec: SocketTimeoutSec})
	util.Check(err)

	// Read from standard input & send packet
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		util.Check(scanner.Err())

		// Send packet
		packet := makePacket(scanner.Text())
		err = syscall.Sendto(socket, packet, 0, proxyAddr)
		util.Check(err)

		// Wait for ack
		for {
			ackPacket := make([]byte, 1500)

			_, _, err = syscall.Recvfrom(socket, ackPacket, 0)
			if err != nil && err != syscall.EAGAIN {
				panic(err)
			}

			if (err == nil && !isAck(ackPacket)) || err == syscall.EAGAIN {
				err = syscall.Sendto(socket, packet, 0, proxyAddr)
				util.Check(err)
			} else {
				break
			}
		}
	}
	err = syscall.Close(socket)
	util.Check(err)
}

func makePacket(msg string) []byte {
	packet := make([]byte, 5)
	// Set Source Port header
	binary.BigEndian.PutUint16(packet[2:4], uint16(SourcePort))
	// Append msg
	packet = append(packet, msg...)
	// Set checksum
	binary.BigEndian.PutUint16(packet[0:2], util.CalcChecksum(packet))
	return packet
}

func isAck(packet []byte) bool {
	return packet[4] == 1
}
