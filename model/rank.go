package model

import (
	"sort"
	"sync"
	"time"
)

// ArticleRankItem 记录每个话题的权重
type ArticleRankItem struct {
	AID    uint64 `json:"a_id"`
	Weight uint64
}

// CategoryRankData 一个分类下的排序数据
type CategoryRankData struct {
	CID       uint64     `json:"c_id"`
	mtx       sync.Mutex // 同一时刻, 只允许一个协程操纵该分类的记录
	maxID     uint64     // 数据库游标, 记录当前已读取数据的最大值
	topicData []ArticleRankItem
}

// RankMap time to live map
type RankMap struct {
	m   map[uint64]*CategoryRankData
	mtx sync.Mutex // 同一时刻, 只允许一个协程操纵map
}

var rankMap *RankMap

func (data *CategoryRankData) resort() {
	// fmt.Println("In sort ", data.CID, data.topicData)
	func() {
		data.mtx.Lock()
		defer data.mtx.Unlock()
		// Sort by age, keeping original order or equal elements.
		sort.SliceStable(data.topicData, func(i, j int) bool {
			return data.topicData[i].Weight >= data.topicData[j].Weight
		})
	}()
	// fmt.Println("After sort ", data.CID, data.topicData)

	var maxID uint64
	for _, d := range data.topicData {
		if d.AID > maxID {
			maxID = d.AID
		}
	}
	func() {
		data.mtx.Lock()
		defer data.mtx.Unlock()
		data.maxID = maxID
	}()
}

func timelyResort(sleeps uint64) {
	m := GetRankMap()
	for _ = range time.Tick(time.Second * time.Duration(sleeps)) { // 每10s刷新一次
		expireItems := []*CategoryRankData{}
		func() {
			m.mtx.Lock()
			defer m.mtx.Unlock()
			for _, v := range m.m {
				expireItems = append(expireItems, v) // 效率低下
			}
		}()

		for _, item := range expireItems {
			t := item
			go t.resort()
		}
	}
}

func newRankMap() (m *RankMap) {
	m = &RankMap{m: make(map[uint64]*CategoryRankData)}
	return m
}

// RankMapInit init a ttl map
func RankMapInit(sleeps uint64) {
	rankMap = newRankMap()
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
	var tmpData []ArticleRankItem
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
		maxIdx := uint64(len(crd.topicData))
		start = min(start, maxIdx)
		end = min(end, maxIdx)
		tmpData = append(tmpData, crd.topicData[start:end]...)
		// fmt.Printf("%p %p", tmpData, crd.topicData[start:end])
	}()

	for _, d := range tmpData {
		retData = append(retData, d.AID)
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
	// fmt.Print("Get data ", cid, rankItems)

	crd := m.m[cid] // categoryRankData
	func() {
		crd.mtx.Lock()
		defer crd.mtx.Unlock()
		crd.topicData = append(crd.topicData, rankItems...) // 直接加入, 不做任何处理
	}()
}

// GetCidArticleMax 获取当前分类的偏移值
func GetCidArticleMax(cid uint64) uint64 {
	m := GetRankMap()
	if _, ok := m.m[cid]; ok {
		return m.m[cid].maxID
	}
	return 0
}