package main

import (
	"blog/data"
	"errors"
	"fmt"
	"log"
	"net/http"
	"path"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate
var latestVer = "/v1"

const adminSign = "whyspacex"

func main() {

	f := createLogFile()
	defer f.Close()

	// log.SetOutput(f)

	// 定制httprouter
	router := NewRouter()
	validate = validator.New()

	router.POST("/v1/users/login", test)

	// router.POST("/v1/users/login", middleware(login))
	// router.POST("/v1/users", middleware(userAdd))
	// router.GET("/v1/users", middleware(users))
	// router.PUT("/v1/users/:id", middleware(usersUpdate))
	// router.DELETE("/v1/users/:id", middleware(usersDelete))

	log.Fatal(http.ListenAndServe(":9090", router))
}

func test(ctx *Context) {
	fmt.Println(ctx.r.Method)

	ctx.JSON(http.StatusOK, H{"code": 200, "msg": "ok"})
}

func middleware(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		var forbidden = false

		defer func() {
			if forbidden {
				http.Redirect(w, r, "/login", http.StatusForbidden)
				err := errors.New("没有访问权限")
				resp := createFailResp(StatusForbidden, err)
				resp.returnJSON(w, r, http.StatusForbidden)
			}
		}()

		log.Printf("[Request]: %s %s %s", r.Host, r.URL.RequestURI(), r.Method)

		r.ParseForm()
		s := r.FormValue("sign")

		if s != adminSign {
			forbidden = true
			return
		}

		//每次请求都判断下session
		p := path.Base(r.URL.Path)

		if p != "login" {
			cookiestr, err := r.Cookie("uid")

			if err != nil || !data.IsLogin(cookiestr.Value) {
				forbidden = true
				return
			}
		}
		h(w, r, ps)
	}
}
