package model

import (
	"time"

	"github.com/astaxie/beego/orm"
)

const (
	TaskRunMode_Server = "server"
	TaskRunMode_Client = "client"
)

const (
	TaskStatus_WaitingPartyConfirm = "awaiting_party_confirm"
	TaskStatus_Created             = "created"
	TaskStatus_Ready               = "ready"
	TaskStatus_Init                = "init"
	TaskStatus_Running             = "running"
	TaskStatus_ClientWaiting       = "client_waiting"
	TaskStatus_Completed           = "completed"
	TaskStatus_Success             = "success"
	TaskStatus_Failed              = "failed"
	TaskStatus_Cancel              = "cancel"
)

// Task is a psi task
type Task struct {
	Id           uint64
	Uuid         string    `orm:"size(64)"`
	JobUid       string    `orm:"size(64)"`
	Name         string    `orm:"size(64)"`
	Desc         string    `orm:"size(256)"`
	Initiator    string    `orm:"size(64)"`
	LocalName    string    `orm:"size(64);column(local_name)"`
	LocalDSName  string    `orm:"size(64);column(local_dsname)"`
	LocalDSIndex int32     `orm:"size(8);column(local_dsindex)"`
	LocalDSCount int64     `orm:"size(16);column(local_dscount)"`
	PartyName    string    `orm:"size(64);column(party_name)"`
	PartyDSName  string    `orm:"size(64);column(party_dsname)"`
	PartyDSIndex int32     `orm:"size(8);column(party_dsindex)"`
	PartyDSCount int64     `orm:"size(16);column(party_dscount)"`
	ServerIP     string    `orm:"size(64)"`
	ServerPort   int32     `orm:"size(8)"`
	Protocol     string    `orm:"size(32)"`
	Mode         string    `orm:"size(8)"`
	Status       string    `orm:"size(64)"`
	Result       string    `orm:"size(8)"`
	ExecStart    time.Time `orm:"null;type(datetime);"`
	ExecEnd      time.Time `orm:"null;type(datetime)"`
	Created      time.Time `orm:"auto_now_add;type(datetime)"`
	Updated      time.Time `orm:"auto_now;type(datetime)"`
}

// TableUnique ...
func (obj *Task) TableUnique() [][]string {
	return [][]string{
		[]string{"Uuid"},
	}
}

// AddTask insert a task
func AddTask(obj *Task) error {
	o := orm.NewOrm()
	_, err := o.Insert(obj)
	return err
}

// UpdateTask update the task
func UpdateTask(obj *Task) error {
	o := orm.NewOrm()
	_, err := o.Update(obj)
	return err
}

// DeleteTask delete a task
func DeleteTask(obj *Task) error {
	o := orm.NewOrm()
	_, err := o.Delete(obj)
	return err
}

// GetTaskById get the task by ID
func GetTaskById(id uint64) (*Task, error) {
	o := orm.NewOrm()
	obj := Task{Id: id}
	err := o.Read(&obj)
	return &obj, err
}

// GetTaskByUuid get the task by UUID
func GetTaskByUuid(uuid string) (*Task, error) {
	o := orm.NewOrm()
	obj := Task{Uuid: uuid}
	err := o.Read(&obj, "Uuid")
	return &obj, err
}

// ListTasks return all tasks
func ListTasks() ([]*Task, error) {
	o := orm.NewOrm()
	var objs []*Task
	_, err := o.QueryTable("task").All(&objs)
	return objs, err
}

// ListTasksOfJob return all tasks of the given job
func ListTasksOfJob(jobUID string) ([]*Task, error) {
	o := orm.NewOrm()
	var objs []*Task
	_, err := o.QueryTable("task").
		Filter("job_uid", jobUID).
		// Filter("mode", TaskRunMode_Server).
		All(&objs)
	return objs, err
}

// ListFailedTasksOfJob return all failed tasks of the given job
func ListFailedTasksOfJob(jobUID string) ([]*Task, error) {
	o := orm.NewOrm()
	var objs []*Task
	_, err := o.QueryTable("task").
		Filter("job_uid", jobUID).
		Filter("status", JobStatus_Failed).
		All(&objs)
	return objs, err
}

// GetOldestTask return the first Ready task in Server Mode by create time
func GetOldestTask() (*Task, error) {
	o := orm.NewOrm()
	var t Task
	err := o.QueryTable("task").
		Filter("status", TaskStatus_Ready).
		Filter("mode", TaskRunMode_Server).
		OrderBy("created").
		Limit(1, 0).
		One(&t)
	return &t, err
}

// GetFailedTasks return the failed tasks in Server Mode by create time
func GetFailedTasks() ([]*Task, error) {
	o := orm.NewOrm()
	var objs []*Task
	_, err := o.QueryTable("task").
		Filter("status", TaskStatus_Failed).
		Filter("mode", TaskRunMode_Server).
		OrderBy("created").
		All(&objs)
	return objs, err
}

// DeleteTasksOfJob delete all tasks of the job
func DeleteTasksOfJob(jobUID string) error {
	o := orm.NewOrm()
	m := &Task{
		JobUid: jobUID,
	}

	// _, err := o.QueryTable("task").
	// 	Filter("job_uid", jobUID).Delete()
	_, err := o.Delete(m, "JobUid")
	return err
}
