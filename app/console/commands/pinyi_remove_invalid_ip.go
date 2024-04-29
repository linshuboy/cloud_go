package commands

import (
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/contracts/queue"
	"github.com/goravel/framework/facades"
	"goravel/app/jobs"
)

type PinyiRemoveInvalidIp struct {
}

// Signature The name and signature of the console command.
func (receiver *PinyiRemoveInvalidIp) Signature() string {
	return "pinyi:remove_invalid_ip"
}

// Description The console command description.
func (receiver *PinyiRemoveInvalidIp) Description() string {
	return "为所有品易账号添加ip池"
}

// Extend The console command extend.
func (receiver *PinyiRemoveInvalidIp) Extend() command.Extend {
	return command.Extend{}
}

// Handle Execute the console command.
func (receiver *PinyiRemoveInvalidIp) Handle(ctx console.Context) error {
	err := facades.Queue().Job(&jobs.PinyiRemoveInvalidIp{}, []queue.Arg{}).Dispatch()
	if err != nil {
	}
	return nil
}
