package model

import (
	"time"

	"gorm.io/gorm"
)

/*
It will record the user relation with the topic
*/
type UserTopic struct {
	gorm.Model
	ID      uint64 `gorm:"primaryKey"`
	UserID  uint64 `gorm:"column:user_id;index"`
	TopicID uint64 `gorm:"column:topic_id;index"`

	LastReadAt         time.Time // last read time
	LastReadPostNumber uint64    // last read post number
	Subscription       string    // subscription status enum(follow, ignore)
}
