package model

import (
	"database/sql"
	"errors"
	"fmt"
	"goyoubbs/model/flarum"
	"goyoubbs/util"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v7"
)

type (
	// Category 帖子分类
	Category = struct {
		ID          uint64 `json:"id"`
		Name        string `json:"name"`
		URLName     string `json:"urlname"`
		Articles    uint64 `json:"articles"`
		About       string `json:"about"`
		ParentID    uint64 `json:"parent_id"`
		Description string `json:"description"`
		Hidden      bool   `json:"hidden"`
	}

	// CategoryMini 帖子分类
	CategoryMini struct {
		ID   uint64 `json:"id"`
		Name string `json:"name"`
	}

	// CategoryPageInfo 显示所有的帖子信息
	CategoryPageInfo struct {
		Items    []Category `json:"items"`
		HasPrev  bool       `json:"hasprev"`
		HasNext  bool       `json:"hasnext"`
		FirstKey uint64     `json:"firstkey"`
		LastKey  uint64     `json:"lastkey"`
	}
)

// SQLGetAllCategory 获取所有分类
func SQLGetAllCategory(db *sql.DB) ([]Category, error) {
	var categories []Category
	rows, err := db.Query("SELECT id, name, urlname, description FROM node order by topic_count desc limit 30")
	defer rowsClose(rows)

	if err != nil {
		fmt.Printf("Query failed,err:%v", err)
	}
	for rows.Next() {
		obj := Category{}
		err = rows.Scan(&obj.ID, &obj.Name, &obj.URLName, &obj.Description) // 不scan会导致连接不释放

		if err != nil {
			fmt.Printf("Scan failed,err:%v", err)
			return categories, errors.New("No result")
		}
		categories = append(categories, obj)
	}

	return categories, nil
}

// SQLGetNotEmptyCategory 获取非空的分类
func SQLGetNotEmptyCategory(db *sql.DB, redisDB *redis.Client) (categories []Category, err error) {
	rows, err := db.Query("SELECT id FROM node where topic_count != 0 and is_hidden = 0")
	defer rowsClose(rows)
	logger := util.GetLogger()

	if err != nil {
		logger.Errorf("Query failed,err:%v", err)
		return
	}
	var categoryList []uint64
	for rows.Next() {
		var item uint64
		err = rows.Scan(&item)
		if err != nil {
			logger.Errorf("Scan failed,err:%v", err)
			continue
		}
		categoryList = append(categoryList, item)
	}
	categories = sqlGetCategoryByList(db, redisDB, categoryList)
	return
}

func sqlGetCategoryByList(db *sql.DB, redisDB *redis.Client, categoryList []uint64) (items []Category) {
	var err error
	var rows *sql.Rows
	var categoryListStr []string
	logger := util.GetLogger()
	if len(categoryList) == 0 {
		logger.Warning("sqlGetCategoryByList: Can't process the category list empty")
		return
	}
	defer rowsClose(rows)
	for _, v := range categoryList {
		categoryListStr = append(categoryListStr, strconv.FormatInt(int64(v), 10))
	}
	qFieldList := []string{
		"id", "name", "urlname",
		"description", "is_hidden", "parent_id",
	}
	sql := fmt.Sprintf("select %s from node where id in (%s)",
		strings.Join(qFieldList, ","),
		strings.Join(categoryListStr, ","))
	rows, err = db.Query(sql)
	if err != nil {
		logger.Errorf("Query failed,err:%v", err)
		return
	}

	for rows.Next() {
		item := Category{}
		err = rows.Scan(
			&item.ID, &item.Name, &item.URLName,
			&item.Description, &item.Hidden,
			&item.ParentID,
		)
		if err != nil {
			logger.Errorf("Scan failed,err:%v", err)
			continue
		}
		items = append(items, item)
	}

	return
}

// SQLCategoryGetByID 通过id获取节点
func SQLCategoryGetByID(db *sql.DB, cid string) (Category, error) {
	return sqlCategoryGet(db, cid, "", "")
}

// SQLCategoryGetByName 通过name获取节点
func SQLCategoryGetByName(db *sql.DB, name string) (Category, error) {
	return sqlCategoryGet(db, "", name, "")
}

// SQLCategoryGetByURLName 通过urlname获取节点
func SQLCategoryGetByURLName(db *sql.DB, urlname string) (Category, error) {
	return sqlCategoryGet(db, "", "", urlname)
}

func sqlCategoryGet(db *sql.DB, cid string, name string, urlname string) (Category, error) {
	obj := Category{}
	var rows *sql.Rows
	var err error
	isAdd := false
	defer rowsClose(rows)

	if cid != "" {
		rows, err = db.Query("SELECT id, name, summary, topic_count FROM node WHERE id =  ?", cid)
	} else if name != "" {
		rows, err = db.Query("SELECT id, name, summary, topic_count FROM node WHERE name =  ?", name)
	} else if urlname != "" {
		rows, err = db.Query("SELECT id, name, summary, topic_count FROM node WHERE urlname =  ?", urlname)
	} else {
		return obj, errors.New("Did not give any category")
	}

	if err != nil {
		fmt.Printf("Query failed,err:%v", err)
	}

	for rows.Next() {
		isAdd = true
		err = rows.Scan(&obj.ID, &obj.Name, &obj.About, &obj.Articles) //不scan会导致连接不释放
		if err != nil {
			fmt.Printf("Scan failed,err:%v", err)
			return obj, errors.New("No result")
		}
	}
	if !isAdd {
		return obj, errors.New("No result")
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

// FlarumGetAllTags flarum下获取分类列表
func FlarumGetAllTags(db *sql.DB) []flarum.Tag {
	var tags []flarum.Tag

	allTags := executeQuery(
		db,
		"SELECT n.id, name, urlname, description, topic_count,position, last_posted_topic_id, nn.parent_id FROM ( SELECT * FROM `node` ) AS n LEFT JOIN nodenode as nn ON n.id = nn.child_id",
	)
	for _, t := range allTags {
		fmt.Println(t, t["id"])
	}

	// var keys [][]byte
	// var hasPrev, hasNext bool
	// var firstKey, lastKey uint64

	// categrayList, err := SQLGetAllCategory(db)
	// if !util.CheckError(err, "获取所有category") {
	// 	for _, cate := range categrayList {
	// 		fullCategory, err := SQLCategoryGetByName(db, cate.Name)
	// 		if !util.CheckError(err, fmt.Sprintf("获取category: %s", cate.Name)) {
	// 			items = append(items, fullCategory)
	// 		}
	// 	}
	// }

	// return CategoryPageInfo{
	// 	Items:    items,
	// 	HasPrev:  hasPrev,
	// 	HasNext:  hasNext,
	// 	FirstKey: firstKey,
	// 	LastKey:  lastKey,
	// }

	return tags
}
