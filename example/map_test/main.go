package main

import (
	"fmt"
)

func main() {

	m := make(map[string]int)

	m["abc"] = 300
	m["egf"] = 500
	m["eiie"] = 100
	for key, val := range m {
		fmt.Printf("key:%s val:%d\n", key, val)
	}
}
