package main

// User 用户数据结构
type User struct {
	Userame  string `json:"name" validate:"required,alphanum"`
	Email    string `json:"email" validate:"required,email"`
	Mobile   string `json:"mobile" validate:"required,numeric,len=11"`
	Password string `json:"password" validate:"required"`
	Checkpwd string `json:"checkpwd" validate:"required,eqfield=Password"`
	Role     int    `json:"role" validate:"numeric"`
}

type result struct {
}
