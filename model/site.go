package model

import (
	"github.com/go-redis/redis/v7"
	"time"
)

type SiteInfo struct {
	Days     uint64
	UserNum  uint64
	NodeNum  uint64
	TagNum   uint64
	PostNum  uint64
	ReplyNum uint64
}

// GetDays 获取从建站开始, 到目前的天数, 用于主页中的显示
func GetDays(redisDB *redis.Client) uint64 {

	siteCreateTime, err := redisDB.Get("site_create_time").Uint64()
	if err != nil {
		siteCreateTime = 1557585456 // 2019-05-11 22:37:36 +0800 HKT
	}
	then := time.Unix(int64(siteCreateTime), 0)
	diff := time.Now().UTC().Sub(then)
	return uint64(diff.Hours()/24) + 1
}
