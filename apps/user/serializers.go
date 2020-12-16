package user

import (
	. "apiserver/utils"
	"errors"
)

type RegisterRep struct {
	Name     string `json:"name" validate:"required" example:"jhon"` // 姓名
	Email    string `json:"email" validate:"required,email" example:"jhon@gmail.com"` //邮箱
	Phone   string `json:"phone" example:"18817281823"` //手机号
	PassWord1 string `json:"password1" validate:"gte=6,lte=12" example:"123456"` //密码
	PassWord2 string `json:"password2" validate:"eqfield=PassWord1" example:"123456" ` //密码
	isOk bool
}

func (r *RegisterRep) isValid() error {
	if hasUser(r.Email){
		return errors.New("邮箱已经注册")
	}
	r.isOk = true
	return nil
}

func (r *RegisterRep) create() error {
	if !r.isOk{
		return errors.New("you must call isValid() before calling create()")
	}
	r.PassWord1 = MD5V([]byte(r.PassWord1))
	user:= User{Name: r.Name, Email: r.Email, Phone: r.Phone, PassWord: r.PassWord1}
	err := createUser(&user)
	if err != nil {
		return err
	}
	return nil
}

//loginReq 登陆请求参数
type LoginReq struct {
	Email string `json:"email" validate:"required,email" example:"wt@com.com"` //邮箱
	PassWord string `json:"password" validate:"required,gte=6,lte=12" example:"123456"` //密码
}

func (l *LoginReq) login()(*LoginResp, error) {
	user:=new(LoginResp)
	return user,userLogin(user,l.Email,MD5V([]byte(l.PassWord)))
}

type LoginResp struct {
	Name string
	Email string
	Phone string
	Authorization string
}

type userUpdateReq struct {
	Name     *string `json:"name" example:"jhon" db:"name"` //姓名
	Phone   *string `json:"phone" example:"18817281823" db:"phone"` //手机号
}

type idUri struct {
	ID   int `uri:"id" validate:"required" db:"id"`
}

type userSerializer struct {
	ID int
	Name string
	Email string
	Phone string
}