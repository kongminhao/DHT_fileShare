package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"bufio"
)

func sender(conn net.Conn, infohash string) {
	words := "get_peers "
	words += infohash
	conn.Write([]byte(words))
	fmt.Println("send over")

}
func showinfo(conn net.Conn)  {
	buffer := make([]byte, 2048)
	words := "get_info"
	conn.Write([]byte(words))
	fmt.Println("send ok")
	n, err := conn.Read(buffer)
	if err !=nil {
		fmt.Println(err)
	}
	str := string(buffer[0:n])
	str_list := strings.Split(str, "_")
	// 123_test.sh_345_kkkk.sh_
	for i :=0 ; i< len(str_list) -1; i+=2{
		fmt.Printf("INFO:%s filename: %s \n", str_list[i], str_list[i+1])
	}
}
func upload(conn net.Conn, path string)  {
	word := "openTcp " + path
	conn.Write([]byte(word))
	fmt.Println("test")
}

func main() {
	server := "223.129.64.13:2333"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
	fmt.Println("connect success")
	fmt.Println("please enter your choice, /q to quit, /help to show help info")
	reader := bufio.NewReader(os.Stdin)
	for   {
		fmt.Print(">>>")
		choice := ""
		strBytes, _, err := reader.ReadLine()
		choice = string(strBytes)
		if err != nil{
			fmt.Println(err)
			break
		}
		if strings.HasPrefix(choice, "/q"){
			break
		}else if strings.HasPrefix(choice,"/help") {
			fmt.Println("/q to quit this program")
			fmt.Println("/help to show this info")
		}else if strings.HasPrefix(choice, "/showinfo") {
			showinfo(conn)
		}else if strings.HasPrefix(choice, "/download"){
			fmt.Println(choice)
			infohash := strings.Split(choice, " ")[1]
			sender(conn, infohash)
		}else if strings.HasPrefix(choice, "/upload"){
			path := strings.Split(choice, " ")[1]
			upload(conn, path)
		}
	}

}
