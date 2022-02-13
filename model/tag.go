package model

import (
	"gorm.io/gorm"
)

// Tag 帖子分类
type Tag struct {
	gorm.Model
	ID   uint64 `gorm:"primaryKey"`
	Name string `json:"name"`

	URLName     string `json:"urlname"`
	Articles    uint64 `json:"articles"`
	About       string `json:"about"`
	ParentID    uint64 `json:"parent_id"`
	Position    uint64 `json:"position"`
	Description string `json:"description"`
	Hidden      bool   `json:"hidden"`
	Color       string `json:"color"`
	IconIMG     string `json:"icon_img"`

	Topics []Topic `gorm:"many2many:topic_tags;"`
}

func SQLGetTags(gormDB *gorm.DB) (tags []Tag, err error) {
	// find all tags
	err = gormDB.Find(&tags).Error
	return
}

func SQLGetTagByUrlName(gormDB *gorm.DB, urlName string) (tag Tag, err error) {
	err = gormDB.First(&tag, "url_name = ?", urlName).Error
	return
}
func SQLGetTagByID(gormDB *gorm.DB, id uint64) (tag Tag, err error) {
	err = gormDB.First(&tag, id).Error
	return
}
