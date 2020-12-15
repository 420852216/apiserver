# apiserver
go web

## JWT全局中间件：
1.登陆接口使用了两个Handler,`router.POST("/login", user.Login,middleware.SetJSONWebToken)`
  第一个是根据账号密码查询数据库获取用户，第二个是根据用户信息生成token，并把token加到用户信息的结构体内，返回json数据

## 接口文档`swagger`
### 安装和使用
1. 下载模块：
```shell script
go get github.com/swaggo/swag/cmd/swag
go get github.com/swaggo/swag/cmd/swag
go get github.com/swaggo/swag/cmd/swag
```
2. 在`$GOPATH/bin/`下会看到多了一个`swag`。把`$GOPATH/bin/`加到`PATH`
3. 在包含`main.go`的Go工程的根目录下执行`swag init`，swag会检索当前工程里的swag注释，生成`docs.go`以及`swagger.json/yaml`
4. 在`router`里面注册路由`router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))`
5. 在`main`文件或者路由文件导入生成的doc目录：`_ "项目目录/docs"`。如果漏掉这步会报错：**not yet registered swag**
6. api文档地址**localhost:port/swagger/index.html**

### 注解说明
1. 入口main函数注解：

常规注解：
名称| 描述 | 例子
---|---|---|---
title |项目标题  |// @title XXXAPI文档
version |项目版本  |// @version x.x.x
description |项目描述  |// @description 项目描述
termsOfService |服务条款  |// @termsOfService http://xxxx.com
contact.name |联系人  |// @contact.name xxx
contact.url |指向联系人的URL，2个必须一起使用  |// @contact.url http://xxxx.com
contact.email |联系人邮箱，3个建议一起使用  |// @contact.email xxxx@163.com
license.name |版权许可协议  |// @license.name MIT
license.url |版权许可协议url  |// @license.url http://xxxx.com
host |提供API服务的主机IP或者域名  |// @host localhost:9527
BasePath |API公共路径  |// @BasePath /

安全验证
名称| 描述 | 参数| 例子
---|---|---|---
securitydefinitions.basic |基本认证  | |// @securityDefinitions.basic BasicAuth
securitydefinitions.apikey |token认证  |in, name |// @securitydefinitions.apikey

验证参数
名称| 描述 | 例子
---|---|---|---
in |位置  |// @in header
name |key  |// @name Authorization

2. API接口注解
名称| 描述 | 例子
---|---|---|---
description |描述  |// @description 接口描述
Tags |标签  |// @Tags
summary |摘要  |// @summary 摘要
Accept |接口接收数据类型  |// @Accept json
produce |接口产生数据类型  |// @produce json
Param |参数  |// @Param val1 query int true "used for calc"(可写多条)。参数详情请参考Param表格
Security |单独为一个接口设置验证规则  |// @Security ApiKeyAuth
success |请求成功的返回  |// @Success 200 {string} string "answer"
failure |请求失败的返回  |// @Failure 400 {string} string "ok" (可写多条)
router |接口地址  |// @Router /examples/attribute [get] (get是请求method)


3. Param参数表例子和说明
* `@Param`后面一般会有5个参数，后面还有一些可选参数，如枚举，默认值
    第一个表示参数的名称key，
    第二个表示参数类型(`query|path|header|body|formData`)
    第三个表示参数的值类型(`string|integer|number|boolean|struct`)
    第四个表示是否是必须(`true|false`)
    第五个是对Param的描述
```go
// 默认值
// @Param val1 query int false "参数1" default(1)

//枚举，且默认选中1.1
// @Param val3 query number false "请选择" Enums(1.1, 1.2, 1.3) default(1.1)

//控制参数值的长度
// @Param val4 query string false "最少5，最大10长度" minlength(5) maxlength(10)

//控制int的最大值和最小值
// @Param int query int false "int valid" minimum(1) maximum(10)

//结构体参数：loginReq是一个结构体
type loginReq struct {
	Email string `json:"email" validate:"required,email" example:"jhon@gmail.com"` //邮箱
	PassWord string `json:"password" validate:"required,gte=6,lte=12" example:"123456"` //密码
}
//example是结构体的默认值
// @Param message body loginReq true "登陆请求参数"

//以上可选参数也适用于结构体：
type Foo struct {
    Bar string `minLength:"4" maxLength:"16"`
    Baz int `minimum:"10" maximum:"20" default:"15"`
    Qux []string `enums:"foo,bar,baz"`
}
```

4. 结构体tag
tag| 描述 | 例子
---|---|---|---
swaggerignore |忽略字段  |`swaggerignore:"true"`
example |例子  |`example:123`
// |结构体字段描述  |//名称 ， 一般在结构体字段后面注释，在上面也可以
minimum |最小值  |`minimum:"10" `
maximum |最大值  |`maximum:"20" `
minlength |最小长度 |`minlength:"20" ` 
maxlength |最大长度  |`maxlength:"20" `
enums |枚举  |`enums:"foo,bar,baz"`
default |默认值  |`default:"15"`

5. success和failure
* // @Success 状态码 {返回数据类型} 返回数据类型对象 "描述信息" 
* success和failure用法一样，但是状态码不能一样。
```go

type Account struct {
    ID   int    `json:"id" example:"1"`
    Name string `json:"name" example:"account name"`
}
//array类型，相当于 `[]Account`
// @Success 200 {array} Account "请求成功数据"

type JSONResult struct {
    Code    int          `json:"code" `
    Message string       `json:"message"`
    Data    interface{}  `json:"data"`
}

type Order struct { //in `proto` package
    ...
}
//结构体类型，并用Order替换JSONResult中的data
// @success 200 {object} JSONResult{data=Order} "desc"

//替换多个字段的类型。如果某字段不存在，将添加该字段。
// @success 200 {object} JSONResult{data1=string,data2=[]string,data3=Order,data4=[]Order} "desc"
```

