package user

import (
	. "apiserver/utils"
	"errors"
)

type registerSerializer struct {
	Name     string `json:"name" validate:"required" example:"jhon"` // 姓名
	Email    string `json:"email" validate:"required,email" example:"jhon@gmail.com"` //邮箱
	Phone   string `json:"phone" example:"18817281823"` //手机号
	PassWord1 string `json:"password1" validate:"gte=6,lte=12" example:"123456"` //密码
	PassWord2 string `json:"password2" validate:"eqfield=PassWord1" example:"123456"` //密码
	isOk bool
}

func (r *registerSerializer) isValid() error {
	if hasUser(r.Email){
		return errors.New("邮箱已经注册")
	}
	r.isOk = true
	return nil
}

func (r *registerSerializer) create()(err error) {
	if !r.isOk{
		return errors.New("you must call isValid() before calling create()")
	}
	r.PassWord1 = MD5V([]byte(r.PassWord1))
	user:= User{Name: r.Name, Email: r.Email, Phone: r.Phone, PassWord: r.PassWord1}
	err = createUser(&user)
	if err != nil {
		return err
	}
	return nil
}