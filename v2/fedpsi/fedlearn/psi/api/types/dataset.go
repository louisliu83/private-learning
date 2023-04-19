package types

import (
	"time"
)

type ChunkCheckRequest struct {
	MD5       string `json:"md5"`
	Chunk     int32  `json:"chunk"`
	ChunkSize int64  `json:"chunkSize"`
}

type ChunkUploadRequest struct {
	MD5      string `json:"md5"`
	Chunk    int32  `json:"chunk"`
	FileData []byte
}

type DatasetMetaInfo struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"desc"`
	ValidDays   int32  `json:"validDays"`
	BizContext  string `json:"bizCtx"`
}

type ChunkMergeRequest struct {
	MD5        string `json:"md5"`
	ChunkSize  int64  `json:"chunkSize"`
	ChunkCount int32  `json:"chunkCount"`
	DatasetMetaInfo
}

type FilePullRequest struct {
	URL string `json:"url"`
	DatasetMetaInfo
}

type FileLocalCopyRequest struct {
	FilePath string `json:"path"`
	DatasetMetaInfo
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

type IntersectListRequest struct {
	DatasetName       string `json:"dataset"`
	DatasetBizContext string `json:"bizContext"`
}

type IntersectDataInfo struct {
	DatasetName       string `json:"dataset"`
	DatasetDesc       string `json:"desc"`
	DatasetBizContext string `json:"bizContext"`
}

type IntersectListResponse struct {
	IntersectDataInfo
	Intersects []string `json:"intersects"`
}
