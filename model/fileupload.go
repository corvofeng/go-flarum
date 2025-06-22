package model

import (
	"time"

	"gorm.io/gorm"
)

type UserFiles struct {
	gorm.Model
	ID         uint64 `gorm:"primaryKey"`
	UUID       string `json:"uuid" gorm:"index:idx_uuid,unique"` // 文件的唯一标识符
	UserID     uint64 `json:"user_id" gorm:"index:idx_user_id"`
	FileName   string `json:"file_name"`
	FilePath   string `json:"file_path"`
	FileSize   int64  `json:"file_size"`
	FileType   string `json:"file_type"`
	Visibility string `json:"visibility"` // 可见性，公开或私有
	// Deleted    bool   `json:"deleted"`    // 是否已删除
	// 其他可能的字段

	// UploadTime int64 `json:"upload_time"` // 上传时间
}

func SQLGetUserFiles(gormDB *gorm.DB, userID uint64) (files []UserFiles, err error) {
	// 获取指定用户的所有文件
	err = gormDB.Where("user_id = ?", userID).Find(&files).Error
	return
}

func getCurDate() string {
	// 获取当前日期的字符串表示, 兼容 MySQL 和 PostgreSQL
	return time.Now().Format("2006-01-02")
}

// 获取用户当日上传的文件数
func SQLGetUserDailyUploads(gormDB *gorm.DB, userID uint64) (count int64, err error) {
	// 获取用户当日上传的文件数
	err = gormDB.Model(&UserFiles{}).
		Where("user_id = ? AND DATE(created_at) = ?", userID, getCurDate()).
		Count(&count).Error
	return
}

func SQLSaveUserFile(gormDB *gorm.DB, file *UserFiles) error {
	// 保存用户文件记录
	return gormDB.Create(file).Error
}
