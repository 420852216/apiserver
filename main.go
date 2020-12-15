package main

import (
	"apiserver/settings"
	"apiserver/utils/db"
	"apiserver/utils/logger"
)

// @title Swagger API 文档
// @version v1
// @description  API 文档
// @contact.name 是
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @host 127.0.0.1:9527
// @BasePath /
func main() {
	settings.InitProject()
	defer db.Sqlx.Close()
	defer logger.Log.Sync()

}

