package prom

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	METRICS_NET_PROXY_BYTES = "net_proxy_bytes"
)

type Metrics struct {
	metrics map[string]*prometheus.Desc
	mutex   sync.Mutex
}

func NewMetrics(namespace string) *Metrics {
	return &Metrics{
		metrics: map[string]*prometheus.Desc{
			METRICS_NET_PROXY_BYTES: prometheus.NewDesc(
				namespace+"_"+METRICS_NET_PROXY_BYTES,
				"net proxy bytes total",
				[]string{"host"},
				nil),
		},
	}
}

func (c *Metrics) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range c.metrics {
		ch <- m
	}
}

func (c *Metrics) Collect(ch chan<- prometheus.Metric) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for host, currentValue := range proxyBytesMap {
		ch <- prometheus.MustNewConstMetric(
			c.metrics[METRICS_NET_PROXY_BYTES],
			prometheus.GaugeValue,
			float64(currentValue),
			host)
	}
}
