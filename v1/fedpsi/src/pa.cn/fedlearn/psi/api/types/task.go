package types

type TaskDataset struct {
	PartyName string `json:"partyName"`
	DSName    string `json:"dsName"`
	DSIndex   int32  `json:"dsIndex"`
	DSCount   int64  `json:"dsCount"`
	DSSize    int64  `json:"dsSize"`
}

type TaskCreateRequest struct {
	Initiator    string        `json:"initiator"`
	TaskUID      string        `json:"taskUID"`
	Mode         string        `json:"mode"`
	Protocol     string        `json:"protocol"`
	Name         string        `json:"name"`
	Desc         string        `json:"desc"`
	TaskDatasets []TaskDataset `json:"dsPair"`
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
	JobUID string `json:"jobUID"`
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
