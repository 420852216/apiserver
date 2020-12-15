package middleware

import (
	"apiserver/apps/user"
	"apiserver/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"time"
)

type JSONWebTokenSerializer struct {
	jwt.StandardClaims
	*user.LoginResp
}

const (
	expireTime = time.Hour * 24
	secretKey = "sunLine"
)

func obtainJSONWebToken(user *user.LoginResp) error {
	expireTime := time.Now().Add(expireTime)
	stdClaims := jwt.StandardClaims{
		ExpiresAt: expireTime.Unix(), //过期时间
		IssuedAt:  time.Now().Unix(), //发行时间
		Issuer:    "sunLine",
	}
	claims := JSONWebTokenSerializer{stdClaims,user}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return err
	}
	user.Authorization = tokenString
	return nil
}

func VerifyJSONWebToken(tokenString string)(*user.LoginResp, error) {
	claims := new(JSONWebTokenSerializer)
	_, err := jwt.ParseWithClaims(tokenString, claims,
		func(token *jwt.Token) (interface{}, error) {return []byte(secretKey), nil})
	if err != nil {
		return nil, err
	}
	return claims.LoginResp, err
}

func JWTMiddleware() gin.HandlerFunc {
	return func (c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		user, err := VerifyJSONWebToken(token)
		if err != nil {
			utils.AuthFailed(c,err.Error())
			c.Abort()
			return
		}
		c.Set("user", *user)
		c.Next()
		// after request
	}
}

func SetJSONWebToken(ctx *gin.Context) {
	if v,ok:=ctx.Get("user");ok{
		user:=v.(*user.LoginResp)
		err:=obtainJSONWebToken(user)
		if err != nil {
			utils.Failed(ctx,"获取token失败",err.Error())
			return
		}
		utils.Success(ctx,"登录成功", user)
		return
	}
	return
}