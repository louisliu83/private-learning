package account

import (
	"net/http"

	"pa.cn/fedlearn/psi/api"
)

func isAdminRequest(req *http.Request) bool {
	return req.Header.Get(api.ReqHeader_PSIUserRole) == "admin"
}
