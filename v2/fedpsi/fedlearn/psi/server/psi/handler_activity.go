package psi

import (
	"fedlearn/psi/api"
	"fedlearn/psi/api/types"
	"fedlearn/psi/common/log"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) ActivityAdd(res http.ResponseWriter, req *http.Request) {
	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "ActivityAdd")
	log.Debugln(ctx, "Server.ActivityAdd is called.")

	var r types.ActivityCreateRequest
	if err := api.ReadObjectFromReqBody(req, &r); err != nil {
		auditContextActionResult(ctx, "ActivityAdd", "Failed")
		api.OutputError(res, err)
		return
	}

	if err := s.TaskMgr.ActivityCreate(ctx, r); err != nil {
		auditContextActionResult(ctx, "ActivityAdd", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "ActivityAdd", "Success")
	api.OutputSuccess(res)
	return
}

func (s *Server) ActivityList(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "ActivityList")
	log.Debugln(ctx, "Server.ActivityList is called.")

	result, err := s.TaskMgr.ActivityList(ctx)
	if err != nil {
		auditContextActionResult(ctx, "ActivityList", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "ActivityList", "Success")
	api.OutputObject(res, result)
	return
}

func (s *Server) ActivityDelete(res http.ResponseWriter, req *http.Request) {
	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "ActivityDelete")
	log.Debugln(ctx, "Server.ActivityDelete is called.")
	vars := mux.Vars(req)
	uuid := vars["uuid"]
	r := types.ActivityDeleteRequest{
		Uuid: uuid,
	}
	if err := s.TaskMgr.ActivityDelete(ctx, r); err != nil {
		auditContextActionResult(ctx, "ActivityDelete", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "ActivityDelete", "Success")
	api.OutputSuccess(res)
	return
}

func (s *Server) ActivityAttachData(res http.ResponseWriter, req *http.Request) {
	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "ActivityAttachData")
	log.Debugln(ctx, "Server.ActivityAttachData is called.")

	var r types.ActivityAttachDataRequest
	if err := api.ReadObjectFromReqBody(req, &r); err != nil {
		auditContextActionResult(ctx, "ActivityAttachData", "Failed")
		api.OutputError(res, err)
		return
	}

	if err := s.TaskMgr.ActivityAttachData(ctx, r); err != nil {
		auditContextActionResult(ctx, "ActivityAttachData", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "ActivityAttachData", "Success")
	api.OutputSuccess(res)
	return
}

func (s *Server) ActivityStart(res http.ResponseWriter, req *http.Request) {
	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "ActivityStart")
	log.Debugln(ctx, "Server.ActivityStart is called.")

	vars := mux.Vars(req)
	uuid := vars["uuid"]

	if err := s.TaskMgr.ActivityStart(ctx, uuid); err != nil {
		auditContextActionResult(ctx, "ActivityStart", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "ActivityStart", "Success")
	api.OutputSuccess(res)
	return
}

func (s *Server) ActivityConfirm(res http.ResponseWriter, req *http.Request) {
	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "ActivityConfirm")
	log.Debugln(ctx, "Server.ActivityConfirm is called.")

	var r types.ActivityConfirmRequest
	if err := api.ReadObjectFromReqBody(req, &r); err != nil {
		auditContextActionResult(ctx, "ActivityConfirm", "Failed")
		api.OutputError(res, err)
		return
	}

	if err := s.TaskMgr.ActivityConfirm(ctx, r); err != nil {
		auditContextActionResult(ctx, "ActivityConfirm", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "ActivityConfirm", "Success")
	api.OutputSuccess(res)
	return
}

///////################
func (s *Server) PartyActivityAdd(res http.ResponseWriter, req *http.Request) {
	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "PartyActivityAdd")
	log.Debugln(ctx, "Server.PartyActivityAdd is called.")

	var r types.ActivityCreateRequest
	if err := api.ReadObjectFromReqBody(req, &r); err != nil {
		auditContextActionResult(ctx, "PartyActivityAdd", "Failed")
		api.OutputError(res, err)
		return
	}

	if err := s.TaskMgr.PartyActivityCreate(ctx, r); err != nil {
		auditContextActionResult(ctx, "PartyActivityAdd", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "PartyActivityAdd", "Success")
	api.OutputSuccess(res)
	return
}

func (s *Server) PartyActivityConfirm(res http.ResponseWriter, req *http.Request) {
	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "PartyActivityConfirm")
	log.Debugln(ctx, "Server.PartyActivityConfirm is called.")

	var r types.ActivityConfirmRequest
	if err := api.ReadObjectFromReqBody(req, &r); err != nil {
		auditContextActionResult(ctx, "PartyActivityConfirm", "Failed")
		api.OutputError(res, err)
		return
	}

	if err := s.TaskMgr.PartyActivityConfirm(ctx, r); err != nil {
		auditContextActionResult(ctx, "PartyActivityConfirm", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "PartyActivityConfirm", "Success")
	api.OutputSuccess(res)
	return
}
