package models

import (
	"gorm.io/gorm"
	"schoolChat/database"
	"strings"
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
	Likes      int    `json:"likes"`      // 点赞数
	Replies    int    `json:"replies"`    // 回复数
	Collects   int    `json:"collects"`   // 收藏数
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
	Email         string `json:"email"`         // 用户名
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
	post.Likes = 0
	post.Replies = 0
	post.Collects = 0
	post.CreateTime = time.Now().UnixMilli()
	post.UpdateTime = time.Now().UnixMilli()

	result := database.GetMySQL().Create(&post)

	if result.Error != nil {
		err = result.Error
	}

	return post, err
}

// GetAllPost 查询所有帖子
func GetAllPost() (posts []PostResponse, err error) {
	// 创建查询语句
	sql := `
			SELECT
				post.id,
				post.user_id,
				user.email,
				user.nickname,
				user.avatar,
				post.is_top,
				post.is_choice,
				post.text,
				post.images,
				post.likes,
				post.replies,
				post.collects,
				post.create_time,
				post.update_time
			FROM post
			LEFT JOIN user ON post.user_id = user.id
	`
	// 执行查询语句
	result := database.GetMySQL().Raw(sql).Scan(&posts)

	if result.Error != nil {
		err = result.Error
	}
	return posts, err
}

// GetAllPostByLogin 查询所有帖子(已登陆)
func GetAllPostByLogin(userId int64) (posts []PostResponse, err error) {
	// 创建查询语句
	sql := `
			SELECT
				post.id,
				post.user_id,
				user.email,
				user.nickname,
				user.avatar,
				post.is_top,
				post.is_choice,
				post.text,
				post.images,
				post.likes,
				post.replies,
				post.collects,
				post.create_time,
				post.update_time,
				CASE WHEN pl.user_id IS NOT NULL THEN 1 ELSE 0 END AS like_status,
				CASE WHEN pc.user_id IS NOT NULL THEN 1 ELSE 0 END AS CollectStatus
	-- 	        (SELECT 1 FROM post_likes WHERE post.id = post_likes.post_id AND post_likes.user_id = ? limit 1) as LikeStatus,
	-- 	        (SELECT 1 FROM post_collects WHERE post.id = post_collects.post_id AND post_collects.user_id = ? limit 1) as CollectStatus
	
			FROM post
			LEFT JOIN user ON post.user_id = user.id
			LEFT JOIN post_likes pl ON post.id = pl.post_id AND pl.user_id = ?
			LEFT JOIN post_collects pc ON post.id = pc.post_id AND pc.user_id = ?
	`
	// 执行查询语句
	result := database.GetMySQL().Raw(sql, userId, userId).Scan(&posts)

	if result.Error != nil {
		err = result.Error
	}

	return posts, err
}

// GetPostByPostId 根据帖子ID获取帖子详情
func GetPostByPostId(userId, postId int64) (post PostResponse, err error) {
	// 创建查询语句
	sql := `
			SELECT
				post.id,
				post.user_id,
				user.email,
				user.nickname,
				user.avatar,
				post.is_top,
				post.is_choice,
				post.text,
				post.images,
				post.likes,
				post.replies,
				post.collects,
				post.create_time,
				post.update_time,
				CASE WHEN pl.user_id IS NOT NULL THEN 1 ELSE 0 END AS like_status,
				CASE WHEN pc.user_id IS NOT NULL THEN 1 ELSE 0 END AS CollectStatus
	-- 	        (SELECT 1 FROM post_likes WHERE post.id = post_likes.post_id AND post_likes.user_id = ? limit 1) as LikeStatus,
	-- 	        (SELECT 1 FROM post_collects WHERE post.id = post_collects.post_id AND post_collects.user_id = ? limit 1) as CollectStatus
	
			FROM post
			LEFT JOIN user ON post.user_id = user.id
			LEFT JOIN post_likes pl ON post.id = pl.post_id AND pl.user_id = ?
			LEFT JOIN post_collects pc ON post.id = pc.post_id AND pc.user_id = ?
			WHERE post.id = ?
	`
	// 执行查询语句
	result := database.GetMySQL().Raw(sql, userId, userId, postId).Scan(&post)

	if result.Error != nil {
		err = result.Error
	}
	return post, err
}

// GetPostByUserId 根据用户ID获取帖子列表
func GetPostByUserId(userId int64) (posts []PostResponse, err error) {
	// 创建查询语句
	sql := `
			SELECT
				post.id,
				post.user_id,
				user.email,
				user.nickname,
				user.avatar,
				post.is_top,
				post.is_choice,
				post.text,
				post.images,
				post.likes,
				post.replies,
				post.collects,
				post.create_time,
				post.update_time,
				CASE WHEN pl.user_id IS NOT NULL THEN 1 ELSE 0 END AS like_status,
				CASE WHEN pc.user_id IS NOT NULL THEN 1 ELSE 0 END AS CollectStatus
			FROM post
			LEFT JOIN user ON post.user_id = user.id
			LEFT JOIN post_likes pl ON post.id = pl.post_id AND pl.user_id = ?
			LEFT JOIN post_collects pc ON post.id = pc.post_id AND pc.user_id = ?
			WHERE post.user_id = ?
	`
	// 执行查询语句
	result := database.GetMySQL().Raw(sql, userId, userId, userId).Scan(&posts)

	if result.Error != nil {
		err = result.Error
	}
	return posts, err
}

// GetPostByLoginUserId 根据已登录用户ID获取帖子列表
func GetPostByLoginUserId(userId int64) (posts []PostResponse, err error) {
	// 创建查询语句
	sql := `
			SELECT
				post.id,
				post.user_id,
				user.email,
				user.nickname,
				user.avatar,
				post.is_top,
				post.is_choice,
				post.text,
				post.images,
				post.likes,
				post.replies,
				post.collects,
				post.create_time,
				post.update_time,
				CASE WHEN pl.user_id IS NOT NULL THEN 1 ELSE 0 END AS like_status,
				CASE WHEN pc.user_id IS NOT NULL THEN 1 ELSE 0 END AS CollectStatus
			FROM post
			LEFT JOIN user ON post.user_id = user.id
			LEFT JOIN post_likes pl ON post.id = pl.post_id AND pl.user_id = ?
			LEFT JOIN post_collects pc ON post.id = pc.post_id AND pc.user_id = ?
			WHERE post.user_id = ?
	`
	// 执行查询语句
	result := database.GetMySQL().Raw(sql, userId, userId, userId).Scan(&posts)

	if result.Error != nil {
		err = result.Error
	}
	return posts, err
}

// GetLikedPostsByUserId 根据已登录用户ID获取点赞的帖子列表
func GetLikedPostsByUserId(userId int64) (posts []PostResponse, err error) {
	// 创建查询语句
	sql := `
			SELECT
				post.id,
				post.user_id,
				user.email,
				user.nickname,
				user.avatar,
				post.is_top,
				post.is_choice,
				post.text,
				post.images,
				post.likes,
				post.replies,
				post.collects,
				post.create_time,
				post.update_time,
				CASE WHEN pl.user_id IS NOT NULL THEN 1 ELSE 0 END AS like_status,
				CASE WHEN pc.user_id IS NOT NULL THEN 1 ELSE 0 END AS CollectStatus
			FROM post
			LEFT JOIN user ON post.user_id = user.id
			LEFT JOIN post_likes pl ON post.id = pl.post_id AND pl.user_id = ?
			LEFT JOIN post_collects pc ON post.id = pc.post_id AND pc.user_id = ?
			WHERE pl.user_id = ?
	`
	// 执行查询语句
	result := database.GetMySQL().Raw(sql, userId, userId, userId).Scan(&posts)

	if result.Error != nil {
		err = result.Error
	}
	return posts, err
}

// GetRepliedPostsByUserId 根据已登录用户ID获取回复的帖子列表
func GetRepliedPostsByUserId(userId int64) (posts []PostResponse, err error) {
	// 创建查询语句
	sql := `
			SELECT
				post.id,
				post.user_id,
				user.email,
				user.nickname,
				user.avatar,
				post.is_top,
				post.is_choice,
				post.text,
				post.images,
				post.likes,
				post.replies,
				post.collects,
				post.create_time,
				post.update_time,
				IFNULL(pl.user_id, 0) as LikeStatus,
		        IFNULL(pc.user_id, 0) as CollectStatus
			FROM post
			LEFT JOIN user ON post.user_id = user.id
			LEFT JOIN post_likes pl ON post.id = pl.post_id AND pl.user_id = ?
			LEFT JOIN post_collects pc ON post.id = pc.post_id AND pc.user_id = ?
			LEFT JOIN reply r on user.id = r.user_id
			WHERE r.user_id = ?
	`
	// 执行查询语句
	result := database.GetMySQL().Raw(sql, userId, userId, userId).Scan(&posts)

	if result.Error != nil {
		err = result.Error
	}
	return posts, err
}

// LikePost 给帖子点赞
func LikePost(userId, postId int64) error {
	// 生成点赞记录
	like := PostLikes{
		PostId: postId,
		UserId: userId,
	}

	// 写入点赞记录
	result := database.GetMySQL().Create(&like)

	if result.Error != nil {
		return result.Error
	}

	// 更新帖子点赞数
	database.GetMySQL().Table("post").Where("id = ?", postId).
		Update("likes", gorm.Expr("likes + ?", 1))

	// 更新用户点赞数
	database.GetMySQL().Table("user").Where("id = (?)", database.GetMySQL().Select("user_id").Where("id = ?", postId).Table("post")).
		Update("likes", gorm.Expr("likes + ?", 1))

	return result.Error
}

// UnlikePost 取消给帖子点赞
func UnlikePost(userId, postId int64) error {
	result := database.GetMySQL().Where("user_id = ? AND post_id = ?", userId, postId).Delete(&PostLikes{})

	if result.Error != nil {
		return result.Error
	}

	// 更新帖子点赞数
	database.GetMySQL().Table("post").Where("id = ?", postId).
		Update("likes", gorm.Expr("likes - ?", 1))

	// 更新用户点赞数
	database.GetMySQL().Table("user").Where("id = (?)", database.GetMySQL().Select("user_id").Where("id = ?", postId).Table("post")).
		Update("likes", gorm.Expr("likes - ?", 1))

	return result.Error
}

// CollectPost 收藏帖子
func CollectPost(userId, postId int64) error {
	collect := PostCollects{
		PostId: postId,
		UserId: userId,
	}

	result := database.GetMySQL().Create(&collect)

	if result.Error != nil {
		return result.Error
	}

	// 更新帖子收藏数
	database.GetMySQL().Table("post").Where("id = ?", postId).
		Update("collects", gorm.Expr("collects + ?", 1))

	return result.Error
}

// UnCollectPost 取消收藏帖子
func UnCollectPost(userId, postId int64) error {
	result := database.GetMySQL().Where("user_id = ? AND post_id = ?", userId, postId).Delete(&PostCollects{})

	if result.Error != nil {
		return result.Error
	}

	// 更新帖子收藏数
	database.GetMySQL().Table("post").Where("id = ?", postId).
		Update("collects", gorm.Expr("collects - ?", 1))

	return result.Error
}

// CheckPostByPostId 根据帖子ID检查帖子是否存在
func CheckPostByPostId(postId int64) (post Post, err error) {
	result := database.GetMySQL().Where("id = ?", postId).First(&post)

	if result.Error != nil {
		err = result.Error
	}

	return post, err
}

// GetPostImagesByPostId 根据帖子ID获取帖子图片
func GetPostImagesByPostId(postId int64) (images []string, err error) {
	var imageString string
	result := database.GetMySQL().Table("post").Where("id = ?", postId).Select("images").Scan(&imageString)

	if result.Error != nil {
		return images, result.Error
	}

	// 将图片字符串转换为数组
	images = strings.Split(imageString, ",")

	return images, nil
}

// DeletePostByPostId 根据帖子ID删除帖子
func DeletePostByPostId(postId int64) error {
	// 删除帖子
	result := database.GetMySQL().Where("id = ?", postId).Delete(&Post{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// DeletePostLikesByPostId 根据帖子ID删除帖子点赞记录
func DeletePostLikesByPostId(postId int64) error {
	result := database.GetMySQL().Exec("DELETE FROM post_likes WHERE post_id = ?", postId)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
