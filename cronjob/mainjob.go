package cronjob

import (
	"database/sql"
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v7"

	"zoe/model"
	"zoe/system"

	"github.com/boltdb/bolt"
	"github.com/ego008/youdb"
	logging "github.com/op/go-logging"
	"github.com/weint/httpclient"
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
		// case <-tick1:
		// 	limit := 10
		// 	timeBefore := uint64(time.Now().UTC().Unix() - daySecond)
		// 	scoreStartB := youdb.I2b(timeBefore)
		// 	zbnList := []string{
		// 		"article_detail_token",
		// 		"user_login_token",
		// 	}
		// 	for _, bn := range zbnList {
		// 		rs := db.Zrscan(bn, []byte(""), scoreStartB, limit)
		// 		if rs.State == "ok" {
		// 			keys := make([][]byte, len(rs.Data)/2)
		// 			j := 0
		// 			for i := 0; i < (len(rs.Data) - 1); i += 2 {
		// 				keys[j] = rs.Data[i]
		// 				j++
		// 			}
		// 			db.Zmdel(bn, keys)
		// 		}
		// 	}

		// case <-tick2:
		// 	if scf.AutoGetTag && len(scf.GetTagApi) > 0 {
		// 		getTagFromTitle(db, scf.GetTagApi)
		// 	}
		// case <-tick3:
		// 	if h.App.Cf.Site.AutoDataBackup {
		// 		dataBackup(db)
		// 	}
		// case <-tick4:
		// 	setArticleTag(db)
		// case <-tickRefreshOrder:
		// saveToRedisSorted(logger, db, redisDB)
		// syncWithMySQL(logger, sqlDB, redisDB)
		// saveToRedisSorted(logger, db, redisDB)
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
	data, _ := redisDB.HGetAll("article_views").Result()
	for aid, clickCnt := range data {
		_aid, _ := strconv.ParseUint(aid, 10, 64)
		_clickCnt, _ := strconv.ParseUint(clickCnt, 10, 64)
		// logger.Debugf("Set %4d with %d", _aid, _clickCnt)
		model.SQLArticleSetClickCnt(sqlDB, _aid, _clickCnt)
	}
	logger.Info("=====  end  sync hits with the mysql =====")
}

func dataBackup(db *youdb.DB) {
	filePath := "databackup/" + time.Now().UTC().Format("2006-01-02") + ".db"
	if _, err := os.Stat(filePath); err != nil {
		// path not exists
		err := db.View(func(tx *bolt.Tx) error {
			return tx.CopyFile(filePath, 0600)
		})
		if err == nil {
			// todo upload to qiniu
		}
	}
}

func getTagFromTitle(db *youdb.DB, apiURL string) {
	rs := db.Hscan("task_to_get_tag", []byte(""), 1)
	if rs.State == "ok" {
		aidB := rs.Data[0][:]

		rs2 := db.Hget("article", aidB)
		if rs2.State != "ok" {
			db.Hdel("task_to_get_tag", aidB)
			return
		}
		aobj := model.Article{}
		json.Unmarshal(rs2.Data[0], &aobj)
		if len(aobj.Tags) > 0 {
			db.Hdel("task_to_get_tag", aidB)
			return
		}

		hc := httpclient.NewHttpClientRequest("POST", apiURL)
		hc.Param("state", "ok")
		hc.Param("ms", string(rs.Data[1]))

		t := struct {
			Code int    `json:"code"`
			Tag  string `json:"tag"`
		}{}
		err := hc.ReplyJson(&t)
		if err != nil {
			return
		}
		if hc.Status() == 200 && t.Code == 200 {
			if len(t.Tag) > 0 {
				tags := strings.Split(t.Tag, ",")
				if len(tags) > 5 {
					tags = tags[:5]
				}

				// get once more
				rs2 := db.Hget("article", youdb.I2b(aobj.ID))
				if rs2.State == "ok" {
					aobj := model.Article{}
					json.Unmarshal(rs2.Data[0], &aobj)
					aobj.Tags = strings.Join(tags, ",")
					jb, _ := json.Marshal(aobj)
					db.Hset("article", youdb.I2b(aobj.ID), jb)

					// tag send task work，自动处理tag与文章id
					at := model.ArticleTag{
						ID:      aobj.ID,
						OldTags: "",
						NewTags: aobj.Tags,
					}
					jb, _ = json.Marshal(at)
					db.Hset("task_to_set_tag", youdb.I2b(at.ID), jb)
				}
			}
			db.Hdel("task_to_get_tag", aidB)
		}
	}
}

func setArticleTag(db *youdb.DB) {
	rs := db.Hscan("task_to_set_tag", nil, 1)
	if rs.OK() {
		info := model.ArticleTag{}
		err := json.Unmarshal(rs.Data[1], &info)
		if err != nil {
			return
		}
		//log.Println("aid", info.ID)

		// set tag
		oldTag := strings.Split(info.OldTags, ",")
		newTag := strings.Split(info.NewTags, ",")

		// remove
		for _, tag1 := range oldTag {
			contains := false
			for _, tag2 := range newTag {
				if tag1 == tag2 {
					contains = true
					break
				}
			}
			if !contains {
				tagLower := strings.ToLower(tag1)
				db.Hdel("tag:"+tagLower, youdb.I2b(info.ID))
				db.Zincr("tag_article_num", []byte(tagLower), -1)
			}
		}

		// add
		for _, tag1 := range newTag {
			contains := false
			for _, tag2 := range oldTag {
				if tag1 == tag2 {
					contains = true
					break
				}
			}
			if !contains {
				tagLower := strings.ToLower(tag1)
				// 记录所有tag，只增不减
				if db.Hget("tag", []byte(tagLower)).State != "ok" {
					db.Hset("tag", []byte(tagLower), []byte(""))
					db.HnextSequence("tag") // 添加这一行
				}
				// check if not exist !important
				if db.Hget("tag:"+tagLower, youdb.I2b(info.ID)).State != "ok" {
					db.Hset("tag:"+tagLower, youdb.I2b(info.ID), []byte(""))
					db.Zincr("tag_article_num", []byte(tagLower), 1)
				}
			}
		}

		db.Hdel("task_to_set_tag", youdb.I2b(info.ID))
	}
}
