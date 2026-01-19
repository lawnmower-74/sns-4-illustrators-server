package main

import (
	"sns-4-illustrators-server/database"
	"github.com/gin-gonic/gin"
)

func main() {
	// セットアップ（マイグレーション）
	database.InitDB()

	r := gin.Default()
	
	// ルーティング設定（ここにアップロード処理などを書く）
	// 今後はここも「handlers」などのフォルダに分けていくとさらに綺麗になります
	r.POST("/upload", func(c *gin.Context) {
		// database.DB を使って保存処理...
	})

	r.Run(":8080")
}