package user

import (
	// "CMS/internal/service"
	// "CMS/pkg/utils"

	"StudentServicePlatform/internal/apiException"
	"StudentServicePlatform/internal/service"
	"StudentServicePlatform/pkg/utils"

	"github.com/gin-gonic/gin"
)

type RegisterData struct{
	Username int    `json:"username" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Sex      string `json:"sex"`
	PhoneNum int    `json:"phone_num" binding:"required"`
	Email    string `json:"email" binding:"required"`
	UserType int    `json:"user_type" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context){
	var data RegisterData
	err:=c.ShouldBindJSON(&data)
    if err!=nil{
		_ = c.AbortWithError(200, apiException.ParamError) //参数错误
		return
	}
	_,err=service.GetUserByUsername(data.Username)
	if err==nil{
		_ = c.AbortWithError(200, apiException.UserAlreadyExisted) //用户名已存在
		return
	}
	// wrong:=service.IsNumericUsername(data.Username)
	// if !wrong{
	// 	utils.JsonFail(c,200502,"用户名只能包含数字")
	// 	return
	// }
	if len(data.Password)<8||len(data.Password)>16{
		_ = c.AbortWithError(200, apiException.PasswordError) //密码长度必须大于8且小于16位
		return
	}
	if data.UserType!=1&&data.UserType!=2&&data.UserType!=3{
		_ = c.AbortWithError(200, apiException.UserTypeError) //用户类型无效
		return
	}
	err=service.Register(data.Username,data.Name,data.Sex,data.PhoneNum,data.Email,data.Password,data.UserType)
	if err!=nil{
		_ = c.AbortWithError(200, apiException.Register) //注册失败
		return
	}
	utils.JsonSuccess(c,nil)
}