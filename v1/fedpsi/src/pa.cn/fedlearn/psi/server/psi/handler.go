package psi

import (
	"context"
	"fmt"
	"net/http"

	"pa.cn/fedlearn/psi/api"
	"pa.cn/fedlearn/psi/api/types"
	"pa.cn/fedlearn/psi/audit"
	"pa.cn/fedlearn/psi/config"
	"pa.cn/fedlearn/psi/log"
	"pa.cn/fedlearn/psi/utils"
	"pa.cn/fedlearn/psi/version"
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

func (s *Server) ConfigInfo(res http.ResponseWriter, req *http.Request) {
	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "ConfigInfo")
	log.Debugln(ctx, "Server.ConfigInfo is called.")

	cfgInfo := types.ConfigInfo{}
	cfgInfo.Listen = config.GetConfig().TProxy.Listen
	cfgInfo.Target = config.GetConfig().TProxy.Target
	cfgInfo.DialTimeout = config.GetConfig().TProxy.DialTimeout
	cfgInfo.KeepAlivePeriod = config.GetConfig().TProxy.KeepAlivePeriod
	cfgInfo.ServerWaitDataTimeout = config.GetConfig().TProxy.ServerWaitDataTimeout
	cfgInfo.BinPath = config.GetConfig().PsiExecutor.BinPath
	cfgInfo.PublicIP = config.GetConfig().PsiExecutor.PublicIP
	cfgInfo.PrivateIP = config.GetConfig().PsiExecutor.PrivateIP
	cfgInfo.PublicPort = config.GetConfig().PsiExecutor.PublicPort
	cfgInfo.PrivatePort = config.GetConfig().PsiExecutor.PrivatePort

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
