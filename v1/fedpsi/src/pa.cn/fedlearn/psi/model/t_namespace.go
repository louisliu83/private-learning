package model

import (
	"time"

	"github.com/astaxie/beego/orm"
)

const (
	NamespaceRoleAdmin  = "admin"
	NamespaceRoleMember = "member"
)

type Namespace struct {
	Id          uint64
	Name        string    `orm:"size(64)"`
	Desc        string    `orm:"size(64)"`
	DisplayName string    `orm:"size(64)"`
	Party       string    `orm:"size(64)"`
	Created     time.Time `orm:"auto_now_add;type(datetime)"`
	Updated     time.Time `orm:"auto_now;type(datetime)"`
}

func (obj *Namespace) TableUnique() [][]string {
	return [][]string{
		[]string{"Name"},
	}
}

type UserNamespace struct {
	Id        uint64
	User      string    `orm:"size(64)"`
	Namespace string    `orm:"size(64)"`
	Role      string    `orm:"size(64)"`
	Created   time.Time `orm:"auto_now_add;type(datetime)"`
	Updated   time.Time `orm:"auto_now;type(datetime)"`
}

func (obj *UserNamespace) TableUnique() [][]string {
	return [][]string{
		[]string{"User", "Namespace"},
	}
}

// DB CRUD

func AddNamespace(obj *Namespace) error {
	o := orm.NewOrm()
	_, err := o.Insert(obj)
	return err
}

func UpdateNamespace(obj *Namespace) error {
	o := orm.NewOrm()
	_, err := o.Update(obj)
	return err
}

func DeleteNamespace(obj *Namespace) error {
	o := orm.NewOrm()
	_, err := o.Delete(obj)
	return err
}

func GetNamespaceById(id uint64) (*Namespace, error) {
	o := orm.NewOrm()
	obj := Namespace{Id: id}
	err := o.Read(&obj)
	return &obj, err
}

func GetNamespaceByName(name string) (*Namespace, error) {
	o := orm.NewOrm()
	obj := Namespace{Name: name}
	err := o.Read(&obj, "Name")
	return &obj, err
}

func AddUserNamespace(obj *UserNamespace) error {
	o := orm.NewOrm()
	_, err := o.Insert(obj)
	return err
}

func UpdateUserNamespace(obj *UserNamespace) error {
	o := orm.NewOrm()
	_, err := o.Update(obj)
	return err
}

func DeleteUserNamespace(obj *UserNamespace) error {
	o := orm.NewOrm()
	_, err := o.Delete(obj)
	return err
}

func GetNamespacesByUser(user string) (map[Namespace]string, error) {
	o := orm.NewOrm()
	var objs []*UserNamespace
	_, err := o.QueryTable("user_namespace").
		Filter("user", user).
		All(&objs)
	if err != nil {
		return nil, err
	}
	nsMap := map[Namespace]string{}
	for _, obj := range objs {
		ns, err := GetNamespaceByName(obj.Namespace)
		if err == nil {
			nsMap[*ns] = obj.Role
		}
	}
	return nsMap, nil
}

func GetUsersByNamespace(ns string) ([]*User, error) {
	o := orm.NewOrm()
	var objs []*UserNamespace
	_, err := o.QueryTable("user_namespace").
		Filter("namespace", ns).
		All(&objs)
	if err != nil {
		return nil, err
	}
	users := make([]*User, 0)
	for _, obj := range objs {
		u, err := GetUserByUserName(obj.User)
		if err == nil {
			users = append(users, u)
		}
	}
	return users, nil
}

func GetNamespacesAll() ([]*Namespace, error) {
	o := orm.NewOrm()
	var objs []*Namespace
	_, err := o.QueryTable("namespace").All(&objs)
	return objs, err
}

/*
* Logic except CRUD
 */

func CheckNamespaceExists(name string) bool {
	ns, err := GetNamespaceByName(name)
	if err == nil && ns != nil {
		return true
	}
	return false
}
