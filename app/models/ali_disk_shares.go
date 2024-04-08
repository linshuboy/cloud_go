package models

import (
	"github.com/goravel/framework/database/orm"
)

type AliDiskShares struct {
	orm.Model
	ShareId  string `gorm:"column:share_id"`
	Password string `gorm:"column:password"`
	SyncFlag string `gorm:"column:flag"`
	orm.SoftDeletes
}
