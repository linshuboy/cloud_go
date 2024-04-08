package models

import (
	"github.com/goravel/framework/database/orm"
)

type XiaoyaStorages struct {
	orm.Model
	Driver   string `gorm:"column:driver"`
	Addition string `gorm:"column:addition"`
	Status   string `gorm:"column:status"`
	Disabled string `gorm:"column:disabled"`
}

func (r *XiaoyaStorages) Connection() string {
	return "xiaoya"
}

func (r *XiaoyaStorages) TableName() string {
	return "x_storages"
}
