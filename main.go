package main

import (
	"sns-4-illustrators-server/database"
	"sns-4-illustrators-server/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	// DBセットアップ（マイグレーション）
	database.InitDB()

	r := gin.Default()
	
	// ルーティング
	r.POST("/upload", handlers.UploadImage)
	r.GET("/images", handlers.GetImages)

	r.Run(":8080")
}