package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
	"yearn/ymonitor/config"
	"yearn/ymonitor/workers"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	godotenv.Load()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.HandleFunc("/block-explorer", blockExplorerRedirect)
		http.ListenAndServe(":8090", nil)
	}()

	// endless loop
	for {
		// parse host config and monitor all hosts with the correct worker
		for hostType, hosts := range parseHosts() {
			hostList := make(chan config.Host, len(hosts))
			var wg sync.WaitGroup
			const NumWorkers = 4
			for i := 0; i < NumWorkers; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					monitor(hostType, hostList)
				}()
			}
			for _, h := range hosts {
				hostList <- h
			}
			close(hostList)
			wg.Wait()
		}
		time.Sleep(10 * time.Second)
	}
}

func monitor(hostType string, hostList chan config.Host) {
	switch hostType {
	case "node":
		workers.NodeMonitor(hostList)
	case "website", "api":
		workers.SimpleMonitor(hostList, hostType)
	case "apy":
		workers.ApyMonitor(hostList)
	default:
		log.Fatal("Unknown host type: " + hostType)
	}
}

func parseHosts() map[string][]config.Host {
	fileName := "/monitors.json"
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal("Can't read input file: " + fileName)
	}
	hosts := make(map[string][]config.Host)
	err = json.Unmarshal([]byte(file), &hosts)
	if err != nil {
		log.Fatal(err)
	}
	return hosts
}

func blockExplorerRedirect(w http.ResponseWriter, req *http.Request) {
	network := req.URL.Query()["var-network"][0]
	block := req.URL.Query()["block"][0]
	var url string
	switch network {
	case "ethereum":
		url = "https://etherscan.io/block/" + block
	case "fantom":
		url = "https://ftmscan.com/block/" + block
	case "optimism":
		url = "https://optimistic.etherscan.io/block/" + block
	case "arbitrum":
		url = "https://arbiscan.io/block/" + block
	default:
		url = "https://etherscan.io/block/" + block
	}

	http.Redirect(w, req, url, 307)
}
