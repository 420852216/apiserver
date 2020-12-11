package user

import (
	. "apiserver/utils"
	"github.com/gin-gonic/gin"
)

func Register(ctx *gin.Context) {
	serializer := new(registerSerializer)
	if err := ctx.ShouldBindJSON(serializer); err != nil {
		Failed(ctx,"参数格式错误", err.Error())
		return
	}
	if err := Validator.Struct(serializer); err != nil {
		Failed(ctx,"参数错误", ErrorTranslate(err))
		return
	}
	if err := serializer.isValid(); err != nil {
		Failed(ctx,"Email已经存在", err.Error())
		return
	}
	if err := serializer.create();err != nil {
		Failed(ctx, "注册失败", err.Error())
		return
	}
	Success(ctx, "注册成功", nil)
}
