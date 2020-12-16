package user

import (
	. "apiserver/utils"
	"github.com/gin-gonic/gin"
)
// @description 注册
// @description 不需要token
// @Tags  用户注册
// @summary 用户注册
// @Accept json
// @produce json
// @Param message body RegisterRep true "请求参数"
// @Success 200 {object} utils.SuccessResp "请求成功"
// @Router /register [post]
func Register(ctx *gin.Context) {
	serializer := new(RegisterRep)
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

// @description 登陆
// @description 不需要token
// @Tags  用户登陆
// @summary 用户登陆
// @Accept json
// @produce json
// @Param message body LoginReq true "请求参数"
// @Success 200 {object} utils.SuccessResp{data=user.LoginResp} "请求成功"
// @Router /login [post]
func Login(ctx *gin.Context) {
	serializer := new(LoginReq)
	if err := ctx.ShouldBindJSON(serializer); err != nil {
		Failed(ctx, "参数格式错误", err.Error())
		return
	}
	if err := Validator.Struct(serializer); err != nil {
		Failed(ctx, "参数错误", ErrorTranslate(err))
		return
	}
	if user,err := serializer.login(); err != nil {
		Failed(ctx, "账号或密码错误", err.Error())
		return
	}else {
		ctx.Set("user",user)
		return
	}
}

// @description 修改用户信息
// @Tags  修改用户信息
// @summary 修改用户信息
// @Security ApiKeyAuth
// @Accept json
// @produce json
// @Param id path int true "user ID"
// @Param id body userUpdateReq true "修改参数"
// @Success 200 {object} utils.SuccessResp "请求成功"
// @Router /api/user/{id} [put]
func userUpdateView(ctx *gin.Context) {
	args := new(idUri)
	if err:=ctx.ShouldBindUri(args); err!=nil{
		Failed(ctx, "用户ID错误", err.Error())
		return
	}
	values := new(userUpdateReq)
	if err:=ctx.ShouldBindJSON(values);err!=nil{
		Failed(ctx, "参数错误", err.Error())
		return
	}
	if err:=updateUser(args,values);err!=nil{
		Failed(ctx, "参数错误", err.Error())
		return
	}
	Success(ctx, "修改成功", nil)
}

// @description 删除用户
// @Tags  删除用户
// @summary 删除用户
// @Security ApiKeyAuth
// @Accept json
// @produce json
// @Param id path int true "user ID"
// @Success 200 {object} utils.SuccessResp "请求成功"
// @Router /api/user/{id} [delete]
func userDestroyView(ctx *gin.Context) {
	args := new(idUri)
	if err:=ctx.ShouldBindUri(args); err!=nil{
		Failed(ctx, "用户ID错误", err.Error())
		return
	}
	if err:=deleteUser(args);err!=nil{
		Failed(ctx, "参数错误", err.Error())
		return
	}
	Success(ctx, "删除成功", nil)
}

// @description 删除用户
// @Tags  删除用户
// @summary 删除用户
// @Security ApiKeyAuth
// @Accept json
// @produce json
// @Param id path int true "user ID"
// @Success 200 {object} utils.SuccessResp "请求成功"
// @Router /api/user/{id} [delete]
func userListView(ctx *gin.Context) {
	filter:=new(userFilter)
	err:=ctx.ShouldBindQuery(filter)
	if err != nil {
		Failed(ctx, "参数错误", err.Error())
	}
	users:=make([]userSerializer,0)
	if err:=filterUser(&users,filter);err!=nil{
		Failed(ctx, "查询失败", err.Error())
		return
	}
	Success(ctx, "查询成功", users)
}