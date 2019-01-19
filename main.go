package main

import (
	"log"
	"net/http"

	"gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate
var latestVer = "/v1"

const token = "pingtouge"

func main() {

	f := createLogFile()
	defer f.Close()

	// log.SetOutput(f)

	// 定制httprouter
	router := NewRouter()
	validate = validator.New()

	// router.GET("/v1/users/auth", checkSession)
	{
		router.POST("/v1/users/login", login)
		router.GET("/v1/users/logout", logout)
		router.POST("/v1/users", userAdd)
		router.GET("/v1/users", users)
		router.PUT("/v1/users/:id", usersUpdate)
		router.DELETE("/v1/users/:id", usersDelete)
		router.GET("/v1/posts", posts)
		router.POST("/v1/posts", postsCreate)
		router.PUT("/v1/posts/{id}", postsUpdate)
		router.DELETE("/v1/posts/{id}", postsDelete)
	}

	log.Fatal(http.ListenAndServe(":9090", router))
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
