package model

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"goyoubbs/util"
	"strconv"
	"strings"

	"github.com/ego008/youdb"
	"github.com/go-redis/redis/v7"
)

// Category 帖子分类
type Category struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	Articles uint64 `json:"articles"`
	About    string `json:"about"`
	Hidden   bool   `json:"hidden"`
}

type CategoryMini struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type CategoryPageInfo struct {
	Items    []Category `json:"items"`
	HasPrev  bool       `json:"hasprev"`
	HasNext  bool       `json:"hasnext"`
	FirstKey uint64     `json:"firstkey"`
	LastKey  uint64     `json:"lastkey"`
}

// SQLGetAllCategory 获取所有分类
func SQLGetAllCategory(db *sql.DB) ([]CategoryMini, error) {
	var categories []CategoryMini
	rows, err := db.Query("SELECT id, name FROM node order by topic_count desc limit 30")
	defer func() {
		if rows != nil {
			rows.Close() // 可以关闭掉未scan连接一直占用
		}
	}()
	if err != nil {
		fmt.Printf("Query failed,err:%v", err)
	}
	for rows.Next() {
		obj := CategoryMini{}
		err = rows.Scan(&obj.ID, &obj.Name) // 不scan会导致连接不释放

		if err != nil {
			fmt.Printf("Scan failed,err:%v", err)
			return categories, errors.New("No result")
		}
		categories = append(categories, obj)
	}

	return categories, nil
}

// SQLCategoryGetByID 通过id获取节点
func SQLCategoryGetByID(db *sql.DB, cid string) (Category, error) {
	return sqlCategoryGet(db, cid, "")
}

// SQLCategoryGetByName 通过name获取节点
func SQLCategoryGetByName(db *sql.DB, name string) (Category, error) {
	return sqlCategoryGet(db, "", name)
}

func sqlCategoryGet(db *sql.DB, cid string, name string) (Category, error) {
	obj := Category{}
	var rows *sql.Rows
	var err error

	if cid != "" {
		rows, err = db.Query("SELECT id, name, summary, topic_count FROM node WHERE id =  ?", cid)
	} else if name != "" {
		rows, err = db.Query("SELECT id, name, summary, topic_count FROM node WHERE name =  ?", name)
	} else {
		return obj, errors.New("Did not give any category")
	}

	defer func() {
		if rows != nil {
			rows.Close() //可以关闭掉未scan连接一直占用
		}
	}()
	if err != nil {
		fmt.Printf("Query failed,err:%v", err)
	}
	for rows.Next() {
		err = rows.Scan(&obj.ID, &obj.Name, &obj.About, &obj.Articles) //不scan会导致连接不释放

		if err != nil {
			fmt.Printf("Scan failed,err:%v", err)
			return obj, errors.New("No result")
		}
	}

	return obj, nil
}

// SQLCategoryList 获取分类列表
func SQLCategoryList(db *sql.DB) CategoryPageInfo {
	// tb := "category"
	var items []Category
	// var keys [][]byte
	var hasPrev, hasNext bool
	var firstKey, lastKey uint64

	categrayList, err := SQLGetAllCategory(db)
	if !util.CheckError(err, "获取所有category") {
		for _, cate := range categrayList {
			fullCategory, err := SQLCategoryGetByName(db, cate.Name)
			if !util.CheckError(err, fmt.Sprintf("获取category: %s", cate.Name)) {
				items = append(items, fullCategory)
			}
		}
	}

	return CategoryPageInfo{
		Items:    items,
		HasPrev:  hasPrev,
		HasNext:  hasNext,
		FirstKey: firstKey,
		LastKey:  lastKey,
	}

}

//  以下代码不再使用
// ============================================================ //

func CategoryGetByID(db *youdb.DB, cid string) (Category, error) {
	obj := Category{}
	rs := db.Hget("category", youdb.DS2b(cid))
	if rs.State == "ok" {
		json.Unmarshal(rs.Data[0], &obj)
		return obj, nil
	}
	return obj, errors.New(rs.State)
}

// GetCategoryNameByCID 通过CID获取该分类的名称
func GetCategoryNameByCID(sqlDB *sql.DB, redisDB *redis.Client, cid uint64) string {
	var cname string
	logger := util.GetLogger()
	rep, err := redisDB.HGet("category", fmt.Sprintf("%d", cid)).Result()
	if err != redis.Nil {
		return rep
	}
	category, err := SQLCategoryGetByID(sqlDB, fmt.Sprintf("%d", cid))
	cname = category.Name
	redisDB.HSet("category", fmt.Sprintf("%d", cid), cname)

	logger.Debugf("category not found for %d %s but we refresh!", category.ID, category.Name)
	return cname
}

func CategoryHot(db *youdb.DB, limit int) []CategoryMini {
	var items []CategoryMini
	rs := db.Zrscan("category_article_num", []byte(""), []byte(""), limit)
	if rs.State == "ok" {
		var keys [][]byte
		for i := 0; i < len(rs.Data)-1; i += 2 {
			keys = append(keys, rs.Data[i])
		}
		if len(keys) > 0 {
			rs2 := db.Hmget("category", keys)
			if rs2.State == "ok" {
				for i := 0; i < len(rs2.Data)-1; i += 2 {
					item := CategoryMini{}
					json.Unmarshal(rs2.Data[i+1], &item)
					items = append(items, item)
				}
			}
		}
	}
	return items
}

func CategoryNewest(db *youdb.DB, limit int) []CategoryMini {
	var items []CategoryMini
	rs := db.Hrscan("category", []byte(""), limit)
	if rs.State == "ok" {
		for i := 0; i < len(rs.Data)-1; i += 2 {
			item := CategoryMini{}
			json.Unmarshal(rs.Data[i+1], &item)
			items = append(items, item)
		}
	}
	return items
}

func CategoryGetMain(db *youdb.DB, currentCobj Category) []CategoryMini {
	var items []CategoryMini
	item := CategoryMini{
		ID:   currentCobj.ID,
		Name: currentCobj.Name,
	}
	items = append(items, item)

	rs := db.Hget("keyValue", []byte("main_category"))
	if rs.State == "ok" {
		var keys [][]byte
		currentCIDStr := strconv.FormatUint(currentCobj.ID, 10)
		cids := strings.Split(rs.String(), ",")
		for _, v := range cids {
			if v != currentCIDStr {
				keys = append(keys, youdb.DS2b(v))
			}
		}
		if len(keys) > 0 {
			rs2 := db.Hmget("category", keys)
			if rs2.State == "ok" {
				for i := 0; i < len(rs2.Data)-1; i += 2 {
					item := CategoryMini{}
					json.Unmarshal(rs2.Data[i+1], &item)
					items = append(items, item)
				}
			}
		}
	}

	return items
}

func CategoryList(db *youdb.DB, cmd, key string, limit int) CategoryPageInfo {
	tb := "category"
	var items []Category
	var keys [][]byte
	var hasPrev, hasNext bool
	var firstKey, lastKey uint64

	keyStart := youdb.DS2b(key)
	if cmd == "hrscan" {
		rs := db.Hrscan(tb, keyStart, limit)
		if rs.State == "ok" {
			for i := 0; i < (len(rs.Data) - 1); i += 2 {
				keys = append(keys, rs.Data[i])
			}
		}
	} else if cmd == "hscan" {
		rs := db.Hscan(tb, keyStart, limit)
		if rs.State == "ok" {
			for i := len(rs.Data) - 2; i >= 0; i -= 2 {
				keys = append(keys, rs.Data[i])
			}
		}
	}

	if len(keys) > 0 {
		rs := db.Hmget("category", keys)
		if rs.State == "ok" {
			for i := 0; i < (len(rs.Data) - 1); i += 2 {
				item := Category{}
				json.Unmarshal(rs.Data[i+1], &item)
				items = append(items, item)
				if firstKey == 0 {
					firstKey = item.ID
				}
				lastKey = item.ID
			}

			rs = db.Hscan(tb, youdb.I2b(firstKey), 1)
			if rs.State == "ok" {
				hasPrev = true
			}
			rs = db.Hrscan(tb, youdb.I2b(lastKey), 1)
			if rs.State == "ok" {
				hasNext = true
			}
		}
	}

	return CategoryPageInfo{
		Items:    items,
		HasPrev:  hasPrev,
		HasNext:  hasNext,
		FirstKey: firstKey,
		LastKey:  lastKey,
	}
}
