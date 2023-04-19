package http

import (
	"net/http"

	"github.com/gorilla/mux"
	psiserver "pa.cn/fedlearn/psi/server/psi"
)

func (ss *Server) setPsiRoutes(r *mux.Router) {
	s := psiserver.New()
	r.HandleFunc("/v1/info", s.Info).Methods(http.MethodGet)
	r.HandleFunc("/v1/config", s.ConfigInfo).Methods(http.MethodGet)

	r.HandleFunc("/v1/parties", s.PartyRegister).Methods(http.MethodPost)
	r.HandleFunc("/v1/parties", s.PartyUpdate).Methods(http.MethodPut)
	r.HandleFunc("/v1/parties", s.PartyList).Methods(http.MethodGet)
	r.HandleFunc("/v1/party/{name}", s.PartyDel).Methods(http.MethodDelete)
	r.HandleFunc("/v1/party/{name}", s.PartyReady).Methods(http.MethodGet)

	r.HandleFunc("/v1/tproxy/start", s.TProxyStart).Methods(http.MethodPost)
	r.HandleFunc("/v1/tproxy/stop", s.TProxyStop).Methods(http.MethodPost)
	r.HandleFunc("/v1/tproxy/status", s.TProxyStatus).Methods(http.MethodGet)

	r.HandleFunc("/v1/chunk/check", s.CheckChunk).Methods(http.MethodPost)
	r.HandleFunc("/v1/chunk/upload", s.UploadChunk).Methods(http.MethodPost)
	r.HandleFunc("/v1/chunk/merge", s.MergeChunk).Methods(http.MethodPost)
	r.HandleFunc("/v1/dataset/pull", s.FilePull).Methods(http.MethodPost)
	r.HandleFunc("/v1/dataset/copy", s.FileCopy).Methods(http.MethodPost)

	r.HandleFunc("/v1/dataset", s.DataSetList).Methods(http.MethodGet)
	r.HandleFunc("/v1/dataset/{name}/{index}", s.DataSetGet).Methods(http.MethodGet)
	r.HandleFunc("/v1/dataset/{name}/{index}", s.DataSetDel).Methods(http.MethodDelete)
	r.HandleFunc("/v1/dataset/shards", s.DataSetShards).Methods(http.MethodGet)

	r.HandleFunc("/v1/dataset/bizcode/{bizcode}/intersects", s.DataSetIntersectList).Methods(http.MethodGet)

	r.HandleFunc("/v1/dataset/grant", s.DataSetGrant).Methods(http.MethodPost)
	r.HandleFunc("/v1/dataset/revoke", s.DataSetRevoke).Methods(http.MethodPost)

	r.HandleFunc("/v1/job", s.JobSubmit).Methods(http.MethodPost)
	r.HandleFunc("/v1/jobs", s.JobList).Methods(http.MethodGet)
	r.HandleFunc("/v1/job/confirm", s.JobConfirm).Methods(http.MethodPost)
	r.HandleFunc("/v1/job/stop", s.JobStop).Methods(http.MethodPost)
	r.HandleFunc("/v1/job/{jobuid}", s.JobDel).Methods(http.MethodDelete)
	r.HandleFunc("/v1/job/intersect", s.JobIntersect).Methods(http.MethodGet)
	r.HandleFunc("/v1/job/intersect/download", s.JobIntersectDownload).Methods(http.MethodGet)

	//r.HandleFunc("/v1/task", s.TaskCreate).Methods(http.MethodPost)
	r.HandleFunc("/v1/task", s.TaskCreateV2).Methods(http.MethodPost)
	r.HandleFunc("/v1/task/confirm", s.TaskConfirm).Methods(http.MethodPost)
	r.HandleFunc("/v1/task", s.TaskGet).Methods(http.MethodGet)
	r.HandleFunc("/v1/task/start", s.TaskStart).Methods(http.MethodPost)
	r.HandleFunc("/v1/task/stop", s.TaskStop).Methods(http.MethodPost)
	r.HandleFunc("/v1/task/rerun", s.TaskRerun).Methods(http.MethodPost)
	r.HandleFunc("/v1/task/intersect", s.TaskIntersect).Methods(http.MethodGet)

	r.HandleFunc("/v1/tasks", s.TaskList).Methods(http.MethodGet)
	r.HandleFunc("/v1/task/intersect/download", s.TaskIntersectDownload).Methods(http.MethodGet)
}
