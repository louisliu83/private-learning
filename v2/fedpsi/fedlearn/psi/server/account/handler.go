package account

import (
	"net/http"

	"fedlearn/psi/api"
)

func isAdminRequest(req *http.Request) bool {
	return req.Header.Get(api.ReqHeader_PSIUserRole) == "admin"
}
