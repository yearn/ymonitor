package workers

import (
	"encoding/json"
	"io"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"time"
	"yearn/ymonitor/config"
	"yearn/ymonitor/prom"
	"yearn/ymonitor/requests"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/tcnksm/go-httpstat"
)

type NodeRequest struct {
	JsonRpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Id      int    `json:"id"`
}

type NodeResponse struct {
	JsonRpc string `json:"jsonrpc"`
	Result  string `json:"result"`
	Id      int    `json:"id"`
}

var blockNumberGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: prom.NS,
	Subsystem: prom.SUB,
	Name:      "block_number",
	Help:      "The latest observed block number",
}, []string{"host", "network", "env", "code", "type", "client", "provider"})

func NodeMonitor(hosts chan config.Host) {
	for host := range hosts {
		url := host.Url.Url
		log.Printf("querying node: %s, network: %s, env: %s\n", host.Name, host.Network, host.Env)

		block, statusCode, stats, err := getBlockNumber(url.String())
		if err != nil {
			continue
		}
		clientVersion, _, _, err := getClientVersion(url.String())

		labels := prometheus.Labels{
			"host":     host.Name,
			"network":  host.Network,
			"env":      host.Env,
			"code":     strconv.Itoa(statusCode),
			"type":     "node",
			"client":   clientVersion,
			"provider": host.Provider,
		}

		blockNumberGauge.With(labels).Set(float64(block.Int64()))
		prom.Observe(stats, labels)
	}

}

func getResponseBody(url string, method string) ([]byte, *http.Response, httpstat.Result, error) {
	request := NodeRequest{JsonRpc: "2.0", Method: method, Id: 1}
	payload, err := json.Marshal(request)
	if err != nil {
		log.Print(err)
		return nil, nil, httpstat.Result{}, err
	}
	res, stats, err := requests.DoPostRequest(url, payload)
	if err != nil {
		log.Print(err)
		return nil, nil, httpstat.Result{}, err
	}

	body, err := io.ReadAll(res.Body)
	stats.End(time.Now())
	if err != nil {
		log.Fatal(err)
		return nil, nil, httpstat.Result{}, err
	}
	return body, res, stats, nil
}

func getBlockNumber(url string) (*big.Int, int, httpstat.Result, error) {
	body, res, stats, err := getResponseBody(url, "eth_blockNumber")
	if err != nil {
		return nil, 500, httpstat.Result{}, err
	}
	block := new(big.Int)
	blockNumberRes := NodeResponse{}
	if res.StatusCode == 200 {
		err = json.Unmarshal(body, &blockNumberRes)
		if err != nil {
			log.Print(err)
			return nil, 500, httpstat.Result{}, err
		}
		block.SetString(blockNumberRes.Result, 0)
	}
	res.Body.Close()
	return block, res.StatusCode, stats, nil
}

func getClientVersion(url string) (string, int, httpstat.Result, error) {
	body, res, stats, err := getResponseBody(url, "web3_clientVersion")
	if err != nil {
		return "", 500, httpstat.Result{}, err
	}
	clientRes := NodeResponse{}
	var clientVersion string
	if res.StatusCode == 200 {
		err = json.Unmarshal(body, &clientRes)
		if err != nil {
			log.Print(err)
			return "", 500, httpstat.Result{}, err
		}
		clientVersion = clientRes.Result
	}
	res.Body.Close()
	return clientVersion, res.StatusCode, stats, nil
}
