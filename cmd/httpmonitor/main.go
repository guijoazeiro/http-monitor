package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	now := time.Now()
	url := os.Args[1]
	get, err := http.Get(url)

	if err != nil {
		fmt.Println("error in get url")
		panic(err)
	}

	elapsed := time.Since(now).Seconds()
	status := get.StatusCode

	fmt.Printf("Status: {%d} Time: {%f}\n", status, elapsed)

}
