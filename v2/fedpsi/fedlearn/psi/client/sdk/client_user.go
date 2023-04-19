package sdk

import (
	"context"
	"encoding/json"
	"fmt"

	"fedlearn/psi/api/types"

	"fedlearn/psi/client/httpc"
)

func (p *PartyClient) DatasetList(withParty string) (map[string][]types.Dataset, error) {
	u := "apis/v2/dataset"
	if withParty != "" {
		u = fmt.Sprintf("%s?party=%s", u, withParty)
	}
	url1 := fmt.Sprintf("%s://%s/%s", p.Scheme, p.Server, u)
	headers := map[string]string{
		"Content-type": "application/json",
		"Accept":       "text/plain",
	}
	if p.Token != "" {
		headers["Authorization"] = fmt.Sprintf("%s %s", "Bearer", p.Token)
	}
	data, err := httpc.DoGet(url1, headers)
	if err != nil {
		return nil, err
	}
	var res PartyDatasetResponse
	if err = json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return res.Result, nil
}

func (p *PartyClient) ConfirmTask(t types.TaskConfirmRequest) (bool, error) {
	u := "apis/v2/task/confirm"
	url1 := fmt.Sprintf("%s://%s/%s", p.Scheme, p.Server, u)
	postData, err := json.Marshal(t)
	if err != nil {
		return false, err
	}
	return p.doPostWithJsonRequestAndReturnNoResultResponse(context.Background(), url1, postData)
}

func (p *PartyClient) StartTask(t types.TaskStartRequest) (bool, error) {
	u := "apis/v2/task/start"
	url1 := fmt.Sprintf("%s://%s/%s", p.Scheme, p.Server, u)
	postData, err := json.Marshal(t)
	if err != nil {
		return false, err
	}
	return p.doPostWithJsonRequestAndReturnNoResultResponse(context.Background(), url1, postData)
}

func (p *PartyClient) StopTask(t types.TaskStopRequest) (bool, error) {
	u := "apis/v2/task/stop"
	url1 := fmt.Sprintf("%s://%s/%s", p.Scheme, p.Server, u)
	postData, err := json.Marshal(t)
	if err != nil {
		return false, err
	}
	return p.doPostWithJsonRequestAndReturnNoResultResponse(context.Background(), url1, postData)
}

func (p *PartyClient) GetTask(taskUID string) (*types.Task, error) {
	u := fmt.Sprintf("apis/v2/task?task_uuid=%s", taskUID)
	url1 := fmt.Sprintf("%s://%s/%s", p.Scheme, p.Server, u)
	headers := map[string]string{
		"Content-type": "application/json",
		"Accept":       "text/plain",
	}
	if p.Token != "" {
		headers["Authorization"] = fmt.Sprintf("%s %s", "Bearer", p.Token)
	}
	data, err := httpc.DoGet(url1, headers)
	if err != nil {
		return nil, err
	}
	var res TaskInfoResponse
	if err = json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return res.Result, nil
}

func (p *PartyClient) ListTasks() ([]*types.Task, error) {
	u := "apis/v2/tasks"
	url1 := fmt.Sprintf("%s://%s/%s", p.Scheme, p.Server, u)
	headers := map[string]string{
		"Content-type": "application/json",
		"Accept":       "text/plain",
	}
	if p.Token != "" {
		headers["Authorization"] = fmt.Sprintf("%s %s", "Bearer", p.Token)
	}
	data, err := httpc.DoGet(url1, headers)
	if err != nil {
		return nil, err
	}
	var res TasksInfoResponse
	if err = json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return res.Result, nil
}

func (p *PartyClient) TaskIntersectResult(taskUID string) (string, error) {
	u := fmt.Sprintf("apis/v2/task/intersect?task_uuid=%s", taskUID)
	url1 := fmt.Sprintf("%s://%s/%s", p.Scheme, p.Server, u)
	headers := map[string]string{
		"Content-type": "application/json",
		"Accept":       "text/plain",
	}
	if p.Token != "" {
		headers["Authorization"] = fmt.Sprintf("%s %s", "Bearer", p.Token)
	}
	data, err := httpc.DoGet(url1, headers)
	if err != nil {
		return "", err
	}
	var res TaskIntersectResponse
	if err = json.Unmarshal(data, &res); err != nil {
		return "", err
	}
	return res.Result, nil
}
