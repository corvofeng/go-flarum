package model

import "gorm.io/gorm"

type ActionRecord struct {
	gorm.Model
	UserID uint64 `gorm:"column:user_id;index"`

	PostData string
}

func CreateActionRecord(gormDB *gorm.DB, userID uint64, postData string) error {
	record := ActionRecord{
		UserID:   userID,
		PostData: postData,
	}
	return gormDB.Create(&record).Error
}
