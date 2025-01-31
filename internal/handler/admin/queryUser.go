package admin

import (
	"StudentServicePlatform/internal/apiException"
	"StudentServicePlatform/internal/service"
	"StudentServicePlatform/pkg/utils"

	"github.com/gin-gonic/gin"
)

type queryUser struct {
	AdminID int `form:"admin_id" binding:"required"`
}

type GetUser struct {
	ID       int    `json:"user_id"`
	Username int    `json:"username"`
	Name     string `json:"name"`
	Sex      string `json:"sex"`
	PhoneNum int    `json:"phone_num"`
	Email    string `json:"email"`
	UserType int    `json:"user_type"`
}

func QueryAdmin(c *gin.Context) {
	var data queryUser
	err := c.ShouldBindQuery(&data)
	if err != nil {
		_ = c.AbortWithError(200, apiException.ParamError)
		return
	}
	// 检验用户存在
	user, err := service.GetUserByUserID(data.AdminID)
	if err != nil {
		_ = c.AbortWithError(200, apiException.AdminNotFind)
		return
	}

	// 检验用户权限
	if user.UserType != 2 {
		_ = c.AbortWithError(200, apiException.LackRight)
		return
	}

	// 获取管理员
	userList, err := service.QueryAdmin()
	if err != nil {
		_ = c.AbortWithError(200, apiException.GetAdminListError)
		return
	}
	var user_list []GetUser
	for _, admin := range userList {
		// 2.获取帖子内容
		user, err := service.GetUserByUserID(admin.ID)
		if err != nil {
			_ = c.AbortWithError(200, apiException.GetUserError)
			return
		}
		
		// 3.返回帖子内容
		user_list = append(user_list, GetUser{
			ID:       user.ID,
			Username: user.Username,
			Name:     user.Name,
			Sex:      user.Sex,
			PhoneNum: user.PhoneNum,
			Email:    user.Email,
			UserType: user.UserType,
		})
	}

	utils.JsonSuccess(c, gin.H{
		"user_list": user_list,
	})
}
