package model

import (
	"time"

	"github.com/astaxie/beego/orm"
)

// created -> data_available(initiator) -> sync to party -> awaiting_party_confirm
// awaiting_party_confirm -> data_available -> confirmed

const (
	ActivityStatus_Created             = "created"
	ActivityStatus_DataAvailable       = "data_available"
	ActivityStatus_WaitingPartyConfirm = "awaiting_party_confirm"
	ActivityStatus_Confirmed           = "confirmed"
	ActivityStatus_Running             = "running"
	ActivityStatus_Completed           = "completed"
	ActivityStatus_failed              = "failed"
)

// Activity is an intersection of biz
type Activity struct {
	Id            uint64
	Uuid          string    `orm:"size(64)"`
	Name          string    `orm:"size(64)"`
	Title         string    `orm:"size(256)"`
	SendID        string    `orm:"size(64);column(send_id)"`
	Desc          string    `orm:"size(256)"`
	InitParty     string    `orm:"size(64)"`
	FollowerParty string    `orm:"size(64)"`
	InitiatorData string    `orm:"size(16384)"`
	FollowerData  string    `orm:"size(16384)"`
	Status        string    `orm:"size(64)"`
	Created       time.Time `orm:"auto_now_add;type(datetime)"`
	Updated       time.Time `orm:"auto_now;type(datetime)"`
}

// TableUnique ...
func (obj *Activity) TableUnique() [][]string {
	return [][]string{
		[]string{"Name"},
	}
}

// AddActivity insert an activity
func AddActivity(obj *Activity) error {
	o := orm.NewOrm()
	_, err := o.Insert(obj)
	return err
}

// UpdateActivity update the activity
func UpdateActivity(obj *Activity) error {
	o := orm.NewOrm()
	_, err := o.Update(obj)
	return err
}

// DeleteActivity delete an activity
func DeleteActivity(obj *Activity) error {
	o := orm.NewOrm()
	_, err := o.Delete(obj)
	return err
}

func DeleteActivityByUUID(uuid string) error {
	o := orm.NewOrm()
	obj := &Activity{Uuid: uuid}
	_, err := o.Delete(obj, "Uuid")
	return err
}

// GetActivityById get the activity by ID
func GetActivityById(id uint64) (*Activity, error) {
	o := orm.NewOrm()
	obj := &Activity{Id: id}
	err := o.Read(obj)
	return obj, err
}

// GetActivityByUuid get the activity by UUID
func GetActivityByUuid(uuid string) (*Activity, error) {
	o := orm.NewOrm()
	obj := &Activity{Uuid: uuid}
	err := o.Read(obj, "Uuid")
	return obj, err
}

// ListActivities return all activities
func ListActivities() ([]*Activity, error) {
	o := orm.NewOrm()
	var objs []*Activity
	_, err := o.QueryTable("activity").All(&objs)
	return objs, err
}

// GetOldestCreatedActivity return the oldest created activity
func GetOldestCreatedActivity() (*Activity, error) {
	o := orm.NewOrm()
	var t Activity
	err := o.QueryTable("activity").
		Filter("status", ActivityStatus_Created).
		OrderBy("created").
		Limit(1, 0).
		One(&t)
	return &t, err
}

// GetOldestConfirmedActivity return the oldest confirmed activity
func GetOldestConfirmedActivity() (*Activity, error) {
	o := orm.NewOrm()
	var t Activity
	err := o.QueryTable("activity").
		Filter("status", ActivityStatus_Confirmed).
		OrderBy("created").
		Limit(1, 0).
		One(&t)
	return &t, err
}
