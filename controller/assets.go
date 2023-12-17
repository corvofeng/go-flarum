package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"zoe/util"

	"goji.io/pat"
)

// GetLocaleData 获取地区对应的语言包
func (h *BaseHandler) GetLocaleData(w http.ResponseWriter, r *http.Request) {
	localeDir := path.Join(h.App.Cf.Main.LocaleDir)
	extDir := path.Join(h.App.Cf.Main.ExtensionsDir)
	flarumDir := path.Join(h.App.Cf.Main.ExtensionsDir, "..", "flarum")
	locale := pat.Param(r, "locale")
	if locale == "en" {
		// 由于最新的flarum中将英文的字段全部放在了每个组件中, 因此简单先进行proxy
		jsURL := "https://discuss.flarum.org/assets/forum-en.js"
		// 从远程获取文件
		resp, err := http.Get(jsURL)
		if err != nil {
			h.App.Logger.Errorf("Failed to fetch JS file: %v", err)
			return
		}
		defer resp.Body.Close()

		jsBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			h.App.Logger.Errorf("Failed to read JS file: %v", err)
			return
		}
		w.Header().Set("Content-Type", "application/javascript")
		w.Write(jsBytes)
	} else {
		localeDataArr := util.FlarumReadLocale(flarumDir, extDir, localeDir, locale)
		arr, _ := json.Marshal(localeDataArr)
		retData := fmt.Sprintf(`flarum.core.app.translator.addTranslations(%s);`, arr)
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		w.Header().Set("Cache-Control", "public, max-age=31536000")

		w.Write([]byte(retData))
	}
}
