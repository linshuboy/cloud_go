package commands

import (
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/contracts/queue"
	"github.com/goravel/framework/facades"
	"goravel/app/jobs"
	"goravel/app/models"
)

type PinyiAddWhiteList struct {
}

// Signature The name and signature of the console command.
func (receiver *PinyiAddWhiteList) Signature() string {
	return "pinyi:add_all_white_list"
}

// Description The console command description.
func (receiver *PinyiAddWhiteList) Description() string {
	return "为所有品易账号添加白名单"
}

// Extend The console command extend.
func (receiver *PinyiAddWhiteList) Extend() command.Extend {
	return command.Extend{}
}

// Handle Execute the console command.
func (receiver *PinyiAddWhiteList) Handle(ctx console.Context) error {
	var proxys []models.ProxyPinYi
	_ = facades.Orm().Query().Get(&proxys)
	for _, proxy := range proxys {
		err := facades.Queue().Job(&jobs.PinyiAddIpPool{}, []queue.Arg{
			{
				Type:  "ProxyPinYi",
				Value: proxy,
			},
		}).Dispatch()
		if err != nil {
			// do something
		}
	}
	return nil
}
