package database

import (
	"fmt"
	"log"
	"os"
	"sns-4-illustrators-server/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ================================================
// テーブルを追加する場合、ここにモデルを追加
// ================================================
var AllModels = []interface{}{
	&models.Image{},
}

func InitDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tokyo",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("ERROR: DB接続失敗: ", err)
	}

	// マイグレーションの実行
	if err := RunAllMigrations(DB); err != nil {
		log.Fatal("ERROR: マイグレーション失敗: ", err)
	}

	fmt.Println("DBセットアップ完了")
}

// ================================================
// AllModels に定義された全モデルをマイグレーション
// ================================================
func RunAllMigrations(db *gorm.DB) error {
	return db.AutoMigrate(AllModels...)
}