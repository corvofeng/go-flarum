package model

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/corvofeng/go-flarum/util"
)

type (
	// xml to struct https://www.onlinetool.io/xmltogo/

	// USERMENTION 引用用户
	USERMENTION struct {
		XMLName     xml.Name `xml:"USERMENTION"`
		Text        string   `xml:",chardata"`
		Displayname string   `xml:"displayname,attr"`
		ID          string   `xml:"id,attr"`
		Username    string   `xml:"username,attr"`
	}

	// POSTMENTION 引用评论
	POSTMENTION struct {
		XMLName      xml.Name `xml:"POSTMENTION"`
		Text         string   `xml:",chardata"`
		Discussionid string   `xml:"discussionid,attr"`
		Displayname  string   `xml:"displayname,attr"`
		ID           string   `xml:"id,attr"`
		Number       string   `xml:"number,attr"`
		Username     string   `xml:"username,attr"`
	}

	// USERHTMLTag 引用用户的tag
	USERHTMLTag struct {
		XMLName xml.Name `xml:"a"`
		Text    string   `xml:",chardata"`
		Href    string   `xml:"href,attr"`
		Class   string   `xml:"class,attr"`
	}

	// POSTHTMLTag  引用其他评论的tag
	POSTHTMLTag struct {
		XMLName xml.Name `xml:"a"`
		Text    string   `xml:",chardata"`
		Href    string   `xml:"href,attr"`
		Class   string   `xml:"class,attr"`
		DataID  string   `xml:"data-id,attr"`
	}
)

// @"corvofeng"#p15 针对15楼的回复
// <r><p><POSTMENTION discussionid="1" displayname="corvofeng" id="15" number="13">@"corvofeng"#p15</POSTMENTION> 针对15楼的回复</p></r>

// @"corvofeng"#1 针对corvofeng的回复
// <r><p><USERMENTION displayname="corvofeng" id="1">@corvofeng</USERMENTION> 针对corvofeng的回复</p></r>

func makeMention(mentionStr []string, comment Comment, user User) string {
	logger := util.GetLogger()
	result := mentionStr[0]
	for {
		userName := mentionStr[1]
		commentID := mentionStr[2]
		if userName != user.Name { // 确保用户信息正确
			logger.Warning("Can't process mention with correct name", mentionStr[0])
			break
		}

		if commentID == "" {
			// for the user mention
			um := USERMENTION{
				ID:          user.StrID(),
				Text:        "@" + user.Name,
				Displayname: user.Nickname,
				Username:    user.Name,
			}
			data, err := xml.Marshal(um)
			if err != nil {
				logger.Warning("Can't create xml data", um)
				break
			}
			result = string(data)
			break
		}

		if commentID == fmt.Sprintf("%d", comment.ID) {
			// for the post mention
			um := POSTMENTION{
				ID:           fmt.Sprintf("%d", comment.ID),
				Text:         "@" + user.Name,
				Discussionid: fmt.Sprintf("%d", comment.AID),
				Number:       fmt.Sprintf("%d", comment.Number),
				Displayname:  user.Nickname,
				Username:     user.Name,
			}
			data, err := xml.Marshal(um)
			if err != nil {
				logger.Warning("Can't create xml data", um)
				break
			}
			result = string(data)
			break
		}
		logger.Warning("Can't get right mention data", mentionStr)
		break
	}
	return result
}

func replaceAllMentions(userData string, mentionDict map[string]string) string {
	for k, v := range mentionDict {
		userData = strings.ReplaceAll(userData, k, v)
	}
	return userData
}

func MentionToHTML(content string) string {
	var um USERMENTION
	var pm POSTMENTION
	var err error
	result := content

	for {
		if err = xml.Unmarshal([]byte(content), &um); err == nil {
			uht := USERHTMLTag{
				Text:  um.Text,
				Href:  fmt.Sprintf("/u/%s", um.Username),
				Class: "UserMention",
			}
			data, err := xml.Marshal(uht)
			if err != nil {
				util.GetLogger().Warning("Can't create xml data", um)
				break
			}
			result = string(data)

			break
		}
		if err = xml.Unmarshal([]byte(content), &pm); err == nil {
			pht := POSTHTMLTag{
				Text:   pm.Text,
				Href:   fmt.Sprintf("/d/%s/%s", pm.Discussionid, pm.Number),
				DataID: pm.ID,
				Class:  "PostMention",
			}
			data, err := xml.Marshal(pht)
			if err != nil {
				util.GetLogger().Warning("Can't create xml data", um)
				break
			}
			result = string(data)
			break
		}
		break
	}
	return result
}
