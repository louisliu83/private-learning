package http

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	_ "pa.cn/fedlearn/psi/prom"
)

func (ss *Server) setPromRoutes(r *mux.Router) {
	r.Handle("/v1/metrics", promhttp.Handler())
}
