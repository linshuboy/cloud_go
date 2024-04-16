package console

import (
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/schedule"
	"github.com/goravel/framework/facades"
	"goravel/app/models"
)

type Kernel struct {
}

func (kernel *Kernel) Schedule() []schedule.Event {
	return []schedule.Event{
		// 每一分钟获取一次白名单地址，如果当前ip 不在的话，则添加
		facades.Schedule().Call(func() {
			var proxys []models.ProxyPinYi
			_ = facades.Orm().Query().Get(&proxys)
		}).Cron("* * * * *").DelayIfStillRunning().OnOneServer().Name("test"),
	}
}

func (kernel *Kernel) Commands() []console.Command {
	return []console.Command{}
}
