package psi

import (
	"net/http"

	"fedlearn/psi/api"
	"fedlearn/psi/common/log"

	"github.com/gorilla/mux"
)

func (s *Server) DataSetIntersectList(res http.ResponseWriter, req *http.Request) {

	ctx := api.ContextFromReq(req)
	auditContextAction(ctx, "DataSetIntersectList")
	log.Debugln(ctx, "Server.DataSetIntersectList is called.")

	vars := mux.Vars(req)
	bizContext := vars["bizcode"]

	result, err := s.TaskMgr.IntersectsOfDatasetBizContext(ctx, bizContext)
	if err != nil {
		auditContextActionResult(ctx, "DataSetIntersectList", "Failed")
		api.OutputError(res, err)
		return
	}

	auditContextActionResult(ctx, "DataSetIntersectList", "Success")
	api.OutputObject(res, result)
	return
}
