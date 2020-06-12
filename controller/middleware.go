package controller

import (
	"context"
	"fmt"
	"net/http"
)

// 与中间件相关的函数

type (
	// HTTPHandleFunc 用于处理http请求的函数
	HTTPHandleFunc func(w http.ResponseWriter, r *http.Request)

	// HTTPMiddleWareFunc 中间件函数
	HTTPMiddleWareFunc func(inner HTTPHandleFunc) HTTPHandleFunc
)

// MiddlewareArrayToChains 中间件整理成链式的函数调用形式
/* 当我们某个使用了多个中间件时, 可以方便的进行整合:
sp.HandleFunc(pat.Get("/"), controller.ArrayToChains(
	[]controller.ReqMiddle{
	controller.TestMiddleware,
	controller.TestMiddleware2,
	},
	h.FlarumIndex,
))

将会返回被中间件包裹的如下形式的函数:
controller.TestMiddleware(controller.TestMiddleware2(h.FlarumIndex))
*/
func MiddlewareArrayToChains(reqProcessFuncs []HTTPMiddleWareFunc, req HTTPHandleFunc) (rp HTTPHandleFunc) {
	rp = req
	rpfs := reqProcessFuncs
	for i := len(rpfs) - 1; i >= 0; i-- {
		rp = rpfs[i](rp)
	}
	return
}

// InitMiddlewareContext 初始化的中间件需要的数据结构
/*
中间件中传递数据依赖于context设计, 当前的context作为结构体而存在, 每次请求时新建一个对应的结构体, 并存储相关信息,
在真正处理请求时, 获取该结构体, 并获得中间件传递的信息.
*/
func (h *BaseHandler) InitMiddlewareContext(inner http.Handler) http.Handler {
	mw := func(w http.ResponseWriter, r *http.Request) {
		reqCtx := &ReqContext{}
		reqCtx.h = h
		r = r.WithContext(
			context.WithValue(r.Context(), ckRequest, reqCtx),
		)
		inner.ServeHTTP(w, r)
	}
	return http.HandlerFunc(mw)
}

// AuthMiddleware 校验用户
func (h *BaseHandler) AuthMiddleware(inner http.Handler) http.Handler {
	mw := func(w http.ResponseWriter, r *http.Request) {
		reqCtx := GetRetContext(r)
		reqCtx.currentUser, _ = h.CurrentUser(w, r)
		inner.ServeHTTP(w, r)
	}
	return http.HandlerFunc(mw)
}

// InAPIMiddleware 被此装饰器修饰表明当前请求为API请求
func InAPIMiddleware(inner HTTPHandleFunc) HTTPHandleFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCtx := GetRetContext(r)
		reqCtx.inAPI = true
		inner(w, r)
	}
}

func TestMiddleware(inner HTTPHandleFunc) HTTPHandleFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("In md 1")
		inner(w, r)
	}
}
func TestMiddleware2(inner HTTPHandleFunc) HTTPHandleFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("In md 2")
		inner(w, r)
	}
}
