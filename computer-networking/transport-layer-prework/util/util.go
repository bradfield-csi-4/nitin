package util

import "syscall"

type Packet struct {
	SourcePort int
	Checksum   uint16
	Ack        uint8
	Msg        []byte
	Payload    []byte
}

func Check(err error) {
	if err != nil {
		if err != syscall.EINTR {
			panic(err)
		}
	}
}

// CalcChecksum TODO: Use actual checksum algorithm
func CalcChecksum(packet []byte) uint16 {
	var checksum uint8
	for i := 2; i < len(packet); i++ {
		checksum += packet[i]
	}

	return uint16(checksum)
}

func IsCorrupt(packet Packet) bool {
	return packet.Checksum != CalcChecksum(packet.Payload)
}
