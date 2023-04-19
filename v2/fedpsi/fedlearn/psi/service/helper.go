package service

import (
	"fmt"

	"fedlearn/psi/api/types"
	"fedlearn/psi/client/sdk"
	"fedlearn/psi/common/config"
	"fedlearn/psi/model"
)

func GetRemoteClient(partyName string) (*sdk.PartyClient, error) {
	r, err := model.GetPartyByName(partyName)
	if err != nil {
		return nil, fmt.Errorf("No party %s exists:%w", partyName, err)
	}
	srcPartyName := config.GetConfig().PartyName
	c := sdk.New(srcPartyName, r.Name, r.Scheme, fmt.Sprintf("%s:%d", r.ControllerServer, r.ControllerPort), r.WorkServer, r.WorkPort, r.Token)
	return c, nil
}

func ToAPIDataset(data *model.Dataset) types.Dataset {
	ds := types.Dataset{
		Id:          data.Id,
		Name:        data.Name,
		Index:       data.Index,
		Desc:        data.Desc,
		Count:       data.Count,
		Size:        data.Size,
		Status:      data.Status,
		ShardsNum:   data.Shards,
		BizContext:  data.BizContext,
		ExpiredDate: data.ExpiredDate,
	}
	return ds
}
