package models

import (
	"schoolChat/database"
	"time"
)

// Reply 回复模型
type Reply struct {
	ID         int64  `json:"id"`         // 回复ID
	DeskType   int    `json:"deskType"`   // 目标类型 1:帖子 2:回复
	DeskId     int64  `json:"deskId"`     // 目标ID
	UserId     int64  `json:"userId"`     // 用户ID
	IsTop      int    `json:"isTop"`      // 是否置顶
	IsChoice   int    `json:"isChoice"`   // 是否精选
	Text       string `json:"text"`       // 回复内容
	Images     string `json:"images"`     // 回复图片
	CreateTime int64  `json:"createTime"` // 创建时间
	UpdateTime int64  `json:"updateTime"` // 更新时间
}

// ReplyLikes 回复点赞模型
type ReplyLikes struct {
	ID      int64 `json:"id"`      // 主键ID
	ReplyId int64 `json:"replyId"` // 回复ID
	UserId  int64 `json:"userId"`  // 用户ID
}

// ReplyCollects 回复收藏模型
type ReplyCollects struct {
	ID      int64 `json:"id"`      // 主键ID
	ReplyId int64 `json:"replyId"` // 回复ID
	UserId  int64 `json:"userId"`  // 用户ID
}

// ReplyResponse 回复响应模型
type ReplyResponse struct {
	ID            int64  `json:"id"`            // 回复ID
	DeskType      int    `json:"deskType"`      // 目标类型 1:帖子 2:回复
	DeskId        int64  `json:"deskId"`        // 目标ID
	UserId        int64  `json:"userId"`        // 用户ID
	Username      string `json:"username"`      // 用户名
	Nickname      string `json:"nickname"`      // 昵称
	Avatar        string `json:"avatar"`        // 头像
	IsTop         int    `json:"isTop"`         // 是否置顶
	IsChoice      int    `json:"isChoice"`      // 是否精选
	Text          string `json:"text"`          // 回复内容
	Images        string `json:"images"`        // 回复图片
	Likes         int    `json:"likes"`         // 点赞数
	Replies       int    `json:"replies"`       // 回复数
	Collects      int    `json:"collects"`      // 收藏数
	CollectStatus int    `json:"collectStatus"` // 收藏状态
	LikeStatus    int    `json:"likeStatus"`    // 点赞状态
	CreateTime    int64  `json:"createTime"`    // 创建时间
	UpdateTime    int64  `json:"updateTime"`    // 更新时间
}

// AddReply 添加回复
func AddReply(reply Reply) (replyData Reply, err error) {

	reply.IsTop = 0
	reply.IsChoice = 0
	reply.CreateTime = time.Now().UnixMilli()
	reply.UpdateTime = time.Now().UnixMilli()

	result := database.GetDB().Create(&reply)

	if result.Error != nil {
		err = result.Error
	}

	return reply, err
}

// GetReplyByPostId 根据帖子ID获取回复列表
func GetReplyByPostId(userId, postId int64) (replies []ReplyResponse, err error) {
	result := database.GetDB().
		Table("reply").
		Select("reply.*, "+
			"user.username, "+
			"user.nickname, "+
			"user.avatar, "+
			"count(distinct reply_likes.id) AS likes, "+
			"count(distinct reply_collects.id) AS collects, "+
			"(select count(reply.id) from reply where reply.desk_type = 2 AND reply.desk_id = reply.id) as replies,"+
			"(SELECT 1 FROM reply_likes WHERE reply.id = reply_likes.reply_id AND reply_likes.user_id = ? limit 1) as LikeStatus, "+
			"(SELECT 1 FROM reply_collects WHERE reply.id = reply_collects.reply_id AND reply_collects.user_id = ? limit 1) as collects", userId, userId).
		Joins("LEFT JOIN user ON user.id = reply.user_id").
		Joins("LEFT JOIN post ON post.id = reply.desk_id and reply.desk_type = 1").
		Joins("LEFT JOIN reply_likes ON reply_likes.reply_id = reply.id").
		Joins("LEFT JOIN reply_collects ON reply_collects.reply_id = reply.id").
		Where("reply.desk_type = 1 AND reply.desk_id = ?", postId).
		Group("reply.id").
		Order("reply.is_top DESC, reply.create_time DESC").
		Find(&replies)

	if result.Error != nil {
		err = result.Error
	}

	return replies, err
}

// GetReplyByReplyId 根据回复ID获取回复详情
func GetReplyByReplyId(userId, replyId int64) (replies []ReplyResponse, err error) {
	result := database.GetDB().
		Table("reply").
		Select("reply.*, "+
			"user.username, "+
			"user.nickname, "+
			"user.avatar, "+
			"count(distinct reply_likes.id) AS likes, "+
			"count(distinct reply_collects.id) AS collects, "+
			"(select count(reply.id) from reply where reply.desk_type = 2 AND reply.desk_id = reply.id) as replies,"+
			"(SELECT 1 FROM reply_likes WHERE reply.id = reply_likes.reply_id AND reply_likes.user_id = ? limit 1) as LikeStatus, "+
			"(SELECT 1 FROM reply_collects WHERE reply.id = reply_collects.reply_id AND reply_collects.user_id = ? limit 1) as collects", userId, userId).
		Joins("LEFT JOIN user ON user.id = reply.user_id").
		Joins("LEFT JOIN reply_likes ON reply_likes.reply_id = reply.id").
		Joins("LEFT JOIN reply_collects ON reply_collects.reply_id = reply.id").
		Where("reply.desk_type = 2 AND reply.desk_id = ?", replyId).
		Group("reply.id").
		Order("reply.is_top DESC, reply.create_time DESC").
		Find(&replies)

	if result.Error != nil {
		err = result.Error
	}

	return replies, err
}

// LikeReply 点赞回复
func LikeReply(userId, replyId int64) (err error) {
	result := database.GetDB().
		Table("reply_likes").
		Create(&ReplyLikes{
			ReplyId: replyId,
			UserId:  userId,
		})

	if result.Error != nil {
		err = result.Error
	}

	return err
}

// UnlikeReply 取消点赞回复
func UnlikeReply(userId, replyId int64) (err error) {
	result := database.GetDB().
		Table("reply_likes").
		Where("reply_id = ? AND user_id = ?", replyId, userId).
		Delete(&ReplyLikes{})

	if result.Error != nil {
		err = result.Error
	}

	return err
}
