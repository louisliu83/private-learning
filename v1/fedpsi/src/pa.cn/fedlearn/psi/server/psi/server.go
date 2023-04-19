package psi

import (
	"pa.cn/fedlearn/psi/server/manager"
)

type Server struct {
	DSMgr     *manager.DataSetMgr
	TaskMgr   *manager.TaskMgr
	UploadMgr *manager.UploaderManager
	PartyMgr  *manager.PartyMgr
	TProxyMgr *manager.TProxyManager
}

func New() *Server {
	return &Server{
		DSMgr:     manager.GetDatasetMgr(),
		TaskMgr:   manager.GetTaskMgr(),
		UploadMgr: manager.GetUploaderManager(),
		PartyMgr:  manager.GetPartyMgr(),
		TProxyMgr: manager.GetTProxyManager(),
	}
}
