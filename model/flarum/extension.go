package flarum

// IExtensionsV1 flarum的扩展
type IExtensionsV1 interface {
	Register()
	SetAttributes(map[string]interface{})
}
