package model

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/ego008/youdb"
)

type Link struct {
	Id    uint64 `json:"id"`
	Name  string `json:"name"`
	Url   string `json:"url"`
	Score int    `json:"score"`
}

func LinkGetById(db *youdb.DB, lid string) Link {
	var item Link
	rs := db.Hget("link", youdb.DS2b(lid))
	if rs.State == "ok" {
		json.Unmarshal(rs.Data[0], &item)
	}
	return item
}

func LinkSet(db *youdb.DB, obj Link) {
	if obj.Id == 0 {
		// add
		newId, _ := db.HnextSequence("link")
		obj.Id = newId
	}
	jb, _ := json.Marshal(obj)
	db.Hset("link", youdb.I2b(obj.Id), jb)
}

func SqlLinkList(db *sql.DB, getAll bool) []Link {
	var items []Link
	// itemMap := map[uint64]Link{}
	rows, err := db.Query("SELECT id, title FROM topic limit 20")
	defer func() {
		if rows != nil {
			rows.Close() //可以关闭掉未scan连接一直占用
		}
	}()
	if err != nil {
		fmt.Printf("Query failed,err:%v", err)
		return items
	}
	for rows.Next() {
		item := Link{}

		// Id    uint64 `json:"id"`
		// Name  string `json:"name"`
		// Url   string `json:"url"`
		err = rows.Scan(&item.Id, &item.Name) //不scan会导致连接不释放
		if err != nil {
			fmt.Printf("Scan failed,err:%v", err)
			return items
		}
		fmt.Print(item)
		items = append(items, item)
	}

	return items
}

func LinkList(db *youdb.DB, getAll bool) []Link {
	var items []Link
	itemMap := map[uint64]Link{}

	startKey := []byte("")

	for {
		rs := db.Hscan("link", startKey, 20)
		if rs.State == "ok" {
			for i := 0; i < len(rs.Data)-1; i += 2 {
				startKey = rs.Data[i]
				item := Link{}
				json.Unmarshal(rs.Data[i+1], &item)
				if getAll {
					// included score == 0
					itemMap[youdb.B2i(rs.Data[i])] = item
				} else {
					if item.Score > 0 {
						itemMap[youdb.B2i(rs.Data[i])] = item
					}
				}
			}
		} else {
			break
		}
	}

	if len(itemMap) > 0 {
		type Kv struct {
			Key   uint64
			Value int
		}

		var ss []Kv
		for k, v := range itemMap {
			ss = append(ss, Kv{k, v.Score})
		}

		sort.Slice(ss, func(i, j int) bool {
			return ss[i].Value > ss[j].Value
		})

		for _, kv := range ss {
			items = append(items, itemMap[kv.Key])
		}
	}

	return items
}
