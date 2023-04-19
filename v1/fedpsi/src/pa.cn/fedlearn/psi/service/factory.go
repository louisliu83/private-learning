package service

import (
	"pa.cn/fedlearn/psi/api/service"
	"pa.cn/fedlearn/psi/config"
)

const (
	SvcTypeLocal  = "Local"
	SvcTypeRemote = "Remote"
)

var (
	dataSetServiceMap map[string]service.DatasetService = make(map[string]service.DatasetService, 0)
)

// GetDatasetService return the dataset service
func GetDatasetService(partyName string) service.DatasetService {
	if partyName == config.GetConfig().PartyName {
		return dataSetServiceMap[SvcTypeLocal]
	}
	return dataSetServiceMap[SvcTypeRemote]
}

func RegisterDatasetService(svcType string, svc service.DatasetService) {
	dataSetServiceMap[svcType] = svc
}
