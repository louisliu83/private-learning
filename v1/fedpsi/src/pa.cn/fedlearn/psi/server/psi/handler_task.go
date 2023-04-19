package psi

import (
	"fmt"
	"net/http"

	"pa.cn/fedlearn/psi/api"
	"pa.cn/fedlearn/psi/api/types"
	log "pa.cn/fedlearn/psi/log"
)

func (s *Server) TaskCreate(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "TaskCreate")
	log.Debugln(ctx, "Server.TaskCreate is called.")

	var r types.TaskCreateRequest
	if err := api.ReadObjectFromReqBody(req, &r); err != nil {
		auditContextActionResult(ctx, "TaskCreate", "Failed")
		api.OutputError(res, err)
		return
	}

	if err := s.TaskMgr.TaskCreate(ctx, r); err != nil {
		auditContextActionResult(ctx, "TaskCreate", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "TaskCreate", "Success")
	api.OutputSuccess(res)
	return
}

func (s *Server) TaskCreateV2(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "TaskCreateV2")
	log.Debugln(ctx, "Server.TaskCreateV2 is called.")

	var r types.TaskCreateRequestV2
	if err := api.ReadObjectFromReqBody(req, &r); err != nil {
		auditContextActionResult(ctx, "TaskCreateV2", "Failed")
		api.OutputError(res, err)
		return
	}

	if err := s.TaskMgr.TaskCreateV2(ctx, r); err != nil {
		auditContextActionResult(ctx, "TaskCreateV2", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "TaskCreateV2", "Success")
	api.OutputSuccess(res)
	return
}

func (s *Server) PartyTaskCreate(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "PartyTaskCreate")
	log.Debugln(ctx, "Server.PartyTaskCreate is called.")

	var r types.TaskCreateRequest
	if err := api.ReadObjectFromReqBody(req, &r); err != nil {
		auditContextActionResult(ctx, "PartyTaskCreate", "Failed")
		api.OutputError(res, err)
		return
	}

	if err := s.TaskMgr.PartyTaskCreate(ctx, r); err != nil {
		auditContextActionResult(ctx, "PartyTaskCreate", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "PartyTaskCreate", "Success")
	api.OutputSuccess(res)
	return
}

func (s *Server) PartyTaskCreateV2(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "PartyTaskCreateV2")
	log.Debugln(ctx, "Server.PartyTaskCreateV2 is called.")

	var r types.TaskCreateRequestV2
	if err := api.ReadObjectFromReqBody(req, &r); err != nil {
		auditContextActionResult(ctx, "PartyTaskCreateV2", "Failed")
		api.OutputError(res, err)
		return
	}

	if err := s.TaskMgr.PartyTaskCreateV2(ctx, r); err != nil {
		auditContextActionResult(ctx, "PartyTaskCreateV2", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "PartyTaskCreateV2", "Success")
	api.OutputSuccess(res)
	return
}

func (s *Server) TaskList(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "TaskList")
	log.Debugln(ctx, "Server.TaskList is called.")

	jobUUID := api.GetQueryValue(req, "job_uuid")
	r := types.TaskListRequest{
		JobUID: jobUUID,
	}
	result, err := s.TaskMgr.TaskList(ctx, r)
	if err != nil {
		auditContextActionResult(ctx, "TaskList", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "TaskList", "Success")
	api.OutputObject(res, result)
	return
}

func (s *Server) TaskStart(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "TaskStart")
	log.Debugln(ctx, "Server.TaskStart is called.")

	var r types.TaskStartRequest
	if err := api.ReadObjectFromReqBody(req, &r); err != nil {
		auditContextActionResult(ctx, "TaskStart", "Failed")
		api.OutputError(res, err)
		return
	}

	if err := s.TaskMgr.TaskStart(ctx, r); err != nil {
		auditContextActionResult(ctx, "TaskStart", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "TaskStart", "Success")
	api.OutputSuccess(res)
	return
}

func (s *Server) PartyTaskStart(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "PartyTaskStart")
	log.Debugln(ctx, "Server.PartyTaskStart is called.")

	var r types.TaskStartRequest
	if err := api.ReadObjectFromReqBody(req, &r); err != nil {
		auditContextActionResult(ctx, "PartyTaskStart", "Failed")
		api.OutputError(res, err)
		return
	}

	if err := s.TaskMgr.PartyTaskStart(ctx, r); err != nil {
		auditContextActionResult(ctx, "PartyTaskStart", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "PartyTaskStart", "Success")
	api.OutputSuccess(res)
	return
}

func (s *Server) TaskConfirm(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "TaskConfirm")
	log.Debugln(ctx, "Server.TaskConfirm is called.")

	var r types.TaskConfirmRequest
	if err := api.ReadObjectFromReqBody(req, &r); err != nil {
		auditContextActionResult(ctx, "TaskConfirm", "Failed")
		api.OutputError(res, err)
		return
	}

	if err := s.TaskMgr.TaskConfirm(ctx, r); err != nil {
		auditContextActionResult(ctx, "TaskConfirm", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "TaskConfirm", "Success")
	api.OutputSuccess(res)
	return
}

func (s *Server) PartyTaskConfirm(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "PartyTaskConfirm")
	log.Debugln(ctx, "Server.PartyTaskConfirm is called.")

	var r types.TaskConfirmRequest
	if err := api.ReadObjectFromReqBody(req, &r); err != nil {
		auditContextActionResult(ctx, "PartyTaskConfirm", "Failed")
		api.OutputError(res, err)
		return
	}

	if err := s.TaskMgr.PartyTaskConfirm(ctx, r); err != nil {
		auditContextActionResult(ctx, "PartyTaskConfirm", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "PartyTaskConfirm", "Success")
	api.OutputSuccess(res)
	return
}

func (s *Server) TaskStop(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "TaskStop")
	log.Debugln(ctx, "Server.TaskStop is called.")

	var r types.TaskStopRequest
	if err := api.ReadObjectFromReqBody(req, &r); err != nil {
		auditContextActionResult(ctx, "TaskStop", "Failed")
		api.OutputError(res, err)
		return
	}

	if err := s.TaskMgr.TaskStop(ctx, r); err != nil {
		auditContextActionResult(ctx, "TaskStop", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "TaskStop", "Success")
	api.OutputSuccess(res)
	return
}

func (s *Server) PartyTaskStop(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "PartyTaskStop")
	log.Debugln(ctx, "Server.PartyTaskStop is called.")

	var r types.TaskStopRequest
	if err := api.ReadObjectFromReqBody(req, &r); err != nil {
		auditContextActionResult(ctx, "PartyTaskStop", "Failed")
		api.OutputError(res, err)
		return
	}

	if err := s.TaskMgr.PartyTaskStop(ctx, r); err != nil {
		auditContextActionResult(ctx, "PartyTaskStop", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "PartyTaskStop", "Success")
	api.OutputSuccess(res)
	return
}

func (s *Server) TaskGet(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "TaskGet")
	log.Debugln(ctx, "Server.TaskGet is called.")

	taskUUID := api.GetQueryValue(req, "task_uuid")

	r := types.TaskGetRequest{
		TaskUID: taskUUID,
	}

	var task *types.Task
	var err error
	if task, err = s.TaskMgr.TaskGet(ctx, r); err != nil {
		auditContextActionResult(ctx, "TaskGet", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "TaskGet", "Success")
	api.OutputObject(res, task)
	return
}

func (s *Server) PartyTaskGet(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "PartyTaskGet")
	log.Debugln(ctx, "Server.PartyTaskGet is called.")

	taskUUID := api.GetQueryValue(req, "task_uuid")

	r := types.TaskGetRequest{
		TaskUID: taskUUID,
	}

	var task *types.Task
	var err error
	if task, err = s.TaskMgr.TaskGet(ctx, r); err != nil {
		auditContextActionResult(ctx, "PartyTaskGet", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "PartyTaskGet", "Success")
	api.OutputObject(res, task)
	return
}

func (s *Server) TaskIntersect(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "TaskIntersect")
	log.Debugln(ctx, "Server.TaskIntersect is called.")

	taskUUID := api.GetQueryValue(req, "task_uuid")

	r := types.TaskIntersectionDownloadRequest{
		TaskUID: taskUUID,
	}

	data, err := s.TaskMgr.TaskIntersectRead(ctx, r)
	if err != nil {
		auditContextActionResult(ctx, "TaskIntersect", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "TaskIntersect", "Success")
	api.OutputObject(res, string(data))
	return
}

func (s *Server) PartyTaskIntersect(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "PartyTaskIntersect")
	log.Debugln(ctx, "Server.PartyTaskIntersect is called.")

	taskUUID := api.GetQueryValue(req, "task_uuid")

	r := types.TaskIntersectionDownloadRequest{
		TaskUID: taskUUID,
	}

	data, err := s.TaskMgr.TaskIntersectRead(ctx, r)
	if err != nil {
		auditContextActionResult(ctx, "PartyTaskIntersect", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "PartyTaskIntersect", "Success")
	api.OutputObject(res, string(data))
	return
}

func (s *Server) TaskIntersectDownload(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "TaskIntersectDownload")
	log.Debugln(ctx, "Server.TaskIntersectDownload is called.")

	taskUUID := api.GetQueryValue(req, "task_uuid")
	r := types.TaskIntersectionDownloadRequest{
		TaskUID: taskUUID,
	}

	data, err := s.TaskMgr.TaskIntersectRead(ctx, r)
	if err != nil {
		auditContextActionResult(ctx, "TaskIntersectDownload", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "TaskIntersectDownload", "Success")
	fileName := fmt.Sprintf("intersect_%s.txt", taskUUID)
	res.Header().Add("Content-Disposition", fmt.Sprintf("attachment;filename=%s", fileName))
	res.Header().Add("Content-Type", "application/octet-stream")
	res.Write(data)
	return
}

func (s *Server) TaskRerun(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "TaskRerun")
	log.Debugln(ctx, "Server.TaskRerun is called.")

	var r types.TaskRerunRequest
	if err := api.ReadObjectFromReqBody(req, &r); err != nil {
		auditContextActionResult(ctx, "TaskRerun", "Failed")
		api.OutputError(res, err)
		return
	}

	if err := s.TaskMgr.TaskRerun(ctx, r); err != nil {
		auditContextActionResult(ctx, "TaskRerun", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "TaskRerun", "Success")
	api.OutputSuccess(res)
	return
}

func (s *Server) PartyTaskRerun(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "PartyTaskRerun")
	log.Debugln(ctx, "Server.PartyTaskRerun is called.")

	var r types.TaskRerunRequest
	if err := api.ReadObjectFromReqBody(req, &r); err != nil {
		auditContextActionResult(ctx, "PartyTaskRerun", "Failed")
		api.OutputError(res, err)
		return
	}

	if err := s.TaskMgr.PartyTaskRerun(ctx, r); err != nil {
		auditContextActionResult(ctx, "PartyTaskRerun", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "PartyTaskRerun", "Success")
	api.OutputSuccess(res)
	return
}
