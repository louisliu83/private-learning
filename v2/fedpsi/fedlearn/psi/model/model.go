package model

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
)

// Initdb...
func Initdb(dbPath string) error {
	if err := orm.RegisterDriver("sqlite", orm.DRSqlite); err != nil {
		return err
	}
	if err := orm.RegisterDataBase("default", "sqlite3", dbPath); err != nil {
		return err
	}
	orm.RegisterModel(new(Dataset),
		new(Project),
		new(Activity),
		new(Job),
		new(Task),
		new(User),
		new(Party),
		new(Namespace),
		new(UserNamespace))
	if err := orm.RunSyncdb("default", false, true); err != nil {
		return err
	}
	return nil
}
