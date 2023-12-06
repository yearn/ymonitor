package workers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strconv"
	"time"
	"yearn/ymonitor/config"
	"yearn/ymonitor/prom"
	"yearn/ymonitor/requests"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var apyGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: prom.NS,
	Subsystem: prom.SUB,
	Name:      "apy_timestamp",
	Help:      "The latest observed apy timestamp in seconds",
}, []string{"host", "network", "env", "code", "type", "client", "provider"})

func ApyMonitor(hosts chan config.Host) {
	for host := range hosts {
		url := host.Url.Url
		log.Printf("querying apy: %s, network: %s, env: %s\n", host.Name, host.Network, host.Env)
		res, stats, err := requests.DoGetRequest(url.String())
		if err != nil {
			log.Print(err)
			continue
		}

		body, err := ioutil.ReadAll(res.Body)
		stats.End(time.Now())

		labels := prometheus.Labels{
			"host":     host.Name,
			"network":  host.Network,
			"env":      host.Env,
			"code":     strconv.Itoa(res.StatusCode),
			"type":     "apy",
			"client":   "n/a",
			"provider": "yexporter",
		}
		prom.Observe(stats, labels)
		if res.StatusCode >= 200 {
			var apyRes interface{}
			err = json.Unmarshal(body, &apyRes)
			if err != nil {
				log.Print(err)
				continue
			}
			firstItem := apyRes.([]interface{})[0]
			updated := firstItem.(map[string]interface{})["updated"]

			apyGauge.With(labels).Set(updated.(float64))
		}
		res.Body.Close()
	}
}
