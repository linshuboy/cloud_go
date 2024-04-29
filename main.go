package main

import (
	"github.com/goravel/framework/contracts/queue"
	"github.com/goravel/framework/facades"

	"goravel/bootstrap"
)

func main() {

	// This bootstraps the framework and gets it ready for use.
	bootstrap.Boot()

	//Start http server by facades.Route().
	go func() {
		if err := facades.Route().Run(); err != nil {
			facades.Log().Errorf("Route run error: %v", err)
		}
	}()

	//启动的时候自动维护
	facades.Artisan().Call("pinyi:add_all_white_list")
	//启动时间自动维护一次
	facades.Artisan().Call("pinyi:add_ip_pool")
	//启动时间自动维护一次
	facades.Artisan().Call("xiaoya:sync")

	// Start schedule by facades.Schedule
	go facades.Schedule().Run()

	go func() {
		if err := facades.Queue().Worker(&queue.Args{
			Connection: "redis",
			Concurrent: 10,
		}).Run(); err != nil {
			facades.Log().Errorf("Queue run error: %v", err)
		}
	}()

	select {}
}
