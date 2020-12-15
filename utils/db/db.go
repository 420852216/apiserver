//connection 的值为：当前表字段-关联表的字段

package db

import (
	"apiserver/utils/logger"
	. "apiserver/utils/logger"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"time"
)

type Config struct {
	ENGINE      string
	NAME        string
	USER        string
	PASSWORD    string
	HOST        string
	PORT        string
	MAXIDLECONN int
	MAXOPENCONN int
	DEBUG bool
}

type modeler interface {
	TableName() string
}

type database struct {
	*sqlx.DB
	showSQL bool
}

var Sqlx *database

func Connect(cfg Config) (err error) {
	//dsn := "root:Initial1@(127.0.0.1:3306)/rbac?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",cfg.USER,
		cfg.PASSWORD,cfg.HOST,cfg.PORT,cfg.NAME)
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		logger.Log.Error("connect db failed, err:%v\n", zap.Error(err))
		return err
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Second * 10)
	Sqlx = new(database)
	Sqlx.DB=db
	Sqlx.showSQL = cfg.DEBUG
	return nil
}

func (db *database) Get(dest interface{}, query string, args ...interface{}) error {
	if db.showSQL{
		Log.Debug(query, zap.Any("args",args))
	}
	return db.Unsafe().Get(dest, query, args...)

}

func (db *database) PrepareNamed(query string) (*sqlx.NamedStmt, error) {
	sql,err:=db.DB.PrepareNamed(query)
	if db.showSQL && err==nil{
		Log.Debug(sql.QueryString)
	}
	return sql,err

}