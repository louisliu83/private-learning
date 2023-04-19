package psi

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"pa.cn/fedlearn/psi/api"
	"pa.cn/fedlearn/psi/api/types"
	"pa.cn/fedlearn/psi/log"
)

func (s *Server) DataSetDel(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "DataSetDel")
	log.Debugln(ctx, "Server.DataSetDel is called.")

	var r types.DataSetDeleteRequest
	vars := mux.Vars(req)
	name := vars["name"]
	strIndex := vars["index"]
	index, err := strconv.Atoi(strIndex)
	if err != nil {
		index = 0
	}
	r.Name = name
	r.Index = int32(index)

	if err := s.DSMgr.DataSetDelete(ctx, r); err != nil {
		auditContextActionResult(ctx, "DataSetDel", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "DataSetDel", "Success")
	api.OutputSuccess(res)
	return
}

func (s *Server) DataSetList(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "DataSetList")
	log.Debugln(ctx, "Server.DataSetList is called.")

	var r types.DataSetListRequest
	partyName := api.GetQueryValue(req, "party")
	r.PartyNames = make([]string, 0)
	if partyName != "" {
		r.PartyNames = append(r.PartyNames, partyName)
	}

	result, err := s.DSMgr.DataSetList(ctx, r)
	if err != nil {
		auditContextActionResult(ctx, "DataSetList", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "DataSetList", "Success")
	api.OutputObject(res, result)
	return
}

func (s *Server) DataSetGet(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "DataSetGet")
	log.Debugln(ctx, "Server.DataSetGet is called.")

	vars := mux.Vars(req)
	name := vars["name"]
	strIndex := vars["index"]
	index, err := strconv.Atoi(strIndex)
	if err != nil {
		index = 0
	}

	result, err := s.DSMgr.DataSetGet(ctx, name, int32(index))
	if err != nil {
		auditContextActionResult(ctx, "DataSetGet", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "DataSetGet", "Success")
	api.OutputObject(res, result)
	return
}

func (s *Server) DataSetShards(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "DataSetShards")
	log.Debugln(ctx, "Server.DataSetShards is called.")

	datasetName := api.GetQueryValue(req, "name")
	result, err := s.DSMgr.DataSetShardsList(ctx, datasetName)
	if err != nil {
		auditContextActionResult(ctx, "DataSetShards", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "DataSetShards", "Success")
	api.OutputObject(res, result)
	return
}

func (s *Server) DataSetGrant(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "DataSetGrant")
	log.Debugln(ctx, "Server.DataSetGrant is called.")

	var r types.DatasetGrantRequest
	if err := api.ReadObjectFromReqBody(req, &r); err != nil {
		api.OutputError(res, err)
		return
	}

	if err := s.DSMgr.DataSetGrant(ctx, r); err != nil {
		api.OutputError(res, err)
		return
	}

	api.OutputSuccess(res)
	return
}

func (s *Server) DataSetRevoke(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "DataSetRevoke")
	log.Debugln(ctx, "Server.DataSetRevoke is called.")

	var r types.DatasetRevokeRequest
	if err := api.ReadObjectFromReqBody(req, &r); err != nil {
		api.OutputError(res, err)
		return
	}

	if err := s.DSMgr.DataSetRevoke(ctx, r); err != nil {
		api.OutputError(res, err)
		return
	}

	api.OutputSuccess(res)
	return
}

//------------------------------ PartyDataset funcs
func (s *Server) PartyDataSetList(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "PartyDataSetList")
	log.Debugln(ctx, "Server.PartyDataSetList is called.")

	var r types.DataSetListRequest
	r.PartyNames = make([]string, 0)

	result, err := s.DSMgr.PartyDataSetList(ctx, r)
	if err != nil {
		auditContextActionResult(ctx, "PartyDataSetList", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "PartyDataSetList", "Success")
	api.OutputObject(res, result)
	return
}

func (s *Server) PartyDataSetGet(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "PartyDataSetGet")
	log.Debugln(ctx, "Server.PartyDataSetGet is called.")

	vars := mux.Vars(req)
	name := vars["name"]
	strIndex := vars["index"]
	index, err := strconv.Atoi(strIndex)
	if err != nil {
		index = 0
	}

	result, err := s.DSMgr.PartyDataSetGet(ctx, name, int32(index))
	if err != nil {
		auditContextActionResult(ctx, "PartyDataSetGet", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "PartyDataSetGet", "Success")
	api.OutputObject(res, result)
	return
}

func (s *Server) PartyDataSetShards(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "PartyDataSetShards")
	log.Debugln(ctx, "Server.PartyDataSetShards is called.")

	datasetName := api.GetQueryValue(req, "name")
	result, err := s.DSMgr.DataSetShardsList(ctx, datasetName)
	if err != nil {
		auditContextActionResult(ctx, "PartyDataSetShards", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "PartyDataSetShards", "Success")
	api.OutputObject(res, result)
	return
}
