package models

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/goravel/framework/database/orm"
)

type XiaoyaStorages struct {
	orm.Model
	Driver   string   `gorm:"column:driver"`
	Addition Addition `gorm:"column:addition;" json:"addition"`
	Status   string   `gorm:"column:status"`
	Disabled string   `gorm:"column:disabled"`
}

func (r *XiaoyaStorages) Connection() string {
	return "xiaoya"
}

func (r *XiaoyaStorages) TableName() string {
	return "x_storages"
}

type Addition struct {
	ShareId  string `json:"share_id"`
	SharePwd string `json:"share_pwd"`
}

func (c *Addition) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

func (c *Addition) Scan(src any) error {
	return json.Unmarshal([]byte(src.(string)), c)
}
