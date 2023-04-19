package http

import (
	"net/http"
	"sync"

	"fedlearn/psi/common/config"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Server is the http server
type Server struct {
}

var (
	server *Server
	once   sync.Once
)

// GetServer return the instance server
func GetServer() *Server {
	once.Do(func() {
		construct()
	})
	return server
}

func construct() *Server {
	return &Server{}
}

// Start  start the controller server
func (ss *Server) Start() {
	r := mux.NewRouter()

	ss.setRoutes(r)
	go func() {
		logrus.Infoln("psi controller server is serving on:", config.GetConfig().Listener.Address)
		if err := http.ListenAndServe(config.GetConfig().Listener.Address, r); err != nil {
			logrus.Fatalln(err)
		}
	}()

	if config.GetConfig().Listener.TLSEnabled {
		go func() {
			logrus.Infoln("psi https server is serving on:", config.GetConfig().Listener.TLSAddress)
			if err := http.ListenAndServeTLS(config.GetConfig().Listener.TLSAddress,
				config.GetConfig().Listener.TLSCertFile, config.GetConfig().Listener.TLSKeyFile, r); err != nil {
				logrus.Fatalln(err)
			}
		}()
	}
}
