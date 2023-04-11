package main

import (
	"flag"
	"fmt"
	"strings"
)

var num int

func init() {
	flag.IntVar(&num, "num", 0, "input num")
}

func main() {
	flag.Parse()
	s := fmt.Sprintf("%08b", num)

	aa := strings.Split(s, "")

	bb := make([]byte, 0, 8)
	for _, b := range aa {
		if b == "0" {
			bb = append(bb, 32)
		} else {
			bb = append(bb, 46)
		}
	}
	fmt.Println(string(bb))
}
