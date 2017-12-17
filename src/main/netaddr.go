package main

import (
	"fmt"
	"network"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	node := network.Init_node()
	node.Init_rpc_server()

	fmt.Println(node)

}
