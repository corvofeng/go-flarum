package flarum

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"strings"
	"zoe/util"
)

type FlarumAdminSettings struct {
	AllowPostEditing                  string `json:"allow_post_editing"`
	AllowRenaming                     string `json:"allow_renaming"`
	AllowSignUp                       string `json:"allow_sign_up"`
	CustomLess                        string `json:"custom_less"`
	DefaultLocale                     string `json:"default_locale"`
	DefaultRoute                      string `json:"default_route"`
	ExtensionsEnabled                 string `json:"extensions_enabled"`
	FlarumMarkdownMdarea              string `json:"flarum-markdown.mdarea"`
	FlarumMentionsAllowUsernameFormat string `json:"flarum-mentions.allow_username_format"`
	FlarumTagsMaxPrimaryTags          string `json:"flarum-tags.max_primary_tags"`
	FlarumTagsMaxSecondaryTags        string `json:"flarum-tags.max_secondary_tags"`
	FlarumTagsMinPrimaryTags          string `json:"flarum-tags.min_primary_tags"`
	FlarumTagsMinSecondaryTags        string `json:"flarum-tags.min_secondary_tags"`
	ForumDescription                  string `json:"forum_description"`
	ForumTitle                        string `json:"forum_title"`
	MailDriver                        string `json:"mail_driver"`
	MailFrom                          string `json:"mail_from"`
	ThemeColoredHeader                string `json:"theme_colored_header"`
	ThemeDarkMode                     string `json:"theme_dark_mode"`
	ThemePrimaryColor                 string `json:"theme_primary_color"`
	ThemeSecondaryColor               string `json:"theme_secondary_color"`
	Version                           string `json:"version"`
	WelcomeMessage                    string `json:"welcome_message"`
	WelcomeTitle                      string `json:"welcome_title"`
}
type FlarumAdminPermissions struct {
	ViewForum                      []string `json:"viewForum"`
	DiscussionFlagPosts            []string `json:"discussion.flagPosts"`
	DiscussionLikePosts            []string `json:"discussion.likePosts"`
	DiscussionReply                []string `json:"discussion.reply"`
	DiscussionReplyWithoutApproval []string `json:"discussion.replyWithoutApproval"`
	DiscussionStartWithoutApproval []string `json:"discussion.startWithoutApproval"`
	SearchUsers                    []string `json:"searchUsers"`
	StartDiscussion                []string `json:"startDiscussion"`
	DiscussionApprovePosts         []string `json:"discussion.approvePosts"`
	DiscussionEditPosts            []string `json:"discussion.editPosts"`
	DiscussionHide                 []string `json:"discussion.hide"`
	DiscussionHidePosts            []string `json:"discussion.hidePosts"`
	DiscussionLock                 []string `json:"discussion.lock"`
	DiscussionRename               []string `json:"discussion.rename"`
	DiscussionSticky               []string `json:"discussion.sticky"`
	DiscussionTag                  []string `json:"discussion.tag"`
	DiscussionViewFlags            []string `json:"discussion.viewFlags"`
	DiscussionViewIpsPosts         []string `json:"discussion.viewIpsPosts"`
	UserSuspend                    []string `json:"user.suspend"`
	UserViewLastSeenAt             []string `json:"user.viewLastSeenAt"`
}

// 这个信息存在于每个扩展的composer.json文件中
// 作为json文件读取进来就可以了
type FlarumAdminExtension interface{}

func extGetLink(extData PluginComposer) map[string]interface{} {
	links := make(map[string]interface{})
	for k, v := range extData.Support {
		links[k] = v
	}
	authors := []string{}
	for _, aObj := range extData.Authors {
		authors = append(authors, aObj.Name)
	}
	links["authors"] = authors
	return links
}

type FlarumExtIcon struct {
	Name            string `json:"name"`
	BackgroundColor string `json:"backgroundColor"`
	Color           string `json:"color"`
	Image           string `json:"image"`
	BackgroundImage string `json:"backgroundImage"`

	BackgroundSize     string `json:"backgroundSize"`
	BackgroundRepeat   string `json:"backgroundRepeat"`
	BackgroundPosition string `json:"backgroundPosition"`
}

type PluginComposer struct {
	// 必须额外加入
	ID    string                 `json:"id"`
	Path  string                 `json:"path"`
	Links map[string]interface{} `json:"links"`
	Icon  FlarumExtIcon          `json:"icon"`

	InstallPath        string                 `json:"install-path"`
	InstallationSource string                 `json:"install-source"`
	HasMigrations      bool                   `json:"hasMigrations"`
	HasAssets          bool                   `json:"hasAssets"`
	Source             map[string]interface{} `json:"source"`

	Name        string   `json:"name"`
	Description string   `json:"description"`
	Keywords    []string `json:"keywords"`
	Type        string   `json:"type"`
	License     string   `json:"license"`
	Authors     []struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Homepage string `json:"homepage,omitempty"`
	} `json:"authors"`
	Support map[string]string `json:"support"`
	Require map[string]string `json:"require"`

	Extra struct {
		FlarumExtension struct {
			Title string        `json:"title"`
			Icon  FlarumExtIcon `json:"icon"`
		} `json:"flarum-extension"`
		Flagrow struct {
			Discuss string `json:"discuss"`
		} `json:"flagrow"`
	} `json:"extra"`
	Autoload interface{} `json:"autoload"`
	Replace  interface{} `json:"replace"`
}

func toBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func ReadExtensionMetadata(extionsDir string) (map[string]FlarumAdminExtension, error) {
	exts := make(map[string]FlarumAdminExtension)
	absPath, err := filepath.Abs(extionsDir)
	if util.CheckError(err, "查找扩展绝对路径") {
		return exts, errors.New("无法找到扩展目录")
	}
	items, err := ioutil.ReadDir(absPath)
	if util.CheckError(err, "读取扩展目录") {
		return exts, errors.New("无法找到扩展目录")
	}
	for _, item := range items {
		// 仅认为目录是真正的扩展
		if item.IsDir() {
			data, err := ioutil.ReadFile(path.Join(extionsDir, item.Name(), "composer.json"))
			if util.CheckError(err, fmt.Sprintf("读取扩展%s配置失败", item.Name())) {
				// 某个扩展读取信息失败, 仅打印日志, 不做处理
				continue
			}
			// var extData map[string]interface{}
			var extData PluginComposer
			_ = json.Unmarshal([]byte(data), &extData)

			// 由于admin页面类似锚点路由, 因此这里需要把/转换一次
			extName := strings.ReplaceAll(extData.Name, "/", "-")
			extData.ID = extName

			extData.Links = extGetLink(extData)
			extData.Path = path.Join(absPath, item.Name()) // 必须增加path变量
			extData.Icon = extData.Extra.FlarumExtension.Icon

			if extData.Icon.Image != "" {
				// 参考flarum/src/Extension/Extension.php, 读取文件作为图标
				bytes, err := ioutil.ReadFile(path.Join(extData.Path, extData.Icon.Image))
				if err != nil {
					log.Fatal(err)
				}
				if util.CheckError(err, "读取扩展的图标") {
					continue
				}
				LOGO_MIMETYPES := map[string]string{
					".svg":  "image/svg+xml",
					".png":  "image/png",
					".jpeg": "image/jpeg",
					".jpg":  "image/jpeg",
				}
				mimeType := LOGO_MIMETYPES[filepath.Ext(extData.Icon.Image)]
				base64Encoding := toBase64(bytes)
				extData.Icon.BackgroundImage = fmt.Sprintf("url('data:%s;base64,%s')", mimeType, base64Encoding)
			}

			// 下面可加可不加
			extData.InstallPath = path.Join(absPath, item.Name()) // 必须增加path变量
			extData.InstallationSource = "dist"
			extData.HasMigrations = true
			extData.HasAssets = false

			extData.Source = map[string]interface{}{
				// "type":      "git",
				// "url":       "https://github.com/flarum/emoji.git",
				// "reference": "f2eb25e9383dd05f9905d2ea087f2e997aae2bb0",
			}
			// extData["time"] = "2021-05-25T20:45:40+00:00"
			// extData["version_normalized"] = "1.0.0.0"
			// extData["version"] = "v1.0.0"
			// extData["extensionDependencyIds"] = []string{}
			// extData["optionalDependencyIds"] = []string{}
			exts[extName] = extData
		}
	}
	return exts, nil
}

type AdminCoreData struct {
	CoreData
	Settings    FlarumAdminSettings    `json:"settings"`
	Permissions FlarumAdminPermissions `json:"permissions"`

	PhpVersion   string `json:"phpVersion"`
	MysqlVersion string `json:"mysqlVersion"`

	Extensions map[string]FlarumAdminExtension `json:"extensions"`

	DisplayNameDrivers []string `json:"displayNameDrivers"`

	SlugDrivers map[string]interface{} `json:"slugDrivers"`

	ModelStatistics struct {
		Users struct {
			Total int `json:"total"`
		} `json:"users"`
	} `json:"modelStatistics"`
}

type AutoGenerated struct {
	SlugDrivers struct {
		// FlarumDiscussionDiscussion []string `json:"Flarum\Discussion\Discussion"`
		// FlarumUserUser             []string `json:"Flarum\User\User"`
	} `json:"slugDrivers"`

	Statistics struct {
		Users struct {
			Total int `json:"total"`
			Timed struct {
				Num1590624000 int `json:"1590624000"`
				Num1590883200 int `json:"1590883200"`
				Num1591747200 int `json:"1591747200"`
			} `json:"timed"`
		} `json:"users"`
		Discussions struct {
			Total int `json:"total"`
			Timed struct {
				Num1590624000 int `json:"1590624000"`
				Num1590883200 int `json:"1590883200"`
				Num1592265600 int `json:"1592265600"`
			} `json:"timed"`
		} `json:"discussions"`
		Posts struct {
			Total int `json:"total"`
			Timed struct {
				Num1590624000 int `json:"1590624000"`
				Num1590883200 int `json:"1590883200"`
				Num1591747200 int `json:"1591747200"`
				Num1592265600 int `json:"1592265600"`
				Num1592697600 int `json:"1592697600"`
				Num1593388800 int `json:"1593388800"`
				Num1594080000 int `json:"1594080000"`
				Num1602201600 int `json:"1602201600"`
			} `json:"timed"`
		} `json:"posts"`
		TimezoneOffset int `json:"timezoneOffset"`
	} `json:"statistics"`
}
