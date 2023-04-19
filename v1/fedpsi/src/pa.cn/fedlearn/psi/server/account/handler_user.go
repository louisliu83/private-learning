package account

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"pa.cn/fedlearn/psi/api/types"
	"pa.cn/fedlearn/psi/auth"
	"pa.cn/fedlearn/psi/config"

	"pa.cn/fedlearn/psi/api"
	"pa.cn/fedlearn/psi/log"
)

var (
	NotPermitError = errors.New("Not permit to operate it.")
)

// UserLogin login the current user and return the User Token of the current user
func (s Server) UserLogin(res http.ResponseWriter, req *http.Request) {
	ctx := api.ContextFromReq(req)
	log.Debugln(ctx, "Server.UserLogin is called.")
	userID := req.Header.Get(api.ReqHeader_PSIUserID)
	t, err := s.TokenMgr.GenToken(ctx, userID, time.Duration(30)*time.Minute, auth.TokenTypeUser)
	if err != nil {
		api.OutputError(res, err)
		return
	}
	api.OutputObject(res, t)
	return
}

// TokenGen  only admin can generate party token
func (s Server) TokenGen(res http.ResponseWriter, req *http.Request) {
	ctx := api.ContextFromReq(req)
	log.Debugln(ctx, "Server.TokenGen is called.")

	timeInHoursStr := api.GetQueryValue(req, "hours")
	timeInHours, err := strconv.Atoi(timeInHoursStr)
	if err != nil {
		timeInHours = int(config.GetConfig().TokenSetting.TokenValidInHours)
	}

	who := api.GetQueryValue(req, "user")
	if who == "" {
		api.OutputError(res, NotPermitError)
		return
	}

	if !isAdminRequest(req) {
		api.OutputError(res, errors.New("Only admin can gen token for other users"))
		return
	}

	t, err := s.TokenMgr.GenToken(ctx, who, time.Duration(timeInHours)*time.Hour, auth.TokenTypeAgent)
	if err != nil {
		api.OutputError(res, err)
		return
	}
	api.OutputObject(res, t)
	return
}

// UserRegister only admin can register user
func (s *Server) UserRegister(res http.ResponseWriter, req *http.Request) {
	ctx := api.ContextFromReq(req)
	log.Debugln(ctx, "Server.UserRegister is called.")
	var r types.UserRegisterRequest
	if err := api.ReadObjectFromReqBody(req, &r); err != nil {
		api.OutputError(res, err)
		return
	}

	if !isAdminRequest(req) {
		api.OutputError(res, NotPermitError)
		return
	}

	if err := s.UserMgr.UserRegister(ctx, r); err != nil {
		api.OutputError(res, err)
		return
	}

	api.OutputSuccess(res)
	return
}

func (s *Server) UserList(res http.ResponseWriter, req *http.Request) {
	ctx := api.ContextFromReq(req)
	log.Debugln(ctx, "Server.UserList is called.")

	if !isAdminRequest(req) {
		api.OutputError(res, NotPermitError)
		return
	}

	if users, err := s.UserMgr.UserList(ctx); err != nil {
		api.OutputError(res, err)
	} else {
		api.OutputObject(res, users)
	}

	return
}

func (s *Server) UserDel(res http.ResponseWriter, req *http.Request) {
	ctx := api.ContextFromReq(req)
	log.Debugln(ctx, "Server.UserDel is called.")

	if !isAdminRequest(req) {
		api.OutputError(res, NotPermitError)
		return
	}

	vars := mux.Vars(req)
	userName := vars["username"]

	if err := s.UserMgr.UserDel(ctx, userName); err != nil {
		api.OutputError(res, err)
	} else {
		api.OutputSuccess(res)
	}

	return
}
