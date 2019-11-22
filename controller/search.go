package controller

import (
	"encoding/json"
	"html/template"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"goyoubbs/model"
)

// SearchDetail 按关键字查找
func (h *BaseHandler) SearchDetail(w http.ResponseWriter, r *http.Request) {
	currentUser, _ := h.CurrentUser(w, r)
	btn, pn := r.FormValue("btn"), r.FormValue("pagenum")
	// if currentUser.ID == 0 {
	// 	http.Redirect(w, r, "/login", http.StatusSeeOther)
	// 	return
	// }
	var pagenum uint64
	var err error

	if len(pn) > 0 {
		pagenum, err = strconv.ParseUint(pn, 10, 64)
		if err != nil {
			w.Write([]byte(`{"retcode":400,"retmsg":"key type err"}`))
			return
		}
	}
	if btn == "" || btn == "next" {
		pagenum++
	} else if btn == "prev" {
		if pagenum > 0 {
			pagenum--
		}
	} else {
		w.Write([]byte(`{"retcode":400,"retmsg":"btn err"}`))
		return
	}

	q := r.FormValue("q")

	if len(q) == 0 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	qLow := strings.ToLower(q)

	db := h.App.Db
	logger := h.App.Logger
	scf := h.App.Cf.Site
	sqlDB := h.App.MySQLdb
	redisDB := h.App.RedisDB

	var body []byte
	resp, err := http.Get(
		fmt.Sprintf("http://172.17.0.1:9192?query=%s&pagenum=%d&pagelen=%d",
			q, pagenum,
			scf.HomeShowNum,
		),
	)
	if err != nil {
		logger.Errorf("Search %s with error: %s ", q, err)
	} else {
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Error("parse body error")
		}
	}
	data := struct {
		Items []struct {
			ID      int    `json:"id"`
			Title   string `json:"title"`
			Content string `json:"content"`
		} `json:"items"`
		IsLastPage bool   `json:"is_last_page"`
		PageNum    uint64 `json:"pagenum"`
		PageLen    int    `json:"pagelen"`
	}{}
	json.Unmarshal(body, &data)

	articleList := make([]uint64, len(data.Items))
	for _, item := range data.Items {
		articleList = append(articleList, uint64(item.ID))
	}
	// Although we get data from the search API, it is necessary to
	// check the data in the database, also we need get article title.
	pageInfo := model.SQLArticleGetByList(sqlDB, db, redisDB, articleList, scf.TimeZone)

	for _, item := range data.Items {
		for idx := range pageInfo.Items {
			if uint64(item.ID) == pageInfo.Items[idx].ID {
				pageInfo.Items[idx].HighlightContent = template.HTML(item.Content)
			}
		}
	}

	type pageData struct {
		PageData
		Q        string
		PageInfo model.ArticlePageInfo
	}
	if data.PageNum > 1 {
		pageInfo.HasPrev = true
	}
	if !data.IsLastPage {
		pageInfo.HasNext = true
	}
	pageInfo.PageNum = data.PageNum

	tpl := h.CurrentTpl(r)

	evn := &pageData{}
	evn.SiteCf = scf
	evn.Title = qLow + " - " + scf.Name
	evn.IsMobile = tpl == "mobile"

	evn.CurrentUser = currentUser
	evn.ShowSideAd = true
	evn.PageName = "category_detail"
	evn.HotNodes = model.CategoryHot(db, scf.CategoryShowNum)
	evn.NewestNodes = model.CategoryNewest(db, scf.CategoryShowNum)

	evn.Q = qLow
	evn.PageInfo = pageInfo
	evn.SiteInfo = model.GetSiteInfo(redisDB, db)

	h.Render(w, tpl, evn, "layout.html", "search.html")
}
