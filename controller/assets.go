package controller

import (
	"encoding/json"
	"fmt"
	"goyoubbs/util"
	"net/http"
	"path"
)

// GetLocaleData 获取地区对应的语言包
func (h *BaseHandler) GetLocaleData(w http.ResponseWriter, r *http.Request) {
	localeDir := path.Join(h.App.Cf.Main.ViewDir, "..", "locale")
	localeDataArr := util.FlarumReadLocale(localeDir, "en")

	arr, _ := json.Marshal(localeDataArr)

	retData := fmt.Sprintf(`flarum.core.app.translator.addTranslations(%s);`, arr)
	w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=31536000")

	w.Write([]byte(retData))
}
