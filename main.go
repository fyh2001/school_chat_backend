package main

import (
	"schoolChat/database"
	"schoolChat/routes"
	"schoolChat/util"
)

func main() {
	database.InitMySQL() // 初始化MySQL数据库连接
	database.InitRedis() // 初始化Redis连接

	util.InitSMSConfig() // 初始化短信配置

	r := routes.InitRouter()

	_ = r.Run(":9112")
}
