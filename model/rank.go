package model

import (
	"database/sql"
	"fmt"
	"goyoubbs/util"
	"sync"

	"github.com/ego008/youdb"
	"github.com/go-redis/redis/v7"

	"strconv"
)

// WeightAble 可以获取权值的一种结构
type WeightAble interface {
	GetWeight() uint64
}

// ArticleRankItem 记录每个话题的权重
type ArticleRankItem struct {
	AID     uint64 `json:"a_id"`
	Weight  uint64
	SQLDB   *sql.DB
	CacheDB *youdb.DB
	RedisDB *redis.Client
}

// CategoryRankData 一个分类下的排序数据
type CategoryRankData struct {
	CID       uint64     `json:"c_id"`
	mtx       sync.Mutex // 同一时刻, 只允许一个协程操纵该分类的记录
	maxID     uint64     // 数据库游标, 记录当前已读取数据的最大值, 从数据库中读取新的数据时使用
	topicData []ArticleRankItem
}

// RankMap time to live map
type RankMap struct {
	m       map[uint64]*CategoryRankData
	mtx     sync.Mutex // 同一时刻, 只允许一个协程操纵map
	SQLDB   *sql.DB
	CacheDB *youdb.DB
	RedisDB *redis.Client
}

func getWeight(rankMap *RankMap, aid uint64) float64 {
	article, err := SQLArticleGetByID(rankMap.SQLDB, rankMap.CacheDB, rankMap.RedisDB, aid)
	if util.CheckError(err, "查询帖子") {
		return 0
	}
	return article.GetWeight(
		rankMap.SQLDB,
		rankMap.CacheDB,
		rankMap.RedisDB,
	)
}

var rankMap *RankMap
var rankRedisDB *redis.Client

// TimelyResort 刷新Redis数据库中每个帖子的权重
func TimelyResort() {
	// 刷新所有节点的排序
	categoryList, err := SQLGetAllCategory(rankMap.SQLDB)
	logger := util.GetLogger()
	if util.CheckError(err, "获取所有节点") {
		return
	}

	for _, v := range categoryList {
		logger.Debugf("Start refresh %s", v.Name)

		// 删除redis中所有无效的帖子
		sqlDataDel, err := sqlGetAllArticleWithCID(rankMap.SQLDB, v.ID, false)
		if util.CheckError(err, fmt.Sprintf("获取%d节点下的无效的帖子列表", v.ID)) {
			return
		}
		for _, t := range sqlDataDel {
			_, err := rankRedisDB.ZRem(fmt.Sprintf("%d", v.ID), fmt.Sprintf("%d", t.ID)).Result()
			util.CheckError(err, "删除无效帖子")
		}

		// 将所有有效帖子更新至redis数据库中
		sqlDataAdd, err := sqlGetAllArticleWithCID(rankMap.SQLDB, v.ID, true)
		if util.CheckError(err, fmt.Sprintf("获取%d节点下的有效的帖子列表", v.ID)) {
			return
		}

		// 	首先从数据库中获取所有有效的ID
		for _, t := range sqlDataAdd {
			_, err := rankRedisDB.ZAddNX(fmt.Sprintf("%d", v.ID), &redis.Z{
				Score:  getWeight(rankMap, t.ID),
				Member: fmt.Sprintf("%d", t.ID)},
			).Result()
			util.CheckError(err, "更新当前帖子")
		}

		// 刷新权重
		rdsData, _ := rankRedisDB.ZRevRange(fmt.Sprintf("%d", v.ID), 0, -1).Result()
		for _, topicID := range rdsData {
			aid, _ := strconv.ParseUint(topicID, 10, 64)
			rankRedisDB.ZAddXX(fmt.Sprintf("%d", v.ID), &redis.Z{
				Score:  float64(getWeight(rankMap, aid)),
				Member: fmt.Sprintf("%d", aid)},
			)
		}
	}
}

func newRankMap() (m *RankMap) {
	m = &RankMap{m: make(map[uint64]*CategoryRankData)}
	return m
}

// RankMapInit init a ttl map
func RankMapInit(sqlDB *sql.DB, cntDB *youdb.DB, redisDB *redis.Client) {
	rankMap = newRankMap()
	rankMap.SQLDB = sqlDB
	rankMap.CacheDB = cntDB
	rankMap.RedisDB = redisDB

	rankRedisDB = redisDB
}

// GetRankMap you can get ttlmap by this.
func GetRankMap() *RankMap {
	return rankMap
}

func min(a, b uint64) uint64 {
	if a <= b {
		return a
	}
	return b
}

// GetTopicListByPageNum 通过给定的页码查找话题的ID值
func GetTopicListByPageNum(cid uint64, page uint64, limit uint64) []uint64 {
	var retData []uint64

	start := (page - 1) * limit
	end := (page)*limit - 1
	data, _ := rankRedisDB.ZRevRange(fmt.Sprintf("%d", cid), int64(start), int64(end)).Result()
	for _, val := range data {
		aid, _ := strconv.ParseUint(val, 10, 64)
		retData = append(retData, aid)
	}
	return retData
}

// AddNewArticleList 为某个分类添加话题
func AddNewArticleList(cid uint64, rankItems []ArticleRankItem) {
	m := GetRankMap()
	if _, ok := m.m[cid]; !ok { // 同一时刻只允许一个协程操作
		func() {
			m.mtx.Lock()
			defer m.mtx.Unlock()
			if _, ok := m.m[cid]; !ok { // 二次检查
				m.m[cid] = &CategoryRankData{CID: cid}
			}
		}()
	}

	var maxID uint64
	// fmt.Printf("Add rank item", rankItems)
	for _, d := range rankItems {
		if d.AID > maxID {
			maxID = d.AID
		}

		rankRedisDB.ZAdd(fmt.Sprintf("%d", cid), &redis.Z{
			Score:  float64(d.Weight),
			Member: fmt.Sprintf("%d", d.AID)})
	}

	crd := m.m[cid] // categoryRankData
	func() {
		crd.mtx.Lock()
		defer crd.mtx.Unlock()
		crd.topicData = append(crd.topicData, rankItems...) // 直接加入, 不做任何处理
		crd.maxID = maxID
	}()
}

// GetCIDArticleMax 获取当前分类的偏移值
func GetCIDArticleMax(cid uint64) uint64 {
	m := GetRankMap()
	if _, ok := m.m[cid]; ok {
		return m.m[cid].maxID
	}
	return 0
}
