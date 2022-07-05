package model

import (
	"html/template"
)

// RestfulAPIBase for restful API
type RestfulAPIBase struct {
	State bool `json:"success"`
}

// RestfulReply Reply for restful API
type RestfulReply struct {
}

// RestfulUser user for restful API
type RestfulUser struct {
	Name   string `json:"loginname"`
	Avatar string `json:"avatar_url"`
}

// RestfulTopic topic for restful API
type RestfulTopic struct {
	ID       uint64        `json:"id"`
	UID      uint64        `json:"author_id"`
	Content  template.HTML `json:"content"`
	Title    string        `json:"title"`
	CreateAt string        `json:"create_at"`
	Author   RestfulUser   `json:"author"`

	ReplyCount uint64 `json:"reply_count"`
	VisitCount uint64 `json:"visit_count"`

	Replies     []RestfulReply `json:"replies"`
	LastReplyAt string         `json:"last_reply_at"`
}
