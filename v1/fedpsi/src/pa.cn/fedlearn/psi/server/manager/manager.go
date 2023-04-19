package manager

import (
	"sync"
)

var (
	datasetMgr       *DataSetMgr
	datasetMgrLocker = sync.Mutex{}
)

func GetDatasetMgr() *DataSetMgr {
	datasetMgrLocker.Lock()
	defer datasetMgrLocker.Unlock()
	if datasetMgr == nil {
		datasetMgr = &DataSetMgr{}
	}
	return datasetMgr
}

var (
	partyMgr       *PartyMgr
	partyMgrLocker = sync.Mutex{}
)

func GetPartyMgr() *PartyMgr {
	partyMgrLocker.Lock()
	defer partyMgrLocker.Unlock()
	if partyMgr == nil {
		partyMgr = &PartyMgr{}
	}
	return partyMgr
}

var (
	taskMgr       *TaskMgr
	taskMgrLocker = sync.Mutex{}
)

func GetTaskMgr() *TaskMgr {
	taskMgrLocker.Lock()
	defer taskMgrLocker.Unlock()
	if taskMgr == nil {
		taskMgr = &TaskMgr{}
	}
	return taskMgr
}

var (
	tokenMgr       *TokenManager
	tokenMgrLocker = sync.Mutex{}
)

func GetTokenManager() *TokenManager {
	tokenMgrLocker.Lock()
	defer tokenMgrLocker.Unlock()
	if tokenMgr == nil {
		tokenMgr = &TokenManager{}
	}
	return tokenMgr
}

var (
	upMgr       *UploaderManager
	upMgrLocker = sync.Mutex{}
)

func GetUploaderManager() *UploaderManager {
	upMgrLocker.Lock()
	defer upMgrLocker.Unlock()
	if upMgr == nil {
		upMgr = &UploaderManager{}
	}
	return upMgr
}

var (
	userMgr       *UserManager
	userMgrLocker = sync.Mutex{}
)

func GetUserManager() *UserManager {
	userMgrLocker.Lock()
	defer userMgrLocker.Unlock()
	if userMgr == nil {
		userMgr = &UserManager{}
	}
	return userMgr
}

var (
	tproxyMgr      *TProxyManager
	tproxyMgrLoker = sync.Mutex{}
)

func GetTProxyManager() *TProxyManager {
	tproxyMgrLoker.Lock()
	defer tproxyMgrLoker.Unlock()
	if tproxyMgr == nil {
		tproxyMgr = &TProxyManager{}
	}
	return tproxyMgr
}
