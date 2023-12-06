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

type GetTransactionCountRequest struct {
	JsonRpc string   `json:"jsonrpc"`
	Method  string   `json:"method"`
	Id      int      `json:"id"`
	Params  []string `json:"params"`
}

type GetTransactionCountResponse struct {
	JsonRpc string `json:"jsonrpc"`
	Result  string `json:"result"`
	Id      int    `json:"id"`
}

var txPoolGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: prom.NS,
	Subsystem: prom.SUB,
	Name:      "tx_pool",
	Help:      "The latest observed nonce for a given address",
}, []string{"host", "network", "env", "code", "type", "client", "provider"})

func TxPoolMonitor(hosts chan config.Host) {
	for host := range hosts {
		url := host.Url.Url
		log.Printf("querying txPool: %s, network: %s, env: %s\n", host.Name, host.Network, host.Env)

		params := []string{"0x7a1057e6e9093da9c1d4c1d049609b6889fc4c67", "pending"}
		request := GetTransactionCountRequest{JsonRpc: "2.0", Method: "eth_getTransactionCount", Id: 1, Params: params}
		payload, err := json.Marshal(request)
		if err != nil {
			log.Print(err)
			continue
		}
		res, stats, err := requests.DoPostRequest(url.String(), payload)
		if err != nil {
			log.Print(err)
			continue
		}

		body, err := ioutil.ReadAll(res.Body)
		stats.End(time.Now())
		if err != nil {
			log.Fatal(err)
		}
		labels := prometheus.Labels{
			"host":     host.Name,
			"network":  host.Network,
			"env":      host.Env,
			"code":     strconv.Itoa(res.StatusCode),
			"type":     "txPool",
			"client":   "n/a",
			"provider": "n/a",
		}
		prom.Observe(stats, labels)

		if res.StatusCode == 200 {
			getTransactionCountRes := GetTransactionCountResponse{}
			err = json.Unmarshal(body, &getTransactionCountRes)
			if err != nil {
				log.Print(err)
				continue
			}
			nonce := new(big.Int)
			nonce.SetString(getTransactionCountRes.Result, 0)
			txPoolGauge.With(labels).Set(float64(nonce.Int64()))
		}
		res.Body.Close()
	}
}
