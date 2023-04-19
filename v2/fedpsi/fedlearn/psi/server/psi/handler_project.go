package psi

import (
	"net/http"
	"strconv"

	"fedlearn/psi/api"
	"fedlearn/psi/api/types"
	log "fedlearn/psi/common/log"

	"github.com/gorilla/mux"
)

func (s *Server) ProjectCreate(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "ProjectCreate")
	log.Debugln(ctx, "Server.ProjectCreate is called.")

	var r types.ProjectCreateRequest
	if err := api.ReadObjectFromReqBody(req, &r); err != nil {
		auditContextActionResult(ctx, "ProjectCreate", "Failed")
		api.OutputError(res, err)
		return
	}

	if err := s.ProjectMgr.ProjectCreate(ctx, &r); err != nil {
		auditContextActionResult(ctx, "ProjectCreate", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "ProjectCreate", "Success")
	api.OutputSuccess(res)
	return
}

func (s *Server) ProjectUpdate(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "ProjectUpdate")
	log.Debugln(ctx, "Server.ProjectUpdate is called.")

	var r types.ProjectUpdateRequest
	if err := api.ReadObjectFromReqBody(req, &r); err != nil {
		auditContextActionResult(ctx, "ProjectUpdate", "Failed")
		api.OutputError(res, err)
		return
	}

	if err := s.ProjectMgr.ProjectUpdate(ctx, &r); err != nil {
		auditContextActionResult(ctx, "ProjectUpdate", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "ProjectUpdate", "Success")
	api.OutputSuccess(res)
	return
}

func (s *Server) ProjectList(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "ProjectList")
	log.Debugln(ctx, "Server.ProjectList is called.")

	name := api.GetQueryValue(req, "name")
	result := s.ProjectMgr.ProjectList(ctx, name)

	auditContextActionResult(ctx, "ProjectList", "Success")
	api.OutputObject(res, result)
	return
}

func (s *Server) ProjectDel(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "ProjectDel")
	log.Debugln(ctx, "Server.ProjectDel is called.")

	vars := mux.Vars(req)
	idstr := vars["id"]

	id, err := strconv.Atoi(idstr)
	if err != nil {
		auditContextActionResult(ctx, "ProjectDel", "id format error")
		api.OutputError(res, err)
		return
	}

	if err := s.ProjectMgr.ProjectDel(ctx, uint64(id)); err != nil {
		auditContextActionResult(ctx, "ProjectDel", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "ProjectDel", "Success")
	api.OutputSuccess(res)
	return
}

func (s *Server) ProjectGet(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "ProjectGet")
	log.Debugln(ctx, "Server.ProjectGet is called.")

	vars := mux.Vars(req)
	idstr := vars["id"]

	id, err := strconv.Atoi(idstr)
	if err != nil {
		auditContextActionResult(ctx, "ProjectGet", "id format error")
		api.OutputError(res, err)
		return
	}

	result, err := s.ProjectMgr.ProjectGet(ctx, uint64(id))
	if err != nil {
		auditContextActionResult(ctx, "ProjectGet", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "ProjectGet", "Success")
	api.OutputObject(res, result)
	return
}

func (s *Server) ProjectJobsGet(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "ProjectJobsGet")
	log.Debugln(ctx, "Server.ProjectJobsGet is called.")

	vars := mux.Vars(req)
	puid := vars["puid"]
	page_num := api.GetQueryValue(req, "page_num")
	pageNum, err := strconv.Atoi(page_num)
	if err != nil {
		pageNum = 1
	}
	page_size := api.GetQueryValue(req, "page_size")
	pageSize, err := strconv.Atoi(page_size)
	if err != nil {
		pageSize = 10
	}
	items, count, pageCount := s.ProjectMgr.ProjectJobsGet(ctx, puid, pageNum, pageSize)
	result := types.PageInfo{
		Items:     items,
		Count:     count,
		PageCount: pageCount,
		PageNum:   pageNum,
		PageSize:  pageSize,
	}
	auditContextActionResult(ctx, "ProjectJobsGet", "Success")
	api.OutputObject(res, result)
	return

}

func (s *Server) PartyProjectCreate(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "PartyProjectCreate")
	log.Debugln(ctx, "Server.PartyProjectCreate is called.")

	var r types.ProjectCreateRequest
	if err := api.ReadObjectFromReqBody(req, &r); err != nil {
		auditContextActionResult(ctx, "PartyProjectCreate", "Failed")
		api.OutputError(res, err)
		return
	}

	if err := s.ProjectMgr.PartyProjectCreate(ctx, &r); err != nil {
		auditContextActionResult(ctx, "PartyProjectCreate", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "PartyProjectCreate", "Success")
	api.OutputSuccess(res)
	return
}
