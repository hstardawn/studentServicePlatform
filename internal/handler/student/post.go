package student

import (
	"StudentServicePlatform/internal/apiException"
	"StudentServicePlatform/internal/model"
	"StudentServicePlatform/internal/service"
	"StudentServicePlatform/pkg/utils"

	"github.com/gin-gonic/gin"
)

type CreatePostData struct {
	UserID      int    `json:"user_id" binding:"required"`
	IsAnonymous int    `json:"is_anonymous" binding:"required"`
	IsUrgent    int    `json:"is_urgent" binding:"required"`
	PostType    int    `json:"post_type" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Content     string `json:"content" binding:"required"`
}

func CreatePost(c *gin.Context) {
	var data CreatePostData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		_ = c.AbortWithError(200, apiException.ParamError) //参数错误
		return
	}
	user, err := service.GetUserByUserID(data.UserID)
	if err != nil {
		_ = c.AbortWithError(200, apiException.UserNotFind) //用户不存在
		return
	}
	if data.PostType != 1 && data.PostType != 2 && data.PostType != 3 && data.PostType != 4 {
		_ = c.AbortWithError(200, apiException.PostTypeError) //反馈类型无效
		return
	}
	err = service.CreatePost(data.UserID, data.IsAnonymous, data.IsUrgent, data.PostType, data.Title, data.Content)
	if err != nil {
		_ = c.AbortWithError(200, apiException.CreatePostError) //提交反馈失败
		return
	}
	
	service.SendMail(user.Email, user.Name, "请您在提交反馈时确保内容的有效性和准确性，感谢您的理解和配合。如有异议，请重新反馈。")
	utils.JsonSuccess(c, nil)
}

type UpdatePostData struct {
	UserID      int    `json:"user_id"`
	ID          int    `json:"post_id"`
	IsAnonymous int    `json:"is_anonymous"`
	IsUrgent    int    `json:"is_urgent"`
	PostType    int    `json:"post_type"`
	Title       string `json:"title"`
	Content     string `json:"content"`
}

func UpdatePost(c *gin.Context) {
	var data UpdatePostData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		_ = c.AbortWithError(200, apiException.ParamError) //参数错误
		return
	}
	_, err = service.GetUserByUserID(data.UserID)
	if err != nil {
		_ = c.AbortWithError(200, apiException.UserNotFind) //用户不存在
		return
	}
	post, err := service.GetPostByID(data.ID)
	if err != nil {
		_ = c.AbortWithError(200, apiException.PostNotFind) //反馈不存在
		return
	}
	if post.UserID != data.UserID {
		_ = c.AbortWithError(200, apiException.UserConnotUpdatePost) //无权修改帖子
		return
	}
	if data.PostType != 1 && data.PostType != 2 && data.PostType != 3 && data.PostType != 4 {
		_ = c.AbortWithError(200, apiException.PostTypeError) //反馈类型无效
		return
	}
	err = service.UpdatePost(data.UserID, data.ID, data.IsAnonymous, data.IsUrgent, data.PostType, data.Title, data.Content)
	if err != nil {
		_ = c.AbortWithError(200, apiException.UpdatePostError) //修改反馈失败
		return
	}
	utils.JsonSuccess(c, nil)
}

type DeletePostData struct {
	UserID int `json:"user_id" binding:"required"`
	ID     int `json:"post_id" binding:"required"`
}

func DeletePost(c *gin.Context) {
	var data DeletePostData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		_ = c.AbortWithError(200, apiException.ParamError) //参数错误
		return
	}
	_, err = service.GetUserByUserID(data.UserID)
	if err != nil {
		_ = c.AbortWithError(200, apiException.UserNotFind) //用户不存在
		return
	}
	post, err := service.GetPostByID(data.ID)
	if err != nil {
		_ = c.AbortWithError(200, apiException.PostNotFind) //反馈不存在
		return
	}
	if post.UserID != data.UserID {
		_ = c.AbortWithError(200, apiException.UserConnotDeletePost) //无权删除帖子
		return
	}
	err = service.DeletePost(data.UserID, data.ID)
	if err != nil {
		_ = c.AbortWithError(200, apiException.DeletePostError) //删除反馈失败
		return
	}
	utils.JsonSuccess(c, nil)
}

func GetPostList(c *gin.Context) {
	var postList []model.Post
	postList, err := service.GetPostList()
	if err != nil {
		_ = c.AbortWithError(200, apiException.GetPostListError) //获取反馈列表失败
		return
	}
	utils.JsonSuccess(c, gin.H{
		"post_list": postList,
	})
}

