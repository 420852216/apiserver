package settings

import (
	"apiserver/apps/user"
	"apiserver/utils"
	"apiserver/utils/db"
	"apiserver/utils/logger"
	"github.com/gin-gonic/gin"
)

var DEBUG = true

func InitProject()  {
	logger.InitLogger(DEBUG)
	dbConfig := db.Config{
		ENGINE:      "mysql",
		NAME:        "apiserver",
		USER:        "root",
		PASSWORD:    "Initial1",
		HOST:        "127.0.0.1",
		PORT:        "3306",
		MAXIDLECONN: 10,
		MAXOPENCONN: 128,
	}
	err:=db.Connect(dbConfig)
	if err != nil {
		return
	}
	utils.InitValidator()
	/* 生产环境
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(middleware.Ginzap(logger.Log,"2006/01/02 15:04:05"),
		middleware.RecoveryWithZap(logger.Log,true))
	*/
	router := gin.Default()
	router.POST("/register", user.Register)
	router.Run(":9527")
}
