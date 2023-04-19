package local

import (
	"context"
	"errors"
	"fmt"
	"sync"

	serviceApi "fedlearn/psi/api/service"

	"fedlearn/psi/api"
	"fedlearn/psi/api/types"
	"fedlearn/psi/common/config"
	"fedlearn/psi/common/log"
	"fedlearn/psi/model"
	"fedlearn/psi/service"
)

func init() {
	service.RegisterDatasetService(service.SvcTypeLocal, NewDatasetService())
}

// DatasetServiceImpl implements DatasetService in local mode
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
	if partyName != config.GetConfig().PartyName {
		return nil, fmt.Errorf("Not local party: %s", partyName)
	}
	index := int32(0)
	ds, err := model.GetDatasetByNameAndIndex(dsName, index)
	if err != nil {
		log.Errorf(ctx, "Cannot load party %s dataset %s-%d, error:%v\n", partyName, dsName, index, err)
		return nil, err
	}
	apid := service.ToAPIDataset(ds)
	return &apid, nil
}

func (s *DatasetServiceImpl) GetDatasetCount(ctx context.Context, partyName, dsName string, index int32) (int64, error) {
	if partyName != config.GetConfig().PartyName {
		return -1, fmt.Errorf("Not local party: %s", partyName)
	}
	ds, err := model.GetDatasetByNameAndIndex(dsName, index)
	if err != nil {
		log.Errorf(ctx, "Cannot load party %s dataset %s-%d, error:%v\n", partyName, dsName, index, err)
		return int64(0), err
	}
	return ds.Count, nil
}

func (s *DatasetServiceImpl) GetDatasetList(ctx context.Context, partyName string) ([]types.Dataset, error) {
	apids := make([]types.Dataset, 0)
	if partyName != config.GetConfig().PartyName {
		return apids, fmt.Errorf("Not local party: %s", partyName)
	}

	listFunc := func() ([]*model.Dataset, error) {
		return model.ListDataset()
	}

	if config.GetConfig().FeatureGate.DatasetPrivate {
		sourceParty := fmt.Sprintf("%s", ctx.Value(api.ReqHeader_PSIUserParty))
		if sourceParty != config.GetConfig().PartyName {
			listFunc = func() ([]*model.Dataset, error) {
				return model.ListDatasetBySrcParty(sourceParty)
			}
		}
	}

	ds, err := listFunc()
	if err != nil {
		log.Errorf(ctx, "Cannot load datasets, error:%v\n", err)
		return apids, errors.New("Load datasets error")
	}

	for _, d := range ds {
		if d.IsValid() {
			apid := service.ToAPIDataset(d)
			apids = append(apids, apid)
		}
	}
	return apids, nil
}
