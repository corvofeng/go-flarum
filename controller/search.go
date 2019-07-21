package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"../model"
)

func (h *BaseHandler) SearchDetail(w http.ResponseWriter, r *http.Request) {
	currentUser, _ := h.CurrentUser(w, r)
	// if currentUser.Id == 0 {
	// 	http.Redirect(w, r, "/login", http.StatusSeeOther)
	// 	return
	// }

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

	// where := "title"
	// if strings.HasPrefix(qLow, "c:") {
	// 	where = "content"
	// 	qLow = qLow[2:]
	// }

	resp, err := http.Get("http://127.0.0.1:9192?query=" + q)
	if err != nil {
		logger.Error("make get error" + q)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("parse body error")
	}
	data := struct {
		Items []struct {
			ID      int    `json:"id"`
			Title   string `json:"title"`
			Content string `json:"content"`
		} `json:"items"`
	}{}
	json.Unmarshal(body, &data)
	articleList := make([]int, len(data.Items))
	for _, item := range data.Items {
		articleList = append(articleList, item.ID)
	}

	// pageInfo := model.ArticleSearchList(db, where, qLow, scf.PageShowNum, scf.TimeZone)
	pageInfo := model.SQLArticleGetByList(sqlDB, db, articleList)

	type pageData struct {
		PageData
		Q        string
		PageInfo model.ArticlePageInfo
	}

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

	h.Render(w, tpl, evn, "layout.html", "search.html")
}
