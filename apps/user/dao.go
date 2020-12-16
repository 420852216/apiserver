package user

import (
	. "apiserver/utils/db"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

type pager interface {
	GetPageSize()int
	GetPage()int
}

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
		if v.IsZero(){
			return errors.New("修改字段不能为空")
		}
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
		if v.IsZero(){
			return errors.New("修改条件不能为空")
		}
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
		return err
	}
	_,err = nstmt.Exec(sqlMap)
	return err
}

func deleteUser(args interface{}) error {
	var sql strings.Builder
	sql.WriteString("DELETE FROM user WHERE ")
	v := reflect.Indirect(reflect.ValueOf(args))
	if v.IsZero(){
		return errors.New("删除条件不能为空")
	}
	first := true
	t := reflect.TypeOf(args).Elem()
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
	}
	nstmt, err := Sqlx.PrepareNamed(sql.String())
	if err != nil {
		return err
	}
	_,err = nstmt.Exec(args)
	return err
}

func filterUser(dest interface{}, args interface{}) error {
	var sql strings.Builder
	sql.WriteString("SELECT * FROM user ")
	v := reflect.Indirect(reflect.ValueOf(args))
	if v.IsZero(){
		return Sqlx.Unsafe().Select(dest, sql.String())
	}
	t := reflect.TypeOf(args).Elem()
	first := true
	for i := 0; i < v.NumField()-2; i++ {
		field,ok:=t.Field(i).Tag.Lookup("db")
		if v.Field(i).IsZero()||!ok{
			continue
		}
		if first {
			sql.WriteString("WHERE ")
		}else {
			sql.WriteString(" AND ")
		}
		first = false
		sql.WriteString(field)
		sql.WriteString(" =:")
		sql.WriteString(field)
	}
	if p,ok:=args.(pager);ok{
		limit:=p.GetPageSize()
		offset:=p.GetPage()*limit
		if limit>0 && offset>0{
			pageClause:=fmt.Sprintf(" LIMIT %d OFFSET %d",limit,(offset-1)*limit)
			sql.WriteString(pageClause)
		}
	}
	nstmt, err := Sqlx.PrepareNamed(sql.String())
	if err != nil {
		return err
	}
	return nstmt.Unsafe().Select(dest, args)
	//return Sqlx.Unsafe().Select(dest, sql.String(),args)
}