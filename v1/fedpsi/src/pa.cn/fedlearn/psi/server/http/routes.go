package http

import (
	"github.com/gorilla/mux"
	"pa.cn/fedlearn/psi/config"
)

// SetRoutes set the routes
func (ss *Server) setRoutes(r *mux.Router) {
	// First, clean all the authenticate HEADER

	// ReqHeader_PSIUserID    = "PSI-User-ID"
	// ReqHeader_PSIUserParty = "PSI-User-Party"
	// ReqHeader_PSIUserRole  = "PSI-User-Role"
	// ReqHeader_PSIUserType  = "PSI-User-Type"

	r.Use(cleanIDHeaderHandler)

	apiRouter := r.PathPrefix("/apis").Subrouter()
	apiRouter.Use(genTraceIDHandler)
	if config.GetConfig().TokenSetting.AuthEnabled {
		apiRouter.Use(authenticationHandler, userAuthorizationHandler)
	}
	apiRouter.Use(auditHandler)
	ss.setAPIRoutes(apiRouter)

	p2pRouter := r.PathPrefix("/p2p").Subrouter()
	p2pRouter.Use(genTraceIDHandler)
	if config.GetConfig().TokenSetting.AuthEnabled {
		p2pRouter.Use(authenticationHandler, partyAuthorizationHandler)
	}
	p2pRouter.Use(auditHandler)
	ss.setParty2PartyRoutes(p2pRouter)
}

func (ss *Server) setAPIRoutes(r *mux.Router) {
	ss.setPsiRoutes(r)
	ss.setAccountRoutes(r)
	ss.setPromRoutes(r)
}
