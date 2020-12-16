package user


import "github.com/gin-gonic/gin"

func Urls(r *gin.RouterGroup)  {
	router:=r.Group("user")
	router.GET("", userListView)
	router.DELETE(":id", userDestroyView)
	router.PUT(":id", userUpdateView)
	//router.POST("resetpwd/:id", resetPwd)
}
