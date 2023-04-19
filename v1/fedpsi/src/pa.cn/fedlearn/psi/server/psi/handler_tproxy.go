package psi

import (
	"net/http"

	"pa.cn/fedlearn/psi/api"
	"pa.cn/fedlearn/psi/api/types"
	"pa.cn/fedlearn/psi/log"
)

func (s *Server) TProxyStart(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "TProxyStart")
	log.Debugln(ctx, "Server.TProxyStart is called.")

	var r types.TProxyStartRequest
	if err := api.ReadObjectFromReqBody(req, &r); err != nil {
		api.OutputError(res, err)
		return
	}

	if err := s.TProxyMgr.Start(ctx, r); err != nil {
		api.OutputError(res, err)
		return
	}

	api.OutputSuccess(res)
	return
}

func (s *Server) TProxyStop(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "TProxyStop")
	log.Debugln(ctx, "Server.TProxyStop is called.")

	var r types.TProxyStopRequest
	if err := api.ReadObjectFromReqBody(req, &r); err != nil {
		api.OutputError(res, err)
		return
	}

	if err := s.TProxyMgr.Stop(ctx, r); err != nil {
		api.OutputError(res, err)
		return
	}

	api.OutputSuccess(res)
	return
}

func (s *Server) TProxyStatus(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "TProxyStatus")
	log.Debugln(ctx, "Server.TProxyStatus is called.")

	proxyStatus := s.TProxyMgr.Status(ctx)
	auditContextActionResult(ctx, "TProxyStatus", "Success")
	api.OutputObject(res, proxyStatus)
	return
}
