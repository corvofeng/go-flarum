package model

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"goyoubbs/util"

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
