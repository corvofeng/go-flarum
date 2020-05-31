package flarum

// Group group信息
type Group struct {
	Type string `json:"type"`

	ID           string `json:"id"`
	NameSingular string `json:"nameSingular"`
	NamePlural   string `json:"namePlural"`
	Color        string `json:"color"`
	Icon         string `json:"icon"`
	IsHidden     int    `json:"isHidden"`
}

// DoInit 初始化Group
func (g *Group) DoInit() {
	g.Type = "group"
}

// GetDefaultAttributes 获取属性
func (g *Group) GetDefaultAttributes(obj interface{}) {
}
