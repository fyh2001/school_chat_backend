package main

import (
	"schoolChat/database"
	"schoolChat/routes"
)

func main() {
	database.Init()
	r := routes.InitRouter()

	_ = r.Run(":9112")
}
