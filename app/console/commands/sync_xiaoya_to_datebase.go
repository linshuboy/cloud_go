package commands

import (
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/contracts/queue"
	"github.com/goravel/framework/facades"
	"goravel/app/jobs"
)

type SyncXiaoyaToDatabase struct {
}

// Signature The name and signature of the console command.
func (receiver *SyncXiaoyaToDatabase) Signature() string {
	return "xiaoya:sync"
}

// Description The console command description.
func (receiver *SyncXiaoyaToDatabase) Description() string {
	return "为所有品易账号添加ip池"
}

// Extend The console command extend.
func (receiver *SyncXiaoyaToDatabase) Extend() command.Extend {
	return command.Extend{}
}

// Handle Execute the console command.
func (receiver *SyncXiaoyaToDatabase) Handle(ctx console.Context) error {
	err := facades.Queue().Job(&jobs.SyncXiaoyaToDatabase{}, []queue.Arg{}).Dispatch()
	if err != nil {
	}
	return nil
}
