package sdk

import (
	"context"
	"encoding/json"
	"fmt"

	"pa.cn/fedlearn/psi/api/types"
	"pa.cn/fedlearn/psi/client/httpc"
)

func (p *PartyClient) PartyStatus(ctx context.Context) (*types.ServerInfo, error) {
	u := fmt.Sprintf("p2p/%s/v1/info", p.SourcePartyName)
	url1 := fmt.Sprintf("%s://%s/%s", p.Scheme, p.Server, u)
	headers := map[string]string{
		"Content-type": "application/json",
		"Accept":       "text/plain",
	}
	setTraceIDHeader(ctx, headers)
	if p.Token != "" {
		headers["Authorization"] = fmt.Sprintf("%s %s", "Bearer", p.Token)
	}
	data, err := httpc.DoGet(url1, headers)
	if err != nil {
		return nil, err
	}
	var res PartyInfoResponse
	if err = json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return res.Result, nil
}

func (p *PartyClient) PartyDatasetGet(ctx context.Context, name string, index int32) (*types.Dataset, error) {
	u := fmt.Sprintf("p2p/%s/v1/dataset/%s/%d", p.SourcePartyName, name, index)
	url1 := fmt.Sprintf("%s://%s/%s", p.Scheme, p.Server, u)
	headers := map[string]string{
		"Content-type": "application/json",
		"Accept":       "text/plain",
	}
	setTraceIDHeader(ctx, headers)
	if p.Token != "" {
		headers["Authorization"] = fmt.Sprintf("%s %s", "Bearer", p.Token)
	}
	data, err := httpc.DoGet(url1, headers)
	if err != nil {
		return nil, err
	}
	var res PartyDatasetGetResponse
	if err = json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return res.Result, nil
}

func (p *PartyClient) PartyDatasetList(ctx context.Context) (map[string][]types.Dataset, error) {
	u := fmt.Sprintf("p2p/%s/v1/dataset", p.SourcePartyName)
	url1 := fmt.Sprintf("%s://%s/%s", p.Scheme, p.Server, u)
	headers := map[string]string{
		"Content-type": "application/json",
		"Accept":       "text/plain",
	}
	setTraceIDHeader(ctx, headers)
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

func (p *PartyClient) PartyDataShards(ctx context.Context, name string) ([]types.Dataset, error) {
	u := fmt.Sprintf("p2p/%s/v1/dataset/shards?name=%s", p.SourcePartyName, name)
	url1 := fmt.Sprintf("%s://%s/%s", p.Scheme, p.Server, u)
	headers := map[string]string{
		"Content-type": "application/json",
		"Accept":       "text/plain",
	}
	setTraceIDHeader(ctx, headers)
	if p.Token != "" {
		headers["Authorization"] = fmt.Sprintf("%s %s", "Bearer", p.Token)
	}
	data, err := httpc.DoGet(url1, headers)
	if err != nil {
		return nil, err
	}
	var res PartyDataShardsResponse
	if err = json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return res.Result, nil
}

func (p *PartyClient) PartyJobSubmit(ctx context.Context, t types.JobSubmitRequest) (bool, error) {
	u := fmt.Sprintf("p2p/%s/v1/job", p.SourcePartyName)
	url1 := fmt.Sprintf("%s://%s/%s", p.Scheme, p.Server, u)
	postData, err := json.Marshal(t)
	if err != nil {
		return false, err
	}
	return p.doPostWithJsonRequestAndReturnNoResultResponse(ctx, url1, postData)
}

func (p *PartyClient) CreatePartyTask(ctx context.Context, t types.TaskCreateRequest) (bool, error) {
	u := fmt.Sprintf("p2p/%s/v1/task", p.SourcePartyName)
	url1 := fmt.Sprintf("%s://%s/%s", p.Scheme, p.Server, u)
	postData, err := json.Marshal(t)
	if err != nil {
		return false, err
	}
	return p.doPostWithJsonRequestAndReturnNoResultResponse(ctx, url1, postData)
}

func (p *PartyClient) CreatePartyTaskV2(ctx context.Context, t types.TaskCreateRequestV2) (bool, error) {
	u := fmt.Sprintf("p2p/%s/v1/task", p.SourcePartyName)
	url1 := fmt.Sprintf("%s://%s/%s", p.Scheme, p.Server, u)
	postData, err := json.Marshal(t)
	if err != nil {
		return false, err
	}
	return p.doPostWithJsonRequestAndReturnNoResultResponse(ctx, url1, postData)
}

func (p *PartyClient) ConfirmPartyJob(ctx context.Context, t types.JobConfirmRequest) (bool, error) {
	u := fmt.Sprintf("p2p/%s/v1/job/confirm", p.SourcePartyName)
	url1 := fmt.Sprintf("%s://%s/%s", p.Scheme, p.Server, u)
	postData, err := json.Marshal(t)
	if err != nil {
		return false, err
	}
	return p.doPostWithJsonRequestAndReturnNoResultResponse(ctx, url1, postData)
}

func (p *PartyClient) ConfirmPartyTask(ctx context.Context, t types.TaskConfirmRequest) (bool, error) {
	u := fmt.Sprintf("p2p/%s/v1/task/confirm", p.SourcePartyName)
	url1 := fmt.Sprintf("%s://%s/%s", p.Scheme, p.Server, u)
	postData, err := json.Marshal(t)
	if err != nil {
		return false, err
	}
	return p.doPostWithJsonRequestAndReturnNoResultResponse(ctx, url1, postData)
}

func (p *PartyClient) StartPartyTask(ctx context.Context, t types.TaskStartRequest) (bool, error) {
	u := fmt.Sprintf("p2p/%s/v1/task/start", p.SourcePartyName)
	url1 := fmt.Sprintf("%s://%s/%s", p.Scheme, p.Server, u)
	postData, err := json.Marshal(t)
	if err != nil {
		return false, err
	}
	return p.doPostWithJsonRequestAndReturnNoResultResponse(ctx, url1, postData)
}

func (p *PartyClient) StopPartyTask(ctx context.Context, t types.TaskStopRequest) (bool, error) {
	u := fmt.Sprintf("p2p/%s/v1/task/stop", p.SourcePartyName)
	url1 := fmt.Sprintf("%s://%s/%s", p.Scheme, p.Server, u)
	postData, err := json.Marshal(t)
	if err != nil {
		return false, err
	}
	return p.doPostWithJsonRequestAndReturnNoResultResponse(ctx, url1, postData)
}

func (p *PartyClient) GetPartyTask(ctx context.Context, taskUID string) (*types.Task, error) {
	u := fmt.Sprintf("p2p/%s/v1/task?task_uuid=%s", p.SourcePartyName, taskUID)
	url1 := fmt.Sprintf("%s://%s/%s", p.Scheme, p.Server, u)
	headers := map[string]string{
		"Content-type": "application/json",
		"Accept":       "text/plain",
	}
	setTraceIDHeader(ctx, headers)
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

func (p *PartyClient) TaskIntersectPartyResult(ctx context.Context, taskUID string) (string, error) {
	u := fmt.Sprintf("p2p/%s/v1/task/intersect?task_uuid=%s", p.SourcePartyName, taskUID)
	url1 := fmt.Sprintf("%s://%s/%s", p.Scheme, p.Server, u)
	headers := map[string]string{
		"Content-type": "application/json",
		"Accept":       "text/plain",
	}
	setTraceIDHeader(ctx, headers)
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

func (p *PartyClient) JobIntersectPartyResult(ctx context.Context, jobUID string) (string, error) {
	u := fmt.Sprintf("p2p/%s/v1/job/intersect?job_uuid=%s", p.SourcePartyName, jobUID)
	url1 := fmt.Sprintf("%s://%s/%s", p.Scheme, p.Server, u)
	headers := map[string]string{
		"Content-type": "application/json",
		"Accept":       "text/plain",
	}
	setTraceIDHeader(ctx, headers)
	if p.Token != "" {
		headers["Authorization"] = fmt.Sprintf("%s %s", "Bearer", p.Token)
	}
	data, err := httpc.DoGet(url1, headers)
	if err != nil {
		return "", err
	}
	var res JobIntersectResponse
	if err = json.Unmarshal(data, &res); err != nil {
		return "", err
	}
	return res.Result, nil
}

func (p *PartyClient) RerunPartyTask(ctx context.Context, t types.TaskRerunRequest) (bool, error) {
	u := fmt.Sprintf("p2p/%s/v1/task/rerun", p.SourcePartyName)
	url1 := fmt.Sprintf("%s://%s/%s", p.Scheme, p.Server, u)
	postData, err := json.Marshal(t)
	if err != nil {
		return false, err
	}
	return p.doPostWithJsonRequestAndReturnNoResultResponse(ctx, url1, postData)
}
