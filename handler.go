package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"time"

	"github.com/kataras/iris"
	validator "gopkg.in/go-playground/validator.v9"
)

// StatusOK 请求成功
const StatusOK = "0000"

// StatusErr 请求失败
const StatusErr = "4000"

// StatusRequestJSONErr 请求json解析错误
const StatusRequestJSONErr = "4001"

// StatusParamsValidErr 请求参数格式错误
const StatusParamsValidErr = "4002"

// StatusServerErr 服务器内部错误
const StatusServerErr = "5000"

func login(ctx iris.Context) {

}

// userAdd 添加用户
func userAdd(ctx iris.Context) {

	defer func() {
		if err := recover(); err != nil {
			newErr := fmt.Errorf("%s", err)
			errReturn(ctx, newErr, StatusRequestJSONErr)
		}
	}()

	var user User
	if err := ctx.ReadJSON(&user); err != nil {
		errReturn(ctx, err, StatusServerErr)
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

	// 密码sha1加密
	hash := sha1.New()
	io.WriteString(hash, user.Password)
	hashpwd := fmt.Sprintf("%x", hash.Sum(nil))

	//save user data
	t := time.Now().Format("2006-01-02 15:04:05")
	stmt, err := db.Prepare("INSERT INTO `users` (`name`, `email`, `mobile`, `password`, `role`, `created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?) ")
	checkErr(err)

	res, err := stmt.Exec(user.Userame, user.Email, user.Mobile, hashpwd, user.Role, t, t)
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	ctx.JSON(iris.Map{
		"status": "successd",
		"data": iris.Map{
			"code":    StatusOK,
			"userid":  id,
			"created": t,
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
