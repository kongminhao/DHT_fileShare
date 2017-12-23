package main

import (
	"net"
	"network"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	node := network.Init_node()
	test_addr, _ := net.ResolveTCPAddr("tcp", "192.168.0.12:7890")
	go node.Announce_peer(123, *test_addr, "test.sh")
	go node.Announce_peer(345, *test_addr,"kkkk.sh")
	node.Init_rpc_server()
}
