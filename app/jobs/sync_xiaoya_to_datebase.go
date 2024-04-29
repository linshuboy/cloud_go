package jobs

import (
	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/facades"
	"goravel/app/models"
	"gorm.io/gorm"
)

type SyncXiaoyaToDatabase struct {
}

// Signature The name and signature of the job.
func (receiver *SyncXiaoyaToDatabase) Signature() string {
	return "SyncXiaoyaToDatabase"
}

// Handle Execute the job.
func (receiver *SyncXiaoyaToDatabase) Handle(args ...any) error {
	err := facades.Orm().Transaction(func(tx orm.Transaction) error {
		var lines []models.XiaoyaStorages
		_ = facades.Orm().Query().Get(&lines)
		shares := models.AliDiskShares{}
		var newFlag = shares.RandomString(10)
		for _, line := range lines {
			if line.Driver == "AliyundriveShare2Open" {
				var aliDiskShare models.AliDiskShares
				_ = facades.Orm().Query().WithTrashed().Where("share_id = ?", line.Addition.ShareId).FirstOr(&aliDiskShare, func() error {
					aliDiskShare.ShareId = line.Addition.ShareId
					aliDiskShare.Password = line.Addition.SharePwd
					return nil
				})
				aliDiskShare.SyncFlag = newFlag
				err := facades.Orm().Query().Save(&aliDiskShare)
				if err != nil {
					return err
				}
			}
		}
		var noDeleteLines []models.AliDiskShares
		_ = facades.Orm().Query().WithTrashed().Where("flag = ?", newFlag).Get(&noDeleteLines)
		for _, line := range noDeleteLines {
			if line.DeletedAt.Valid {
				line.DeletedAt = gorm.DeletedAt{}
				_ = facades.Orm().Query().Save(&line)
			}
		}
		var deleteLines []models.AliDiskShares
		_ = facades.Orm().Query().Where("flag != ?", newFlag).Get(&deleteLines)
		for _, line := range deleteLines {
			_, _ = facades.Orm().Query().Delete(&line)
		}
		//
		return nil
	})
	if err != nil {
	}
	return nil
}
