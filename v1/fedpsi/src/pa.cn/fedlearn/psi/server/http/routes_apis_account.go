package http

import (
	"net/http"

	"github.com/gorilla/mux"
	accountserver "pa.cn/fedlearn/psi/server/account"
)

func (ss *Server) setAccountRoutes(r *mux.Router) {
	s := accountserver.New()
	//r.HandleFunc("/v1/user/register", s.UserRegister).Methods(http.MethodPost)
	r.HandleFunc("/v1/user/login", s.UserLogin).Methods(http.MethodPost)
	r.HandleFunc("/v1/user/token", s.TokenGen).Methods(http.MethodGet)
	r.HandleFunc("/v1/users", s.UserList).Methods(http.MethodGet)
	r.HandleFunc("/v1/user/{username}", s.UserDel).Methods(http.MethodDelete)

}
