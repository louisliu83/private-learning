package http

import (
	"net/http"

	psiserver "fedlearn/psi/server/psi"

	"github.com/gorilla/mux"
)

func (ss *Server) setParty2PartyRoutes(r *mux.Router) {
	s := psiserver.New()
	r.HandleFunc("/{srcParty}/v2/info", s.Info).Methods(http.MethodGet)
	r.HandleFunc("/{srcParty}/v2/worker/info", s.WorkerInfo).Methods(http.MethodGet)
	r.HandleFunc("/{srcParty}/v2/config", s.ConfigInfo).Methods(http.MethodGet)
	r.HandleFunc("/{srcParty}/v2/dataset", s.PartyDataSetList).Methods(http.MethodGet)
	r.HandleFunc("/{srcParty}/v2/dataset/{name}/{index}", s.PartyDataSetGet).Methods(http.MethodGet)
	r.HandleFunc("/{srcParty}/v2/dataset/shards", s.DataSetShards).Methods(http.MethodGet)
	r.HandleFunc("/{srcParty}/v2/job", s.PartyJobSubmit).Methods(http.MethodPost)
	r.HandleFunc("/{srcParty}/v2/job/confirm", s.PartyJobConfirm).Methods(http.MethodPost)
	r.HandleFunc("/{srcParty}/v2/job/intersect", s.PartyJobIntersect).Methods(http.MethodGet)
	r.HandleFunc("/{srcParty}/v2/task", s.PartyTaskCreate).Methods(http.MethodPost)
	r.HandleFunc("/{srcParty}/v2/task/confirm", s.PartyTaskConfirm).Methods(http.MethodPost)
	r.HandleFunc("/{srcParty}/v2/task/start", s.PartyTaskStart).Methods(http.MethodPost)
	r.HandleFunc("/{srcParty}/v2/task/stop", s.PartyTaskStop).Methods(http.MethodPost)
	r.HandleFunc("/{srcParty}/v2/task/rerun", s.PartyTaskRerun).Methods(http.MethodPost)
	r.HandleFunc("/{srcParty}/v2/task", s.PartyTaskGet).Methods(http.MethodGet)
	r.HandleFunc("/{srcParty}/v2/task/intersect", s.PartyTaskIntersect).Methods(http.MethodGet)
	r.HandleFunc("/{srcParty}/v2/activity", s.PartyActivityAdd).Methods(http.MethodPost)
	r.HandleFunc("/{srcParty}/v2/activity/confirm", s.PartyActivityConfirm).Methods(http.MethodPost)
	r.HandleFunc("/{srcParty}/v2/project", s.PartyProjectCreate).Methods(http.MethodPost)
}
