package model

import (
	"fmt"
	"sync"
	"testing"

	"github.com/astaxie/beego/orm"
)

var (
	once sync.Once
)

const (
	sqlitedb = "D:\\psi.db"
)

func Setup() {
	once.Do(func() {
		Initdb(sqlitedb)
	})
}

func TestInit(t *testing.T) {
	Setup()
	o := orm.NewOrm()
	s := o.DBStats()
	if s == nil {
		t.Fail()
	}
}

func TestAddDataset(t *testing.T) {
	Setup()
	name := "user_ids"
	index := int32(0)
	ds := &Dataset{
		Name:  name,
		Index: index,
		Path:  fmt.Sprintf("/home/jack/psi/dataset/user_ids/%s_%d", name, index),
	}
	err := AddDataset(ds)
	if err != nil {
		t.Fail()
	}
}

func TestGetTaskByUuid(t *testing.T) {
	Setup()
	task, err := GetTaskByUuid("0fbab3b5-c3ee-448a-90fd-9c0c1172f0c5")
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%v", task)
}

func TestAddTask(t *testing.T) {
	Setup()
	task := &Task{
		Uuid:   "hh3",
		Status: TaskStatus_Init,
	}
	err := AddTask(task)
	if err != nil {
		t.Fatalf("%v", err)
	}
}

func TestGetOldestTask(t *testing.T) {
	Setup()
	task, err := GetOldestTask()
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%v", task)
}
