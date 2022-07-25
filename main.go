package main

import (
	"fmt"
	"net"
)

func main() {
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8080})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Local: <%s> \n", listener.LocalAddr().String())

	data := make([]byte, 1024)
	for {
		// 等待连接
		n, remoteAddr, err := listener.ReadFromUDP(data)
		if err != nil || n != 10 || data[9]&(1<<1) == 0 {
			fmt.Printf("ERROR %s\n", err.Error())
		}

		// 发送ACK和SYN
		serverSerial := byte(0)
		clientSerial := data[3]
		start := make([]byte, 10)
		start[3] = serverSerial
		start[7] = clientSerial + 1
		start[9] = 1<<1 + 1<<4
		_, err = listener.WriteToUDP(start, remoteAddr)
		if err != nil {
			fmt.Printf("ERROR %s\n", err.Error())
		}

		// 读取ACK
		n, err = listener.Read(data)
		if err != nil || n != 10 || data[9]&(1<<4) == 0 {
			fmt.Printf("ERROR %s\n", err.Error())
		}
		clientSerial = data[3]

		// 读取信息
		n, err = listener.Read(data)
		if err != nil {
			fmt.Printf("ERROR %s\n", err.Error())
		}
		fmt.Printf("<%s> %s\n", remoteAddr, data[:n])

		// 做出回复
		_, err = listener.WriteToUDP([]byte("world"), remoteAddr)
		if err != nil {
			fmt.Printf("ERROR %s\n", err.Error())
		}
	}
}
