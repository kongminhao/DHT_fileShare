package network

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

// simple DHT network implementation
type node struct {
	id      uint16      // 16位的id， 保证小型局域网内不重复
	ip_addr net.UDPAddr // 节点的ip地址， 包含端口号
}

type info_hash struct {
	infohash uint64 // md5值取前16位
	filename string // 储存文件名
}

type peer struct {
	ip net.TCPAddr // ip && port , 实现分片下载用
}

type peer_list struct {
	info       info_hash
	peer_lists []peer // 规定下大小， 最多5个，否则丢弃
}

type route_table struct { // 简化的DHT, 只存储前者和后者节点的信息
	pre_node   node
	after_node node
}

var broadcast_addr net.UDPAddr = net.UDPAddr{
	IP:   net.IPv4(255, 255, 255, 255),
	Port: 8765,
}
var minnode node = node{
	id:      0,
	ip_addr: broadcast_addr,
}
var maxnode node = node{
	id:      0xffff,
	ip_addr: broadcast_addr,
}
var node_route_table route_table = route_table{
	pre_node:   minnode,
	after_node: maxnode,
}
var peer_lists []peer_list = make([]peer_list, 10) // 一个节点存10条记录，emmmmm
// 构建一个无缓冲的node信道， 便于构建路由表
var nodech chan node = make(chan node)
var infolist []info_hash

// 构建一个无缓冲的info_hash信道，便于自动的获取info_hash
var infoch chan info_hash = make(chan info_hash)

func get_localip() (ip net.IP) {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Println(ipnet.IP.String())
				ip = ipnet.IP
			}

		}

	}
	return
}

func (Info info_hash) String() string {
	return strconv.Itoa(int(Info.infohash)) + "_" + Info.filename
}
