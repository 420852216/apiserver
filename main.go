package main

import (
	"apiserver/settings"
	"apiserver/utils/db"
	"apiserver/utils/logger"
)

func main() {
	settings.InitProject()
	defer db.Sqlx.Close()
	defer logger.Log.Sync()

}
