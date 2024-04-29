package routes

import (
	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"
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
			shares := models.AliDiskShares{}
			var newFlag = shares.RandomString(10)
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
		AliDiskShareFile := models.AliDiskShareFile{}
		facades.Orm().Query().Where("id = ?", 103).With("Share").Find(&AliDiskShareFile)
		if AliDiskShareFile.CompletedAt.IsValid() {
			facades.Log().Info(AliDiskShareFile.CompletedAt)
			facades.Log().Info("已经整完了")
			return nil
		} else {
			facades.Log().Info(AliDiskShareFile.CompletedAt)
		}
		AliDiskShareFile.Share.InitProxyIp()
		var list []models.AliDiskShareFile
		var nextMarker string
		list, nextMarker = AliDiskShareFile.Share.GetFileListByFileId(AliDiskShareFile.FileId, AliDiskShareFile.NextMarker)
		facades.Log().Info(AliDiskShareFile.Share.ProxyIp.Ip)
		facades.Log().Info(len(list))
		facades.Log().Info(nextMarker)
		for _, line := range list {
			facades.Log().Info(line.Name)
			facades.Orm().Query().Save(&line)
		}
		if nextMarker != "" {
			AliDiskShareFile1 := models.AliDiskShareFile{}
			AliDiskShareFile1.FileId = AliDiskShareFile.FileId
			AliDiskShareFile1.DriveId = AliDiskShareFile.DriveId
			AliDiskShareFile1.DomainId = AliDiskShareFile.DomainId
			AliDiskShareFile1.AliDiskShareId = AliDiskShareFile.AliDiskShareId
			AliDiskShareFile1.Name = AliDiskShareFile.Name
			AliDiskShareFile1.Type = AliDiskShareFile.Type
			AliDiskShareFile1.ParentFileId = AliDiskShareFile.ParentFileId
			AliDiskShareFile1.NextMarker = nextMarker
			facades.Orm().Query().Save(&AliDiskShareFile1)
		}
		AliDiskShareFile.CompletedAt = carbon.DateTime{Carbon: carbon.Now()}
		facades.Orm().Query().Save(&AliDiskShareFile)
		return ctx.Response().Success().Json(http.Json{"massage": AliDiskShareFile.ID})
	})
}
