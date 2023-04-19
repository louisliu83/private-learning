package service

import (
	"fmt"

	"pa.cn/fedlearn/psi/api/types"
	"pa.cn/fedlearn/psi/client/sdk"
	"pa.cn/fedlearn/psi/config"
	"pa.cn/fedlearn/psi/model"
)

func GetRemoteClient(partyName string) (*sdk.PartyClient, error) {
	r, err := model.GetPartyByName(partyName)
	if err != nil {
		return nil, fmt.Errorf("No party %s exists:%w", partyName, err)
	}
	srcPartyName := config.GetConfig().PartyName
	c := sdk.New(srcPartyName, r.Name, r.Scheme, fmt.Sprintf("%s:%d", r.ControllerServer, r.ControllerPort), r.Token)
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
