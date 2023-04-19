package model

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type Party struct {
	Id               uint64
	Name             string    `orm:"size(64)"`
	Scheme           string    `orm:"size(32)"`
	ControllerServer string    `orm:"size(64)"`
	ControllerPort   int32     `orm:"size(8)"`
	WorkServer       string    `orm:"size(64)"`
	WorkPort         int32     `orm:"size(8)"`
	Status           string    `orm:"size(32)"`
	Token            string    `orm:"size(8192)"`
	Created          time.Time `orm:"auto_now_add;type(datetime)"`
	Updated          time.Time `orm:"auto_now;type(datetime)"`
}

func (obj *Party) TableUnique() [][]string {
	return [][]string{
		[]string{"Name"},
	}
}

func AddParty(obj *Party) error {
	o := orm.NewOrm()
	_, err := o.Insert(obj)
	return err
}

func UpdateParty(obj *Party) error {
	o := orm.NewOrm()
	_, err := o.Update(obj)
	return err
}

func DeleteParty(obj *Party) error {
	o := orm.NewOrm()
	_, err := o.Delete(obj)
	return err
}

func GetPartyByName(name string) (*Party, error) {
	o := orm.NewOrm()
	obj := Party{Name: name}
	err := o.Read(&obj, "Name")
	return &obj, err
}

func ListParties() ([]*Party, error) {
	o := orm.NewOrm()
	var objs []*Party
	_, err := o.QueryTable("party").
		All(&objs)
	return objs, err
}
