package admin

import (
	"StudentServicePlatform/internal/apiException"
	"StudentServicePlatform/internal/service"
	"StudentServicePlatform/pkg/utils"

	"github.com/gin-gonic/gin"
)

type handleTrash struct {
	AdminID  int `json:"admin_id" binding:"required"`
	PostID   int `json:"post_id" binding:"required"`
	Approval int `json:"approval" binding:"required"`
}

func HandleTrash(c *gin.Context) {
	var data handleTrash
	err := c.ShouldBindJSON(&data)
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

	//检验反馈存在
	post, err := service.GetPostByID(data.PostID)
	if err != nil {
		_ = c.AbortWithError(200, apiException.PostNotFind)
		return
	}

	// 检验用户权限
	if user.UserType != 2 {
		_ = c.AbortWithError(200, apiException.LackRight)
		return
	}

	// 检查反馈状态
	if post.Status == 1 {
		_ = c.AbortWithError(200, apiException.ReatHandle)
		return
	} else if post.Status == 0 {
		_ = c.AbortWithError(200, apiException.PostNotHandle)
		return
	}

	// 处理
	err = service.HandleTrash(data.AdminID, data.PostID, data.Approval)
	if err != nil {
		_ = c.AbortWithError(200, apiException.HandleError)
		return
	}
	
	if data.Approval==1 {
		service.SendMail(user.Email, user.Name, "请您在提交反馈时确保内容的有效性和准确性，感谢您的理解和配合。如有异议，请重新反馈。")
	}
	utils.JsonSuccess(c, nil)
}
