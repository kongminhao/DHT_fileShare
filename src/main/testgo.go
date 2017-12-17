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

var test chan string = make(chan  string)
func main() {
	test <- "1213"
	test <- "sdasda"
	go say("world")
	say("hello")
}
