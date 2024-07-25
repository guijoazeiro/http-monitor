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
	FailureDate   string
}

func createServerList(serverList *os.File) []Server {
	csvReader := csv.NewReader(serverList)
	data, err := csvReader.ReadAll()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
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

func checkServer(servers []Server) []Server {
	var downServers []Server
	for _, server := range servers {
		now := time.Now()
		get, err := http.Get(server.ServerUrl)
		if err != nil {
			fmt.Printf("Server %s is down[%s]\n", server.ServerUrl, err.Error())
			server.Status = 0
			server.FailureDate = now.Format("2006-01-02 15:04:05")
			downServers = append(downServers, server)
			continue
		}
		server.Status = get.StatusCode
		if server.Status != 200 {
			server.FailureDate = now.Format("2006-01-02 15:04:05")
			downServers = append(downServers, server)
		}
		server.ExecutionTime = time.Since(now).Seconds()
		fmt.Printf("Status: [%d] Time: [%f] Url: [%s]\n", server.Status, server.ExecutionTime, server.ServerUrl)
	}
	return downServers
}

func generateDownTime(downtimeList *os.File, downServers []Server) {
	csvWriter := csv.NewWriter(downtimeList)

	for _, server := range downServers {
		line := []string{server.ServerName, server.ServerUrl, server.FailureDate, fmt.Sprintf("%f", server.ExecutionTime), fmt.Sprintf("%d", server.Status)}
		csvWriter.Write(line)
	}
	csvWriter.Flush()
}

func openFiles(serverListFile string, downtimeFile string) (*os.File, *os.File) {
	serverList, err := os.OpenFile(serverListFile, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println(err)
	}

	downtimeList, err := os.OpenFile(downtimeFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)
	}

	return serverList, downtimeList

}

func main() {

	serverList, downtimeList := openFiles(os.Args[1], os.Args[2])
	defer serverList.Close()
	defer downtimeList.Close()

	servers := createServerList(serverList)

	downServers := checkServer(servers)

	generateDownTime(downtimeList, downServers)

}
