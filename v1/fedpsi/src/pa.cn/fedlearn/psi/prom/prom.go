package prom

import (
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	col1 := NewMetrics("psi")
	prometheus.MustRegister(col1)
}
