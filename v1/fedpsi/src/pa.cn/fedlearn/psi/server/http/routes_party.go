package http

import (
	"net/http"

	"github.com/gorilla/mux"
	psiserver "pa.cn/fedlearn/psi/server/psi"
)

func (ss *Server) setParty2PartyRoutes(r *mux.Router) {
	s := psiserver.New()
	r.HandleFunc("/{srcParty}/info", s.Info).Methods(http.MethodGet)
	r.HandleFunc("/{srcParty}/v1/info", s.Info).Methods(http.MethodGet)
	r.HandleFunc("/{srcParty}/v1/dataset", s.PartyDataSetList).Methods(http.MethodGet)
	r.HandleFunc("/{srcParty}/v1/dataset/{name}/{index}", s.PartyDataSetGet).Methods(http.MethodGet)
	r.HandleFunc("/{srcParty}/v1/dataset/shards", s.DataSetShards).Methods(http.MethodGet)
	r.HandleFunc("/{srcParty}/v1/job", s.PartyJobSubmit).Methods(http.MethodPost)
	r.HandleFunc("/{srcParty}/v1/job/confirm", s.PartyJobConfirm).Methods(http.MethodPost)
	r.HandleFunc("/{srcParty}/v1/job/intersect", s.PartyJobIntersect).Methods(http.MethodGet)
	//r.HandleFunc("/{srcParty}/v1/task", s.PartyTaskCreate).Methods(http.MethodPost)
	r.HandleFunc("/{srcParty}/v1/task", s.PartyTaskCreateV2).Methods(http.MethodPost)
	r.HandleFunc("/{srcParty}/v1/task/confirm", s.PartyTaskConfirm).Methods(http.MethodPost)
	r.HandleFunc("/{srcParty}/v1/task/start", s.PartyTaskStart).Methods(http.MethodPost)
	r.HandleFunc("/{srcParty}/v1/task/stop", s.PartyTaskStop).Methods(http.MethodPost)
	r.HandleFunc("/{srcParty}/v1/task/rerun", s.PartyTaskRerun).Methods(http.MethodPost)
	r.HandleFunc("/{srcParty}/v1/task", s.PartyTaskGet).Methods(http.MethodGet)
	r.HandleFunc("/{srcParty}/v1/task/intersect", s.PartyTaskIntersect).Methods(http.MethodGet)
}
