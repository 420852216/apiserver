package user

import (
	. "apiserver/utils/db"
	"fmt"
	"reflect"
	"strings"
	"time"
)

func createUser(user *User) error {
	sql := "INSERT INTO user (password,name,email,phone,created_at,updated_at) VALUES(:password,:name,:email,:phone,:created_at,:updated_at)"
	nstmt, err := Sqlx.PrepareNamed(sql)
	if err != nil {
		return err
	}
	user.UpdatedAt=time.Now()
	user.CreatedAt=time.Now()
	_,err = nstmt.Exec(user)
	return err
}

func hasUser(email string) bool {
	sql := "SELECT 1 FROM user WHERE email=? LIMIT 1"
	var count int
	Sqlx.Get(&count,sql,email)
	if count==1{
		return true
	}
	return false
}

func userLogin(dest interface{}, email, pwd string) error {
	sql := "SELECT name,email,phone FROM user WHERE email=? AND password=? LIMIT 1"
	return Sqlx.Get(dest,sql,email,pwd)
}

func updateUser(args interface{}, values interface{}) error {
	var sql strings.Builder
	sqlMap := make(map[string]interface{})
	sql.WriteString("UPDATE user SET ")
	{
		v := reflect.Indirect(reflect.ValueOf(values))
		t := reflect.TypeOf(values).Elem()
		first := true
		for i := 0; i < v.NumField(); i++ {
			if v.Field(i).IsZero() {
				continue
			}
			if !first {
				sql.WriteString(",")
			}
			first = false
			field:=t.Field(i).Tag.Get("db")
			sql.WriteString(field)
			sql.WriteString(" =:")
			sql.WriteString(field)
			sqlMap[field] = reflect.Indirect(v.Field(i)).Interface()
		}
	}
	{
		v := reflect.Indirect(reflect.ValueOf(args))
		t := reflect.TypeOf(args).Elem()
		first := true
		sql.WriteString(" WHERE ")
		for i := 0; i < v.NumField(); i++ {
			if v.Field(i).IsZero() {
				continue
			}
			if !first {
				sql.WriteString(" AND ")
			}
			first = false
			field:=t.Field(i).Tag.Get("db")
			sql.WriteString(field)
			sql.WriteString(" =:")
			sql.WriteString(field)
			sqlMap[field] = reflect.Indirect(v.Field(i)).Interface()
		}
	}
	nstmt, err := Sqlx.PrepareNamed(sql.String())
	if err != nil {
		fmt.Println(err)
		return err
	}
	rows,err := nstmt.Exec(sqlMap)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(rows.RowsAffected())
	fmt.Println(rows.LastInsertId())
	return nil
}