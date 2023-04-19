package psi

import (
	"fedlearn/psi/server/manager"
)

type Server struct {
	DSMgr      *manager.DataSetMgr
	TaskMgr    *manager.TaskMgr
	UploadMgr  *manager.UploaderManager
	PartyMgr   *manager.PartyMgr
	ProjectMgr *manager.ProjectMgr
}

func New() *Server {
	return &Server{
		DSMgr:      manager.GetDatasetMgr(),
		TaskMgr:    manager.GetTaskMgr(),
		UploadMgr:  manager.GetUploaderManager(),
		PartyMgr:   manager.GetPartyMgr(),
		ProjectMgr: manager.GetProjectMgr(),
	}
}
