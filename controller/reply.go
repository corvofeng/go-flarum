package controller

import (
	"encoding/json"
	"goyoubbs/model/flarum"
	"goyoubbs/util"
	"html/template"
	"net/http"
)

// ContentPreviewPost 预览主题以及评论
func (h *BaseHandler) ContentPreviewPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	token := h.GetCookie(r, "token")
	if len(token) == 0 {
		w.Write([]byte(`{"retcode":400,"retmsg":"token cookie missed"}`))
		return
	}

	currentUser, _ := h.CurrentUser(w, r)
	if !currentUser.CanCreateTopic() || !currentUser.CanReply() {
		w.Write([]byte(`{"retcode":403,"retmsg":"forbidden"}`))
		return
	}

	type recForm struct {
		Act     string `json:"act"`
		Link    string `json:"link"`
		Content string `json:"content"`
	}

	type response struct {
		normalRsp
		Content string        `json:"content"`
		Html    template.HTML `json:"html"`
	}

	decoder := json.NewDecoder(r.Body)
	var rec recForm
	err := decoder.Decode(&rec)
	if err != nil {
		w.Write([]byte(`{"retcode":400,"retmsg":"json Decode err:` + err.Error() + `"}`))
		return
	}
	defer r.Body.Close()

	// db := h.App.Db
	rsp := response{}

	if rec.Act == "preview" && len(rec.Content) > 0 {
		rsp.Retcode = 200
		rsp.Html = template.HTML(util.ContentFmt(rec.Content))
	}
	json.NewEncoder(w).Encode(rsp)

}

// FlarumAPICreatePost flarum进行评论的接口
func FlarumAPICreatePost(w http.ResponseWriter, r *http.Request) {
	ctx := GetRetContext(r)
	h := ctx.h

	rsp := response{}
	decoder := json.NewDecoder(r.Body)
	post := flarum.NewResource(flarum.EPost, 0)
	err := decoder.Decode(&post)
	if err != nil {
		rsp = response{400, "json Decode err:" + err.Error()}
		h.jsonify(w, rsp)
		return
	}
	defer r.Body.Close()
	rsp.Retcode = 404
	w.WriteHeader(http.StatusBadGateway)
	h.jsonify(w, rsp)

	// if rec.Act == "preview" && len(rec.Content) > 0 {
	// 	rsp.Retcode = 200
	// 	rsp.HTML = template.HTML(util.ContentFmt(rec.Content))
	// }
	// json.NewEncoder(w).Encode(rsp)

}
