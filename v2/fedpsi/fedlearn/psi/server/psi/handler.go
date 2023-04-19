package psi

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"fedlearn/psi/api"
	"fedlearn/psi/api/types"
	"fedlearn/psi/audit"
	"fedlearn/psi/common/config"
	"fedlearn/psi/common/log"
	"fedlearn/psi/common/utils"
	"fedlearn/psi/common/version"
)

func (s *Server) Info(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "Info")
	log.Debugln(ctx, "Server.Info is called.")

	serverInfo := types.ServerInfo{
		Version:   version.Version.String(),
		Status:    "UP",
		PartyName: config.GetConfig().PartyName,
		Protocols: []string{"OT", "DH"},
	}

	auditContextActionResult(ctx, "Info", "Success")
	api.OutputObject(res, &serverInfo)
	return
}

func (s *Server) WorkerInfo(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "WorkerInfo")
	log.Debugln(ctx, "Server.WorkerInfo is called.")

	serverInfo := types.ServerInfo{
		Version:   version.Version.String(),
		PartyName: config.GetConfig().PartyName,
		Protocols: []string{"OT", "DH"},
	}
	worker := api.GetQueryValue(req, "worker")
	port := api.GetQueryValue(req, "port")
	address := net.JoinHostPort(worker, port)
	conn, err := net.DialTimeout("tcp", address, 3*time.Second)
	if err != nil {
		serverInfo.Status = "DOWN"
	} else {
		if conn != nil {
			serverInfo.Status = "UP"
			_ = conn.Close()
		} else {
			serverInfo.Status = "DOWN"
		}
	}
	auditContextActionResult(ctx, "WorkerInfo", "Success")
	api.OutputObject(res, &serverInfo)
	return
}

func (s *Server) ConfigInfo(res http.ResponseWriter, req *http.Request) {
	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "ConfigInfo")
	log.Debugln(ctx, "Server.ConfigInfo is called.")

	cfgInfo := types.ConfigInfo{}

	auditContextActionResult(ctx, "ConfigInfo", "Success")
	api.OutputObject(res, &cfgInfo)
	return
}

func auditContextAction(ctx context.Context, action string) {
	auditContextActionResult(ctx, action, "")
}

func auditContextActionResult(ctx context.Context, action, result string) {
	auditContextActionResultMessage(ctx, action, result, "")
}

func auditContextActionResultMessage(ctx context.Context, action, result, msg string) {
	partyName := ctx.Value(api.ReqHeader_PSIUserParty)
	user := ctx.Value(api.ReqHeader_PSIUserID)
	who := utils.GetWho(fmt.Sprintf("%s", user), fmt.Sprintf("%s", partyName))
	audit.Log(who, action, result, msg)
}
