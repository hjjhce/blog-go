package main

import (
	"blog/data"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	validator "gopkg.in/go-playground/validator.v9"
)

// StatusOK 请求成功
const StatusOK = "success"

// StatusErr 请求失败
const StatusErr = "failed"

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

type result struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  errResp     `json:"error,omitempty"`
	Time   string      `json:"time"`
}

type errResp struct {
	Code   string `json:"code,omitempty"`
	Errmsg string `json:"errmsg,omitempty"`
}

var session *data.Session

func (res *result) returnJSON(w http.ResponseWriter, r *http.Request, code int) {
	log.Printf("[response]: %s %s [%s]%s", r.Method, r.URL.Path, res.Error.Code, res.Error.Errmsg)

	w.WriteHeader(code)
	// json.MarshalIndent(res);w.Write(res)
	json.NewEncoder(w).Encode(res)
}

func createFailResp(code string, err error) result {
	return result{Status: "failed", Error: errResp{Code: code, Errmsg: err.Error()}, Time: time.Now().Format("2006-01-02 15:03:04")}
}

func createSuccessResp(data interface{}) result {
	return result{Status: "success", Data: data, Time: time.Now().Format("2006-01-02 15:03:04")}
}

// 用户后台登录
func login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var auth data.Auth
	err := json.NewDecoder(r.Body).Decode(&auth)
	if err != nil {
		resp := createFailResp(StatusRequestJSONErr, err)
		resp.returnJSON(w, r, http.StatusBadRequest)
	}

	err = validate.Struct(auth)
	if err != nil {
		var errmsg string
		for _, err := range err.(validator.ValidationErrors) {
			errmsg += err.Field() + ":" + fmt.Sprintf("%s", err.Value()) + ", the type should be " + err.Type().String() + "(" + err.Tag() + "); "
		}
		resp := createFailResp(StatusParamsValidErr, errors.New(errmsg))
		resp.returnJSON(w, r, http.StatusBadRequest)
		return
	}

	var user *data.User
	user, err = auth.Login()
	if err != nil {
		resp := createFailResp(StatusAuthErr, err)
		resp.returnJSON(w, r, http.StatusUnauthorized)
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

	http.SetCookie(w, &c)

	type mdata struct {
		ID int64 `json:"id"`
	}

	resp := createSuccessResp(mdata{ID: user.ID})
	resp.returnJSON(w, r, http.StatusOK)
}

// userAdd 添加用户
func userAdd(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var user data.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		resp := createFailResp(StatusRequestJSONErr, err)
		resp.returnJSON(w, r, http.StatusBadRequest)
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
		newErr := fmt.Errorf("%s", errmsg)
		resp := createFailResp(StatusParamsValidErr, newErr)
		resp.returnJSON(w, r, http.StatusBadRequest)
		return
	}

	err = user.Create()
	if err != nil {
		resp := createFailResp(StatusServerErr, err)
		resp.returnJSON(w, r, http.StatusServiceUnavailable)
		return
	}

	type d struct {
		ID int64 `json:"id"`
	}

	resp := createSuccessResp(d{ID: user.ID})
	resp.returnJSON(w, r, http.StatusOK)
	return
}

// 获取用户列表
func users(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	rows, err := data.Users()
	if err != nil {
		resp := createFailResp(StatusServerErr, err)
		resp.returnJSON(w, r, http.StatusServiceUnavailable)
		return
	}

	resp := createSuccessResp(rows)
	resp.returnJSON(w, r, http.StatusOK)
}

func usersUpdate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func usersDelete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	s := ps.ByName("id")
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		resp := createFailResp(StatusParamsValidErr, err)
		resp.returnJSON(w, r, http.StatusBadRequest)
		return
	}

	err = data.DeleteUser(id)
	if err != nil {
		resp := createFailResp(StatusServerErr, err)
		resp.returnJSON(w, r, http.StatusServiceUnavailable)
		return
	}

	// w.WriteHeader(http.StatusOK)
	resp := createSuccessResp(nil)
	resp.returnJSON(w, r, http.StatusOK)
}
