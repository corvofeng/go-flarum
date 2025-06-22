package model

import (
	"fmt"

	"gorm.io/gorm"
)

type BlogMeta struct {
	gorm.Model
	ID      uint64 `gorm:"primaryKey"`
	TopicID uint64 `gorm:"column:topic_id;index"` // 关联的主题ID

	Summary         string `json:"summary"`           // 摘要
	FeaturedImage   string `json:"featured_image"`    // 特色图片
	IsFeatured      bool   `json:"is_featured"`       // 是否是特色文章
	IsPendingReview bool   `json:"is_pending_review"` // 是否待审核
	IsSized         bool   `json:"is_sized"`          // 是否已调整大小

	// TopicData Topic `gorm:"foreignKey:ID;references:TopicID"` // 关联的主题
}

func (blogMeta *BlogMeta) CreateFlarumBlogMeta(gormDB *gorm.DB) (bool, error) {
	// 创建或更新博客元数据
	if err := gormDB.Create(blogMeta).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (blogMeta *BlogMeta) CreateOrUpdate(gormDB *gorm.DB) error {
	// 创建或更新博客元数据
	if blogMeta.ID == 0 {
		return gormDB.Create(blogMeta).Error
	}
	return gormDB.Save(blogMeta).Error
}

func (blogMeta *BlogMeta) GetFormatedString() string {
	return fmt.Sprintf(
		"[BlogMeta] ID: %d, TopicID: %d, Summary: %s",
		blogMeta.ID,
		blogMeta.TopicID,
		blogMeta.Summary,
	)
}

func SQLGetAllBlogMeta(gormDB *gorm.DB) (metas []BlogMeta, err error) {
	// 获取所有博客元数据
	// err = gormDB.Preload("TopicData").Find(&metas).Error
	err = gormDB.Find(&metas).Error
	if err != nil {
		return nil, err
	}
	return metas, nil
}

func SQLGetBlogMetaByTopicID(gormDB *gorm.DB, topicID uint64) (meta BlogMeta, err error) {
	// 获取指定主题的博客元数据
	err = gormDB.First(&meta, "topic_id = ?", topicID).Error
	if err != nil {
		return meta, err
	}
	return meta, nil
}

func SQLGetBlogMetaByID(gormDB *gorm.DB, id uint64) (meta BlogMeta, err error) {
	// 获取指定ID的博客元数据
	err = gormDB.First(&meta, "id = ?", id).Error
	if err != nil {
		return meta, err
	}
	return meta, nil
}

func SQLSaveBlogMeta(gormDB *gorm.DB, meta *BlogMeta) error {
	// 保存或更新博客元数据
	if meta.ID == 0 {
		return gormDB.Create(meta).Error
	}
	return gormDB.Save(meta).Error
}
