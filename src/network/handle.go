package network

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

func (Node node) recvUDPMsg(conn *net.UDPConn) {

	// 以协程方式启动两个更新程序
	go Node.update_route_table()
	go update_infolist()
	var buf [200]byte
	// 循环服务
	for {
		n, raddr, err := conn.ReadFromUDP(buf[0:])
		if err != nil {
			return
		}
		fmt.Println("msg is ", string(buf[0:n]))
		fmt.Println("from", raddr)
		if bytes.HasPrefix(buf[0:], []byte("ping")) {
			msg := Node.handle_ping(conn, buf[0:], n)
			conn.WriteToUDP([]byte(msg), raddr)
		} else if bytes.HasPrefix(buf[0:], []byte("findnode")) {
			handle_find_node()
		} else if bytes.HasPrefix(buf[0:], []byte("announcepeer")) {
			msg := handle_announce_peer(conn, buf[0:], n)
			conn.WriteToUDP([]byte(msg), raddr)
		} else if bytes.HasPrefix(buf[0:], []byte("getpeers")) {
			handle_get_peer(conn, buf[0:], n)
		} else if bytes.HasPrefix(buf[0:], []byte("broadcastinfo")) {
			_, err := handle_broadcastinfo(conn, buf[0:], n, raddr)
			if err != nil {
				go Node.Ping_all() // 以协程方式启动，防止阻塞
			}
		}
		//WriteToUDP
		//func (c *UDPConn) WriteToUDP(b []byte, addr *UDPAddr) (int, error)
		n, err = conn.WriteToUDP([]byte("nice to see u"), raddr)
		fmt.Println(n)
		checkError(err)
	}
}

func (Node node) handle_ping(conn *net.UDPConn, buf []byte, n int) string {
	//todo: 处理ping请求
	// msg := "ping:" + string(Node.id) + ":" + string(Node.ip_addr.IP) + ":" + string(Node.ip_addr.Port)
	msg := string(buf[0:n])
	var str_list []string = strings.Split(msg, ":")
	fmt.Println(str_list)
	id64, err := strconv.ParseUint(str_list[1], 10, 16)
	checkError(err)
	id := uint16(id64)
	ipaddr := str_list[2] + ":" + str_list[3]
	node_addr, err := net.ResolveUDPAddr("udp", ipaddr)
	checkError(err)
	ping_node := node{
		id:      id,
		ip_addr: *node_addr,
	}
	nodech <- ping_node
	return_msg := "pingresp:" + strconv.Itoa(int(Node.id)) + ":" + Node.ip_addr.IP.String() + ":" + strconv.Itoa(Node.ip_addr.Port)
	return return_msg
}

func handle_find_node() {
	// todo: 处理find_node请求
}

func handle_get_peer(conn *net.UDPConn, buf []byte, n int) string {
	// todo: 处理get peers 请求
	// 判断是否在自己的peer_list中， 如在，返回对应的peer_list
	// 否则，向下级路由转发， 同时，将路由的跳数-1
	// msg := "getpeers:" + string(info_hash) + ":" + string(Node.id) + ":" + string(Node.ip_addr.IP) + ":" + string(Node.ip_addr.Port) + ":" + string(255)
	msg := string(buf[0:n])
	var str_list []string = strings.Split(msg, ":")
	infohash64, err := strconv.ParseUint(str_list[1], 10, 64)
	checkError(err)
	infohash := uint64(infohash64)
	// 查找是否存在该info_hash对应的peer_list
	flag := 0
	for _, peerlist := range peer_lists {
		if peerlist.info.infohash == infohash {
			flag = 1
			// todo: return peer_list
			return_msg := "peer_list"
			for _, peer := range peerlist.peer_lists {
				return_msg += peer.ip.String()
				return_msg += ":"
				fmt.Println(peer.ip.String())
			}
			// todo: check info stream
			remote := str_list[3] + str_list[4]
			raddr, err := net.ResolveUDPAddr("udp", remote)
			conn.WriteToUDP([]byte(return_msg), raddr)

			checkError(err)
		}
	}
	if flag == 0 {
		nodeid := uint16(infohash % 0xffff) // 计算映射节点
		ttl, err := strconv.Atoi(str_list[5])
		checkError(err)
		ttl -= 1
		forward_msg := "getpeers:" + str_list[1] + ":" + str_list[2] + ":" + str_list[3] + ":" + str_list[4] + ":" + strconv.Itoa(ttl)
		if distance(nodeid, node_route_table.pre_node.id) < distance(nodeid, node_route_table.after_node.id) {
			// 转发给pre_node
			forward_conn, err := net.DialUDP("udp", nil, &node_route_table.pre_node.ip_addr)
			checkError(err)
			forward_conn.Write([]byte(forward_msg))

		} else {
			// 转发给after_node
			forward_conn, err := net.DialUDP("udp", nil, &node_route_table.after_node.ip_addr)
			checkError(err)
			forward_conn.Write([]byte(forward_msg))
		}
	}
	return "success"
}

func handle_announce_peer(conn *net.UDPConn, buf []byte, n int) string {
	// todo: 处理announce peer请求， 向下级路由转发也要做，妈个鸡，好烦啊, 判断是否有节点比自己离info_hash近，没有的话， 再加入节点。
	msg := string((buf[0:n]))
	var str_list []string = strings.Split(msg, ":")
	// msg := "announcepeer:" + string(info_hash) + ":" + string(Node.id) + ":" + string(tcpaddr.IP) + ":" + string(tcpaddr.Port)
	s64, err := strconv.Atoi(str_list[1])
	checkError(err)
	s := uint64(s64)
	info := info_hash{
		infohash: s,
		filename: str_list[5],
	}
	ipaddr := str_list[3] + ":" + str_list[4]
	tcpaddr, err := net.ResolveTCPAddr("tcp", ipaddr)
	checkError(err)
	PEER := peer{
		ip: *tcpaddr,
	}
	// 如果该peer有对应的info_hash, 则将其加入peer_list
	flag := 0
	for _, peerlist := range peer_lists {
		if peerlist.info.infohash == info.infohash {
			if len(peerlist.peer_lists) < 5 { // 确保peer_list大小小于5
				flag = 1
				peerlist.peer_lists = append(peerlist.peer_lists, PEER) // append 采用副作用编程
			}
		}
	}
	// init empty peer , 如果没有对应的info_hash
	if flag == 0 {
		var p []peer = []peer{}
		p = append(p, PEER)
		newpeerlist := peer_list{
			info:       info,
			peer_lists: p,
		}
		peer_lists = append(peer_lists, newpeerlist)
	}
	// todo: 返回值，啦啦啦
	msg = "success"
	return msg
}

func handle_broadcastinfo(conn *net.UDPConn, buf []byte, n int, faddr *net.UDPAddr) (return_msg string, error error) {
	// raddr 朝相反方向转发
	msg := string(buf[0:n])
	str_list := strings.Split(msg, "_")
	raddr, err := net.ResolveUDPAddr("udp", str_list[1])
	checkError(err)
	rconn, err := net.DialUDP("udp", nil, raddr)
	checkError(err)

	defer rconn.Close()
	for _, peer := range peer_lists {
		infohash := peer.info
		return_msg := "infohash" + "_" + infohash.String()
		rconn.Write([]byte(return_msg))
	}
	if faddr.String() == node_route_table.pre_node.ip_addr.String() {
		fconn, err := net.DialUDP("udp", nil, &node_route_table.after_node.ip_addr)
		checkError(err)
		fconn.Write(buf[0:n]) //朝同一个方向转发
	} else if faddr.String() == node_route_table.after_node.ip_addr.String() {
		fconn, err := net.DialUDP("udp", nil, &node_route_table.pre_node.ip_addr)
		checkError(err)
		fconn.Write(buf[0:n])
	} else {
		//todo: 出错了,emmmmm, 重新ping全局路由
		return_msg := "fail"
		error := errors.New("fail_to_forwardmsg")
		return return_msg, error
	}
	return return_msg, nil
}

func update_infolist() {
	for {
		flag := 0
		info_hashrecv := <-infoch
		for _, infohash := range infolist {
			if infohash.infohash == info_hashrecv.infohash {
				flag = 1
				break
			}
		}
		if flag == 0 {
			infolist = append(infolist, info_hashrecv)
		}
	}
}

// 构建一个全局的nch, 便于做路由表的更新
func (N node) update_route_table() {
	// 测试更新正常
	for {
		Node := <-nodech
		// 更新路由表
		if Node.id < N.id {
			if distance(Node.id, N.id) < distance(node_route_table.pre_node.id, N.id) {
				node_route_table.pre_node = Node
			}
		} else {
			if distance(Node.id, N.id) < distance(node_route_table.after_node.id, N.id) {
				node_route_table.after_node = Node
			}
		}
		fmt.Println(Node)
	}
}
func handle_ping_resp(buf []byte) {
	// complete
	msg := string(buf)
	var str_list []string = strings.Split(msg, ":")
	// 构建出node
	id64, err := strconv.Atoi(str_list[1])
	checkError(err)
	id := uint16(id64)
	ipaddr := str_list[2] + ":" + str_list[3]
	node_addr, err := net.ResolveUDPAddr("udp", ipaddr)
	checkError(err)
	Node := node{
		id:      id,
		ip_addr: *node_addr,
	}
	nodech <- Node
}
