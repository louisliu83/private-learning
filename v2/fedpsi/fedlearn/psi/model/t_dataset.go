package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

const (
	DatasetOK             = "loaded" // means the dataset is uploaded to the server
	DatasetSharding       = "sharding"
	DatasetAvailable      = "available" // means the dataset can be used as intersect task source
	DatasetExpired        = "expired"   // means the dataste expired
	DatasetDeleted        = "deleted"
	DatasetDownloading    = "downloading"
	DatasetDownloadFailed = "download_failed"
	DatasetCopyFailed     = "copy_failed"
)

const (
	DatasetIsPublic  = "public"
	TimeFormatString = "2006-01-02 15:04:05"
)

func GetAnchorTime() time.Time {
	a, _ := time.Parse(TimeFormatString, "2021-01-01 08:00:00")
	return a
}

type Dataset struct {
	Id          uint64
	Name        string    `orm:"size(64)"`
	Type        string    `orm:"type(64)"`
	Desc        string    `orm:"size(256)"`
	Index       int32     `orm:"size(8)"`
	Shards      int32     `orm:"size(8)"`
	Path        string    `orm:"size(1024)"`
	Count       int64     `orm:"size(16)"`       // line of file
	Size        int64     `orm:"size(16)"`       // size of file
	Md5         string    `orm:"size(64)"`       // Md5 of file
	URL         string    `orm:"size(1024)"`     // url of file
	Status      string    `orm:"null;type(32)"`  // status of file, maybe file download failed.
	Parties     string    `orm:"null;size(256)"` // Comma-Sperated Parties that can access this dataset
	BizContext  string    `orm:"null;size(1024)"`
	ExpiredDate time.Time `orm:"null;type(datetime)"`
	Created     time.Time `orm:"auto_now_add;type(datetime)"`
	Updated     time.Time `orm:"auto_now;type(datetime)"`
}

func (obj *Dataset) IsValid() bool {
	a := GetAnchorTime()
	return (obj.ExpiredDate.Before(a) || obj.ExpiredDate.After(time.Now())) &&
		obj.Status != DatasetDeleted
}

func (obj *Dataset) TableUnique() [][]string {
	return [][]string{
		[]string{"Name", "Index"},
	}
}

func AddDataset(obj *Dataset) error {
	o := orm.NewOrm()
	_, err := o.Insert(obj)
	return err
}

func UpdateDataset(obj *Dataset) error {
	o := orm.NewOrm()
	_, err := o.Update(obj)
	return err
}

func DeleteDataset(obj *Dataset) error {
	o := orm.NewOrm()
	_, err := o.Delete(obj)
	return err
}

func DeleteDatasetsByName(name string) error {
	o := orm.NewOrm()
	m := &Dataset{
		Name: name,
	}
	_, err := o.Delete(m, "Name")
	// _, err := o.QueryTable("dataset").Filter("name", name).Delete()
	return err
}

func GetDatasetById(id uint64) (*Dataset, error) {
	o := orm.NewOrm()
	obj := Dataset{Id: id}
	err := o.Read(&obj)
	return &obj, err
}

func GetDatasetByBizContext(bizContext string) (*Dataset, error) {
	o := orm.NewOrm()
	obj := Dataset{BizContext: bizContext}
	err := o.Read(&obj, "BizContext")
	return &obj, err
}

func GetDatasetByNameAndIndex(name string, index int32) (*Dataset, error) {
	o := orm.NewOrm()
	obj := Dataset{Name: name, Index: index}
	err := o.Read(&obj, "Name", "Index")
	return &obj, err
}

func GetDatasetShards(name string) ([]*Dataset, error) {
	shards, err := GetDatasetByName(name)
	if err != nil {
		return nil, err
	}
	if len(shards) <= 1 {
		return shards, nil
	}
	retShards := shards[1:]
	return retShards, nil
}

func GetDatasetByName(name string) ([]*Dataset, error) {
	o := orm.NewOrm()
	var objs []*Dataset
	_, err := o.QueryTable("dataset").
		Filter("name", name).
		OrderBy("index").
		All(&objs)
	return objs, err
}

func ListDataset() ([]*Dataset, error) {
	o := orm.NewOrm()
	var objs []*Dataset
	_, err := o.QueryTable("dataset").
		Filter("index", int32(0)).
		All(&objs)
	return objs, err
}

func ListDatasetBySrcParty(partyName string) ([]*Dataset, error) {
	o := orm.NewOrm()
	var objs []*Dataset
	_, err := o.QueryTable("dataset").
		Filter("index", int32(0)).
		All(&objs)
	if err != nil {
		return nil, err
	}
	retObjs := make([]*Dataset, 0)
	for _, obj := range objs {
		if strings.Contains(obj.Parties, fmt.Sprintf("%s,", partyName)) ||
			strings.Contains(obj.Parties, fmt.Sprintf(",%s", partyName)) {
			retObjs = append(retObjs, obj)
		}
	}
	return retObjs, nil
}
