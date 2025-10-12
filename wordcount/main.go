//go:build !solution

package main

import (
	"fmt"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	args := os.Args[1:]
	mp := make(map[string]int)
	for _, name := range args {
		data, err := os.ReadFile(name)
		check(err)
		words := strings.Split(string(data), "\n")
		for _, word := range words {
			mp[word]++
		}
	}
	for k, v := range mp {
		if v >= 2 {
			fmt.Printf("%d\t%s\n", v, k)
		}
	}
}
