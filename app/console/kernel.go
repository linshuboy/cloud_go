package console

import (
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/schedule"
	"github.com/goravel/framework/facades"
	"goravel/app/console/commands"
)

type Kernel struct {
}

func (kernel *Kernel) Schedule() []schedule.Event {
	return []schedule.Event{
		// 每一分钟获取一次白名单地址，如果当前ip 不在的话，则添加
		facades.Schedule().Command("pinyi:add_all_white_list").Cron("* * * * *").DelayIfStillRunning().OnOneServer().Name("add_all_white_list"),
		facades.Schedule().Command("pinyi:add_ip_pool").Cron("* * * * *").DelayIfStillRunning().OnOneServer().Name("add_all_white_list"),
		facades.Schedule().Command("pinyi:remove_invalid_ip").Cron("* * * * *").DelayIfStillRunning().OnOneServer().Name("add_all_white_list"),
		facades.Schedule().Command("xiaoya:sync").DailyAt("2:00").DelayIfStillRunning().OnOneServer().Name("add_all_white_list"),
	}
}

func (kernel *Kernel) Commands() []console.Command {
	return []console.Command{
		&commands.PinyiAddWhiteList{},
		&commands.PinyiAddIpPool{},
		&commands.PinyiRemoveInvalidIp{},
		&commands.SyncXiaoyaToDatabase{},
	}
}
