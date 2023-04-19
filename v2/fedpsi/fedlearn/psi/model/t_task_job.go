package model

import (
	"time"

	"fedlearn/psi/common/config"

	"github.com/astaxie/beego/orm"
)

const (
	JobStatus_WaitingPartyConfirm = "awaiting_party_confirm"
	JobStatus_Confirmed           = "confirmed"
	JobStatus_Ready               = "ready"
	JobStatus_Running             = "running"
	JobStatus_Completed           = "completed"
	JobStatus_Success             = "success"
	JobStatus_Failed              = "failed"
	JobStatus_Cancel              = "cancel"
)

// Job is a collection of tasks
type Job struct {
	Id             uint64
	Uuid           string    `orm:"size(64)"`
	ActUid         string    `orm:"size(64)"`
	Name           string    `orm:"size(64)"`
	Desc           string    `orm:"size(256)"`
	Initiator      string    `orm:"size(64)"`
	LocalName      string    `orm:"size(64);column(local_name)"`
	LocalDSName    string    `orm:"size(64);column(local_dsname)"`
	LocalDSCount   int64     `orm:"size(16);column(local_dscount)"`
	PartyName      string    `orm:"size(64);column(party_name)"`
	PartyDSName    string    `orm:"size(64);column(party_dsname)"`
	PartyDSCount   int64     `orm:"size(16);column(party_dscount)"`
	Protocol       string    `orm:"size(32)"`
	Mode           string    `orm:"size(8)"`
	Status         string    `orm:"size(64)"`
	Result         string    `orm:"size(8)"`
	IntersectCount int64     `orm:"size(64)"`
	ExecStart      time.Time `orm:"null;type(datetime);"`
	ExecEnd        time.Time `orm:"null;type(datetime)"`
	Created        time.Time `orm:"auto_now_add;type(datetime)"`
	Updated        time.Time `orm:"auto_now;type(datetime)"`
}

// TableUnique ...
func (obj *Job) TableUnique() [][]string {
	return [][]string{
		[]string{"Uuid"},
	}
}

// AddJob insert a job
func AddJob(obj *Job) error {
	o := orm.NewOrm()
	_, err := o.Insert(obj)
	return err
}

// UpdateJob update job
func UpdateJob(obj *Job) error {
	o := orm.NewOrm()
	_, err := o.Update(obj)
	return err
}

// DeleteJob delete job
func DeleteJob(obj *Job) error {
	o := orm.NewOrm()
	_, err := o.Delete(obj)
	return err
}

// GetJobById get job by id
func GetJobById(id uint64) (*Job, error) {
	o := orm.NewOrm()
	obj := Job{Id: id}
	err := o.Read(&obj)
	return &obj, err
}

// GetJobByUuid get job by UUID
func GetJobByUuid(uuid string) (*Job, error) {
	o := orm.NewOrm()
	obj := Job{Uuid: uuid}
	err := o.Read(&obj, "Uuid")
	return &obj, err
}

// ListJobs return all jobs
func ListJobs() ([]*Job, error) {
	o := orm.NewOrm()
	var objs []*Job
	_, err := o.QueryTable("job").
		OrderBy("-created").
		All(&objs)
	return objs, err
}

func ListJobsOfActivityByPage(actUID string, pageNum int, pageSize int) ([]*Job, int64, error) {
	o := orm.NewOrm()
	var objs []*Job
	count, _ := o.QueryTable("job").
		Filter("act_uid", actUID).
		Count()
	_, err := o.QueryTable("job").
		Filter("act_uid", actUID).
		OrderBy("-created").
		Limit(pageSize).
		Offset((pageNum - 1) * pageSize).
		All(&objs)
	return objs, count, err
}

func ListJobsOfActivity(actUID string) ([]*Job, error) {
	o := orm.NewOrm()
	var objs []*Job
	_, err := o.QueryTable("job").
		Filter("act_uid", actUID).
		OrderBy("-created").
		All(&objs)
	return objs, err
}

// ListJobsByDataset return all jobs of local dataset name
func ListJobsByLocalDataset(name string) ([]*Job, error) {
	o := orm.NewOrm()
	var objs []*Job
	_, err := o.QueryTable("job").
		Filter("local_name", config.GetConfig().PartyName).
		Filter("local_dsname", name).
		All(&objs)
	return objs, err
}

// GetOldestConfirmedJob return the first confirmed job order by create_time
func GetOldestConfirmedJob() (*Job, error) {
	o := orm.NewOrm()
	var t Job
	err := o.QueryTable("job").
		Filter("status", JobStatus_Confirmed).
		OrderBy("created").
		Limit(1, 0).
		One(&t)
	return &t, err
}

// ListConfirmedJobs return all confirmed jobs
func ListConfirmedJobs() ([]*Job, error) {
	o := orm.NewOrm()
	var objs []*Job
	_, err := o.QueryTable("job").
		Filter("status", JobStatus_Confirmed).
		All(&objs)
	return objs, err
}
