package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"sns-4-illustrators-server/database"
	"sns-4-illustrators-server/models"
	"time"

	"github.com/gin-gonic/gin"
)

func UploadImage(c *gin.Context) {
	// アップロードデータの抽出
	illustratorName := c.PostForm("illustrator_name")

	// image（キー）として送信された全画像データを取得
    form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": "画像データの解析に失敗しました"})
		return
	}
    files := form.File["image"]

	// ※結果を格納するリスト（レスポンス用）
    var successRecords []models.Image

    // ループで1枚ずつ登録
    for _, file := range files {
        
        // すでに登録されている場合はスキップ（画像名で重複チェック）
        var existingImage models.Image
        if err := database.DB.Where("file_name = ?", file.Filename).First(&existingImage).Error; err == nil {
            continue 
        }

        // 画像データ自体は uploads/ 配下に保存（※DBにはファイルパスを保存）
        filename := filepath.Base(file.Filename)
        savePath := filepath.Join("uploads", filename)
        if err := c.SaveUploadedFile(file, savePath); err != nil {
            continue
        }

        // DB登録
        now := time.Now()
        imageRecord := models.Image{
            IllustratorName: illustratorName,
            FileName:        filename,
            FileSize:        file.Size,
            MimeType:        file.Header.Get("Content-Type"),
            StoragePath:     savePath,
            ShotAt:          &now,
        }
        database.DB.Create(&imageRecord)
        
        successRecords = append(successRecords, imageRecord)
    }

    c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("%d件のアップロードに成功しました", len(successRecords)),
		"data":    successRecords,
	})
}