package model

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"goyoubbs/util"

	"github.com/go-redis/redis/v7"
)

// ISQLLoader dict结果的loader
type ISQLLoader interface {
	LoadDictData(map[string]interface{})
}

// executeQuery 执行SQL语句, 以dictCursor的形式返回数据
/*
 * 使用时要注意:
 *   1. 一次取特别多数据会很影响性能, 并且会占用较多的空间
 *   2. 此函数只能用于查询语句
 */
func executeQuery(db *sql.DB, query string, args ...interface{}) []map[string]interface{} {
	var err error
	var rows *sql.Rows
	var dictData []map[string]interface{}

	rows, err = db.Query(query, args...)
	defer func() {
		if rows != nil {
			rows.Close() // 未scan, 连接会一直占用, 需要关闭
		}
	}()

	if util.CheckError(err, fmt.Sprintf("Query (%s) %+v failed", query, args)) {
		return dictData
	}

	dictData, err = dataGetByRows(rows)
	if err != nil {
		return dictData
	}
	return dictData
}

// rowsClose scan没有结束或是没有进行scan操作时, 需要手动释放连接
func rowsClose(rows *sql.Rows) {
	if rows != nil {
		rows.Close()
	}
}

// dataGetByRows 从数据库返回结果中获取数据, 解析成为dict的形式
func dataGetByRows(rows *sql.Rows) ([]map[string]interface{}, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	size := len(columns)
	var obj []map[string]interface{}

	colData := make([]interface{}, size)
	container := make([]interface{}, size)
	for i := range colData {
		colData[i] = &container[i]
	}

	for rows.Next() {
		err := rows.Scan(colData...)
		if err != nil {
			return nil, err
		}
		var r = make(map[string]interface{}, size)
		for i, column := range columns {
			r[column] = colData[i]
		}

		obj = append(obj, r)
	}

	return obj, nil
}

func rSet(redisDB *redis.Client, bucket string, key string, value interface{}) error {
	p, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return redisDB.HSet(bucket, key, p).Err()
}
func rDel(redisDB *redis.Client, bucket, key string) error {
	return redisDB.HDel(bucket, key).Err()
}

func rGet(redisDB *redis.Client, bucket string, key string, dest interface{}) error {
	p, err := redisDB.HGet(bucket, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(p), dest)
}
