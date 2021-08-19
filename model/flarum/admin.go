package flarum

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path"
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

//       "dist": {
//         "type": "zip",
//         "url": "https://api.github.com/repos/flarum/emoji/zipball/f2eb25e9383dd05f9905d2ea087f2e997aae2bb0",
//         "reference": "f2eb25e9383dd05f9905d2ea087f2e997aae2bb0",
//         "shasum": ""
//       },
func ReadExtensionMetadata(extionsDir string) (map[string]FlarumAdminExtension, error) {
	getIcon := func() map[string]string {
		return map[string]string{
			"name":            "fas fa-flag",
			"backgroundColor": "#D659B5",
			"color":           "#fff",
		}
	}

	exts := make(map[string]FlarumAdminExtension)
	items, err := ioutil.ReadDir(extionsDir)
	if util.CheckError(err, "读取配置") {
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
			var extData map[string]interface{}
			_ = json.Unmarshal([]byte(data), &extData)
			// 由于admin页面类似锚点路由, 因此这里需要把/转换一次
			extName := strings.ReplaceAll(extData["name"].(string), "/", "-")
			extData["id"] = extName
			extData["links"] = map[string]interface{}{
				"source":  "https://github.com/flarum/emoji.git",
				"discuss": "https://discuss.flarum.org",
				"website": "https://flarum.org",
				"donate":  "https://flarum.org/donate/",
				"authors": []string{}, // 必须给定authors
			}
			extData["path"] = "/flarum/app/vendor/composer/../flarum/likes" // 必须增加path变量
			extData["icon"] = getIcon()

			// 下面可加可不加
			extData["install-path"] = "../flarum/flags"
			extData["type"] = "flarum-extension"
			extData["installation-source"] = "dist"
			extData["hasAssets"] = false
			extData["hasMigrations"] = true
			extData["source"] = map[string]interface{}{
				"type":      "git",
				"url":       "https://github.com/flarum/emoji.git",
				"reference": "f2eb25e9383dd05f9905d2ea087f2e997aae2bb0",
			}
			extData["time"] = "2021-05-25T20:45:40+00:00"
			extData["version_normalized"] = "1.0.0.0"
			extData["version"] = "v1.0.0"
			extData["extensionDependencyIds"] = []string{}
			extData["optionalDependencyIds"] = []string{}
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
