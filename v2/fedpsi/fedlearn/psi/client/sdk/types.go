package sdk

import (
	"fedlearn/psi/api/types"
)

type PartyConfigResponse struct {
	ErrorCode int               `json:"errcode"`
	ErrorMsg  string            `json:"errmsg"`
	Result    *types.ConfigInfo `json:"result"`
}

type PartyInfoResponse struct {
	ErrorCode int               `json:"errcode"`
	ErrorMsg  string            `json:"errmsg"`
	Result    *types.ServerInfo `json:"result"`
}

type PartyDatasetResponse struct {
	ErrorCode int                        `json:"errcode"`
	ErrorMsg  string                     `json:"errmsg"`
	Result    map[string][]types.Dataset `json:"result"`
}

type PartyDataShardsResponse struct {
	ErrorCode int             `json:"errcode"`
	ErrorMsg  string          `json:"errmsg"`
	Result    []types.Dataset `json:"result"`
}

type PartyDatasetGetResponse struct {
	ErrorCode int            `json:"errcode"`
	ErrorMsg  string         `json:"errmsg"`
	Result    *types.Dataset `json:"result"`
}

type TaskInfoResponse struct {
	ErrorCode int         `json:"errcode"`
	ErrorMsg  string      `json:"errmsg"`
	Result    *types.Task `json:"result"`
}

type TasksInfoResponse struct {
	ErrorCode int           `json:"errcode"`
	ErrorMsg  string        `json:"errmsg"`
	Result    []*types.Task `json:"result"`
}

type TaskIntersectResponse struct {
	ErrorCode int    `json:"errcode"`
	ErrorMsg  string `json:"errmsg"`
	Result    string `json:"result"`
}

type JobIntersectResponse struct {
	ErrorCode int    `json:"errcode"`
	ErrorMsg  string `json:"errmsg"`
	Result    string `json:"result"`
}

type NoResultResponse struct {
	ErrorCode int    `json:"errcode"`
	ErrorMsg  string `json:"errmsg"`
}
