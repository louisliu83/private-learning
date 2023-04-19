package remote

import (
	"context"
	"fmt"
	"sync"

	serviceApi "pa.cn/fedlearn/psi/api/service"
	"pa.cn/fedlearn/psi/api/types"
	"pa.cn/fedlearn/psi/config"
	"pa.cn/fedlearn/psi/log"
	service "pa.cn/fedlearn/psi/service"
)

func init() {
	service.RegisterDatasetService(service.SvcTypeRemote, NewDatasetService())
}

// DatasetServiceImpl implements DatasetService in remote mode
type DatasetServiceImpl struct {
}

var (
	datasetService serviceApi.DatasetService
	datasetLock    = sync.Mutex{}
)

func NewDatasetService() serviceApi.DatasetService {
	datasetLock.Lock()
	defer datasetLock.Unlock()
	if datasetService == nil {
		datasetService = &DatasetServiceImpl{}
	}
	return datasetService
}

func (s *DatasetServiceImpl) GetDataset(ctx context.Context, partyName, dsName string) (*types.Dataset, error) {
	if partyName == config.GetConfig().PartyName {
		return nil, fmt.Errorf("it is local party: %s", partyName)
	}

	remoteClient, err := service.GetRemoteClient(partyName)
	if err != nil {
		log.Errorf(ctx, "Create party client for %s error:%v\n", partyName, err)
		return nil, fmt.Errorf("Remote Error:%v", err)
	}

	index := int32(0)
	ds, err := remoteClient.PartyDatasetGet(ctx, dsName, index)
	if err != nil {
		log.Errorf(ctx, "Cannot get party %s dataset %s-%d, error:%v\n", partyName, dsName, index, err)
		return nil, err
	}

	return ds, nil
}

func (s *DatasetServiceImpl) GetDatasetCount(ctx context.Context, partyName, dsName string, index int32) (int64, error) {
	if partyName == config.GetConfig().PartyName {
		return -1, fmt.Errorf("it is local party: %s", partyName)
	}
	remoteClient, err := service.GetRemoteClient(partyName)
	if err != nil {
		log.Errorf(ctx, "Create party client for %s error:%v\n", partyName, err)
		return int64(0), fmt.Errorf("Remote Error:%v", err)
	}
	ds, err := remoteClient.PartyDatasetGet(ctx, dsName, index)
	if err != nil {
		log.Errorf(ctx, "Cannot get party %s dataset %s-%d, error:%v\n", partyName, dsName, index, err)
		return int64(0), err
	}
	return ds.Count, nil

}

func (s *DatasetServiceImpl) GetDatasetList(ctx context.Context, partyName string) ([]types.Dataset, error) {
	if partyName == config.GetConfig().PartyName {
		return nil, fmt.Errorf("it is local party: %s", partyName)
	}

	remoteClient, err := service.GetRemoteClient(partyName)
	if err != nil {
		log.Errorf(ctx, "Create party client for %s error:%v\n", partyName, err)
		return nil, fmt.Errorf("Remote Error:%v", err)
	}

	remoteDataSet, err := remoteClient.PartyDatasetList(ctx)
	if err != nil {
		log.Errorf(ctx, "Get dataset of party %s error:%v\n", partyName, err)
	}

	if datasetList, ok := remoteDataSet[partyName]; ok {
		return datasetList, nil
	}

	return make([]types.Dataset, 0), nil
}
