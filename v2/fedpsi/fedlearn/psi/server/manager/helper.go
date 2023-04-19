package manager

import (
	"context"
	"strings"

	"fedlearn/psi/api/types"
	"fedlearn/psi/client/sdk"
	"fedlearn/psi/model"
	"fedlearn/psi/service"
	_ "fedlearn/psi/service/local"
	_ "fedlearn/psi/service/remote"
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

func toAPIJob(j *model.Job) *types.Job {
	if j == nil {
		return nil
	}
	job := &types.Job{}
	job.ActUID = j.ActUid
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

func ToAPIActivity(a *model.Activity) *types.Activity {
	if a == nil {
		return nil
	}

	ta := types.Activity{
		Uuid:          a.Uuid,
		Name:          a.Name,
		Title:         a.Title,
		SendID:        a.SendID,
		InitParty:     a.InitParty,
		FollowerParty: a.FollowerParty,
		Status:        a.Status,
	}

	if len(a.InitiatorData) != 0 {
		ta.Dataset = strings.Split(a.InitiatorData, ",")
	} else {
		ta.Dataset = make([]string, 0)
	}

	if len(a.FollowerData) != 0 {
		ta.FollowerDataset = strings.Split(a.FollowerData, ",")
	} else {
		ta.FollowerDataset = make([]string, 0)
	}

	return &ta
}

func ToAPIProject(j *model.Project) *types.Project {
	if j == nil {
		return nil
	}
	proj := &types.Project{}
	proj.Id = j.Id
	proj.Uuid = j.Uuid
	proj.Name = j.Name
	proj.Status = j.Status
	proj.Type = j.Type
	proj.Desc = j.Desc
	proj.InitParty = j.InitParty
	proj.FollowerParty = j.FollowerParty
	proj.Creator = j.Creator
	proj.UpdateUser = j.UpdateUser
	proj.Created = j.Created
	proj.Updated = j.Updated
	return proj
}
