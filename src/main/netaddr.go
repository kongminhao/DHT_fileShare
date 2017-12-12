package main

import (
	"fmt"
	"network"
)

func main() {

	node := network.Init_node()
	node.Init_rpc_server()

	fmt.Println(node)

}
