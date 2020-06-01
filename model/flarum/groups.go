package flarum

// Group group信息
type Group struct {
	BaseResources

	NameSingular string `json:"nameSingular"`
	NamePlural   string `json:"namePlural"`
	Color        string `json:"color"`
	Icon         string `json:"icon"`
	IsHidden     bool   `json:"isHidden"`
}

// DoInit 初始化Group
func (g *Group) DoInit() {
	g.Type = "groups"
}

// GetType 获取类型
func (g *Group) GetType() string {
	return g.Type
}

// GetDefaultAttributes 获取属性
func (g *Group) GetDefaultAttributes(obj interface{}) {
}
