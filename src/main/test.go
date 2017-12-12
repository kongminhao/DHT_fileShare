package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	// 这里设置发送者的IP地址，自己查看一下自己的IP自行设定
	laddr := net.UDPAddr{
		IP:   net.IPv4(10, 132, 24, 246),
		Port: 32919,
	}
	// 这里设置接收者的IP地址为广播地址
	raddr := net.UDPAddr{
		IP: net.IPv4(255, 255, 255, 255),
		//IP: net.IPv4(192,168,191,1),
		Port: 8765,
	}
	fmt.Println(raddr.String())
	fmt.Println(laddr)
	conn, err := net.ListenUDP("udp", &laddr)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
	var msg [200]byte
	var i int = 10
	for i = 0; i < 10; i++ {
		conn.WriteToUDP([]byte("ping:123:119.29.16.200:6789"), &raddr)
	}

	n, _, err := conn.ReadFromUDP(msg[0:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
	t := msg[0:n]
	fmt.Println(string(msg[0:n]))
	fmt.Println(string(t))
	conn.Close()
}
