package model

import (
	"database/sql"
	"fmt"
	"github.com/ego008/youdb"
	"github.com/go-redis/redis/v7"
	"sort"
	"sync"
	"time"

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

// GetWeight 获取权重
func (articleItem *ArticleRankItem) GetWeight() uint64 {

	if articleItem.CacheDB != nil {
		articleItem.Weight = GetArticleCntFromRedisDB(
			articleItem.SQLDB,
			articleItem.CacheDB,
			articleItem.RedisDB,
			articleItem.AID,
		)
	}

	return articleItem.Weight
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

func getWeight(rankMap *RankMap, aid uint64) uint64 {

	weight := GetArticleCntFromRedisDB(
		rankMap.SQLDB,
		rankMap.CacheDB,
		rankMap.RedisDB,
		aid,
	)

	return weight
}

var rankMap *RankMap
var rankRedisDB *redis.Client

func (data *CategoryRankData) resort() {
	// fmt.Println("In sort ", data.CID, data.topicData)
	func() {
		data.mtx.Lock()
		defer data.mtx.Unlock()
		// Sort by age, keeping original order or equal elements.
		sort.SliceStable(data.topicData, func(i, j int) bool {
			return data.topicData[i].GetWeight() >= data.topicData[j].GetWeight()
		})
	}()
}

func timelyResort(sleeps uint64) {
	m := GetRankMap()
	for range time.Tick(time.Second * time.Duration(sleeps)) { // 每10s刷新一次
		func() {
			m.mtx.Lock()
			defer m.mtx.Unlock()
			for _, v := range m.m {

				data, _ := rankRedisDB.ZRevRange(fmt.Sprintf("%d", v.CID), 0, -1).Result()
				for _, topicID := range data {

					aid, _ := strconv.ParseUint(topicID, 10, 64)

					rankRedisDB.ZAddXX(fmt.Sprintf("%d", v.CID), &redis.Z{
						Score:  float64(getWeight(m, aid)),
						Member: fmt.Sprintf("%d", aid)},
					)
				}
			}
		}()
	}
}

func newRankMap() (m *RankMap) {
	m = &RankMap{m: make(map[uint64]*CategoryRankData)}
	return m
}

// RankMapInit init a ttl map
func RankMapInit(sleeps uint64, sqlDB *sql.DB, cntDB *youdb.DB, redisDB *redis.Client) {
	rankMap = newRankMap()
	rankMap.SQLDB = sqlDB
	rankMap.CacheDB = cntDB
	rankMap.RedisDB = redisDB

	rankRedisDB = redisDB
	go timelyResort(sleeps)
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
	// var tmpData []ArticleRankItem
	m := GetRankMap()
	if _, ok := m.m[cid]; !ok {
		return retData
	}

	crd := m.m[cid] // categoryRankData

	func() {
		crd.mtx.Lock()
		defer crd.mtx.Unlock()
		start := (page - 1) * limit
		end := (page) * limit
		// maxIDx := uint64(len(crd.topicData))
		// start = min(start, maxIDx)
		// end = min(end, maxIDx)
		// tmpData = append(tmpData, crd.topicData[start:end]...)
		data, _ := rankRedisDB.ZRevRange(fmt.Sprintf("%d", cid), int64(start), int64(end)).Result()
		for _, val := range data {
			// fmt.Println("Get val", val)
			aid, _ := strconv.ParseUint(val, 10, 64)
			retData = append(retData, aid)
		}
		// fmt.Printf("%p %p", tmpData, crd.topicData[start:end])
		// fmt.Println("Get ret data", retData)
	}()
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
