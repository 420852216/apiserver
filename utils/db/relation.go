package db

import (
	. "apiserver/utils/logger"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"reflect"
	"strings"
)

//subQuery 嵌套子查询
func (db database) getQuery(dest interface{}, field reflect.StructField, relations []string) (err error) {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	switch kind := field.Type.Kind(); kind {
	case reflect.Slice:
		subDest := reflect.New(field.Type)
		element := reflect.Indirect(subDest).Type().Elem()
		if element.Kind()==reflect.Ptr{
			element = element.Elem()
		}
		var tableName string
		if m, ok := reflect.New(element).Interface().(modeler); ok {
			tableName = m.TableName()
		}else {
			err := fmt.Errorf("%s: missing the TableName method", element.String())
			return err
		}
		sqlStr := fmt.Sprintf("SELECT * FROM %s WHERE %s = ?", tableName, relations[0])
		err := db.RelationSelect(subDest.Interface(), sqlStr, db.Mapper.FieldByName(destValue, relations[1]).Interface())
		if err != nil {
			return err
		}
		subDest = reflect.Indirect(subDest)
		if subDest.Len()==0{
			//如果没有值，需要插入一个空切片
			destValue.FieldByName(field.Name).Set(reflect.MakeSlice(field.Type,0,0))
		}else {
			destValue.FieldByName(field.Name).Set(subDest)
		}

	case reflect.Ptr:
		subDest := reflect.New(field.Type.Elem())
		var tableName string
		if m,ok:=subDest.Interface().(modeler);ok{
			tableName = m.TableName()
		}else {
			err := fmt.Errorf("%s: missing the TableName method", subDest.Type().Elem().String())
			return err
		}
		sqlStr := fmt.Sprintf("SELECT * FROM %s WHERE %s = ?", tableName, relations[0])
		err := db.RelationGet(subDest.Interface(), sqlStr, db.Mapper.FieldByName(destValue, relations[1]).Interface())
		if err != nil {
			return err
		}
		destValue.FieldByName(field.Name).Set(subDest)
	default:
		return fmt.Errorf("invalid type:%s", field.Type.String())
	}
	return nil
}

//RelationGet 查询一条数据，数据内部嵌套
func (db database) RelationGet(dest interface{}, query string, args ...interface{}) (err error) {
	err = db.Unsafe().Get(dest, query, args...)
	if err != nil {
		return
	}
	Log.Debug(query, zap.Any("-->",args))
	destType := reflect.TypeOf(dest).Elem()
	for i := 0; i < destType.NumField(); i++ {
		if destType.Field(i).Anonymous {
			continue
		}
		connection := destType.Field(i).Tag.Get("connection")
		if connection == "" {
			continue
		}
		relations := strings.Split(connection, "-")
		if len(relations) != 2 {
			return errors.New(fmt.Sprintf("connection tag error, must be 'string-string', but get '%v'", connection))
		}
		err = db.getQuery(dest, destType.Field(i), relations)
		if err != nil {
			return
		}
	}
	if cb:=reflect.ValueOf(dest).MethodByName("AfterQuery");cb.IsValid(){
		cb.Call(nil)
	}
	return nil
}
//RelationSelect 查询多条数据，数据内部嵌套
func (db database) RelationSelect(dest interface{}, query string, args ...interface{}) (err error) {
	err = db.Unsafe().Select(dest, query, args...)
	if err != nil {
		return
	}
	Log.Debug(query, zap.Any("-->",args))
	destVal := reflect.Indirect(reflect.ValueOf(dest))
	count := destVal.Len()
	if count == 0 {
		return nil
	}
	elementType := reflect.Indirect(destVal.Index(0)).Type()
	for i := 0; i < elementType.NumField(); i++ {
		if elementType.Field(i).Anonymous {
			continue
		}
		connection := elementType.Field(i).Tag.Get("connection")
		if connection == "" {
			continue
		}
		relations := strings.Split(connection, "-")
		if len(relations) != 2 {
			return errors.New(fmt.Sprintf("connection tag error, must be 'string-string', but get '%v'", connection))
		}
		err = db.selectQuery(dest,elementType.Field(i),relations)
		if err != nil {
			return
		}
	}

	if _,has:=reflect.PtrTo(elementType).MethodByName("AfterQuery");has{
		for i := 0; i < destVal.Len(); i++ {
			reflect.Indirect(destVal.Index(i)).Addr().MethodByName("AfterQuery").Call(nil)
		}
	}
	return nil
}

func (db database) selectQuery(dest interface{}, field reflect.StructField, relations []string) (err error){
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	sqlValMap := make(map[interface{}]interface{}, 0)
	for i := 0; i < destValue.Len(); i++ {
		k:=db.Mapper.FieldByName(destValue.Index(i), relations[1]).Interface()
		sqlValMap[k] = nil
	}
	sqlVals := make([]interface{},0,len(sqlValMap))
	for k:= range sqlValMap {
		sqlVals = append(sqlVals,k)
	}
	var tableName string
	var subDest reflect.Value
	kind := field.Type.Kind()
	if kind == reflect.Slice{
		subDest = reflect.New(field.Type)
		element := reflect.Indirect(subDest).Type().Elem()
		if element.Kind()==reflect.Ptr{
			element = element.Elem()
		}
		if m, ok := reflect.New(element).Interface().(modeler); ok {
			tableName = m.TableName()
		}else {
			err := fmt.Errorf("%s: missing the TableName method", element.String())
			return err
		}
	}else if kind == reflect.Ptr{
		subDest = reflect.New(field.Type.Elem())
		if m,ok:=subDest.Interface().(modeler);ok{
			tableName = m.TableName()
		}else {
			err := fmt.Errorf("%s: missing the TableName method", subDest.Type().Elem().String())
			return err
		}
		subDest = reflect.New(reflect.SliceOf(subDest.Type()))
	}else {
		return fmt.Errorf("invalid type:%s", field.Type.String())
	}
	sqlStr := fmt.Sprintf("SELECT * FROM %s WHERE %s IN (?);",tableName, relations[0])
	query, args, err := sqlx.In(sqlStr, sqlVals)
	if err != nil {
		return err
	}
	err = db.RelationSelect(subDest.Interface(), query, args...)
	if err != nil {
		return err
	}
	indirectSubDest:=reflect.Indirect(subDest)
	if indirectSubDest.Len()==0{return nil}//没有查到数据，直接返回
	subDestMap := make(map[interface{}]reflect.Value)
	if kind == reflect.Slice{
		for i := 0; i < indirectSubDest.Len(); i++ {
			val := indirectSubDest.Index(i)
			myId := db.Mapper.FieldByName(val, relations[0])
			//如果没有这条数据，则创建一个slice并把这条数据加到slice，然后把整个slice放进map
			if _, has := subDestMap[myId.Interface()]; !has {
				subDestMap[myId.Interface()] = reflect.New(reflect.SliceOf(field.Type.Elem())).Elem()
			}
			subDestMap[myId.Interface()] = reflect.Append(subDestMap[myId.Interface()], val)
		}
		for i := 0; i < destValue.Len(); i++ {
			//fmt.Println(destValue.Index(i))
			fid := db.Mapper.FieldByName(destValue.Index(i), relations[1])
			if sub, has := subDestMap[fid.Interface()]; has {
				reflect.Indirect(destValue.Index(i)).FieldByName(field.Name).Set(sub)
			}else {
				reflect.Indirect(destValue.Index(i)).FieldByName(field.Name).Set(reflect.MakeSlice(field.Type,0,0))
			}

		}
	}else {
		for i := 0; i < indirectSubDest.Len(); i++ {
			val := indirectSubDest.Index(i)
			myId := db.Mapper.FieldByName(val, relations[0])
			subDestMap[myId.Interface()] = val
		}
		for i := 0; i < destValue.Len(); i++ {
			fid := db.Mapper.FieldByName(destValue.Index(i), relations[1])
			if sub, has := subDestMap[fid.Interface()]; has {
				reflect.Indirect(destValue.Index(i)).FieldByName(field.Name).Set(sub)
			}
		}


	}
	return nil
}
