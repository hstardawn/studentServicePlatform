package user

import (
	"StudentServicePlatform/internal/apiException"
	"StudentServicePlatform/internal/service"
	"StudentServicePlatform/pkg/utils"

	"github.com/gin-gonic/gin"
)

type SendCodeData struct {
	Email string `json:"email"`
}

func SendCode(c *gin.Context) {
	var data SendCodeData
	if err := c.ShouldBindJSON(&data); err != nil {
		_ = c.AbortWithError(200, apiException.ParamError)
		return
	}

	code := service.GenerateVerificationCode()
	if err := service.SendVerifyCode( data.Email, code); err != nil {
		_ = c.AbortWithError(200, apiException.SendError)
	    return
    }
	utils.JsonSuccess(c, nil)
}

type RegisterData struct{
	Username int    `json:"username" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Code string `json:"code" binding:"required"`
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
		_ = c.AbortWithError(200, apiException.UserAlreadyExisted) //用户已存在
		return
	}
	if data.Username<1e11||data.Username>1e12{
		_ = c.AbortWithError(200, apiException.UsernameError) //用户名长度必须为12位
		return
	}
	if len(data.Password)<8||len(data.Password)>16{
		_ = c.AbortWithError(200, apiException.PasswordError) //密码长度必须大于8且小于16位
		return
	}
	// if data.UserType!=1&&data.UserType!=2&&data.UserType!=3{
	// 	_ = c.AbortWithError(200, apiException.UserTypeError) //用户类型无效
	// 	return
	// }
	// 校验验证码
	if !service.VerifyVerificationCode(data.Email, data.Code) {
		_ = c.AbortWithError(200, apiException.VartiyError)
		return
	}

	// 加密密码
	hashpassword , err:= utils.HashPassword(data.Password)
	if err != nil{
		_ = c.AbortWithError(200, apiException.EncryptionFailed)
		return
	}

	err=service.Register(data.Username,data.Name,data.Email,hashpassword)
	if err!=nil{
		_ = c.AbortWithError(200, apiException.Register) //注册失败
		return
	}
	utils.JsonSuccess(c,nil)
}