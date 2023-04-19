package types

import (
	"time"
)

type TaskDataset struct {
	PartyName string `json:"partyName"`
	DSName    string `json:"dsName"`
	DSIndex   int32  `json:"dsIndex"`
	DSCount   int64  `json:"dsCount"`
	DSSize    int64  `json:"dsSize"`
}

type TaskCreateRequestV2 struct {
	Initiator    string      `json:"initiator"`
	TaskUID      string      `json:"taskUID"`
	JobUID       string      `json:"jobUID"`
	Mode         string      `json:"mode"`
	Protocol     string      `json:"protocol"`
	Name         string      `json:"name"`
	Desc         string      `json:"desc"`
	LocalDataset TaskDataset `json:"localDS"`
	PartyDataset TaskDataset `json:"partyDS"`
}

type TaskListRequest struct {
	JobUID   string `json:"jobUID"`
	PageNum  int    `json:"pageNum"`
	PageSize int    `json:"pageSize"`
}

type Task struct {
	Uuid          string `json:"uuid"`
	JobUID        string `json:"jobUid"`
	Name          string `json:"name"`
	Desc          string `json:"desc"`
	Initiator     string `json:"initiator"`
	LocalName     string `json:"localName"`
	LocalDSName   string `json:"localDSName"`
	LocalDSIndex  int32  `json:"localDSIndex"`
	LocalDSCount  int64  `json:"localDSCount"`
	PartyName     string `json:"partyName"`
	PartyDSName   string `json:"partyDSName"`
	PartyDSIndex  int32  `json:"partyDSIndex"`
	PartyDSCount  int64  `json:"partyDSCount"`
	Mode          string `json:"mode"`
	PSIServerIP   string `json:"psiServerIP"`
	PSIServerPort int32  `json:"psiServerPort"`
	Status        string `json:"status"`
}

type TaskConfirmRequest struct {
	TaskUID       string `json:"taskUID"`
	PartyDSCount  int64  `json:"partyDSCount"`
	PSIServerIP   string `json:"psiServerIP"`
	PSIServerPort int32  `json:"psiServerPort"`
}

type TaskStartRequest struct {
	TaskUID string `json:"taskUID"`
}

type TaskStopRequest struct {
	TaskUID string `json:"taskUID"`
}

type TaskGetRequest struct {
	TaskUID string `json:"taskUID"`
}

type TaskIntersectionDownloadRequest struct {
	TaskUID string `json:"taskUID"`
}

type TaskRerunRequest struct {
	TaskUID string `json:"taskUID"`
}

type BatchJobSubmitRequest []JobSubmitRequest

type JobSubmitRequest struct {
	Initiator    string      `json:"initiator"`
	ActUID       string      `json:"actUID"`
	JobUID       string      `json:"jobUID"`
	Mode         string      `json:"mode"`
	Protocol     string      `json:"protocol"`
	Name         string      `json:"name"`
	Desc         string      `json:"desc"`
	LocalDataset TaskDataset `json:"localDS"`
	PartyDataset TaskDataset `json:"partyDS"`
}

type JobConfirmRequest struct {
	JobUID string `json:"jobUID"`
}

type JobStopRequest struct {
	JobUID string `json:"jobUID"`
}

type JobDelRequest struct {
	JobUID string `json:"jobUID"`
}

type JobDownloadRequest struct {
	JobUID string `json:"jobUID"`
}

type Job struct {
	Uuid           string    `json:"uuid"`
	ActUID         string    `json:"actUID"`
	Name           string    `json:"name"`
	Desc           string    `json:"desc"`
	Initiator      string    `json:"initiator"`
	LocalName      string    `json:"localName"`
	LocalDSName    string    `json:"localDSName"`
	LocalDSCount   int64     `json:"localDSCount"`
	PartyName      string    `json:"partyName"`
	PartyDSName    string    `json:"partyDSName"`
	PartyDSCount   int64     `json:"partyDSCount"`
	Protocol       string    `json:"protocol"`
	Mode           string    `json:"mode"`
	Status         string    `json:"status"`
	Result         string    `json:"result"`
	IntersectCount int64     `json:"intersectCount"`
	ExecStart      time.Time `json:"execStart"`
	ExecEnd        time.Time `json:"execEnd"`
}
