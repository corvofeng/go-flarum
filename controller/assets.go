package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"

	"github.com/corvofeng/go-flarum/util"

	"goji.io/pat"
)

// GetLocaleData 获取地区对应的语言包
func (h *BaseHandler) GetLocaleData(w http.ResponseWriter, r *http.Request) {
	localeDir := path.Join(h.App.Cf.Main.LocaleDir)
	extDir := path.Join(h.App.Cf.Main.ExtensionsDir)
	flarumDir := path.Join(h.App.Cf.Main.ExtensionsDir, "..", "framework")
	extDirs := []string{
		extDir,
		path.Join(flarumDir, "extensions"),
	}
	locale := pat.Param(r, "locale")
	localeDataArr := util.FlarumReadLocale(
		path.Join(flarumDir, "framework", "core"),
		extDirs, localeDir, locale)
	arr, _ := json.Marshal(localeDataArr)
	retData := fmt.Sprintf(`flarum.core.app.translator.addTranslations(%s);`, arr)
	w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=31536000")

	w.Write([]byte(retData))
}
