package api

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type PsiServerError struct {
	ErrorCode int         `json:"errcode"`
	ErrorMsg  string      `json:"errmsg"`
	Result    interface{} `json:"result"`
}

func (e PsiServerError) ErrorBytes() []byte {
	data, _ := json.MarshalIndent(e, "", "\t")
	return data
}

func (e PsiServerError) Error() string {
	data, _ := json.MarshalIndent(e, "", "\t")
	return string(data)
}

func OutputObject(res http.ResponseWriter, resObj interface{}) {
	pse := PsiServerError{
		ErrorCode: 0,
		ErrorMsg:  "success",
		Result:    resObj,
	}
	res.Write(pse.ErrorBytes())
}

func OutputError(res http.ResponseWriter, err error) {
	var pe PsiServerError
	switch err.(type) {
	case PsiServerError:
		pe, _ = err.(PsiServerError)
	default:
		pe = PsiServerError{
			ErrorCode: -1,
			ErrorMsg:  err.Error(),
		}
	}
	res.Write(pe.ErrorBytes())
}

func OutputSuccess(res http.ResponseWriter) {
	pe := PsiServerError{
		ErrorCode: 0,
		ErrorMsg:  "success",
	}
	res.Write(pe.ErrorBytes())
}

func ReadObjectFromReqBody(req *http.Request, obj interface{}) error {
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, obj)
}

func GetQueryValue(req *http.Request, key string) string {
	return req.URL.Query().Get(key)
}

func ContextFromReq(req *http.Request) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, Trace_ID, req.Header.Get(Trace_Req_ID))
	user := req.Header.Get(ReqHeader_PSIUserID)
	if user == "" {
		user = "anonymous"
	}
	ctx = context.WithValue(ctx, ReqHeader_PSIUserID, user)
	party := req.Header.Get(ReqHeader_PSIUserParty)
	if party == "" {
		party = "anonymous_party"
	}
	ctx = context.WithValue(ctx, ReqHeader_PSIUserParty, party)
	return ctx
}
