package main

import (
	"fmt"
	"strconv"
)

func main() {
	i64, err := strconv.Atoi("6131483275329207190")
	fmt.Println(err)
	fmt.Println(i64)
}
