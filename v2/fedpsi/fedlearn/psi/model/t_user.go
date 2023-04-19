package model

import (
	"time"

	"github.com/astaxie/beego/orm"
)

const (
	LockIfFailedCountPerDay uint32 = 5
)

type User struct {
	Id          uint64
	UserName    string    `orm:"size(64)"`
	UserPass    string    `orm:"size(64)"`
	DisplayName string    `orm:"size(64)"`
	Party       string    `orm:"size(64)"`
	IsRoot      bool      `orm:"null;size(16)"`
	FailedCount uint32    `orm:"size(8)"`
	FailedDate  string    `orm:"size(16)"`
	Created     time.Time `orm:"auto_now_add;type(datetime)"`
	Updated     time.Time `orm:"auto_now;type(datetime)"`
}

func (obj *User) TableUnique() [][]string {
	return [][]string{
		[]string{"UserName"},
	}
}

func AddUser(obj *User) error {
	o := orm.NewOrm()
	_, err := o.Insert(obj)
	return err
}

func UpdateUser(obj *User) error {
	o := orm.NewOrm()
	_, err := o.Update(obj)
	return err
}

func DeleteUser(obj *User) error {
	o := orm.NewOrm()
	_, err := o.Delete(obj)
	return err
}

func GetUserById(id uint64) (*User, error) {
	o := orm.NewOrm()
	obj := User{Id: id}
	err := o.Read(&obj)
	return &obj, err
}

func GetUserByUserName(userName string) (*User, error) {
	o := orm.NewOrm()
	obj := User{UserName: userName}
	err := o.Read(&obj, "UserName")
	return &obj, err
}

func ListUsers() ([]*User, error) {
	o := orm.NewOrm()
	var ts []*User
	_, err := o.QueryTable("user").
		All(&ts)
	return ts, err
}
