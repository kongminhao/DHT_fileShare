package network

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"
)

// 初始化node
func Init_node() node {
	addr := net.UDPAddr{IP: get_localip(), Port: 8765}
	Node := node{uint16(rand.Uint32()), addr}
	fmt.Println(Node)
	//return_msg := "pingresp:" + strconv.Itoa(int(Node.id)) + ":" + Node.ip_addr.IP.String() + ":" + strconv.Itoa(Node.ip_addr.Port)
	return Node
}

func (Node node) Init_rpc_server() {
	ludpaddr := net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 8765,
	}
	conn, err := net.ListenUDP("udp", &ludpaddr)
	checkError(err)
	defer conn.Close()
	Node.recvUDPMsg(conn)

}

func (Node node) Ping_all() {
	// todo: 处理ping的返回结果
	/*
		ping 局域网内所有主机， 来进行路由表的初始化
		写入的结构
		node的ping:id:node-ipport
	*/
	// 重新初始化我们的节点
	node_route_table.pre_node = minnode
	node_route_table.after_node = maxnode
	ip := Node.ip_addr.IP
	laddr := net.UDPAddr{
		IP:   ip,
		Port: 6789,
	}
	// 这里设置接收者的IP地址为广播地址
	raddr := net.UDPAddr{
		IP:   net.IPv4(255, 255, 255, 255),
		Port: 8765,
	}
	conn, err := net.ListenUDP("udp", &laddr)
	checkError(err)
	defer conn.Close()
	msg := "ping:" + strconv.Itoa(int(Node.id)) + ":" + Node.ip_addr.IP.String() + ":" + strconv.Itoa(Node.ip_addr.Port)
	conn.WriteToUDP([]byte(msg), &raddr)
	buf := make(chan []byte)
	var recv [200]byte
	defer close(buf)
	go func() {
		n, _, err := conn.ReadFromUDP(recv[0:])
		checkError(err)
		buf <- recv[0:n]
	}()

	for {
		select {
		case <-time.After(3 * time.Second):
			// ok
			return
		case ch := <-buf:
			handle_ping_resp(ch)
			go func() {
				n, _, err := conn.ReadFromUDP(recv[0:])
				checkError(err)
				buf <- recv[0:n]
			}()
		}
	}
}

func (Node node) get_all_info() {
	// 由需求导出的处理函数。 自动的得到局域网内的info_hash
	// 其实只是强行使用DHT做一个路由而已
	// 不难，明天做， 我不管那么多了，先实现需求再说。妈个鸡
	msg := 	"broadcastinfo" + "_" + Node.ip_addr.String()
	// 这里向两侧路由转发
	conn_pre, err := net.DialUDP("udp",nil,&node_route_table.pre_node.ip_addr)
	checkError(err)
	conn_pre.Write([]byte(msg))
	conn_after, err := net.DialUDP("udp", nil, &node_route_table.after_node.ip_addr)
	checkError(err)
	conn_after.Write([]byte(msg))

}

func (Node node) ping() {
	// complete: 超时更新路由表, ping的返回处理, 错误处理
	ip := Node.ip_addr.IP
	laddr := net.UDPAddr{
		IP:   ip,
		Port: 32200,
	}

	pre_addr := net.UDPAddr{
		IP:   node_route_table.pre_node.ip_addr.IP,
		Port: 8765,
	}
	after_addr := net.UDPAddr{
		IP:   node_route_table.after_node.ip_addr.IP,
		Port: 8765,
	}
	conn, err := net.ListenUDP("udp", &laddr)
	checkError(err)
	msg := "ping:" + strconv.Itoa(int(Node.id)) + ":" + Node.ip_addr.IP.String() + ":" + strconv.Itoa(Node.ip_addr.Port)
	// 关闭连接 nice
	defer conn.Close()
	conn.WriteToUDP([]byte(msg), &pre_addr)
	conn.WriteToUDP([]byte(msg), &after_addr)
	// 以下做超时处理
	buf := make(chan []byte)
	var recv [200]byte
	go func() {
		n, _, err := conn.ReadFromUDP(recv[0:])
		checkError(err)
		buf <- recv[0:n]
	}()

	for i := 0; i < 1; i++ {
		select {
		case <-time.After(3 * time.Second):
			// todo: 超时处理, 重建路由表，删除超时项
			Node.Ping_all() // 执行发现节点。判断，是否是maxnode超时或者是否有一个为minnode, 因为这样肯定超时
		case ch := <-buf:
			handle_ping_resp(ch)
			go func() {
				n, _, err := conn.ReadFromUDP(recv[0:])
				checkError(err)
				buf <- recv[0:n]
			}()
		}
	}
	close(buf)
}
func (Node node) Find_node(id uint16) {
	// 255， ttl
	// todo: 查找指定id // 虽说我觉得这个真没什么用处
	var raddr net.UDPAddr
	laddr := net.UDPAddr{
		IP:   Node.ip_addr.IP,
		Port: 32221,
	}
	msg := "findnode:" + string(id) + ":" + string(Node.id) + ":" + string(Node.ip_addr.IP) + ":" + string(Node.ip_addr.Port) + ":" + string(255)
	if id > node_route_table.after_node.id {
		raddr = node_route_table.after_node.ip_addr
	}
	if id < node_route_table.pre_node.id {
		raddr = node_route_table.pre_node.ip_addr
	}
	conn, err := net.ListenUDP("udp", &laddr)
	checkError(err)
	conn.WriteToUDP([]byte(msg), &raddr)
}
func (Node node) Get_peers(info_hash uint64) {
	// todo: 获取正在下载文件的节点list
	var raddr net.UDPAddr
	msg := "getpeers:" + string(info_hash) + ":" + string(Node.id) + ":" + string(Node.ip_addr.IP) + ":" + string(Node.ip_addr.Port) + ":" + string(255)
	laddr := net.UDPAddr{
		IP:   Node.ip_addr.IP,
		Port: 32222,
	}
	// 计算info_hash映射到的节点.
	target_id := uint16(info_hash % 0xffff)
	// 寻找与info_hash 映射节点更近的节点
	if distance(target_id, node_route_table.pre_node.id) > distance(target_id, Node.id) && distance(target_id, node_route_table.after_node.id) > distance(target_id, Node.id) {
		raddr = Node.ip_addr
	} else if distance(target_id, node_route_table.pre_node.id) > distance(target_id, node_route_table.after_node.id) {
		raddr = node_route_table.after_node.ip_addr
	} else {
		raddr = node_route_table.pre_node.ip_addr
	}
	conn, err := net.ListenUDP("udp", &laddr)
	checkError(err)

	_, err = conn.WriteToUDP([]byte(msg), &raddr)
	checkError(err)
}
func (Node node) Announce_peer(info_hash uint64, tcpaddr net.TCPAddr, target_addr net.UDPAddr, filename string) {
	// 完成
	laddr := net.UDPAddr{
		IP:   Node.ip_addr.IP,
		Port: 32223,
	}
	msg := "announcepeer:" + strconv.Itoa(int(info_hash)) + ":" + strconv.Itoa(int(Node.id)) + ":" + tcpaddr.IP.String() + ":" + strconv.Itoa(tcpaddr.Port) + ":" + filename
	conn, err := net.DialUDP("udp", &laddr, &target_addr)
	checkError(err)
	_, err = conn.Write([]byte(msg))
	checkError(err)
}
func checkError(err error) {
	if err != nil {
		fmt.Println("Error: %s", err.Error())
		os.Exit(1)
	}
}

func distance(id1 uint16, id2 uint16) uint16 {
	// 简单的距离计算，异或
	return id1 ^ id2
}
