package manager

import (
	"fmt"
	"os"

	"pa.cn/fedlearn/psi/config"

	"context"

	"pa.cn/fedlearn/psi/api/types"
	"pa.cn/fedlearn/psi/log"
	"pa.cn/fedlearn/psi/model"
	"pa.cn/fedlearn/psi/service"
	"pa.cn/fedlearn/psi/utils"
)

type DataSetMgr struct {
}

const (
	TMPDIR = "tmp"
)

func (m *DataSetMgr) DataSetGet(ctx context.Context, name string, index int32) (ds *types.Dataset, err error) {
	log.Debugln(ctx, "DataSetMgr.DataSetGet is called")
	return m.getDatasetByNameAndIndex(ctx, name, index)
}

func (m *DataSetMgr) DataSetShardsList(ctx context.Context, name string) (ds []types.Dataset, err error) {
	log.Debugln(ctx, "DataSetMgr.DataSetShardsList is called")
	dataShards, err := model.GetDatasetShards(name)
	if err != nil {
		log.Errorf(ctx, "Cannot load datashards of %s, error:%v\n", name, err)
		return nil, fmt.Errorf("Cannot load datashards of %s", name)
	}
	apids := make([]types.Dataset, 0)
	for _, d := range dataShards {
		apid := service.ToAPIDataset(d)
		apids = append(apids, apid)
	}
	return apids, nil
}

func (m *DataSetMgr) DataSetList(ctx context.Context, r types.DataSetListRequest) (data map[string][]types.Dataset, err error) {
	log.Debugln(ctx, "DataSetMgr.DataSetList is called")

	if r.PartyNames == nil {
		r.PartyNames = make([]string, 0)
	}
	r.PartyNames = append(r.PartyNames, config.GetConfig().PartyName)

	data = map[string][]types.Dataset{}
	for _, partyName := range r.PartyNames {
		datasetList, err := getDataSetList(ctx, partyName)
		if err != nil {
			log.Errorf(ctx, "Get dataset list of %s error:%v\n", partyName, err)
			continue
		}
		data[partyName] = datasetList
	}

	return data, nil
}

func (m *DataSetMgr) DataSetDelete(ctx context.Context, r types.DataSetDeleteRequest) (err error) {
	log.Debugln(ctx, "DataSetMgr.DataSetDelete is called")
	ds, err := model.GetDatasetByNameAndIndex(r.Name, r.Index)
	if err != nil {
		log.Errorf(ctx, "Get dataset %s_%d from db failed %v\n", r.Name, r.Index, err)
		return err
	}

	shards, err := model.GetDatasetShards(r.Name)
	if err != nil {
		log.Warningf(ctx, "Get dataset shards of  %s from db failed %v\n", r.Name, err)
	}

	err = model.DeleteDatasetsByName(r.Name)
	if err != nil {
		log.Errorf(ctx, "Delete dataset %s from db failed %v\n", r.Name, err)
		return err
	}
	// Remove dataset from filesystem
	if len(shards) > 0 {
		for _, shard := range shards {
			targetPath := DatasetPath(shard.Md5, shard.Name, shard.Index)
			if err := os.Remove(targetPath); err != nil {
				log.Warningf(ctx, "remove dataset shard %s failed:%v", targetPath, err)
			}
		}
		if len(shards) > 1 {
			targetPath := DatasetPath(ds.Md5, ds.Name, int32(0))
			if err := os.Remove(targetPath); err != nil {
				log.Warningf(ctx, "remove dataset full file %s failed:%v", targetPath, err)
			}
		}
	}

	return nil
}

//------------------ party funcs -----------------------------------

func (m *DataSetMgr) PartyDataSetList(ctx context.Context, r types.DataSetListRequest) (data map[string][]types.Dataset, err error) {
	log.Debugln(ctx, "DataSetMgr.PartyDataSetList is called")
	r.PartyNames = make([]string, 0)
	return m.DataSetList(ctx, r)
}

func (m *DataSetMgr) PartyDataSetGet(ctx context.Context, name string, index int32) (ds *types.Dataset, err error) {
	log.Debugln(ctx, "DataSetMgr.PartyDataSetGet is called")
	return m.getDatasetByNameAndIndex(ctx, name, index)
}

func (m *DataSetMgr) DataSetGrant(ctx context.Context, r types.DatasetGrantRequest) error {
	log.Debugln(ctx, "DataSetMgr.DataSetGrant is called")
	if len(r.PartyList) == 0 {
		return nil
	}
	ds, err := model.GetDatasetByNameAndIndex(r.Name, int32(0))
	if err != nil {
		log.Errorf(ctx, "Dataset %s not found %v\n", r.Name, err)
		return err
	}
	originPartyList := utils.CommaSeperatedStringToSlice(ds.Parties)
	partyList := utils.Union(originPartyList, r.PartyList)

	ds.Parties = utils.SliceToCommaSeperatedString(partyList)
	if err := model.UpdateDataset(ds); err != nil {
		log.Errorf(ctx, "Dataset %s grant update parties error %v\n", r.Name, err)
		return err
	}
	return nil
}

func (m *DataSetMgr) DataSetRevoke(ctx context.Context, r types.DatasetRevokeRequest) error {
	log.Debugln(ctx, "DataSetMgr.DataSetRevoke is called")
	if len(r.PartyList) == 0 {
		return nil
	}
	ds, err := model.GetDatasetByNameAndIndex(r.Name, int32(0))
	if err != nil {
		log.Errorf(ctx, "Dataset %s not found %v\n", r.Name, err)
		return err
	}
	originPartyList := utils.CommaSeperatedStringToSlice(ds.Parties)
	partyList := utils.Intersect(originPartyList, r.PartyList)

	ds.Parties = utils.SliceToCommaSeperatedString(partyList)
	if err := model.UpdateDataset(ds); err != nil {
		log.Errorf(ctx, "Dataset %s revoke update parties error %v\n", r.Name, err)
		return err
	}
	return nil
}

// ----helper---
func (m *DataSetMgr) getDatasetByNameAndIndex(ctx context.Context, name string, index int32) (ds *types.Dataset, err error) {
	log.Debugln(ctx, "DataSetMgr.getDatasetByNameAndIndex is called")
	data, err := model.GetDatasetByNameAndIndex(name, index)
	if err != nil {
		log.Errorf(ctx, "Cannot load dataset %s-%d, error:%v\n", name, index, err)
		return nil, fmt.Errorf("Cannot load dataset %s-%d", name, index)
	}
	d := service.ToAPIDataset(data)
	return &d, nil
}
