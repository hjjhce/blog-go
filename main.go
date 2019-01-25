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
		// v1.GET("/test", func(ctx *Context){ ctx.JSON(http.StatusOK, M{"hello":"world"})})
		v1.POST("/users/login", login)
		v1.GET("/users/logout", logout)
		v1.POST("/users", userAdd)
		v1.GET("/users", users)
		v1.PUT("/users/:id", usersUpdate)
		v1.DELETE("/users/:id", usersDelete)

		v1.GET("/posts", posts)
		v1.POST("/posts", postsCreate)
		v1.PUT("/posts/{id}", postsUpdate)
		v1.DELETE("/posts/{id}", postsDelete)
	}

	c.Run(":9090")
}

func middleware(ctx *Context) {
	token, ok := ctx.r.Header["token"]
	if !ok || token[0] != ctx.core.token {
		// http.Error(ctx.w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		panic("访问限制")
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
		ctx.index = -1
		return ctx
	}
	c.token = "rulesaremeanttobebroken"
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
			c := core.createContext(w, req)
			for i := 0; ; i++ {
				pc, file, line, ok := runtime.Caller(i)
				if ok {
					log.Printf("%s:%d (0x%x) ", file, line, pc)
				}
			}

			log.Printf("Panic: %s", rev)
			c.String(http.StatusServiceUnavailable, fmt.Sprintf("%s", rev))
			core.pool.Put(c)
		}
	}()

	log.Printf("[req] %s %s %s", req.Host, req.Method, req.URL.Path)
	core.httprouter.ServeHTTP(w, req)
}

/*
func middleware(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		log.Printf("[Request]: %s %s %s", r.Host, r.URL.RequestURI(), r.Method)

		r.ParseForm()
		s := r.FormValue("sign")

		if s != adminSign {
			err := errors.New("没有访问权限")
			resp := createFailResp(StatusForbidden, err)
			resp.returnJSON(w, r, http.StatusForbidden)
			return
		}

		//每次请求都判断下session
		p := path.Base(r.URL.Path)

		if p != "login" {
			cookiestr, err := r.Cookie("uid")

			if err != nil || !data.IsLogin(cookiestr.Value) {
				err := errors.New("没有访问权限")
				resp := createFailResp(StatusForbidden, err)
				resp.returnJSON(w, r, http.StatusForbidden)
				return
			}
		}
		// res := h(w, r, ps)
		// w.WriteHeader(res.Code)
		// json.NewEncoder(w).Encode(res.Data)
	}
}
*/
