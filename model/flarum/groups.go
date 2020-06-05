package flarum

// Group group信息
type Group struct {
	BaseResources

	NameSingular string `json:"nameSingular"`
	NamePlural   string `json:"namePlural"`
	Color        string `json:"color"`
	Icon         string `json:"icon"`
	IsHidden     bool   `json:"isHidden"`

	// FlarumExtensions []IExtensions
}

// DoInit 初始化Group
func (g *Group) DoInit() {
	g.Type = "groups"
}

// GetType 获取类型
func (g *Group) GetType() string {
	return g.Type
}

// GetAttributes 获取属性
func (g *Group) GetAttributes() map[string]interface{} {
	return nil
}
