package model

// RestfulAPIBase for restful API
type RestfulAPIBase struct {
	State string `json:"success"`
}

type RestfulReply struct {
}

// RestfulTopic topic for restful API
type RestfulTopic struct {
	ID      uint64 `json:"id"`
	UID     uint64 `json:"author_id"`
	Content string `json:"content"`
	Title   string `json:"title"`

	ReplyCount uint64 `json:"reply_count"`
	VisitCount uint64 `json:"visit_count"`

	Replies []RestfulReply `json:"replies"`
}
