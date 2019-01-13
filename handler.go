package main

import (
	"blog/data"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	validator "gopkg.in/go-playground/validator.v9"
)

// StatusOK 请求成功
const StatusOK = "0000"

// StatusRequestJSONErr 请求json解析错误
const StatusRequestJSONErr = "4001"

// StatusParamsValidErr 请求参数格式错误
const StatusParamsValidErr = "4002"

// StatusForbidden 没有访问的权限
const StatusForbidden = "4003"

// StatusAuthErr 登录认证失败
const StatusAuthErr = "4005"

// StatusServerErr 服务器内部错误
const StatusServerErr = "5000"

var session *data.Session
var sessList map[string]data.Session

// M 响应数据结构
type M map[string]interface{}

// JSON 响应JSON数据
func (ctx *Context) JSON(code int, res M) {

	if code == http.StatusOK {
		log.Printf("[response] %s %d", ctx.r.URL.Path, code)
	} else {
		log.Printf("[response] %s %d %s", ctx.r.URL.Path, code, res)
	}

	ctx.w.WriteHeader(code)
	json.NewEncoder(ctx.w).Encode(res)
}

func test(ctx *Context) {
	fmt.Println(ctx.r.Method)
	ctx.JSON(http.StatusOK, M{"code": 200, "msg": "ok", "hello": "world"})
}

// 用户后台登录
func login(ctx *Context) {

	var auth data.Auth
	err := json.NewDecoder(ctx.r.Body).Decode(&auth)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, M{"code": StatusRequestJSONErr, "msg": fmt.Sprintf("%s", err)})
		return
	}

	err = validate.Struct(auth)
	if err != nil {
		var errmsg string
		for _, err := range err.(validator.ValidationErrors) {
			errmsg += err.Field() + ":" + fmt.Sprintf("%s", err.Value()) + ", the type should be " + err.Type().String() + "(" + err.Tag() + "); "
		}

		ctx.JSON(http.StatusBadRequest, M{"code": StatusParamsValidErr, "msg": errmsg})
		return
	}

	var user *data.User
	user, err = auth.Login()
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, M{"code": StatusAuthErr, "msg": fmt.Sprintf("%s", err)})
		return
	}
	session = user.CreateSession()
	fmt.Println(session)
	//设置cookie
	c := http.Cookie{
		Name:     "uid",
		Value:    session.UID,
		HttpOnly: true,
	}

	http.SetCookie(ctx.w, &c)
	ctx.JSON(http.StatusOK, M{"code": StatusOK, "msg": "ok"})
}

// logout 登出
func logout(ctx *Context) {

	defer func() {
		if p := recover(); p != nil {
			ctx.JSON(http.StatusForbidden, M{"code": StatusForbidden, "msg": fmt.Sprintf("%s", p)})
			return
		}
	}()

	uid, err := ctx.r.Cookie("uid")
	if err != nil {
		panic(err)
	}

	if _, ok := sessList[uid.Value]; !ok {
		panic(err)
	}

	delete(sessList, uid.Value)
	ctx.JSON(http.StatusOK, M{"code": StatusOK, "msg": "ok"})
}

// userAdd 添加用户
func userAdd(ctx *Context) {

	var user data.User

	err := json.NewDecoder(ctx.r.Body).Decode(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, M{"code": StatusRequestJSONErr, "msg": fmt.Sprintf("%s", err)})
		return
	}

	err = validate.Struct(user)
	if err != nil {

		var errmsg = make(map[string]string)

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

		ctx.JSON(http.StatusBadRequest, M{"code": StatusParamsValidErr, "msg": errmsg})
		return
	}

	err = user.Create()
	if err != nil {
		ctx.JSON(http.StatusServiceUnavailable, M{"code": StatusServerErr, "msg": fmt.Sprintf("%s", err)})
		return
	}

	ctx.JSON(http.StatusOK, M{"code": StatusOK, "msg": "ok"})
	return
}

// 获取用户列表
func users(ctx *Context) {

	log.Printf("haha")
	_, err := data.Users()
	if err != nil {
		ctx.JSON(http.StatusServiceUnavailable, M{"code": StatusServerErr, "msg": fmt.Sprintf("%s", err)})
		return
	}

	ctx.JSON(http.StatusOK, M{"code": StatusOK, "msg": "ok"})
}

func usersUpdate(ctx *Context) {

}

func usersDelete(ctx *Context) {

	s := ctx.r.FormValue("id")
	if len(s) == 0 {
		ctx.JSON(http.StatusBadRequest, M{"code": StatusParamsValidErr, "msg": "id错误"})
	}

	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, M{"code": StatusParamsValidErr, "msg": fmt.Sprintf("%s", err)})
		return
	}

	err = data.DeleteUser(id)
	if err != nil {
		ctx.JSON(http.StatusServiceUnavailable, M{"code": StatusServerErr, "msg": fmt.Sprintf("%s", err)})
		return
	}

	ctx.JSON(http.StatusOK, M{"code": StatusOK, "msg": "ok"})
}

func posts(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	rows, err := data.Posts()
	if err != nil {
		resp := createFailResp(StatusServerErr, err)
		resp.returnJSON(w, r, http.StatusServiceUnavailable)
		return
	}

	resp := createSuccessResp(rows)
	resp.returnJSON(w, r, http.StatusOK)
}

func postsCreate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func postsUpdate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func postsDelete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}
