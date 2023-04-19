package account

import (
	"pa.cn/fedlearn/psi/server/manager"
)

type Server struct {
	UserMgr  *manager.UserManager
	TokenMgr *manager.TokenManager
}

func New() *Server {
	return &Server{
		UserMgr:  manager.GetUserManager(),
		TokenMgr: manager.GetTokenManager(),
	}
}
