package flarum

type FlarumFoFFiles struct {
	BaseResources

	BaseName    string `json:"baseName"`    // 文件的基本名称
	UUID        string `json:"uuid"`        // 文件的唯一标识符
	BBCode      string `json:"bbcode"`      // 文件的BBCode格式
	URL         string `json:"url"`         // 文件的URL
	Shared      bool   `json:"shared"`      // 是否共享
	CanViewInfo bool   `json:"canViewInfo"` // 是否可以查看文件信息
	CanHide     bool   `json:"canHide"`     // 是否可以隐藏文件
	CanDelete   bool   `json:"canDelete"`   // 是否可以删除文件
	FileType    string `json:"fileType"`    // 文件类型
}

// DoInit 初始化tags
func (t *FlarumFoFFiles) DoInit(id uint64) {
	t.setID(id)
	t.setType("files")
}

// GetType 获取类型
func (t *FlarumFoFFiles) GetType() string {
	return t.Type
}

// GetAttributes
func (t *FlarumFoFFiles) GetAttributes() (map[string]interface{}, error) {
	// 返回属性值
	return map[string]interface{}{
		"uuid":        t.UUID,
		"bbcode":      t.BBCode,
		"url":         t.URL,
		"shared":      t.Shared,
		"canViewInfo": t.CanViewInfo,
		"canHide":     t.CanHide,
		"canDelete":   t.CanDelete,
		"type":        t.FileType,
	}, nil
}
