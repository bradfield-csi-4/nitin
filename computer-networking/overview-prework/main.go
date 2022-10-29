package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
)

const (
	PcapFileHeaderLength      = 24
	PcapPerPacketHeaderLength = 16
	EthernetHeaderLength      = 14
	HttpSourcePort            = uint16(80)
)

func main() {
	f, err := os.Open("computer-networking/overview-prework/net.cap")
	if err != nil {
		panic(err)
	}

	// Parse capture file header
	fileHeader := make([]byte, PcapFileHeaderLength)
	_, err = f.Read(fileHeader)
	if err != nil {
		panic(err)
	}

	packetHeader := make([]byte, PcapPerPacketHeaderLength)

	//var segments []segment

	responseParts := map[uint32][]byte{}

	for {
		// Parse packet header
		_, err = f.Read(packetHeader)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("error reading packet header %v", err)
		}

		packetLength := binary.LittleEndian.Uint32(packetHeader[8:12])

		frame := parseEthernetFrame(f, packetLength)
		datagram := parseIPDatagram(frame)
		if !isDestinationLocalHost(datagram) {
			continue
		}

		segment := parseTCPSegment(datagram)
		responseParts[segment.sequenceNumber] = segment.payload
		//printLayers(frame, datagram, segment)
	}

	// Sort by sequence number
	keys := make([]uint32, 0, len(responseParts))
	values := make([][]byte, len(responseParts))
	for k := range responseParts {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	for i, k := range keys {
		values[i] = responseParts[k]
	}

	// Combine fragments
	messageBytes := bytes.Join(values, []byte{})

	//var messageBytes []byte
	//for _, segment := range segments {
	//	// - The header length of the 1st packet is 40 bytes (which matches Wireshark)
	//	// - That leaves 2 bytes ("b2 f8") that I thought were the HTTP message, but they're too short and so must be part of the TCP segment
	//	// - But I don't get it, because they shouldn't be part of the header because they're over 40 bytes
	//	// - Decided to only consider HTTP payloads larger than 2 bytes
	//	if len(segment.payload) > 2 {
	//		messageBytes = append(messageBytes, segment.payload...)
	//	}
	//}

	endOfHeader := bytes.Index(messageBytes, []byte{'\r', '\n', '\r', '\n'})
	var message string
	var i int
	for i < endOfHeader {
		message += string(messageBytes[i])
		i++
	}

	//fmt.Println()
	//fmt.Println(message[:endOfHeader])

	// Skip 4 bytes used to identify end-of-header
	body := messageBytes[i+4:]
	//fmt.Println(toHexString(body))

	outputFile, err := os.Create("computer-networking/overview-prework/output.jpg")
	if err != nil {
		panic(err)
	}

	_, err = outputFile.Write(body)
	if err != nil {
		panic(err)
	}
	outputFile.Close()
}

func sortBySequenceNumber(segments []segment) {
	//segments = removeDuplicateValues(segments)
	sort.Slice(segments, func(i, j int) bool {
		return segments[i].sequenceNumber < segments[j].sequenceNumber
	})
}

func removeDuplicateValues(segments []segment) []segment {
	keys := make(map[int]bool)
	var result []segment

	for _, segment := range segments {
		if _, value := keys[int(segment.sequenceNumber)]; !value {
			keys[int(segment.sequenceNumber)] = true
			result = append(result, segment)
		}
	}
	return result
}

func printLayers(frame frame, datagram datagram, segment segment) {
	fmt.Println("ETHERNET FRAME")
	fmt.Println(toHexString(frame.raw))
	fmt.Println()
	fmt.Println("IP DATAGRAM")
	fmt.Println(toHexString(frame.payload))
	fmt.Println()
	fmt.Println("TCP SEGMENT")
	fmt.Println(toHexString(datagram.payload))
	fmt.Println()
	fmt.Println("HTTP MESSAGE")
	fmt.Println(toHexString(segment.payload))
	fmt.Println()
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

type segment struct {
	sourcePort      uint16
	destinationPort uint16
	sequenceNumber  uint32
	headerLength    int64
	payload         []byte
}

func parseTCPSegment(datagram datagram) segment {
	payload := datagram.payload

	tcpHeaderLength := getTCPHeaderLength(payload)

	return segment{
		sourcePort:      binary.BigEndian.Uint16(payload[0:2]),
		destinationPort: binary.BigEndian.Uint16(payload[2:4]),
		sequenceNumber:  binary.BigEndian.Uint32(payload[4:8]),
		headerLength:    tcpHeaderLength,
		payload:         payload[tcpHeaderLength:],
	}
}

func isDestinationLocalHost(datagram datagram) bool {
	return bytes.Compare(datagram.destination, []byte{192, 168, 0, 101}) == 0
}

type datagram struct {
	headerLength int64
	length       uint16
	source       []byte
	destination  []byte
	protocol     string
	payload      []byte
}

func parseIPDatagram(frame frame) datagram {
	payload := frame.payload

	ipHeaderLength := getIPHeaderLength(payload)

	return datagram{
		headerLength: ipHeaderLength,
		length:       binary.BigEndian.Uint16(payload[2:4]),
		source:       payload[12:16],
		destination:  payload[16:20],
		protocol:     getNetworkProtocol(payload[9]),
		payload:      payload[ipHeaderLength:],
	}

}

func getIPHeaderLength(payload []byte) int64 {
	headerLengthRawBinary := strconv.FormatUint(uint64(int64(payload[0])), 2)
	// Could also have done a logical AND against '00001111' to get the last 4 bits
	lastFourBits := headerLengthRawBinary[len(headerLengthRawBinary)-4:]
	headerLengthValue, _ := strconv.ParseInt(lastFourBits, 2, 4)
	return headerLengthValue * 4
}

func getTCPHeaderLength(payload []byte) int64 {
	headerLengthRawBinary := strconv.FormatUint(uint64(int64(payload[12])), 2)
	// Could have also done this by right shifting by 4 bits.
	firstFourBits := headerLengthRawBinary[0:4]
	dataOffset, _ := strconv.ParseInt(firstFourBits, 2, 8)
	return dataOffset * 4
}

type frame struct {
	destination string
	source      string
	ethertype   string
	payload     []byte
	raw         []byte
}

func parseEthernetFrame(f *os.File, length uint32) frame {
	var frame frame
	frameBytes := make([]byte, length)
	_, err := f.Read(frameBytes)
	if err != nil {
		fmt.Println("error reading ethernet frame")
	}

	frame.destination = formatMacAddress(hex.EncodeToString(frameBytes[0:6]))
	frame.source = formatMacAddress(hex.EncodeToString(frameBytes[6:12]))
	frame.ethertype = getEthertype(hex.EncodeToString(frameBytes[12:14]))
	frame.payload = frameBytes[EthernetHeaderLength:]
	frame.raw = frameBytes
	return frame
}

func getEthertype(etherTypeHex string) string {
	if etherTypeHex == "0800" {
		return "IPv4"
	}
	return ""
}

func getNetworkProtocol(byte byte) string {
	if byte == 6 {
		return "TCP"
	}
	return ""
}

func formatMacAddress(unformatted string) string {
	var result string
	for i, ch := range unformatted {
		result += string(ch)
		if i != len(unformatted)-1 && i%2 > 0 {
			result += ":"
		}
	}
	return result
}
