package http

import (
	"net/http"

	accountserver "fedlearn/psi/server/account"

	"github.com/gorilla/mux"
)

func (ss *Server) setAccountRoutes(r *mux.Router) {
	s := accountserver.New()
	//r.HandleFunc("/v2/user/register", s.UserRegister).Methods(http.MethodPost)
	r.HandleFunc("/v2/user/login", s.UserLogin).Methods(http.MethodPost)
	r.HandleFunc("/v2/user/token", s.TokenGen).Methods(http.MethodGet)
	r.HandleFunc("/v2/users", s.UserList).Methods(http.MethodGet)
	r.HandleFunc("/v2/user/{username}", s.UserDel).Methods(http.MethodDelete)

}
