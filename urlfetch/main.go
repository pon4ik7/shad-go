//go:build !solution

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	args := os.Args[1:]
	for _, arg := range args {
		resp, err := http.Get(arg)
		if err != nil {
			fmt.Println(err)
			panic(1)
		}
		defer resp.Body.Close()
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			panic(1)
		}
		fmt.Println(string(data))
	}
}
