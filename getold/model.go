package getold

type OldArticle struct {
	Code         int    `json:"code"`
	ID           string `json:"id"`
	UID          string `json:"uid"`
	CID          string `json:"cid"`
	RUID         string `json:"ruid"`
	Title        string `json:"title"`
	Content      string `json:"content"`
	Tags         string `json:"tags"`
	AddTime      string `json:"addtime"`
	EditTime     string `json:"edittime"`
	Comments     string `json:"comments"`
	CloseComment string `json:"closecomment"`
	Visible      string `json:"visible"`
	Views        string `json:"views"`
}
type OldUser struct {
	Code          int    `json:"code"`
	ID            string `json:"id"`
	Name          string `json:"name"`
	Flag          string `json:"flag"`
	Avatar        string `json:"avatar"`
	Password      string `json:"password"`
	Email         string `json:"email"`
	URL           string `json:"url"`
	Articles      string `json:"articles"`
	Replies       string `json:"replies"`
	RegTime       string `json:"regtime"`
	LastPostTime  string `json:"lastposttime"`
	LastReplyTime string `json:"lastreplytime"`
	LastLoginTime string `json:"lastLoginTime"`
	About         string `json:"about"`
	Notice        string `json:"notic"`
}
type OldCategory struct {
	Code     int    `json:"code"`
	ID       string `json:"id"`
	Name     string `json:"name"`
	Articles string `json:"articles"`
	About    string `json:"about"`
	Comments string `json:"comments"`
	Hidden   bool   `json:"hidden"`
}
type OldComment struct {
	Code    int    `json:"code"`
	ID      string `json:"id"`
	AID     string `json:"articleid"`
	UID     string `json:"uid"`
	Content string `json:"content"`
	AddTime string `json:"addtime"`
}
type OldQQ struct {
	Code   int    `json:"code"`
	ID     string `json:"id"`
	UID    string `json:"uid"`
	Name   string `json:"name"`
	Openid string `json:"openid"`
}
type OldWeibo struct {
	Code   int    `json:"code"`
	ID     string `json:"id"`
	UID    string `json:"uid"`
	Name   string `json:"name"`
	Openid string `json:"openid"`
}
type OldTag struct {
	Code     int    `json:"code"`
	ID       string `json:"id"`
	Name     string `json:"name"`
	Articles string `json:"articles"`
	IDs      string `json:"ids"`
}
