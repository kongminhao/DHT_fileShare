package main

import "net"

func main() {
	raddr, _ := net.ResolveUDPAddr("udp", "119.29.16.200:687")
	taddr, _ := net.ResolveUDPAddr("udp", "119.29.16.200:687")
	if (*raddr).String() == (taddr).String() {
		print(123)
	}
}
