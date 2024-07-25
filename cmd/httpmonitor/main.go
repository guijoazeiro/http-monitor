package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Server struct {
	ServerName    string
	ServerUrl     string
	ExecutionTime float64
	Status        int
}

func createServerList(data [][]string) []Server {
	var servers []Server

	for i, line := range data {
		if i > 0 {
			server := Server{
				ServerName: line[0],
				ServerUrl:  line[1],
			}
			servers = append(servers, server)
		}

	}
	return servers
}

func checkServer(servers []Server) {
	for _, server := range servers {
		now := time.Now()
		get, err := http.Get(server.ServerUrl)
		if err != nil {
			fmt.Println(err)
		}
		server.ExecutionTime = time.Since(now).Seconds()
		server.Status = get.StatusCode
		fmt.Printf("Status: [%d] Time: [%f] Url: [%s]\n", server.Status, server.ExecutionTime, server.ServerUrl)
	}
}

func main() {

	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer f.Close()

	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	servers := createServerList(data)
	checkServer(servers)

}
