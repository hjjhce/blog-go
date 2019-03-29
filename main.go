package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"sync"

	"github.com/julienschmidt/httprouter"
)

// Core xyz
type Core struct {
	*Routers
	pool       sync.Pool
	httprouter *httprouter.Router
	token      string
}

// HandlerFunc 处理器函数
type HandlerFunc func(*Context)

func main() {

	f := createLogFile()
	defer f.Close()

	// log.SetOutput(f)

	c := New()
	c.Use(middleware)
	v1 := c.Group("/v1")
	// v1.Use(checkAuth)
	{
		v1.POST("/users/login", login)
		v1.GET("/users/logout", logout)
		v1.POST("/users", userAdd)
		v1.GET("/users", users)
		v1.PUT("/users/:id", usersUpdate)
		v1.DELETE("/users/:id", usersDelete)

		v1.GET("/posts", posts)
		v1.GET("/posts/:id", getPostRow)
		v1.POST("/posts", postsCreate)
		v1.PUT("/posts/:id", postsUpdate)
		v1.DELETE("/posts/:id", postsDelete)
	}

	c.Run(":9090")
}

func middleware(ctx *Context) {
	// fmt.Println(ctx.r.Header)
	token, ok := ctx.r.Header["Token"]
	if !ok || token[0] != ctx.core.token {
		panic("forbiden")
	}
}

func checkAuth(ctx *Context) {
	fmt.Println("Auth", ctx.r.URL.Path)
}

// New 创建框架实例
func New() *Core {
	c := &Core{}
	c.Routers = &Routers{
		handlers: nil,
		basePath: "/",
		core:     c,
	}
	c.httprouter = httprouter.New()
	c.pool.New = func() interface{} {
		ctx := &Context{}
		return ctx
	}
	c.token = "pingtouge"
	return c
}

// Run serve
func (core *Core) Run(addr string) {
	if err := http.ListenAndServe(addr, core); err != nil {
		log.Fatal(err)
	}
}

func (core *Core) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if rev := recover(); rev != nil {
			core.HandlerPanic(w, req, rev)
		}
	}()

	log.Printf("[req] %s %s %s", req.Host, req.Method, req.URL.Path)
	core.httprouter.ServeHTTP(w, req)
}

// HandlerPanic 异常处理
func (core *Core) HandlerPanic(w http.ResponseWriter, req *http.Request, rev interface{}) {
	c := core.createContext(w, req)
	for i := 3; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		log.Printf("%s:%d (0x%x) ", file, line, pc)
	}

	log.Printf("Panic: %s", rev)
	c.String(http.StatusServiceUnavailable, fmt.Sprintf("%s", rev))
	core.pool.Put(c)
}
