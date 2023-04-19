package psi

import (
	"net/http"

	"fedlearn/psi/api"
	"fedlearn/psi/api/types"
	log "fedlearn/psi/common/log"

	"github.com/gorilla/mux"
)

func (s *Server) PartyRegister(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "PartyRegister")
	log.Debugln(ctx, "Server.PartyRegister is called.")

	var r types.PartyRegisterRequest
	if err := api.ReadObjectFromReqBody(req, &r); err != nil {
		auditContextActionResult(ctx, "PartyRegister", "Failed")
		api.OutputError(res, err)
		return
	}

	if err := s.PartyMgr.PartyRegister(ctx, &r); err != nil {
		auditContextActionResult(ctx, "PartyRegister", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "PartyRegister", "Success")
	api.OutputSuccess(res)
	return
}

func (s *Server) PartyUpdate(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "PartyUpdate")
	log.Debugln(ctx, "Server.PartyUpdate is called.")

	var r types.PartyUpdateRequest
	if err := api.ReadObjectFromReqBody(req, &r); err != nil {
		auditContextActionResult(ctx, "PartyUpdate", "Failed")
		api.OutputError(res, err)
		return
	}

	if err := s.PartyMgr.PartyUpdate(ctx, &r); err != nil {
		auditContextActionResult(ctx, "PartyUpdate", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "PartyUpdate", "Success")
	api.OutputSuccess(res)
	return
}

func (s *Server) PartyList(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "PartyList")
	log.Debugln(ctx, "Server.PartyList is called.")

	result := s.PartyMgr.PartyList(ctx)

	auditContextActionResult(ctx, "PartyList", "Success")
	api.OutputObject(res, result)
	return
}

func (s *Server) PartyDel(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "PartyDel")
	log.Debugln(ctx, "Server.PartyDel is called.")

	vars := mux.Vars(req)
	name := vars["name"]

	if err := s.PartyMgr.PartyDel(ctx, name); err != nil {
		auditContextActionResult(ctx, "PartyDel", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "PartyDel", "Success")
	api.OutputSuccess(res)
	return
}

func (s *Server) PartyReady(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "PartyReady")
	log.Debugln(ctx, "Server.PartyReady is called.")

	vars := mux.Vars(req)
	name := vars["name"]

	if err := s.PartyMgr.PartyReady(ctx, name); err != nil {
		auditContextActionResult(ctx, "PartyReady", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "PartyReady", "Success")
	api.OutputSuccess(res)
	return
}

func (s *Server) PartyWorkerReady(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "PartyWorkerReady")
	log.Debugln(ctx, "Server.PartyWorkerReady is called.")

	vars := mux.Vars(req)
	name := vars["name"]

	if err := s.PartyMgr.PartyWorkerReady(ctx, name); err != nil {
		auditContextActionResult(ctx, "PartyWorkerReady", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "PartyWorkerReady", "Success")
	api.OutputSuccess(res)
	return
}
