package http

import (
	"net/http"

	psiserver "fedlearn/psi/server/psi"

	"github.com/gorilla/mux"
)

func (ss *Server) setPsiRoutes(r *mux.Router) {
	s := psiserver.New()
	r.HandleFunc("/v2/info", s.Info).Methods(http.MethodGet)
	r.HandleFunc("/v2/config", s.ConfigInfo).Methods(http.MethodGet)

	r.HandleFunc("/v2/project", s.ProjectCreate).Methods(http.MethodPost)
	r.HandleFunc("/v2/project", s.ProjectUpdate).Methods(http.MethodPut)
	r.HandleFunc("/v2/projects", s.ProjectList).Methods(http.MethodGet)
	r.HandleFunc("/v2/project/{id}", s.ProjectDel).Methods(http.MethodDelete)
	r.HandleFunc("/v2/project/{id}", s.ProjectGet).Methods(http.MethodGet)
	r.HandleFunc("/v2/project/jobs/{puid}", s.ProjectJobsGet).Methods(http.MethodGet)

	r.HandleFunc("/v2/parties", s.PartyRegister).Methods(http.MethodPost)
	r.HandleFunc("/v2/parties", s.PartyUpdate).Methods(http.MethodPut)
	r.HandleFunc("/v2/parties", s.PartyList).Methods(http.MethodGet)
	r.HandleFunc("/v2/party/{name}", s.PartyDel).Methods(http.MethodDelete)
	r.HandleFunc("/v2/party/{name}", s.PartyReady).Methods(http.MethodGet)
	r.HandleFunc("/v2/party/worker/{name}", s.PartyWorkerReady).Methods(http.MethodGet)

	r.HandleFunc("/v2/chunk/check", s.CheckChunk).Methods(http.MethodPost)
	r.HandleFunc("/v2/chunk/upload", s.UploadChunk).Methods(http.MethodPost)
	r.HandleFunc("/v2/chunk/merge", s.MergeChunk).Methods(http.MethodPost)
	r.HandleFunc("/v2/dataset/pull", s.FilePull).Methods(http.MethodPost)
	r.HandleFunc("/v2/dataset/copy", s.FileCopy).Methods(http.MethodPost)

	r.HandleFunc("/v2/dataset", s.DataSetList).Methods(http.MethodGet)
	r.HandleFunc("/v2/dataset/{name}/{index}", s.DataSetGet).Methods(http.MethodGet)
	r.HandleFunc("/v2/dataset/{name}/{index}", s.DataSetDel).Methods(http.MethodDelete)
	r.HandleFunc("/v2/dataset/shards", s.DataSetShards).Methods(http.MethodGet)

	r.HandleFunc("/v2/dataset/bizcode/{bizcode}/intersects", s.DataSetIntersectList).Methods(http.MethodGet)

	r.HandleFunc("/v2/dataset/grant", s.DataSetGrant).Methods(http.MethodPost)
	r.HandleFunc("/v2/dataset/revoke", s.DataSetRevoke).Methods(http.MethodPost)

	r.HandleFunc("/v2/activity", s.ActivityAdd).Methods(http.MethodPost)
	r.HandleFunc("/v2/activity", s.ActivityAttachData).Methods(http.MethodPut)
	r.HandleFunc("/v2/activity/start/{uuid}", s.ActivityStart).Methods(http.MethodPost)
	r.HandleFunc("/v2/activity/confirm", s.ActivityConfirm).Methods(http.MethodPost)

	r.HandleFunc("/v2/activities", s.ActivityList).Methods(http.MethodGet)
	r.HandleFunc("/v2/activity/{uuid}", s.ActivityDelete).Methods(http.MethodDelete)

	r.HandleFunc("/v2/job", s.JobSubmit).Methods(http.MethodPost)
	r.HandleFunc("/v2/batch/job", s.BatchJobSubmit).Methods(http.MethodPost)
	r.HandleFunc("/v2/jobs", s.JobList).Methods(http.MethodGet)
	r.HandleFunc("/v2/job/confirm", s.JobConfirm).Methods(http.MethodPost)
	r.HandleFunc("/v2/job/stop", s.JobStop).Methods(http.MethodPost)
	r.HandleFunc("/v2/job/{jobuid}", s.JobDel).Methods(http.MethodDelete)
	r.HandleFunc("/v2/job/intersect", s.JobIntersect).Methods(http.MethodGet)
	r.HandleFunc("/v2/job/intersect/download", s.JobIntersectDownload).Methods(http.MethodPost)

	r.HandleFunc("/v2/task", s.TaskCreate).Methods(http.MethodPost)
	r.HandleFunc("/v2/tasks", s.TaskList).Methods(http.MethodGet)
	r.HandleFunc("/v2/task/confirm", s.TaskConfirm).Methods(http.MethodPost)
	r.HandleFunc("/v2/task", s.TaskGet).Methods(http.MethodGet)
	r.HandleFunc("/v2/task/start", s.TaskStart).Methods(http.MethodPost)
	r.HandleFunc("/v2/task/stop", s.TaskStop).Methods(http.MethodPost)
	r.HandleFunc("/v2/task/rerun", s.TaskRerun).Methods(http.MethodPost)
	r.HandleFunc("/v2/task/intersect", s.TaskIntersect).Methods(http.MethodGet)
	r.HandleFunc("/v2/task/intersect/download", s.TaskIntersectDownload).Methods(http.MethodGet)
}
