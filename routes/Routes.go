package routes

import (
	"github.com/gin-gonic/gin"
	"schoolChat/app/controllers"
	"schoolChat/app/middleware/cors"
	"schoolChat/app/middleware/jwt"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	r.Static("/public", "./public") // 静态文件服务

	r.Use(cors.Cors())

	//路由组件
	root := r.Group("/api")
	{
		user := root.Group("/user")
		{
			user.GET("/getAll", jwt.JWT(), controllers.GetAllUser)
			user.GET("/getUserByUserId", jwt.JWT(), controllers.GetUserByUserId)
			user.POST("/loginOrRegisterByPhone", controllers.LoginOrRegisterByPhone)
			user.POST("/loginOrRegisterByMail", controllers.LoginOrRegisterByMail)
			user.PUT("/update", jwt.JWT(), controllers.UpdateUser)
			user.PUT("/updateNickname", jwt.JWT(), controllers.UpdateNickname)
			user.POST("/getMailCode", controllers.GetMailCode)
			user.POST("/getPhoneCode", controllers.GetPhoneCode)
		}
		post := root.Group("/post")
		{
			post.POST("/addPost", jwt.JWT(), controllers.AddPost)
			post.GET("/getAllPost", controllers.GetAllPost)
			post.GET("/getAllPostByLogin", jwt.JWT(), controllers.GetAllPostByLogin)
			post.GET("/getPostByPostId", jwt.JWT(), controllers.GetPostByPostId)
			post.GET("/getPostByUserId", jwt.JWT(), controllers.GetPostByUserId)
			post.GET("/getPostByLoginUserId", jwt.JWT(), controllers.GetPostByLoginUserId)
			post.GET("/getLikedPostsByUserId", jwt.JWT(), controllers.GetLikedPostsByUserId)
			post.GET("/getRepliedPostsByUserId", jwt.JWT(), controllers.GetRepliedPostsByUserId)
			post.POST("/likePost", jwt.JWT(), controllers.LikePost)
			post.POST("/unlikePost", jwt.JWT(), controllers.UnlikePost)
			post.POST("collectPost", jwt.JWT(), controllers.CollectPost)
			post.POST("unCollectPost", jwt.JWT(), controllers.UnCollectPost)
		}
		reply := root.Group("/reply")
		{
			reply.POST("/addReply", jwt.JWT(), controllers.AddReply)
			reply.GET("/getReplyByPostId", jwt.JWT(), controllers.GetReplyByPostId)
			reply.GET("/getReplyByReplyId", jwt.JWT(), controllers.GetReplyByReplyId)
			reply.GET("/getSecondReplyByReplyId", jwt.JWT(), controllers.GetSecondReplyByReplyId)
			reply.POST("likeReply", jwt.JWT(), controllers.LikeReply)
			reply.POST("unlikeReply", jwt.JWT(), controllers.UnlikeReply)
		}
		file := root.Group("/file")
		{
			file.POST("/upload", jwt.JWT(), controllers.UploadHandler)
			file.GET("/download", controllers.DownloadHandler)
		}
	}

	return r
}
