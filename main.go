package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate
var latestVer = "/v1"

const adminSign = "whyspacex"

func ilog(h httprouter.Handle) httprouter.Handle {
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
		h(w, r, ps)
	}
}

func main() {

	f := createLogFile()
	defer f.Close()

	log.SetOutput(f)

	// 重新写个路由
	router := httprouter.New()
	validate = validator.New()

	router.GET("/ping", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		fmt.Fprintln(w, p)
	})
	router.POST("/v1/login", ilog(login))
	router.POST("/v1/user/add", ilog(userAdd))
	router.GET("/v1/users", ilog(users))

	log.Fatal(http.ListenAndServe(":9090", router))
}

/*
	app := iris.New()
	app.Logger().SetOutput(f) //日志写入文件

	app.Logger().SetLevel("debug")
	app.Use(recover.New())
	app.Use(logger.New())

	app.Get("/ping", func(ctx iris.Context) {
		ctx.JSON(iris.Map{
			"message": "pong",
		})
	})

	//接口版本
	version := app.Party(latestVer)
	{
		version.Post("/login", login)
		version.Post("/user-add", userAdd)
	}

	app.Run(iris.Addr(":9090"))

	// userAdd 添加用户
	func userAdd(ctx iris.Context) {

		defer func() {
			if err := recover(); err != nil {
				newErr := fmt.Errorf("%s", err)
				errReturn(ctx, newErr, StatusErr)
			}
		}()

		var user data.User

		if err := ctx.ReadJSON(&user); err != nil {
			errReturn(ctx, err, StatusRequestJSONErr)
			return
		}

		err := validate.Struct(user)
		if err != nil {

			var errmsg = make(map[string]string)

			ctx.StatusCode(iris.StatusBadRequest)
			for _, err := range err.(validator.ValidationErrors) {

				fmt.Println(err.Namespace())
				fmt.Println(err.Field())
				fmt.Println(err.StructNamespace()) // can differ when a custom TagNameFunc is registered or
				fmt.Println(err.StructField())     // by passing alt name to ReportError like below
				fmt.Println(err.Tag())
				fmt.Println(err.ActualTag())
				fmt.Println(err.Kind())
				fmt.Println(err.Type())
				fmt.Println(err.Value())
				fmt.Println(err.Param())

				errmsg[err.Field()] = fmt.Sprintf("%s", err.Value()) + "; but type should be " + err.Type().String() + " and " + err.Tag() + ";"
			}
			newErr := fmt.Errorf("%s", errmsg)
			errReturn(ctx, newErr, StatusParamsValidErr)
			return
		}

		id, err := user.Create()
		checkErr(err)

		ctx.JSON(iris.Map{
			"status": "successd",
			"data": iris.Map{
				"code":    StatusOK,
				"userid":  id,
				"created": time.Now().Format("2006-01-02 15:03:04"),
			},
		})
		return
	}

	func errReturn(ctx iris.Context, err error, errCode string) {
		ctx.Application().Logger().Errorf("hanlder:%s\r\n\terrmsg:%s", ctx.HandlerName(), err.Error())
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{
			"status": "failed",
			"errors": iris.Map{
				"code":    errCode,
				"message": err.Error(),
			},
		})
	}
*/
