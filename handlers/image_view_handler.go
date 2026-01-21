package handlers

import (
    "fmt"
	"net/http"
	"sns-4-illustrators-server/database"
	"sns-4-illustrators-server/models"

	"github.com/gin-gonic/gin"
)

func GetImages(c *gin.Context) {
	var images []models.Image

	// 全レコードを取得（※最新のものから順に取得）
	if err := database.DB.Order("created_at desc").Find(&images).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": "データの取得に失敗しました"})
		return
	}

	// JSONで返却
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("%d件のデータ取得に成功しました", len(images)),
		"data":    images,
	})
}