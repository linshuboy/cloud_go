package routes

import (
	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"goravel/app/models"
	"gorm.io/gorm"

	"goravel/app/http/controllers"
)

func Api() {
	userController := controllers.NewUserController()
	facades.Route().Get("/users/{id}", userController.Show)
	facades.Route().Get("test", func(ctx http.Context) http.Response {
		err := facades.Orm().Transaction(func(tx orm.Transaction) error {
			var lines []models.XiaoyaStorages
			_ = facades.Orm().Query().Get(&lines)
			var newFlag = models.AliDiskShares{}.RandomString(10)
			for _, line := range lines {
				if line.Driver == "AliyundriveShare2Open" {
					var aliDiskShare models.AliDiskShares
					_ = facades.Orm().Query().WithTrashed().Where("share_id", line.Addition.ShareId).FirstOr(&aliDiskShare, func() error {
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
			return ctx.Response().Success().Json(http.Json{"massage": "失败"})
		}
		return ctx.Response().Success().Json(http.Json{"massage": "成功"})
	})
	facades.Route().Get("test2", func(ctx http.Context) http.Response {
		aliDiskShares := models.AliDiskShares{}
		_ = facades.Orm().Query().Where("id", 1).First(&aliDiskShares)
		shareFiles := aliDiskShares.GetFileListByFileId("root", "")
		for _, shareFile := range shareFiles {
			err := facades.Orm().Query().Save(&shareFile)
			if err != nil {
				return nil
			}
		}
		return ctx.Response().Success().Json(http.Json{"massage": aliDiskShares.ShareId})
	})
}
