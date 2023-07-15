package controllers

import (
	"github.com/gin-gonic/gin"
	"schoolChat/app/models"
	Results "schoolChat/app/result"
	"strconv"
)

// ReplyResponse 回复响应模型
type ReplyResponse struct {
	ID            int64                  `json:"id"`            // 回复ID
	DeskType      int                    `json:"deskType"`      // 目标类型 1:帖子 2:回复
	PostID        int64                  `json:"postId"`        // 所属帖子ID
	DeskId        int64                  `json:"deskId"`        // 目标ID
	DeskSecondId  int64                  `json:"deskSecondId"`  // 目标二级ID
	UserId        int64                  `json:"userId"`        // 用户ID
	Email         string                 `json:"email"`         // 邮箱
	Nickname      string                 `json:"nickname"`      // 昵称
	Avatar        string                 `json:"avatar"`        // 头像
	IsTop         int                    `json:"isTop"`         // 是否置顶
	IsChoice      int                    `json:"isChoice"`      // 是否精选
	Text          string                 `json:"text"`          // 回复内容
	Images        string                 `json:"images"`        // 回复图片
	SecondReplies []models.SecondReplies `json:"secondReplies"` // 二级回复
	Likes         int                    `json:"likes"`         // 点赞数
	Replies       int                    `json:"replies"`       // 回复数
	LikeStatus    int                    `json:"likeStatus"`    // 点赞状态
	CreateTime    int64                  `json:"createTime"`    // 创建时间
	UpdateTime    int64                  `json:"updateTime"`    // 更新时间
}

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

	var replyResponse []ReplyResponse

	for _, reply := range result {
		// 获取二级回复
		secondReplies, _ := models.GetSecondReplyByReplyId(userId.(int64), reply.ID)

		replyResponse = append(replyResponse, ReplyResponse{
			ID:            reply.ID,
			DeskType:      reply.DeskType,
			PostID:        reply.PostID,
			DeskId:        reply.DeskId,
			DeskSecondId:  reply.DeskSecondId,
			UserId:        reply.UserId,
			Email:         reply.Email,
			Nickname:      reply.Nickname,
			Avatar:        reply.Avatar,
			IsTop:         reply.IsTop,
			IsChoice:      reply.IsChoice,
			Text:          reply.Text,
			Images:        reply.Images,
			SecondReplies: secondReplies,
			Likes:         reply.Likes,
			Replies:       reply.Replies,
			LikeStatus:    reply.LikeStatus,
			CreateTime:    reply.CreateTime,
			UpdateTime:    reply.UpdateTime,
		})
	}

	if err != nil {
		c.JSON(200, Results.Err.Fail("查询失败"))
		return
	}

	c.JSON(200, Results.Ok.Success(replyResponse))
}

// GetReplyByReplyId 根据回复ID获取回复详情
func GetReplyByReplyId(c *gin.Context) {
	// 获取参数
	replyId, _ := strconv.ParseInt(c.Query("replyId"), 10, 64)
	// 获取用户id
	userId, _ := c.Get("id")

	result, err := models.GetReplyByReplyId(userId.(int64), replyId)

	secondReplies, _ := models.GetSecondReplyByReplyId(userId.(int64), result.ID)

	replyResponse := ReplyResponse{
		ID:            result.ID,
		DeskType:      result.DeskType,
		PostID:        result.PostID,
		DeskId:        result.DeskId,
		DeskSecondId:  result.DeskSecondId,
		UserId:        result.UserId,
		Email:         result.Email,
		Nickname:      result.Nickname,
		Avatar:        result.Avatar,
		IsTop:         result.IsTop,
		IsChoice:      result.IsChoice,
		Text:          result.Text,
		Images:        result.Images,
		SecondReplies: secondReplies,
		Likes:         result.Likes,
		Replies:       result.Replies,
		LikeStatus:    result.LikeStatus,
		CreateTime:    result.CreateTime,
		UpdateTime:    result.UpdateTime,
	}

	if err != nil {
		c.JSON(200, Results.Err.Fail("查询失败"))
		return
	}

	c.JSON(200, Results.Ok.Success(replyResponse))
}

// GetSecondReplyByReplyId 根据回复ID获取二级回复列表
func GetSecondReplyByReplyId(c *gin.Context) {
	// 获取参数
	replyId, _ := strconv.ParseInt(c.Query("replyId"), 10, 64)
	// 获取用户id
	userId, _ := c.Get("id")

	result, err := models.GetSecondReplyByReplyId(userId.(int64), replyId)

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
