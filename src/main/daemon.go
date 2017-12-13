package main

import (
	"net"
	"network"
)

func main() {
	node := network.Init_node()
	test_addr, _ := net.ResolveTCPAddr("tcp", "192.168.0.12:7890")
	udpaddr := net.UDPAddr{
		IP:   net.IPv4(223, 129, 64, 108),
		Port: 8765,
	}
	node.Announce_peer(123, *test_addr, udpaddr, "test.sh")
	node.Init_rpc_server()
}
