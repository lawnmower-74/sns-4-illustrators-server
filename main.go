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
	
	// 静的ファイルの配信（ブラウザから画像を見れるようにする）
	r.Static("/view-images", "./uploads")

	// ルーティング
	r.POST("/upload", handlers.UploadImage)

	r.Run(":8080")
}