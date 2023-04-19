package http

import (
	"fmt"
	"net/http"
	"strings"

	"fedlearn/psi/common/config"
	"fedlearn/psi/common/utils"

	"fedlearn/psi/api"
	"fedlearn/psi/audit"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func auditHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		logrus.Infoln("Middleware auditHandler")
		partyName := req.Header.Get(api.ReqHeader_PSIUserParty)
		if partyName == "" {
			if strings.Contains(req.RequestURI, "p2p") {
				vars := mux.Vars(req)
				partyName = vars["srcParty"]
			} else {
				partyName = config.GetConfig().PartyName
			}
			req.Header.Set(api.ReqHeader_PSIUserParty, partyName)
		}
		user := req.Header.Get(api.ReqHeader_PSIUserID)
		if user == "" {
			user = "anonymous"
		}
		who := utils.GetWho(user, partyName)
		action := fmt.Sprintf("%s %s", req.Method, req.RequestURI)
		audit.Log(who, action, "", "Accept request")
		h.ServeHTTP(w, req)
		audit.Log(who, action, "", "Return response")
	})
}

func genTraceIDHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		logrus.Infoln("Middleware genTraceIDHandler")

		traceID := req.Header.Get(api.Trace_Req_ID)
		if traceID == "" {
			traceID = fmt.Sprintf("Trace_Http_%s", utils.UUIDStr())
		}
		logrus.Debugln("Use passed-in traceID ", traceID)
		req.Header.Set(api.Trace_Req_ID, traceID)
		h.ServeHTTP(w, req)
	})
}

func cleanIDHeaderHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		logrus.Infoln("Middleware cleanIDHeaderHandler")

		req.Header.Del(api.ReqHeader_PSIUserID)
		req.Header.Del(api.ReqHeader_PSIUserParty)
		req.Header.Del(api.ReqHeader_PSIUserRole)
		req.Header.Del(api.ReqHeader_PSIUserType)

		h.ServeHTTP(w, req)
	})
}
