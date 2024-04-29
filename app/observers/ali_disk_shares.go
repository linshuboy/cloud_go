package observers

import (
	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/contracts/queue"
	"github.com/goravel/framework/facades"
	"goravel/app/jobs"
	"goravel/app/models"
)

type AliDiskShareFileObservers struct{}

func (u *AliDiskShareFileObservers) Retrieved(event orm.Event) error {
	return nil
}

func (u *AliDiskShareFileObservers) Creating(event orm.Event) error {
	return nil
}

func (u *AliDiskShareFileObservers) Updating(event orm.Event) error {
	return nil
}

func (u *AliDiskShareFileObservers) Updated(event orm.Event) error {
	return nil
}

func (u *AliDiskShareFileObservers) Saving(event orm.Event) error {
	return nil
}

func (u *AliDiskShareFileObservers) Saved(event orm.Event) error {
	return nil
}

func (u *AliDiskShareFileObservers) Deleting(event orm.Event) error {
	return nil
}

func (u *AliDiskShareFileObservers) Deleted(event orm.Event) error {
	return nil
}

func (u *AliDiskShareFileObservers) ForceDeleting(event orm.Event) error {
	return nil
}

func (u *AliDiskShareFileObservers) ForceDeleted(event orm.Event) error {
	return nil
}

func (u *AliDiskShareFileObservers) Created(event orm.Event) error {
	if event.GetAttribute("Type") == "folder" {
		share := &models.AliDiskShares{}
		_ = event.Query().Where("id = ?", event.GetAttribute("ID")). /*.With("Files")*/ Find(&share)
		err := facades.Queue().Job(&jobs.GetAliShareFile{}, []queue.Arg{
			{
				Type:  "AliDiskShareFile",
				Value: share,
			},
		}).Dispatch()
		if err != nil {
			// do something
		}
	}
	return nil
}
