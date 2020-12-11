//全局响应API请求的返回数据格式
package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	SUCCESS = 0
	FAILED = 1
)

type Response struct {
	Code int `json:"code"`
	Msg interface{} `json:"msg"`
	Data interface{} `json:"data"`
}

func Success(ctx *gin.Context, msg string, v interface{}) {
	ctx.JSON(http.StatusOK, Response{
		Code: SUCCESS,
		Msg:  msg,
		Data: v,
	})
}

func Failed(ctx *gin.Context, msg string, v interface{}) {
	ctx.JSON(http.StatusOK, Response{
		Code: FAILED,
		Msg:  msg,
		Data: v,
	})
}

