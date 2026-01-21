package main

import (
	"sns-4-illustrators-server/database"
	"sns-4-illustrators-server/handlers"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/pprof"
)

func main() {
	// DBセットアップ（マイグレーション）
	database.InitDB()

	router := gin.Default()

	pprof.Register(router)

	// http://127.0.0.1:8080/uploads/file_name にアクセスすると、./uploads/file_name のファイルを返す
	router.Static("/uploads", "./uploads")
	
	// ルーティング
	router.POST("/upload", handlers.UploadImage)
	router.GET("/images", handlers.GetImages)

	router.Run(":8080")
}