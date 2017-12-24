package network

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"path"
	"strconv"
	"strings"
	"time"
)

func (Node node) recvtcp_msg(Listenconn net.Listener) {
	for {
		conn, err := Listenconn.Accept()
		if err != nil {
			continue
		}
		defer conn.Close()
		go Node.handleConnection(conn)
	}

}

func (Node node) handleConnection(conn net.Conn) {

	buffer := make([]byte, 2048)

	for {

		n, err := conn.Read(buffer)

		if err != nil {
			fmt.Println(conn.RemoteAddr().String(), " connection error: ", err)
			return
		}
		if bytes.HasPrefix(buffer[:n], []byte("get_peers")) {
			infohash := bytes.Split(buffer[:n], []byte(" "))[1]
			i64, err := strconv.Atoi(string(infohash))
			checkError(err)
			var temp_peer_list peer_list
			fmt.Println(local_peer_lists)
			for _, peer_list := range local_peer_lists {
				if peer_list.info.infohash == uint64(i64) {
					temp_peer_list = peer_list
					break
				}
			}
			data := ""
			for _, peer := range temp_peer_list.peer_lists {
				if peer.ip.Port == 0 {
					continue
				}
				data += peer.ip.String() + "_"
			}
			conn.Write([]byte(data))
		}
		if bytes.HasPrefix(buffer[:n], []byte("get_route")) {
			// 调试用
			fmt.Println(node_route_table)
		}
		if bytes.HasPrefix(buffer[:n], []byte("get_info")) {
			data := ""
			for _, info := range infolist {
				if info.infohash == 0 {
					continue
				}
				Node.Get_peers(info.infohash)
				fmt.Println(info.String())
				data += info.String() + "_"
			}
			fmt.Println(data)
			conn.Write([]byte(data))
		}
		if bytes.HasPrefix(buffer[:n], []byte("openTcp")) { // openTcp filepath
			msg := string(buffer[:n])
			filepath := strings.Split(msg, " ")[1]
			tcpstring := openTcpPort(filepath)
			tcpaddr, err := net.ResolveTCPAddr("tcp", tcpstring)
			checkError(err)
			fmt.Println(tcpaddr)
			filename := path.Base(filepath)
			infohash := md5toinfohash(filepath)
			infoHash := info_hash{
				infohash: infohash,
				filename: filename,
			}
			go Node.doAnnouncePeer(infoHash, *tcpaddr)
		}
	}

}
func (Node node) doAnnouncePeer(infoHash info_hash, tcpaddr net.TCPAddr) {
	for {
		Node.Announce_peer(infoHash.infohash, tcpaddr, infoHash.filename)
		time.Sleep(10 * time.Second)
	}
}
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
		//fmt.Println("from", raddr)

		if bytes.HasPrefix(buf[0:], []byte("ping")) {
			msg := Node.handle_ping(conn, buf[0:], n)
			conn.WriteToUDP([]byte(msg), raddr)
		} else if bytes.HasPrefix(buf[0:], []byte("findnode")) {
			handle_find_node()
		} else if bytes.HasPrefix(buf[0:], []byte("announcepeer")) {
			msg := Node.handle_announce_peer(conn, buf[0:], n)
			conn.WriteToUDP([]byte(msg), raddr)
		} else if bytes.HasPrefix(buf[0:], []byte("getpeers")) {
			handle_get_peer(buf[0:], n)
		} else if bytes.HasPrefix(buf[0:], []byte("broadcastinfo")) {
			_, err := handle_broadcastinfo(conn, buf[0:], n, raddr)
			if err != nil {
				go Node.Ping_all() // 以协程方式启动，防止阻塞
				fmt.Println(err)
			}
		} else if bytes.HasPrefix(buf[0:], []byte("infohash")) {
			handle_infohash(buf[0:], n)
		} else if bytes.HasPrefix(buf[0:], []byte("peer")) {
			handlepeer(buf[0:], n)
		}
		//WriteToUDP
		//func (c *UDPConn) WriteToUDP(b []byte, addr *UDPAddr) (int, error)
		//n, err = conn.WriteToUDP([]byte("nice to see u"), raddr)
		//fmt.Println(n)
		//checkError(err)
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
func handlepeer(buf []byte, n int) {
	// 处理get_peer的返回值
	//   peer_192.168.0.12:7890_
	msg := string(buf[0:n])
	str_list := strings.Split(msg, "_")
	i64, err := strconv.Atoi(str_list[1])
	checkError(err)
	str_list = str_list[2:]
	ui64 := uint64(i64)
	fmt.Println(i64, ui64)
	info := info_hash{
		infohash: ui64,
		filename: "test.txt",
	}
	var temp_peerlist peer_list
	fmt.Println(local_peer_lists)
	for _, each := range local_peer_lists {
		if each.info.infohash == info.infohash {
			temp_peerlist = each
		}
	}
	temp_peerlist.info = info
	for _, straddr := range str_list {
		tcpaddr, err := net.ResolveTCPAddr("tcp", straddr)
		checkError(err)
		peer := peer{
			ip: (*tcpaddr),
		}
		temp_peerlist.peer_lists = append(temp_peerlist.peer_lists, peer)
	}
	local_peer_lists = append(local_peer_lists, temp_peerlist)
}

func handle_get_peer(buf []byte, n int) string {
	// 处理get peers 请求
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
			return_msg := "peer_" + str_list[1] + "_"
			for _, peer := range peerlist.peer_lists {
				return_msg += peer.ip.String()
				return_msg += "_"
				fmt.Println(peer.ip.String())
			}
			// todo: check info stream
			remote := str_list[3] + ":" + str_list[4]
			raddr, err := net.ResolveUDPAddr("udp", remote)

			checkError(err)
			conn, err := net.DialUDP("udp", nil, raddr)
			checkError(err)
			defer conn.Close()
			n, err := conn.Write([]byte(return_msg))
			checkError(err)
			fmt.Println(n)
		}
	}
	if flag == 0 {
		nodeid := uint16(infohash % 0xffff) // 计算映射节点
		ttl, err := strconv.Atoi(str_list[5])
		checkError(err)
		ttl -= 1
		if ttl < 0 {
			return "fail"
		}
		forward_msg := "getpeers:" + str_list[1] + ":" + str_list[2] + ":" + str_list[3] + ":" + str_list[4] + ":" + strconv.Itoa(ttl)
		if distance(nodeid, node_route_table.pre_node.id) < distance(nodeid, node_route_table.after_node.id) {
			// 转发给pre_node
			forward_conn, err := net.DialUDP("udp", nil, &node_route_table.pre_node.ip_addr)
			checkError(err)
			defer forward_conn.Close()
			forward_conn.Write([]byte(forward_msg))

		} else {
			// 转发给after_node
			forward_conn, err := net.DialUDP("udp", nil, &node_route_table.after_node.ip_addr)
			checkError(err)
			defer forward_conn.Close()
			forward_conn.Write([]byte(forward_msg))
		}
	}
	return "success"
}

func (Node node) handle_announce_peer(conn *net.UDPConn, buf []byte, n int) string {
	// 处理announce_peer请求，完成！
	msg := string((buf[0:n]))
	var str_list []string = strings.Split(msg, ":")
	// msg := "announcepeer:" + string(info_hash) + ":" + string(Node.id) + ":" + string(tcpaddr.IP) + ":" + string(tcpaddr.Port)
	s64, err := strconv.Atoi(str_list[1])
	checkError(err)
	s := uint64(s64)
	infohash := s
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

	nodeid := uint16(infohash % 0xffff) // 计算映射节点

	if distance(nodeid, Node.id) > distance(nodeid, node_route_table.pre_node.id) && node_route_table.pre_node.ip_addr.String() != broadcast_addr.String() {
		forwardconn, err := net.DialUDP("udp", nil, &node_route_table.pre_node.ip_addr)
		checkError(err)
		forwardconn.Write([]byte(msg))
	} else if distance(nodeid, Node.id) > distance(nodeid, node_route_table.after_node.id) && node_route_table.after_node.ip_addr.String() != broadcast_addr.String() {
		forwardconn, err := net.DialUDP("udp", nil, &node_route_table.after_node.ip_addr)
		checkError(err)
		forwardconn.Write([]byte(msg))
	} else {
		// 如果该peer有对应的info_hash, 则将其加入peer_list
		flag := 0
		for _, peerlist := range peer_lists {
			if peerlist.info.infohash == info.infohash {
				flag = 1
				if len(peerlist.peer_lists) < 5 { // 确保peer_list大小小于5
					for _, Peer := range peerlist.peer_lists {
						if Peer.ip.String() == PEER.ip.String() {
							flag = 2
						}
					}
					if flag == 1 {
						peerlist.peer_lists = append(peerlist.peer_lists, PEER) // append 采用副作用编程
					}
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
		msg = "success"
	}
	return msg
}
func handle_infohash(buf []byte, n int) {
	msg := string(buf[0:n])
	str_list := strings.Split(msg, "_")
	id64, err := strconv.ParseInt(str_list[1], 10, 64)
	if err != nil {
		fmt.Println(err)
	}
	infoid := uint64(id64)
	info := info_hash{
		infohash: infoid,
		filename: str_list[2],
	}
	infoch <- info
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
	for i, peer := range peer_lists {
		if i == 0 {
			continue
		}
		infohash := peer.info
		return_msg := "infohash" + "_" + infohash.String()
		if return_msg == "infohash_0_" { // 暴力解决
			continue
		}
		rconn.Write([]byte(return_msg))
	}
	//fmt.Println(faddr.String())
	//fmt.Println(node_route_table.pre_node.ip_addr.String())
	//fmt.Println(node_route_table.after_node.ip_addr.String())
	if faddr.IP.String() == get_localip().String() {
		return "ok", nil
	}
	if faddr.IP.String() == node_route_table.pre_node.ip_addr.IP.String() {
		fconn, err := net.DialUDP("udp", nil, &node_route_table.after_node.ip_addr)
		checkError(err)
		fconn.Write(buf[0:n]) //朝同一个方向转发
	} else if faddr.IP.String() == node_route_table.after_node.ip_addr.IP.String() {
		fconn, err := net.DialUDP("udp", nil, &node_route_table.pre_node.ip_addr)
		checkError(err)
		fconn.Write(buf[0:n])
	} else {
		//出错了,emmmmm, 重新ping全局路由
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
		if N.id == Node.id { // 防止本地环回，辣鸡
			continue
		}
		if Node.id < N.id {
			if distance(Node.id, N.id) < distance(node_route_table.pre_node.id, N.id) {
				node_route_table.pre_node = Node
			}
		} else {
			if distance(Node.id, N.id) < distance(node_route_table.after_node.id, N.id) {
				node_route_table.after_node = Node
			}
		}
		fmt.Println("update")
		fmt.Println(Node)
	}
}
func handle_ping_resp(buf []byte) {
	// complete
	msg := string(buf)
	//fmt.Println(msg)
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

func openTcpPort(path string) string {
	laddr := net.TCPAddr{
		IP:   get_localip(),
		Port: 0, // to random port
	}
	Listen_conn, err := net.ListenTCP("tcp", &laddr)
	if err != nil {
		panic(err)
	}
	//defer Listen_conn.Close()
	// open file to transport
	go openFileDownload(*Listen_conn, path)
	return Listen_conn.Addr().String()
}
func openFileDownload(listener net.TCPListener, path string) {
	for {
		//fmt.Println("test donwload")
		conn, err := listener.AcceptTCP()
		if err != nil {
			continue
		}
		defer conn.Close()
		go func() {
			b, err := ioutil.ReadFile(path)
			if err != nil {
				fmt.Println(err)
				return
			}
			conn.Write(b)
			conn.Close()
		}()
	}
}
func md5toinfohash(file string) uint64 {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return 0
	}
	b := md5.Sum(data)
	b_buf := bytes.NewBuffer(b[0:8])
	var x uint64
	binary.Read(b_buf, binary.BigEndian, &x)
	x = x % 0xffffff
	return x
}
