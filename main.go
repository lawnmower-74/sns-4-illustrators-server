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
	
	// ルーティング
	router.POST("/upload", handlers.UploadImage)
	router.GET("/images", handlers.GetImages)

	router.Run(":8080")
}