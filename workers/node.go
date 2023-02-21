package workers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/big"
	"strconv"
	"time"
	"yearn/ymonitor/config"
	"yearn/ymonitor/prom"
	"yearn/ymonitor/requests"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type BlockNumberRequest struct {
	JsonRpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Id      int    `json:"id"`
}

type BlockNumberResponse struct {
	JsonRpc string `json:"jsonrpc"`
	Result  string `json:"result"`
	Id      int    `json:"id"`
}

var blockNumberGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: prom.NS,
	Subsystem: prom.SUB,
	Name:      "block_number",
	Help:      "The latest observed block number",
}, []string{"host", "network", "code", "type"})

func NodeMonitor(hosts chan config.Host) {
	for host := range hosts {
		url := host.Url.Url
		log.Printf("querying node %s\n", host.Name)

		blockNumberRequest := BlockNumberRequest{JsonRpc: "2.0", Method: "eth_blockNumber", Id: 1}
		payload, err := json.Marshal(blockNumberRequest)
		if err != nil {
			log.Fatal(err)
		}
		res, stats := requests.DoPostRequest(url.String(), payload)

		body, err := ioutil.ReadAll(res.Body)
		stats.End(time.Now())
		if err != nil {
			log.Fatal(err)
		}
		labels := prometheus.Labels{
			"host":    host.Name,
			"network": host.Network,
			"code":    strconv.Itoa(res.StatusCode),
			"type":    "node",
		}
		prom.Observe(stats, labels)

		if res.StatusCode == 200 {
			blockNumberRes := BlockNumberResponse{}
			err = json.Unmarshal(body, &blockNumberRes)
			if err != nil {
				log.Fatal(err)
			}
			block := new(big.Int)
			block.SetString(blockNumberRes.Result, 0)
			blockNumberGauge.With(labels).Set(float64(block.Int64()))
		}
		res.Body.Close()
	}
}
