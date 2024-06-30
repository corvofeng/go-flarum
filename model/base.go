package model

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v7"
	"gorm.io/gorm"
)

// ISQLLoader dict结果的loader
type ISQLLoader interface {
	LoadDictData(map[string]interface{})
}

// rowsClose scan没有结束或是没有进行scan操作时, 需要手动释放连接
func rowsClose(rows *sql.Rows) {
	if rows != nil {
		rows.Close()
	}
}

// clearGormTransaction gorm的事务清理
// 使用方式:
// tx := gormDB.Begin()
// defer clearGormTransaction(tx)
func clearGormTransaction(tx *gorm.DB) {
	if err := recover(); err != nil {
		tx.Rollback()
	}
}

func clearTransaction(tx *sql.Tx) {
	err := tx.Rollback()
	if err != sql.ErrTxDone && err != nil {
		fmt.Println("error in transaction", err)
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
