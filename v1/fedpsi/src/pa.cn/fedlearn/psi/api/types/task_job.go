package types

import (
	"time"
)

type JobSubmitRequest struct {
	Initiator    string      `json:"initiator"`
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

type Job struct {
	Uuid           string    `json:"uuid"`
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
