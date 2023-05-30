package controllers

import (
	"github.com/gin-gonic/gin"
	"schoolChat/app/models"
	Results "schoolChat/app/result"
	"strconv"
)

// AddReply 新增回复
func AddReply(c *gin.Context) {
	var reply models.Reply
	err := c.BindJSON(&reply)

	if err != nil {
		c.JSON(200, Results.ErrBind.Fail("获取参数失败"))
		return
	}

	userId, _ := c.Get("id")

	reply.UserId = userId.(int64)

	reply, err = models.AddReply(reply)

	if err != nil {
		c.JSON(200, Results.Err.Fail("新增失败"))
		return
	}

	c.JSON(200, Results.Ok.Success(reply))
}

// GetReplyByPostId 根据帖子ID获取回复列表
func GetReplyByPostId(c *gin.Context) {
	// 获取参数
	postId, _ := strconv.ParseInt(c.Query("postId"), 10, 64)
	// 获取用户id
	userId, _ := c.Get("id")

	result, err := models.GetReplyByPostId(userId.(int64), postId)

	if err != nil {
		c.JSON(200, Results.Err.Fail("查询失败"))
		return
	}

	c.JSON(200, Results.Ok.Success(result))
}

// GetReplyByReplyId 根据回复ID获取回复详情
func GetReplyByReplyId(c *gin.Context) {
	// 获取参数
	replyId, _ := strconv.ParseInt(c.Query("replyId"), 10, 64)
	// 获取用户id
	userId, _ := c.Get("id")

	result, err := models.GetReplyByReplyId(userId.(int64), replyId)

	if err != nil {
		c.JSON(200, Results.Err.Fail("查询失败"))
		return
	}

	c.JSON(200, Results.Ok.Success(result))
}

// LikeReply 点赞回复
func LikeReply(c *gin.Context) {
	// 获取参数
	var replyId int64
	err := c.BindJSON(&replyId)

	// 获取用户id
	userId, _ := c.Get("id")

	err = models.LikeReply(userId.(int64), replyId)

	if err != nil {
		c.JSON(200, Results.Err.Fail("点赞失败"))
		return
	}

	c.JSON(200, Results.Ok.Success("点赞成功"))
}

// UnlikeReply 取消点赞回复
func UnlikeReply(c *gin.Context) {
	// 获取参数
	var replyId int64
	err := c.BindJSON(&replyId)

	// 获取用户id
	userId, _ := c.Get("id")

	err = models.UnlikeReply(userId.(int64), replyId)

	if err != nil {
		c.JSON(200, Results.Err.Fail("取消点赞失败"))
		return
	}

	c.JSON(200, Results.Ok.Success("取消点赞成功"))
}
