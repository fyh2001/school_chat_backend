package models

import (
	"schoolChat/database"
	"time"
)

// Post 帖子模型
type Post struct {
	ID         int64  `json:"id"`         // 帖子ID
	UserId     int64  `json:"userId"`     // 用户ID
	IsTop      int    `json:"isTop"`      // 是否置顶
	IsChoice   int    `json:"isChoice"`   // 是否精选
	Text       string `json:"text"`       // 帖子内容
	Images     string `json:"images"`     // 帖子图片
	CreateTime int64  `json:"createTime"` // 创建时间
	UpdateTime int64  `json:"updateTime"` // 更新时间
}

// PostLikes 帖子点赞模型
type PostLikes struct {
	ID     int64 `json:"id"`     // 主键ID
	PostId int64 `json:"postId"` // 帖子ID
	UserId int64 `json:"userId"` // 用户ID
}

// PostCollects 帖子收藏模型
type PostCollects struct {
	ID     int64 `json:"id"`     // 主键ID
	PostId int64 `json:"postId"` // 帖子ID
	UserId int64 `json:"userId"` // 用户ID
}

// PostResponse 帖子响应模型
type PostResponse struct {
	ID            int64  `json:"id"`            // 帖子ID json
	UserId        int64  `json:"userId"`        // 用户ID
	email         string `json:"email"`         // 用户名
	Nickname      string `json:"nickname"`      // 昵称
	Avatar        string `json:"avatar"`        // 头像
	IsTop         int    `json:"isTop"`         // 是否置顶
	IsChoice      int    `json:"isChoice"`      // 是否精选
	Text          string `json:"text"`          // 帖子内容
	Images        string `json:"images"`        // 帖子图片
	Likes         int    `json:"likes"`         // 点赞数
	Replies       int    `json:"replies"`       // 回复数
	Collects      int    `json:"collects"`      // 收藏数
	CollectStatus int    `json:"collectStatus"` // 收藏状态
	LikeStatus    int    `json:"likeStatus"`    // 点赞状态
	CreateTime    int64  `json:"createTime"`    // 创建时间
	UpdateTime    int64  `json:"updateTime"`    // 更新时间
}

// AddPost @Title 新增帖子
func AddPost(post Post) (postData Post, err error) {

	post.IsTop = 0
	post.IsChoice = 0
	post.CreateTime = time.Now().UnixMilli()
	post.UpdateTime = time.Now().UnixMilli()

	result := database.GetDB().Create(&post)

	if result.Error != nil {
		err = result.Error
	}

	return post, err
}

// GetAllPost 查询所有帖子
func GetAllPost() (posts []PostResponse, err error) {
	result := database.GetDB().
		Select("post.id, " +
			"post.user_id, " +
			"user.email, " +
			"user.nickname, " +
			"user.avatar, " +
			"post.is_top, " +
			"post.is_choice, " +
			"post.text, " +
			"post.images, " +
			"post.create_time, " +
			"post.update_time, " +
			"count(distinct post_likes.id) as likes, " +
			"count(distinct reply.id) as replies, " +
			"count(distinct post_collects.id) as collects").
		Table("post").
		Joins("left join user on post.user_id = user.id").
		Joins("left join post_likes on post.id = post_likes.post_id").
		Joins("left join reply on post.id = reply.desk_id and reply.desk_type = 1").
		Joins("left join post_collects on post.id = post_collects.post_id").
		Group("post.id").
		Scan(&posts)

	if result.Error != nil {
		err = result.Error
	}
	return posts, err
}

// GetAllPostByLogin 查询所有帖子(已登陆)
func GetAllPostByLogin(userId int64) (posts []PostResponse, err error) {
	result := database.GetDB().
		Select("post.id, "+
			"post.user_id, "+
			"user.email, "+
			"user.nickname, "+
			"user.avatar, "+
			"post.is_top, "+
			"post.is_choice, "+
			"post.text, "+
			"post.images, "+
			"post.create_time, "+
			"post.update_time, "+
			"count(distinct post_likes.id) as likes, "+
			"(select count(reply.id) from reply where reply.desk_type = 1 AND reply.desk_id = post.id) + "+
			"(select count(reply.id) from reply where reply.desk_id in (select reply.id from reply where reply.desk_type = 1 AND reply.desk_id = post.id)) as replies, "+
			"count(distinct post_collects.id) as collects, "+
			//"if (post_likes.user_id = ?, 1, 0) AS likeStatus", id).
			"(SELECT 1 FROM post_likes WHERE post.id = post_likes.post_id AND post_likes.user_id = ? limit 1) as LikeStatus, "+
			"(SELECT 1 FROM post_collects WHERE post.id = post_collects.post_id AND post_collects.user_id = ? limit 1) as collects", userId, userId).
		Table("post").
		Joins("left join user on post.user_id = user.id").
		Joins("left join post_likes on post.id = post_likes.post_id").
		Joins("left join reply on post.id = reply.desk_id and reply.desk_type = 1").
		Joins("left join post_collects on post.id = post_collects.post_id").
		Group("post.id").
		Scan(&posts)

	if result.Error != nil {
		err = result.Error
	}
	return posts, err
}

// GetPostByPostId 根据帖子ID获取帖子详情
func GetPostByPostId(userId, postId int64) (post PostResponse, err error) {
	result := database.GetDB().
		Select("post.id, "+
			"post.user_id, "+
			"user.email, "+
			"user.nickname, "+
			"user.avatar, "+
			"post.is_top, "+
			"post.is_choice, "+
			"post.text, "+
			"post.images, "+
			"post.create_time, "+
			"post.update_time, "+
			"count(distinct post_likes.id) as likes, "+
			"(select count(reply.id) from reply where reply.desk_type = 1 AND reply.desk_id = post.id) + "+
			"(select count(reply.id) from reply where reply.desk_id in (select reply.id from reply where reply.desk_type = 1 AND reply.desk_id = post.id)) as replies, "+
			"count(distinct post_collects.id) as collects, "+
			"(SELECT 1 FROM post_likes WHERE post.id = post_likes.post_id AND post_likes.user_id = ? limit 1) as LikeStatus, "+
			"(SELECT 1 FROM post_collects WHERE post.id = post_collects.post_id AND post_collects.user_id = ? limit 1) as collects", userId, userId).
		Table("post").
		Joins("left join user on post.user_id = user.id").
		Joins("left join post_likes on post.id = post_likes.post_id").
		Joins("left join reply on post.id = reply.desk_id and reply.desk_type = 1").
		Joins("left join post_collects on post.id = post_collects.post_id").
		Where("post.id = ?", postId).
		Group("post.id").
		Scan(&post)

	if result.Error != nil {
		err = result.Error
	}
	return post, err
}

// GetPostByUserId 根据用户ID获取帖子列表
func GetPostByUserId(userId int64) (posts []PostResponse, err error) {
	result := database.GetDB().
		Select("post.id, "+
			"post.user_id, "+
			"user.email, "+
			"user.nickname, "+
			"user.avatar, "+
			"post.is_top, "+
			"post.is_choice, "+
			"post.text, "+
			"post.images, "+
			"post.create_time, "+
			"post.update_time, "+
			"count(distinct post_likes.id) as likes, "+
			"(select count(reply.id) from reply where reply.desk_type = 1 AND reply.desk_id = post.id) + "+
			"(select count(reply.id) from reply where reply.desk_id in (select reply.id from reply where reply.desk_type = 1 AND reply.desk_id = post.id)) as replies, "+
			"count(distinct post_collects.id) as collects, "+
			"(SELECT 1 FROM post_likes WHERE post.id = post_likes.post_id AND post_likes.user_id = ? limit 1) as LikeStatus, "+
			"(SELECT 1 FROM post_collects WHERE post.id = post_collects.post_id AND post_collects.user_id = ? limit 1) as collects", userId, userId).
		Table("post").
		Joins("left join user on post.user_id = user.id").
		Joins("left join post_likes on post.id = post_likes.post_id").
		Joins("left join reply on post.id = reply.desk_id and reply.desk_type = 1").
		Joins("left join post_collects on post.id = post_collects.post_id").
		Where("post.user_id = ?", userId).
		Group("post.id").
		Scan(&posts)

	if result.Error != nil {
		err = result.Error
	}
	return posts, err
}

// GetPostByLoginUserId 根据已登录用户ID获取帖子列表
func GetPostByLoginUserId(userId int64) (posts []PostResponse, err error) {
	result := database.GetDB().
		Select("post.id, "+
			"post.user_id, "+
			"user.email, "+
			"user.nickname, "+
			"user.avatar, "+
			"post.is_top, "+
			"post.is_choice, "+
			"post.text, "+
			"post.images, "+
			"post.create_time, "+
			"post.update_time, "+
			"count(distinct post_likes.id) as likes, "+
			"(select count(reply.id) from reply where reply.desk_type = 1 AND reply.desk_id = post.id) + "+
			"(select count(reply.id) from reply where reply.desk_id in (select reply.id from reply where reply.desk_type = 1 AND reply.desk_id = post.id)) as replies, "+
			"count(distinct post_collects.id) as collects, "+
			"(SELECT 1 FROM post_likes WHERE post.id = post_likes.post_id AND post_likes.user_id = ? limit 1) as LikeStatus, "+
			"(SELECT 1 FROM post_collects WHERE post.id = post_collects.post_id AND post_collects.user_id = ? limit 1) as collects", userId, userId).
		Table("post").
		Joins("left join user on post.user_id = user.id").
		Joins("left join post_likes on post.id = post_likes.post_id").
		Joins("left join reply on post.id = reply.desk_id and reply.desk_type = 1").
		Joins("left join post_collects on post.id = post_collects.post_id").
		Where("post.user_id = ?", userId).
		Group("post.id").
		Scan(&posts)

	if result.Error != nil {
		err = result.Error
	}
	return posts, err
}

// GetLikedPostsByUserId 根据已登录用户ID获取点赞的帖子列表
func GetLikedPostsByUserId(userId int64) (posts []PostResponse, err error) {
	result := database.GetDB().
		Select("post.id, "+
			"post.user_id, "+
			"user.email, "+
			"user.nickname, "+
			"user.avatar, "+
			"post.is_top, "+
			"post.is_choice, "+
			"post.text, "+
			"post.images, "+
			"post.create_time, "+
			"post.update_time, "+
			"count(distinct post_likes.id) as likes, "+
			"(select count(reply.id) from reply where reply.desk_type = 1 AND reply.desk_id = post.id) + "+
			"(select count(reply.id) from reply where reply.desk_id in (select reply.id from reply where reply.desk_type = 1 AND reply.desk_id = post.id)) as replies, "+
			"count(distinct post_collects.id) as collects, "+
			"(SELECT 1 FROM post_likes WHERE post.id = post_likes.post_id AND post_likes.user_id = ? limit 1) as LikeStatus, "+
			"(SELECT 1 FROM post_collects WHERE post.id = post_collects.post_id AND post_collects.user_id = ? limit 1) as collects", userId, userId).
		Table("post").
		Joins("left join user on post.user_id = user.id").
		Joins("left join post_likes on post.id = post_likes.post_id").
		Joins("left join reply on post.id = reply.desk_id and reply.desk_type = 1").
		Joins("left join post_collects on post.id = post_collects.post_id").
		Where("post_likes.user_id = ?", userId).
		Group("post.id").
		Scan(&posts)

	if result.Error != nil {
		err = result.Error
	}
	return posts, err
}

// GetRepliedPostsByUserId 根据已登录用户ID获取回复的帖子列表
func GetRepliedPostsByUserId(userId int64) (posts []PostResponse, err error) {
	result := database.GetDB().
		Select("post.id, "+
			"post.user_id, "+
			"user.email, "+
			"user.nickname, "+
			"user.avatar, "+
			"post.is_top, "+
			"post.is_choice, "+
			"post.text, "+
			"post.images, "+
			"post.create_time, "+
			"post.update_time, "+
			"count(distinct post_likes.id) as likes, "+
			"(select count(reply.id) from reply where reply.desk_type = 1 AND reply.desk_id = post.id) + "+
			"(select count(reply.id) from reply where reply.desk_id in (select reply.id from reply where reply.desk_type = 1 AND reply.desk_id = post.id)) as replies, "+
			"count(distinct post_collects.id) as collects, "+
			"(SELECT 1 FROM post_likes WHERE post.id = post_likes.post_id AND post_likes.user_id = ? limit 1) as LikeStatus, "+
			"(SELECT 1 FROM post_collects WHERE post.id = post_collects.post_id AND post_collects.user_id = ? limit 1) as collects", userId, userId).
		Table("post").
		Joins("left join user on post.user_id = user.id").
		Joins("left join post_likes on post.id = post_likes.post_id").
		Joins("left join reply on post.id = reply.desk_id and reply.desk_type = 1").
		Joins("left join post_collects on post.id = post_collects.post_id").
		Where("reply.user_id = ?", userId).
		Group("post.id").
		Scan(&posts)

	if result.Error != nil {
		err = result.Error
	}
	return posts, err
}

// LikePost 给帖子点赞
func LikePost(userId, postId int64) (err error) {
	like := PostLikes{
		PostId: postId,
		UserId: userId,
	}

	result := database.GetDB().Create(&like)

	if result.Error != nil {
		err = result.Error
	}
	return err
}

// UnlikePost 取消给帖子点赞
func UnlikePost(userId, postId int64) (err error) {
	result := database.GetDB().Where("user_id = ? AND post_id = ?", userId, postId).Delete(&PostLikes{})

	if result.Error != nil {
		err = result.Error
	}
	return err
}

// CollectPost 收藏帖子
func CollectPost(userId, postId int64) (err error) {
	collect := PostCollects{
		PostId: postId,
		UserId: userId,
	}

	result := database.GetDB().Create(&collect)

	if result.Error != nil {
		err = result.Error
	}
	return err
}

// UnCollectPost 取消收藏帖子
func UnCollectPost(userId, postId int64) (err error) {
	result := database.GetDB().Where("user_id = ? AND post_id = ?", userId, postId).Delete(&PostCollects{})

	if result.Error != nil {
		err = result.Error
	}
	return err
}
