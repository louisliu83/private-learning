package types

import (
	"time"
)

//  !!!!Dreprecated!!!!
type DataSetUploadRequest struct {
	PackageID string //If user pass this packageID, we use this one
	Data      []byte
}

//  !!!!Dreprecated!!!!
type DataSetUploadResponse struct {
	PackageID string `json:"packageID"`
}

//  !!!!Dreprecated!!!!
type DataSetAddRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	PackageID   string `json:"packageID"`
}

type DataSetDeleteRequest struct {
	Name  string `json:"name"`
	Index int32  `json:"index"`
}

type DataSetListRequest struct {
	PartyNames []string `json:"partyNames"`
}

type Dataset struct {
	Id          uint64    `json:"id"`
	Name        string    `json:"name"`
	Index       int32     `json:"index"`
	Desc        string    `json:"desc"`
	Count       int64     `json:"count"`
	Size        int64     `json:"size"`
	ShardsNum   int32     `json:"shardsNum"`
	BizContext  string    `json:"bizCtx"`
	Status      string    `json:"status"`
	ExpiredDate time.Time `json:"expiredDate"`
}

type DatasetGrantRequest struct {
	Name      string   `json:"name"`
	PartyList []string `json:"partyList"`
}

type DatasetRevokeRequest struct {
	DatasetGrantRequest
}
