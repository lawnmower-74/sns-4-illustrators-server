package models

import "time"

type Image struct {
	ID              uint       `gorm:"primaryKey"`
	IllustratorName string     `gorm:"column:illustrator_name;size:255;not null;index"`
	FileName        string     `gorm:"column:file_name;size:255;not null;unique;index"`
	FileSize        int64      `gorm:"column:file_size;not null"`
	MimeType        string     `gorm:"column:mime_type;size:50;not null"`
	StoragePath     string     `gorm:"column:storage_path;size:500;not null"`
	ShotAt          *time.Time `gorm:"column:shot_at;index"`
	CreatedAt       time.Time  `gorm:"column:created_at;autoCreateTime"`
}