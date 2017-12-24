package main

import (
	"network"
	"runtime"
	"fmt"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	node := network.Init_node()
	node.Init_rpc_server()
	fmt.Println(node)
}
