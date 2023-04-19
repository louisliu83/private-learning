package psi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"fedlearn/psi/api"
	"fedlearn/psi/api/types"
	log "fedlearn/psi/common/log"

	"github.com/gorilla/mux"
)

func (s *Server) JobSubmit(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "JobSubmit")
	log.Debugln(ctx, "Server.JobSubmit is called.")

	var r types.JobSubmitRequest
	if err := api.ReadObjectFromReqBody(req, &r); err != nil {
		auditContextActionResult(ctx, "JobSubmit", "Failed")
		api.OutputError(res, err)
		return
	}

	if err := s.TaskMgr.JobSubmit(ctx, r); err != nil {
		auditContextActionResult(ctx, "JobSubmit", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "JobSubmit", "Success")
	api.OutputSuccess(res)
	return
}

func (s *Server) BatchJobSubmit(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "BatchJobSubmit")
	log.Debugln(ctx, "Server.BatchJobSubmit is called.")

	var r types.BatchJobSubmitRequest
	if err := api.ReadObjectFromReqBody(req, &r); err != nil {
		auditContextActionResult(ctx, "BatchJobSubmit", "Failed")
		api.OutputError(res, err)
		return
	}

	errors := s.TaskMgr.BatchJobSubmit(ctx, r)
	if len(errors) > 0{
		err, _ := json.Marshal(errors) 
		auditContextActionResult(ctx, "BatchJobSubmit", "Failed")
		api.OutputObject(res, string(err))
		return
	}

	auditContextActionResult(ctx, "BatchJobSubmit", "Success")
	api.OutputSuccess(res)
	return
}

func (s *Server) PartyJobSubmit(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "PartyJobSubmit")
	log.Debugln(ctx, "Server.PartyJobSubmit is called.")

	var r types.JobSubmitRequest
	if err := api.ReadObjectFromReqBody(req, &r); err != nil {
		auditContextActionResult(ctx, "PartyJobSubmit", "Failed")
		api.OutputError(res, err)
		return
	}

	if err := s.TaskMgr.PartyJobSubmit(ctx, r); err != nil {
		auditContextActionResult(ctx, "PartyJobSubmit", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "PartyJobSubmit", "Success")
	api.OutputSuccess(res)
	return
}

func (s *Server) JobList(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "JobList")
	log.Debugln(ctx, "Server.JobList is called.")

	listFunc := func() ([]*types.Job, error) {
		return s.TaskMgr.JobList(ctx)
	}

	dsName := api.GetQueryValue(req, "dsname")
	if dsName != "" {
		listFunc = func() ([]*types.Job, error) {
			return s.TaskMgr.JobListByDataset(ctx, dsName)
		}
	}

	result, err := listFunc()

	if err != nil {
		auditContextActionResult(ctx, "JobList", "Failed")
		api.OutputError(res, err)
		return
	}
	auditContextActionResult(ctx, "JobList", "Success")
	api.OutputObject(res, result)
	return
}

func (s *Server) JobStop(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "JobStop")
	log.Debugln(ctx, "Server.JobStop is called.")

	var r types.JobStopRequest
	if err := api.ReadObjectFromReqBody(req, &r); err != nil {
		auditContextActionResult(ctx, "JobStop", "Failed")
		api.OutputError(res, err)
		return
	}
	if err := s.TaskMgr.JobStop(ctx, r); err != nil {
		auditContextActionResult(ctx, "JobStop", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "JobStop", "Success")
	api.OutputSuccess(res)
	return
}

func (s *Server) JobConfirm(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "JobConfirm")
	log.Debugln(ctx, "Server.JobConfirm is called.")

	var r types.JobConfirmRequest
	if err := api.ReadObjectFromReqBody(req, &r); err != nil {
		auditContextActionResult(ctx, "JobConfirm", "Failed")
		api.OutputError(res, err)
		return
	}
	if err := s.TaskMgr.JobConfirm(ctx, r); err != nil {
		auditContextActionResult(ctx, "JobConfirm", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "JobConfirm", "Success")
	api.OutputSuccess(res)
	return
}

func (s *Server) PartyJobConfirm(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "PartyJobConfirm")
	log.Debugln(ctx, "Server.PartyJobConfirm is called.")

	var r types.JobConfirmRequest
	if err := api.ReadObjectFromReqBody(req, &r); err != nil {
		auditContextActionResult(ctx, "PartyJobConfirm", "Failed")
		api.OutputError(res, err)
		return
	}
	if err := s.TaskMgr.PartyJobConfirm(ctx, r); err != nil {
		auditContextActionResult(ctx, "PartyJobConfirm", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "PartyJobConfirm", "Success")
	api.OutputSuccess(res)
	return
}

func (s *Server) JobIntersect(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "JobIntersect")
	log.Debugln(ctx, "Server.JobIntersect is called.")

	jobUUID := api.GetQueryValue(req, "job_uuid")
	data, err := s.TaskMgr.JobIntersectRead(ctx, jobUUID)
	if err != nil {
		auditContextActionResult(ctx, "JobIntersect", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "JobIntersect", "Success")
	api.OutputObject(res, string(data))
	return
}

func (s *Server) PartyJobIntersect(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "PartyJobIntersect")
	log.Debugln(ctx, "Server.PartyJobIntersect is called.")

	jobUUID := api.GetQueryValue(req, "job_uuid")
	data, err := s.TaskMgr.JobIntersectRead(ctx, jobUUID)
	if err != nil {
		auditContextActionResult(ctx, "PartyJobIntersect", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "PartyJobIntersect", "Success")
	api.OutputObject(res, string(data))
	return
}

func (s *Server) JobIntersectDownload(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "JobIntersectDownload")
	log.Debugln(ctx, "Server.JobIntersectDownload is called.")

	//jobUUID := api.GetQueryValue(req, "job_uuid")
	var r types.JobDownloadRequest
	if err := api.ReadObjectFromReqBody(req, &r); err != nil {
		auditContextActionResult(ctx, "JobIntersectDownload", "Failed")
		api.OutputError(res, err)
		return
	}

	jobUUID := r.JobUID
	data, err := s.TaskMgr.JobIntersectRead(ctx, jobUUID)
	if err != nil {
		auditContextActionResult(ctx, "JobIntersectDownload", "Failed")
		api.OutputError(res, err)
		return
	}
	auditContextActionResult(ctx, "JobIntersectDownload", "Success")
	fileName := fmt.Sprintf("intersect_%s.txt", jobUUID)
	//res.Header().Add("Content-Length", data.)
	res.Header().Add("Content-Disposition", fmt.Sprintf("attachment;filename=%s", fileName))
	res.Header().Add("Content-Type", "application/octet-stream")
	res.Write(data)
	return
}

func (s *Server) JobDel(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "JobDel")
	log.Debugln(ctx, "Server.JobDel is called.")

	vars := mux.Vars(req)
	jobUID := vars["jobuid"]

	r := types.JobDelRequest{
		JobUID: jobUID,
	}

	if err := s.TaskMgr.JobDel(ctx, r); err != nil {
		auditContextActionResult(ctx, "JobDel", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "JobDel", "Success")
	api.OutputSuccess(res)
	return
}
