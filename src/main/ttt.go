package main

import (
	"fmt"
	"strconv"
)

func main() {
	var s uint8 = 10
	fmt.Println(string(s))
	fmt.Println(strconv.Itoa(int(s)))
	a := []int{1, 2, 3, 4, 5}
	a = append(a, 10)
	fmt.Println(cap(a))
	for _, num := range a {
		fmt.Println(num)
	}

}
