//go:build !solution

package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

var initTime time.Time = time.Now()

func fetch(u string, wg *sync.WaitGroup) {
	defer wg.Done()

	parsed, err := url.Parse(u)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		fmt.Println(time.Since(initTime), u, err)
		return
	}

	resp, err := http.Get(u)
	if err != nil {
		fmt.Println(time.Since(initTime), u, err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(time.Since(initTime), err)
			return
		}
	}(resp.Body)
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(time.Since(initTime), err)
		return
	}
	fmt.Println(time.Since(initTime), len(data), u)
}

func main() {
	args := os.Args[1:]
	var wg sync.WaitGroup
	for _, arg := range args {
		wg.Add(1)
		go fetch(arg, &wg)
	}
	wg.Wait()
}
