package settings

import (
	"apiserver/apps/user"
	_ "apiserver/docs"
	"apiserver/settings/middleware"
	"apiserver/utils"
	"apiserver/utils/db"
	"apiserver/utils/logger"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
		DEBUG: DEBUG,
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
	router.Use(middleware.Cors())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.POST("/register", user.Register)
	router.POST("/login", user.Login,middleware.SetJSONWebToken)
	v1:=router.Group("/api")
	v1.Use(middleware.JWTMiddleware())
	user.Urls(v1)
	router.Run(":9527")
}
