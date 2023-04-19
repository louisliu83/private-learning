package sdk

import (
	"context"
	"fmt"

	"encoding/json"

	"fedlearn/psi/api"

	"fedlearn/psi/client/httpc"
)

type PartyClient struct {
	SourcePartyName string
	PartyName       string
	Server          string
	WorkAddr        string
	WorkPort        int32
	Scheme          string
	Token           string
}

func New(srcPartyName, partyName, scheme, server, workAddr string, workPort int32, token string) *PartyClient {
	return &PartyClient{
		SourcePartyName: srcPartyName,
		PartyName:       partyName,
		Server:          server,
		WorkAddr:        workAddr,
		WorkPort:        workPort,
		Scheme:          scheme,
		Token:           token,
	}
}

func (p *PartyClient) doPostWithJsonRequestAndReturnNoResultResponse(ctx context.Context, url1 string, postData []byte) (bool, error) {
	headers := map[string]string{
		"Content-type": "application/json",
		"Accept":       "text/plain",
	}
	setTraceIDHeader(ctx, headers)
	if p.Token != "" {
		headers["Authorization"] = fmt.Sprintf("%s %s", "Bearer", p.Token)
	}

	data, err := httpc.DoPostWithJson(url1, headers, postData)
	if err != nil {
		return false, err
	}
	var res NoResultResponse
	if err = json.Unmarshal(data, &res); err != nil {
		return false, err
	}
	if res.ErrorCode == 0 {
		return true, nil
	}
	return false, fmt.Errorf("%s", res.ErrorMsg)
}

func setTraceIDHeader(ctx context.Context, headers map[string]string) {
	traceID := ctx.Value(api.Trace_ID)
	if traceID == nil {
		return
	}

	traceIDStr := fmt.Sprintf("%s", traceID)
	if headers == nil {
		headers = make(map[string]string, 0)
	}

	headers[api.Trace_Req_ID] = traceIDStr
}
