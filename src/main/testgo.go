package main

import (
	"fmt"
)

func say(s string) {
	for {
		//runtime.Gosched() //让执行此函数的goroutine进行go的轮询调度
		fmt.Println(s)
	}
}

func main() {
	go say("world")
	say("hello")
}
