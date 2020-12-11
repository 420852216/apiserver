package user

import (
	. "apiserver/utils/db"
	"apiserver/utils/logger"
	"go.uber.org/zap"
	"time"
)

func createUser(user *User) error {
	sql := "INSERT INTO user (password,name,email,phone,created_at,updated_at) VALUES(:password,:name,:email,:phone,:created_at,:updated_at)"
	//nstmt, err := Sqlx.PrepareNamed(sql)
	nstmt, err := Sqlx.PrepareNamed(sql)
	if err != nil {
		return err
	}
	user.UpdatedAt=time.Now()
	user.CreatedAt=time.Now()
	_,err = nstmt.Exec(user)
	logger.Log.Debug(sql,zap.Any("-->", *user))
	return err
}

func hasUser(email string) bool {
	sql := "SELECT 1 FROM user WHERE email=? LIMIT 1"
	var count int
	Sqlx.Get(&count,sql,email)
	if count==1{
		return true
	}
	logger.Log.Debug(sql,zap.String("-->", email))
	return false
}