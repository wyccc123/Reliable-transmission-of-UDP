package main

import (
	"fmt"
	"net"
)

func main() {
	ip := net.ParseIP("127.0.0.1")
	srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	dstAddr := &net.UDPAddr{IP: ip, Port: 8080}
	conn, err := net.DialUDP("udp", srcAddr, dstAddr)
	if err != nil {
		fmt.Printf("ERROR %s\n", err.Error())
	}
	defer conn.Close()

	// 4 byte sequence | 4 byte received | 2 byte sign
	clientSerial := byte(0)
	start := make([]byte, 10)
	start[3] = clientSerial
	start[9] = 1 << 1
	_, err = conn.Write(start)
	if err != nil {
		fmt.Printf("ERROR %s\n", err.Error())
	}

	data := make([]byte, 1024)
	n, err := conn.Read(data)
	if err != nil || n != 10 || data[9]&(1<<4) == 0 {
		fmt.Printf("ERROR %s\n", err.Error())
	}
	serverSerial := data[3]
	clientSerial += 1
	start[3] = clientSerial
	start[7] = serverSerial
	start[9] = 1 << 4
	_, err = conn.Write(start)

	_, err = conn.Write([]byte("hello"))
	n, err = conn.Read(data)
	if err != nil {
		fmt.Printf("ERROR %s\n", err.Error())
	}
	fmt.Printf("<%s> %s\n", conn.RemoteAddr(), data[:n])
}
