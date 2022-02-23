package cronjob

import (
	"database/sql"
	"time"

	"github.com/go-redis/redis/v7"

	"zoe/model"
	"zoe/system"

	logging "github.com/op/go-logging"
)

// BaseHandler I do not know
type BaseHandler struct {
	App *system.Application
}

// MainCronJob my job
func (h *BaseHandler) MainCronJob() {
	sqlDB := h.App.MySQLdb
	// cntDB := h.App.Db
	// scf := h.App.Cf.Site
	logger := h.App.Logger
	redisDB := h.App.RedisDB
	// tick1 := time.Tick(3600 * time.Second)
	// tick2 := time.Tick(120 * time.Second)
	// tick3 := time.Tick(30 * time.Minute)
	// tick4 := time.Tick(31 * time.Second)
	// tick5 := time.Tick(1 * time.Minute)
	// tickRefreshOrder := time.Tick(3 * time.Second)
	// daySecond := int64(3600 * 24)

	// 每3小时将Redis中的数据库中的排序数据刷新
	tickResortRankMap := time.Tick(3 * time.Hour)

	// 每小时将点击量刷到数据库中
	tickStoreHitToMySQL := time.Tick(1 * time.Hour)

	// 每十分钟刷新排序
	// tickRefreshOrder := time.Tick(10 * time.Minute)

	logger.Info("Start cron job")
	syncWithMySQL(logger, sqlDB, redisDB)
	refreshRankMap(logger)

	for {
		select {
		case <-tickStoreHitToMySQL:
			syncWithMySQL(logger, sqlDB, redisDB)
		case <-tickResortRankMap:
			refreshRankMap(logger)
		}
	}
}

// refreshRankMap 刷新redis中存储的排序数据
func refreshRankMap(logger *logging.Logger) {
	logger.Info("===== Start refresh rank map =====")
	model.TimelyResort()
	logger.Info("=====  End  refresh rank map =====")
}

// syncWithMySQL 将Redis中的统计数据同步给mysql
func syncWithMySQL(logger *logging.Logger, sqlDB *sql.DB, redisDB *redis.Client) {
	logger.Info("===== start sync hits with the mysql =====")
	// data, _ := redisDB.HGetAll("article_views").Result()
	// for aid, clickCnt := range data {
	// _aid, _ := strconv.ParseUint(aid, 10, 64)
	// _clickCnt, _ := strconv.ParseUint(clickCnt, 10, 64)
	// logger.Debugf("Set %4d with %d", _aid, _clickCnt)
	// model.SQLArticleSetClickCnt(sqlDB, _aid, _clickCnt)
	// }
	logger.Info("=====  end  sync hits with the mysql =====")
}
