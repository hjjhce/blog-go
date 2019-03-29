package main

import (
	"blog/data"
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Context 请求和响应数据
type Context struct {
	w        http.ResponseWriter
	r        *http.Request
	hanlders []HandlerFunc
	sessions *data.Session
	index    int8
	core     *Core
	params   httprouter.Params
}

func (core *Core) createContext(w http.ResponseWriter, r *http.Request) *Context {
	c := core.pool.Get().(*Context)
	c.w = w
	c.r = r
	c.core = core
	return c
}

// M 响应数据结构
type M map[string]interface{}

// JSON 响应JSON数据
func (ctx *Context) JSON(code int, res M) {

	if code == http.StatusOK {
		log.Printf("[response] %s %d", ctx.r.URL.Path, code)
	} else {
		log.Printf("[response] %s %d %s", ctx.r.URL.Path, code, res)
	}

	ctx.w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	ctx.w.WriteHeader(code)
	json.NewEncoder(ctx.w).Encode(res)
}

func (ctx *Context) String(code int, str string) {

	ctx.w.Header().Set("Content-Type", "text/html;charset=UTF-8")
	ctx.w.WriteHeader(code)
	ctx.w.Write([]byte(str))

}

// Next 遍历执行handlers,包括中间件的handler
func (ctx *Context) Next() {
	// for k, h := range ctx.hanlders {
	// 	fmt.Printf("%d, %v\r\n", k, h)
	// }
	count := len(ctx.hanlders)
	for index := 0; index < count; index++ {
		ctx.hanlders[index](ctx)
	}
}
