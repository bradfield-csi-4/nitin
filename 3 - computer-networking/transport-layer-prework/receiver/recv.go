package main

import (
	"encoding/binary"
	"fmt"
	"syscall"
	"transport/util"
)

const SourcePort int = 9999

func main() {
	// Create UDP Socket
	socket := startUDPServer()
	for {
		// Receive packets
		data := make([]byte, 1500)
		_, _, err := syscall.Recvfrom(socket, data, 0)
		util.Check(err)

		var packet util.Packet
		var ack []byte
		packet = extractData(data)

		if util.IsCorrupt(packet) {
			ack = makeAck(0)
		} else {
			ack = makeAck(1)
			deliverData(packet.Msg)
		}

		err = syscall.Sendto(socket, ack, 0, &syscall.SockaddrInet4{Port: packet.SourcePort, Addr: [4]byte{127, 0, 0, 1}})
	}
}

func startUDPServer() int {
	socket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
	util.Check(err)

	serverAddr := &syscall.SockaddrInet4{Port: SourcePort, Addr: [4]byte{127, 0, 0, 1}}

	err = syscall.Bind(socket, serverAddr)
	util.Check(err)

	return socket
}

func extractData(data []byte) util.Packet {
	sourcePort := binary.BigEndian.Uint16(data[2:4])
	checksum := binary.BigEndian.Uint16(data[0:2])
	ack := data[4]
	msg := data[5:]

	return util.Packet{SourcePort: int(sourcePort), Checksum: checksum, Ack: ack, Msg: msg, Payload: data}
}

func makeAck(val int) []byte {
	packet := make([]byte, 5)

	// Set Source Port header
	binary.BigEndian.PutUint16(packet[2:4], uint16(SourcePort))
	// Set Ack value
	packet[4] = uint8(val)
	return packet
}

func deliverData(data []byte) {
	fmt.Printf("Received Data: %v\n", string(data))
}
