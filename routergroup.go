package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Routers 路由组
type Routers struct {
	handlers []HandlerFunc
	basePath string
	core     *Core
}

// Use 添加中间件
func (rs *Routers) Use(middlewares ...HandlerFunc) {
	for _, m := range middlewares {
		rs.handlers = append(rs.handlers, m)
	}
}

// Group 分组
func (rs *Routers) Group(path string, handlers ...HandlerFunc) *Routers {
	return &Routers{handlers: rs.appendHandlers(handlers), basePath: rs.joinPath(path), core: rs.core}
}

// Handle 处理器
func (rs *Routers) Handle(method, path string, handlers []HandlerFunc) {

	handlers = rs.appendHandlers(handlers)
	finalPath := rs.joinPath(path)

	// 添加到httprouter的tree节点，并在serveHTTP的handle(w,req,ps)里执行c.Next
	rs.core.httprouter.Handle(method, finalPath, func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		c := rs.core.createContext(w, req)
		c.hanlders = handlers
		c.params = p
		c.Next() //遍历执行handlers
		rs.core.pool.Put(c)
	})
}

// appendHandlers 合并处理器函数
func (rs *Routers) appendHandlers(handlers []HandlerFunc) []HandlerFunc {
	al := len(rs.handlers)
	bl := len(handlers)
	length := al + bl

	if al > 0 && bl > 0 {
		returnHandlers := make([]HandlerFunc, length)
		copy(returnHandlers, rs.handlers)
		copy(returnHandlers[al:], handlers) //如果使用append，则会在切片末尾加上handlers, 这样会导致之前make的切片实际没填满，存在nil
		return returnHandlers
	} else if al > 0 {
		return rs.handlers
	}
	return handlers

}

// GET handle http get method
func (rs *Routers) GET(path string, handlers ...HandlerFunc) {
	rs.Handle("GET", path, handlers)
}

// POST method
func (rs *Routers) POST(path string, handlers ...HandlerFunc) {
	rs.Handle("POST", path, handlers)
}

// PUT method
func (rs *Routers) PUT(path string, handlers ...HandlerFunc) {
	rs.Handle("PUT", path, handlers)
}

// PATCH method
func (rs *Routers) PATCH(path string, handlers ...HandlerFunc) {
	rs.Handle("PATCH", path, handlers)
}

// DELETE method
func (rs *Routers) DELETE(path string, handlers ...HandlerFunc) {
	rs.Handle("DELETE", path, handlers)
}

// OPTIONS method
func (rs *Routers) OPTIONS(path string, handlers ...HandlerFunc) {
	rs.Handle("OPTIONS", path, handlers)
}

// HEAD method
func (rs *Routers) HEAD(path string, handlers ...HandlerFunc) {
	rs.Handle("HEAD", path, handlers)
}

func (rs *Routers) joinPath(relativePath string) string {
	return JoinPath(rs.basePath, relativePath)
}
