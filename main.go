package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql" // 使用mysql驱动
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	"gopkg.in/go-playground/validator.v9"
)

var db *sql.DB
var validate *validator.Validate
var latestVer = "/v1"

func init() {
	var err error
	db, err = sql.Open("mysql", "root:root@/blog")
	if err != nil {
		panic(err)
	}
}

func main() {

	validate = validator.New()

	f := createLogFile()
	defer f.Close()

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
}
