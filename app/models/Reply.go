package models

import (
	"fmt"
	"gorm.io/gorm"
	"schoolChat/database"
	"time"
)

// Reply 回复模型
type Reply struct {
	ID           int64  `json:"id"`           // 回复ID
	DeskType     int    `json:"deskType"`     // 目标类型 1:帖子 2:回复 3:回复的回复
	PostID       int64  `json:"postId"`       // 所属帖子ID
	DeskId       int64  `json:"deskId"`       // 目标ID
	DeskSecondId int64  `json:"deskSecondId"` // 目标二级ID
	UserId       int64  `json:"userId"`       // 用户ID
	IsTop        int    `json:"isTop"`        // 是否置顶
	IsChoice     int    `json:"isChoice"`     // 是否精选
	Text         string `json:"text"`         // 回复内容
	Images       string `json:"images"`       // 回复图片
	Likes        int    `json:"likes"`        // 点赞数
	Replies      int    `json:"replies"`      // 回复数
	CreateTime   int64  `json:"createTime"`   // 创建时间
	UpdateTime   int64  `json:"updateTime"`   // 更新时间
}

// SecondReplies 二级回复模型
type SecondReplies struct {
	ID           int64  `json:"id"`           // 回复ID
	DeskType     int    `json:"deskType"`     // 目标类型 1:帖子 2:回复 3:回复的回复
	PostID       int64  `json:"postId"`       // 所属帖子ID
	DeskId       int64  `json:"deskId"`       // 目标ID
	DeskSecondId int64  `json:"deskSecondId"` // 目标二级ID
	UserId       int64  `json:"userId"`       // 用户ID
	Email        string `json:"email"`        // 邮箱
	Nickname     string `json:"nickname"`     // 昵称
	Avatar       string `json:"avatar"`       // 头像
	IsTop        int    `json:"isTop"`        // 是否置顶
	IsChoice     int    `json:"isChoice"`     // 是否精选
	Text         string `json:"text"`         // 回复内容
	Images       string `json:"images"`       // 回复图片
	Likes        int    `json:"likes"`        // 点赞数
	Replies      int    `json:"replies"`      // 回复数
	CreateTime   int64  `json:"createTime"`   // 创建时间
	UpdateTime   int64  `json:"updateTime"`   // 更新时间
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
	ID           int64  `json:"id"`           // 回复ID
	DeskType     int    `json:"deskType"`     // 目标类型 1:帖子 2:回复
	PostID       int64  `json:"postId"`       // 所属帖子ID
	DeskId       int64  `json:"deskId"`       // 目标ID
	DeskSecondId int64  `json:"deskSecondId"` // 目标二级ID
	UserId       int64  `json:"userId"`       // 用户ID
	Email        string `json:"email"`        // 邮箱
	Nickname     string `json:"nickname"`     // 昵称
	Avatar       string `json:"avatar"`       // 头像
	IsTop        int    `json:"isTop"`        // 是否置顶
	IsChoice     int    `json:"isChoice"`     // 是否精选
	Text         string `json:"text"`         // 回复内容
	Images       string `json:"images"`       // 回复图片
	Likes        int    `json:"likes"`        // 点赞数
	Replies      int    `json:"replies"`      // 回复数
	LikeStatus   int    `json:"likeStatus"`   // 点赞状态
	CreateTime   int64  `json:"createTime"`   // 创建时间
	UpdateTime   int64  `json:"updateTime"`   // 更新时间
}

// AddReply 添加回复
func AddReply(reply Reply) (replyData Reply, err error) {

	reply.IsTop = 0
	reply.IsChoice = 0
	reply.CreateTime = time.Now().UnixMilli()
	reply.UpdateTime = time.Now().UnixMilli()

	result := database.GetMySQL().Create(&reply)

	if result.Error != nil {
		err = result.Error
	}

	if reply.DeskType == 1 {
		// 更新帖子回复数
		database.GetMySQL().Table("post").Where("id = ?", reply.DeskId).Update("replies", gorm.Expr("replies + ?", 1))
	}

	if reply.DeskType == 2 {
		// 更新回复回复数
		database.GetMySQL().Table("reply").Where("id = ?", reply.DeskId).Update("replies", gorm.Expr("replies + ?", 1))
		// 更新帖子回复数
		database.GetMySQL().Table("post").Where("id = ?", reply.PostID).Update("replies", gorm.Expr("replies + ?", 1))
	}

	if reply.DeskType == 3 {
		// 更新回复回复数
		database.GetMySQL().Table("reply").Where("id = ?", reply.DeskId).Update("replies", gorm.Expr("replies + ?", 1))
		// 更新回复回复数
		database.GetMySQL().Table("reply").Where("id = ?", reply.DeskSecondId).Update("replies", gorm.Expr("replies + ?", 1))
		// 更新帖子回复数
		database.GetMySQL().Table("post").Where("id = ?", reply.PostID).Update("replies", gorm.Expr("replies + ?", 1))
	}

	return reply, err
}

// GetReplyByPostId 根据帖子ID获取回复列表
func GetReplyByPostId(userId, postId int64) (replies []ReplyResponse, err error) {
	// 创建查询语句
	sql := `
			SELECT
				reply.id,
				reply.desk_type,
				reply.user_id,
				reply.desk_id,
				reply.desk_second_id,
				reply.user_id,
				user.email,
				user.nickname,
				user.avatar,
		        reply.is_top,
		        reply.is_choice,
				reply.text,
		        reply.images,
				reply.likes,
		        reply.replies,
		        reply.create_time,
		        reply.update_time,
		        CASE WHEN rl.user_id IS NOT NULL THEN 1 ELSE 0 END AS like_status
			FROM reply
			LEFT JOIN user ON user.id = reply.user_id
			LEFT JOIN reply_likes rl on reply.id = rl.reply_id AND rl.user_id = ?
			WHERE reply.desk_type = 1 AND reply.desk_id = ?
	`

	// 执行查询语句
	result := database.GetMySQL().Raw(sql, userId, postId).Scan(&replies)

	fmt.Printf("replies: %+v\n", replies)

	if result.Error != nil {
		err = result.Error
	}

	return replies, err
}

// GetReplyByReplyId 根据回复ID获取回复详情
func GetReplyByReplyId(userId, replyId int64) (reply ReplyResponse, err error) {
	// 创建查询语句
	sql := `
			SELECT
				reply.id,
				reply.desk_type,
				reply.user_id,
				reply.desk_id,
				reply.desk_second_id,
				reply.user_id,
				user.email,
				user.nickname,
				user.avatar,
		        reply.is_top,
		        reply.is_choice,
				reply.text,
		        reply.images,
				reply.likes,
		        reply.replies,
		        reply.create_time,
		        reply.update_time,
		        CASE WHEN rl.user_id IS NOT NULL THEN 1 ELSE 0 END AS like_status
			FROM reply
			LEFT JOIN user ON user.id = reply.user_id
			LEFT JOIN reply_likes rl on reply.id = rl.reply_id AND rl.user_id = ?
			WHERE reply.id = ?
	`

	// 执行查询语句
	result := database.GetMySQL().Raw(sql, userId, replyId).Scan(&reply)

	if result.Error != nil {
		err = result.Error
	}

	return reply, err
}

// GetSecondReplyByReplyId 根据回复ID获取二级回复
func GetSecondReplyByReplyId(userId, replyId int64) (replies []SecondReplies, err error) {

	// 创建查询语句
	sql := `
			SELECT
				reply.id,
				reply.desk_type,
				reply.user_id,
				reply.desk_id,
				reply.desk_second_id,
				reply.user_id,
				user.email,
				user.nickname,
				user.avatar,
		        reply.is_top,
		        reply.is_choice,
				reply.text,
		        reply.images,
				reply.likes,
		        reply.replies,
		        reply.create_time,
		        reply.update_time,
		        CASE WHEN rl.user_id IS NOT NULL THEN 1 ELSE 0 END AS like_status
			FROM reply
			LEFT JOIN user ON user.id = reply.user_id
			LEFT JOIN reply_likes rl on reply.id = rl.reply_id AND rl.user_id = ?
			WHERE (reply.desk_type = 3 OR reply.desk_type = 2 ) AND reply.desk_id = ?
	`

	// 执行查询语句
	result := database.GetMySQL().Raw(sql, userId, replyId).Scan(&replies)

	if result.Error != nil {
		err = result.Error
	}

	return replies, err
}

// LikeReply 点赞回复
func LikeReply(userId, replyId int64) (err error) {
	// 生成点赞记录
	result := database.GetMySQL().
		Table("reply_likes").
		Create(&ReplyLikes{
			ReplyId: replyId,
			UserId:  userId,
		})

	if result.Error != nil {
		err = result.Error
	}

	// 更新回复点赞数
	database.GetMySQL().Table("reply").Where("id = ?", replyId).
		Update("likes", gorm.Expr("likes + ?", 1))

	// 更新用户点赞数
	database.GetMySQL().Table("user").Where("id = (?)", database.GetMySQL().Select("user_id").Where("id = ?", replyId).Table("reply")).
		Update("likes", gorm.Expr("likes + ?", 1))
	return err
}

// UnlikeReply 取消点赞回复
func UnlikeReply(userId, replyId int64) (err error) {
	result := database.GetMySQL().
		Table("reply_likes").
		Where("reply_id = ? AND user_id = ?", replyId, userId).
		Delete(&ReplyLikes{})

	if result.Error != nil {
		err = result.Error
	}

	// 更新回复点赞数
	database.GetMySQL().Table("reply").Where("id = ?", replyId).
		Update("likes", gorm.Expr("likes - ?", 1))

	// 更新用户点赞数
	database.GetMySQL().Table("user").Where("id = (?)", database.GetMySQL().Select("user_id").Where("id = ?", replyId).Table("reply")).
		Update("likes", gorm.Expr("likes - ?", 1))

	return err
}

// DeleteReplyByPostId 根据帖子ID删除所有回复
func DeleteReplyByPostId(postId int64) (err error) {
	// 删除回复
	result := database.GetMySQL().Exec("DELETE FROM reply WHERE post_id = ?", postId)

	if result.Error != nil {
		err = result.Error
	}

	return err
}
