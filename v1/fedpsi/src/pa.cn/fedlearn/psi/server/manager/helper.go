package manager

import (
	"context"

	"pa.cn/fedlearn/psi/api/types"
	"pa.cn/fedlearn/psi/client/sdk"
	"pa.cn/fedlearn/psi/config"
	"pa.cn/fedlearn/psi/model"
	service "pa.cn/fedlearn/psi/service"
	_ "pa.cn/fedlearn/psi/service/local"
	_ "pa.cn/fedlearn/psi/service/remote"
)

func getDatasetCount(ctx context.Context, partyName, dsName string, index int32) (int64, error) {
	datasetService := service.GetDatasetService(partyName)
	return datasetService.GetDatasetCount(ctx, partyName, dsName, index)
}

func getDataSetList(ctx context.Context, partyName string) (data []types.Dataset, err error) {
	datasetService := service.GetDatasetService(partyName)
	return datasetService.GetDatasetList(ctx, partyName)
}

func getRemoteClient(partyName string) (*sdk.PartyClient, error) {
	return service.GetRemoteClient(partyName)
}

func toTaskCreateRequestV2(t types.TaskCreateRequest) types.TaskCreateRequestV2 {
	v2 := types.TaskCreateRequestV2{
		Initiator: t.Initiator,
		TaskUID:   t.TaskUID,
		Mode:      t.Mode,
		Protocol:  t.Protocol,
		Name:      t.Name,
		Desc:      t.Desc,
	}
	for _, d := range t.TaskDatasets {
		taskDataset := types.TaskDataset{
			PartyName: d.PartyName,
			DSName:    d.DSName,
			DSIndex:   d.DSIndex,
			DSCount:   d.DSCount,
			DSSize:    d.DSSize,
		}
		if taskDataset.PartyName == config.GetConfig().PartyName {
			v2.LocalDataset = taskDataset
		} else {
			v2.PartyDataset = taskDataset
		}
	}
	return v2
}

func toTaskCreateRequest(v2 types.TaskCreateRequestV2) types.TaskCreateRequest {
	v1 := types.TaskCreateRequest{
		Initiator: v2.Initiator,
		TaskUID:   v2.TaskUID,
		Mode:      v2.Mode,
		Protocol:  v2.Protocol,
		Name:      v2.Name,
		Desc:      v2.Desc,
	}
	v1.TaskDatasets = make([]types.TaskDataset, 0)
	v1.TaskDatasets = append(v1.TaskDatasets, v2.LocalDataset, v2.PartyDataset)
	return v1
}

func toAPIJob(j *model.Job) *types.Job {
	if j == nil {
		return nil
	}
	job := &types.Job{}
	job.Uuid = j.Uuid
	job.Name = j.Name
	job.Desc = j.Desc
	job.Status = j.Status
	job.Initiator = j.Initiator
	job.Mode = j.Mode
	job.LocalName = j.LocalName
	job.LocalDSName = j.LocalDSName
	job.LocalDSCount = j.LocalDSCount
	job.PartyName = j.PartyName
	job.PartyDSName = j.PartyDSName
	job.PartyDSCount = j.PartyDSCount
	job.Protocol = j.Protocol
	job.Result = j.Result
	job.IntersectCount = j.IntersectCount
	job.ExecStart = j.ExecStart
	job.ExecEnd = j.ExecEnd
	return job
}
