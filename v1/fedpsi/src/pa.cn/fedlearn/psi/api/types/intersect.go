package types

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
