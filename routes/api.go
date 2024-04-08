package routes

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"goravel/app/models"

	"goravel/app/http/controllers"
)

func Api() {
	userController := controllers.NewUserController()
	facades.Route().Get("/users/{id}", userController.Show)
	facades.Route().Get("test", func(ctx http.Context) http.Response {
		var lines []models.XiaoyaStorages
		err := facades.Orm().Query().Get(&lines)
		if err != nil {
			return nil
		}
		//foreach ($s as $value){
		//    (new AliDiskShare())->firstOrCreate(['share_id'=>$value->addition['share_id']],['password'=>$value->addition['share_pwd'] ?? '']);
		//}
		facades.Log().Info(lines)
		facades.Log().Info("test")
		for _, line := range lines {
			facades.Log().Info(line.Addition)
			//var aliDiskShare models.AliDiskShares
			//facades.Orm().Query().Where("share_id", lines.Addition["share_id"]).FirstOr(&aliDiskShare, func() error {
			//	aliDiskShare.Name = "goravel"
			//	return nil
			//})
		}
		return ctx.Response().Success().Json(http.Json{})
	})
}
