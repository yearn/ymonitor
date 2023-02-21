package workers

import (
	"log"
	"strconv"
	"yearn/ymonitor/config"
	"yearn/ymonitor/prom"
	"yearn/ymonitor/requests"

	"github.com/prometheus/client_golang/prometheus"
)

func WebsiteMonitor(hosts chan config.Host) {
	for host := range hosts {
		url := host.Url.Url
		log.Printf("querying website %s\n", host.Name)
		res, stats := requests.DoGetRequest(url.String())

		prom.Observe(stats, prometheus.Labels{
			"host":    host.Name,
			"network": host.Network,
			"code":    strconv.Itoa(res.StatusCode),
			"type":    "website",
		})
	}
}
