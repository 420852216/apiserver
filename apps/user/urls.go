package user


import "github.com/gin-gonic/gin"

func Urls(r *gin.RouterGroup)  {
	router:=r.Group("user")
	//router.GET("", userList)
	//router.DELETE("delete/:id", userDelete)
	router.PUT(":id", userUpdate)
	//router.POST("resetpwd/:id", resetPwd)
}
