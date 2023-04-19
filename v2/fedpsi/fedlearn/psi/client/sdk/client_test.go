package sdk

import (
	"testing"

	"fedlearn/psi/api/types"
)

var (
	c = PartyClient{
		PartyName: "pab",
		Server:    "127.0.0.1:8080",
		Scheme:    "http",
	}
	taskUUID = "0fbab3b5-c3ee-448a-90fd-9c0c1172f0c5"
)

func TestPartyInfo(t *testing.T) {
	info, err := c.PartyInfo()
	if err != nil {
		t.Logf("%v", err)
		t.Fail()
	} else {
		t.Log(info.PartyName, info.Status)
	}
}

func TestPartyDataset(t *testing.T) {
	ds, err := c.Dataset("jd")
	if err != nil {
		t.Logf("%v", err)
		t.Fail()
	} else {
		t.Logf("%v", ds)
	}
}

func TestPartyStartTask(t *testing.T) {
	r := types.TaskStartRequest{
		TaskUID: taskUUID,
	}
	if ok, err := c.StartTask(r); err != nil {
		t.Logf("%v", err)
		t.Fail()
	} else {
		t.Logf("%v", ok)
	}
}

func TestPartyStopTask(t *testing.T) {
	r := types.TaskStopRequest{
		TaskUID: taskUUID,
	}
	if ok, err := c.StopTask(r); err != nil {
		t.Logf("%v", err)
		t.Fail()
	} else {
		t.Logf("%v", ok)
	}
}

func TestTaskGet(t *testing.T) {
	c = PartyClient{
		PartyName: "jd",
		Server:    "192.168.56.101:9090",
		Scheme:    "http",
	}
	taskUUID = "9e7e9de2-8bc4-478f-a4e1-3c70fe43b23b"
	task, err := c.GetTask(taskUUID)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%v", task)
}

func TestTaskList(t *testing.T) {
	tasks, err := c.ListTasks()
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Log(len(tasks))
	for _, task := range tasks {
		t.Logf("%v", task)
	}
}
func TestTaskIntersect(t *testing.T) {
	c = PartyClient{
		PartyName: "jd",
		Server:    "192.168.56.101:9090",
		Scheme:    "http",
	}
	taskUUID = "9e7e9de2-8bc4-478f-a4e1-3c70fe43b23b"
	intersect, err := c.TaskIntersectResult(taskUUID)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%v", intersect)
}
