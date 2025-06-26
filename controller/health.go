package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/corvofeng/go-flarum/model"
)

func FlarumHealthFunc(w http.ResponseWriter, r *http.Request) {
	started := time.Now()
	ctx := GetRetContext(r)
	h := ctx.h
	redisDB := h.App.RedisDB
	model.GetSiteInfo(redisDB)
	gormDB := h.App.GormDB
	gormDB.Exec("SELECT 1") // Test the database connection

	model.SQLGetTagByUrlName(gormDB, "r_root")

	duration := time.Since(started)
	if duration.Seconds() > 10 {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("error: %v", duration.Seconds())))
	} else {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	}
}
