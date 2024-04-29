package jobs

import (
	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"
	"goravel/app/models"
)

type GetAliShareFile struct {
}

// Signature The name and signature of the job.
func (receiver *GetAliShareFile) Signature() string {
	return "GetAliShareFile"
}

// Handle Execute the job.
func (receiver *GetAliShareFile) Handle(args ...any) error {
	AliDiskShareFile := models.AliDiskShareFile{}
	for _, arg := range args {
		AliDiskShareFile = arg.(models.AliDiskShareFile)
	}
	err := facades.Orm().Transaction(func(tx orm.Transaction) error {
		facades.Orm().Query().Where("id = ?", AliDiskShareFile.ID).With("Share").Find(&AliDiskShareFile)
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
		return nil
	})
	if err != nil {
	}
	return nil
}
