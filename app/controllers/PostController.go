package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"schoolChat/app/models"
	Results "schoolChat/app/result"
	"strconv"
	"strings"
)

// AddPost 新增帖子
func AddPost(c *gin.Context) {
	var post models.Post
	err := c.BindJSON(&post)

	if err != nil {
		c.JSON(200, Results.Err.Fail("获取参数失败"))
		return
	}

	userId, _ := c.Get("id")

	post.UserId = userId.(int64)

	post, err = models.AddPost(post)

	if err != nil {
		c.JSON(200, Results.Err.Fail("新增失败"))
		return
	}

	c.JSON(200, Results.Ok.Success(post))
}

// GetAllPost 获取所有帖子
func GetAllPost(c *gin.Context) {

	result, err := models.GetAllPost()

	if err != nil {
		c.JSON(200, Results.Err.Fail("查询失败"))
		return
	}

	c.JSON(200, Results.Ok.Success(result))
}

// GetAllPostByLogin 获取所有帖子
func GetAllPostByLogin(c *gin.Context) {

	userId, _ := c.Get("id")

	result, err := models.GetAllPostByLogin(userId.(int64))

	if err != nil {
		c.JSON(200, Results.Err.Fail("查询失败"))
		return
	}

	c.JSON(200, Results.Ok.Success(result))
}

// GetPostByPostId 根据帖子ID获取帖子详情
func GetPostByPostId(c *gin.Context) {
	// 获取参数
	postId, _ := strconv.ParseInt(c.Query("postId"), 10, 64)

	// 获取用户id
	userId, _ := c.Get("id")

	result, err := models.GetPostByPostId(userId.(int64), postId)

	if err != nil {
		c.JSON(200, Results.Err.Fail("查询失败"))
		return
	}

	c.JSON(200, Results.Ok.Success(result))
}

// GetPostByUserId 根据用户ID获取帖子列表
func GetPostByUserId(c *gin.Context) {
	// 获取参数
	userId, _ := strconv.ParseInt(c.Query("userId"), 10, 64)

	result, err := models.GetPostByUserId(userId)

	if err != nil {
		c.JSON(200, Results.Err.Fail("查询失败"))
		return
	}

	c.JSON(200, Results.Ok.Success(result))
}

// GetPostByLoginUserId 根据已登录用户ID获取帖子列表
func GetPostByLoginUserId(c *gin.Context) {
	// 获取用户id
	userId, _ := c.Get("id")

	result, err := models.GetPostByLoginUserId(userId.(int64))

	if err != nil {
		c.JSON(200, Results.Err.Fail("查询失败"))
		return
	}

	c.JSON(200, Results.Ok.Success(result))
}

// GetLikedPostsByUserId 根据已登录用户ID获取点赞的帖子列表
func GetLikedPostsByUserId(c *gin.Context) {
	// 获取用户id
	userId, _ := c.Get("id")

	result, err := models.GetLikedPostsByUserId(userId.(int64))

	if err != nil {
		c.JSON(200, Results.Err.Fail("查询失败"))
		return
	}

	c.JSON(200, Results.Ok.Success(result))
}

// GetRepliedPostsByUserId 根据已登录用户ID获取回复的帖子列表
func GetRepliedPostsByUserId(c *gin.Context) {
	// 获取用户id
	userId, _ := c.Get("id")

	result, err := models.GetRepliedPostsByUserId(userId.(int64))

	if err != nil {
		c.JSON(200, Results.Err.Fail("查询失败"))
		return
	}

	c.JSON(200, Results.Ok.Success(result))
}

// LikePost 给帖子点赞
func LikePost(c *gin.Context) {
	// 获取参数
	var postId int64
	err := c.BindJSON(&postId)

	// 获取用户id
	userId, _ := c.Get("id")

	err = models.LikePost(userId.(int64), postId)

	if err != nil {
		c.JSON(200, Results.Err.Fail("点赞失败"))
		return
	}

	c.JSON(200, Results.Ok.Success("点赞成功"))
}

// UnlikePost 取消给帖子点赞
func UnlikePost(c *gin.Context) {
	// 获取参数
	var postId int64
	err := c.BindJSON(&postId)

	// 获取用户id
	userId, _ := c.Get("id")

	err = models.UnlikePost(userId.(int64), postId)

	if err != nil {
		c.JSON(200, Results.Err.Fail("取消点赞失败"))
		return
	}

	c.JSON(200, Results.Ok.Success("取消点赞成功"))
}

// CollectPost 收藏帖子
func CollectPost(c *gin.Context) {
	// 获取参数
	var postId int64
	err := c.BindJSON(&postId)

	// 获取用户id
	userId, _ := c.Get("id")

	err = models.CollectPost(userId.(int64), postId)

	if err != nil {
		c.JSON(200, Results.Err.Fail("收藏失败"))
		return
	}

	c.JSON(200, Results.Ok.Success("收藏成功"))
}

// UnCollectPost 取消收藏帖子
func UnCollectPost(c *gin.Context) {
	// 获取参数
	var postId int64
	err := c.BindJSON(&postId)

	// 获取用户id
	userId, _ := c.Get("id")

	err = models.UnCollectPost(userId.(int64), postId)

	if err != nil {
		c.JSON(200, Results.Err.Fail("取消收藏失败"))
		return
	}

	c.JSON(200, Results.Ok.Success("取消收藏成功"))
}

// DeletePostByPostId 根据帖子ID删除帖子
func DeletePostByPostId(c *gin.Context) {
	// 获取参数
	// 获取参数
	postId, err := strconv.ParseInt(c.Query("postId"), 10, 64)

	if err != nil {
		c.JSON(200, Results.Err.Fail("删除失败。原因:参数错误"))
		return
	}

	userId, _ := c.Get("id") // 获取用户id

	post, err := models.CheckPostByPostId(postId) // 检查帖子是否存在

	if post.ID == 0 {
		c.JSON(200, Results.Err.Fail("删除失败。原因:帖子不存在"))
		return
	}

	if post.UserId != userId.(int64) {
		c.JSON(200, Results.Err.Fail("删除失败。原因:无权限"))
		return
	}

	// 将图片字符串转换为数组
	images := strings.Split(post.Images, ",")

	fmt.Printf("images: %v\n", images)

	err = models.DeletePostByPostId(postId) // 删除帖子

	if err != nil {
		c.JSON(200, Results.Err.Fail("删除失败"))
		return
	}

	err = models.DeletePostLikesByPostId(postId) // 删除帖子点赞

	if err != nil {
		c.JSON(200, Results.Err.Fail("删除失败。帖子点赞记录删除失败"))
		return
	}

	err = models.DeleteReplyByPostId(postId) // 删除帖子回复

	if err != nil {
		c.JSON(200, Results.Err.Fail("删除失败。帖子回复删除失败"))
		return
	}

	err = DeleteImage(images) // 删除图片

	if err != nil {
		c.JSON(200, Results.Err.Fail(err.Error()))
		return
	}

	c.JSON(200, Results.Ok.Success("删除成功"))
}
