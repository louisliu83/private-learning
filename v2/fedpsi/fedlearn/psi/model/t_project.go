package model

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type Project struct {
	Id            uint64
	Uuid          string    `orm:"size(64)"`
	Name          string    `orm:"size(64)"`
	Type          string    `orm:"type(64)"`
	Desc          string    `orm:"size(256)"`
	InitParty     string    `orm:"size(64)"`
	FollowerParty string    `orm:"size(64)"`
	Status        string    `orm:"size(64)"`
	Creator       string    `orm:"size(64)"`
	UpdateUser    string    `orm:"size(64)"`
	Created       time.Time `orm:"auto_now_add;type(datetime)"`
	Updated       time.Time `orm:"auto_now;type(datetime)"`
}

func (obj *Project) TableUnique() [][]string {
	return [][]string{
		[]string{"Id"},
	}
}

func AddProject(obj *Project) error {
	o := orm.NewOrm()
	_, err := o.Insert(obj)
	return err
}

func UpdateProject(obj *Project) error {
	o := orm.NewOrm()
	_, err := o.Update(obj)
	return err
}

func DeleteProject(obj *Project) error {
	o := orm.NewOrm()
	_, err := o.Delete(obj)
	return err
}

func GetProjectById(id uint64) (*Project, error) {
	o := orm.NewOrm()
	obj := Project{Id: id}
	err := o.Read(&obj, "Id")
	return &obj, err
}

func GetProjectByUid(uid string) (*Project, error) {
	o := orm.NewOrm()
	obj := Project{Uuid: uid}
	err := o.Read(&obj, "Uuid")
	return &obj, err
}

func ListProjects(name string) ([]*Project, error) {
	o := orm.NewOrm()
	var objs []*Project
	_, err := o.QueryTable("project").
		Filter("name__icontains", name).
		OrderBy("-created").
		All(&objs)
	return objs, err
}
