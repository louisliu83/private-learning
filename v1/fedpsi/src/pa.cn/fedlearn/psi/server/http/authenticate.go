package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"pa.cn/fedlearn/psi/api"
	"pa.cn/fedlearn/psi/auth"
)

func authenticationHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		logrus.Infoln("Middleware authenticationHandler")
		credString := req.Header.Get("Authorization")
		if credString == "" {
			if strings.Contains(req.RequestURI, "download") { // process download
				jwtToken := req.URL.Query().Get("jwt")
				if jwtToken == "" {
					logrus.Errorln("Authenticate failed.")
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				credString = fmt.Sprintf("Bearer %s", jwtToken)
			} else {
				logrus.Errorln("Authenticate failed.")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
		userAuthInfo, ok := auth.Authenticate(credString)
		if !ok {
			logrus.Errorln("Authenticate failed.")
			// w.Header().Add("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, ""))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		req.Header.Add(api.ReqHeader_PSIUserID, userAuthInfo.UserName)
		req.Header.Add(api.ReqHeader_PSIUserParty, userAuthInfo.Party)
		req.Header.Add(api.ReqHeader_PSIUserType, userAuthInfo.Type)
		if userAuthInfo.IsAdmin {
			req.Header.Add(api.ReqHeader_PSIUserRole, "admin")
		} else {
			req.Header.Add(api.ReqHeader_PSIUserRole, "member")
		}
		h.ServeHTTP(w, req)
	})
}

func userAuthorizationHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		logrus.Infoln("Middleware userAuthorizationHandler")
		userType := req.Header.Get(api.ReqHeader_PSIUserType)
		if userType != auth.TokenTypeUser {
			logrus.Errorf("Authorization failed: wrong user type")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, req)
	})
}

func partyAuthorizationHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		logrus.Infoln("Middleware partyAuthorizationHandler")
		userType := req.Header.Get(api.ReqHeader_PSIUserType)
		if userType != auth.TokenTypeAgent {
			logrus.Errorf("Authorization failed: wrong user type")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		partyName := req.Header.Get(api.ReqHeader_PSIUserParty)
		vars := mux.Vars(req)
		srcParty := vars["srcParty"]
		if srcParty != partyName {
			logrus.Errorf("Authorization for party %s failed.", srcParty)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, req)
	})
}
