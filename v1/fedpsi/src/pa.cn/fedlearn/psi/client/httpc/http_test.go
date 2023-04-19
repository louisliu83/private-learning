package httpc

import (
	"encoding/json"
	"testing"

	"pa.cn/fedlearn/psi/api/types"
)

func TestGet(t *testing.T) {
	url1 := "http://127.0.0.1:8080/apis/info"
	headers := map[string]string{
		"Content-type": "application/json",
		"Accept":       "text/plain",
	}
	data, err := DoGet(url1, headers)
	if err != nil {
		t.Errorf("%v", err)
	}
	t.Log(string(data))
}

func TestPost(t *testing.T) {
	url1 := "http://127.0.0.1:8080/apis/v1/dataset"
	headers := map[string]string{
		"Content-type": "application/json",
		"Accept":       "text/plain",
	}

	ds := types.DataSetAddRequest{
		Name:        "testPostName",
		Description: "testPostDesc",
		PackageID:   "30bcaf2f-5ae3-494a-9898-65b9676c5e8b",
	}
	dsData, err := json.MarshalIndent(ds, "", "\t")
	if err != nil {
		t.Errorf("%v", err)
	}
	data, err := DoPostWithJson(url1, headers, dsData)
	if err != nil {
		t.Errorf("%v", err)
	}
	t.Log(string(data))
}
