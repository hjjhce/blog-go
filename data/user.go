package data

import (
	"errors"
	"log"
	"time"
)

// IDNotFound 无法获取用户ID
const IDNotFound int64 = -1

// User 用户数据结构
type User struct {
	ID        int64
	Name      string `json:"name" validate:"required,alphanum"`
	Email     string `json:"email" validate:"required,email"`
	Mobile    string `json:"mobile" validate:"required,numeric,len=11"`
	Password  string `json:"password" validate:"required"`
	Checkpwd  string `json:"checkpwd" validate:"required,eqfield=Password"`
	Role      int    `json:"role" validate:"numeric"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Auth 用户登录认证结构
type Auth struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// Login 用户登录认证
func (m *Auth) Login() (int64, error) {

	var id int64
	err := Conn.QueryRow("SELECT id FROM users WHERE email=? and password=?", m.Email, Encrypt(m.Password)).Scan(&id)
	if err != nil {
		log.Printf("%s:%s", m.Email, err)
		return IDNotFound, errors.New("用户名或密码错误")
	}

	return id, nil
}

// Create 创建用户
func (user *User) Create() (int64, error) {

	//save user data

	stmt, err := Conn.Prepare("INSERT INTO `users` (`name`, `email`, `mobile`, `password`, `role`, `created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?) ")
	if err != nil {
		return IDNotFound, err
	}

	hashpwd := Encrypt(user.Password)
	t := time.Now().Format("2006-01-02 15:04:05")
	res, err := stmt.Exec(user.Name, user.Email, user.Mobile, hashpwd, user.Role, t, t)
	if err != nil {
		return IDNotFound, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return IDNotFound, err
	}
	return id, nil
}

// Update 更新用户信息
func Update() {

}

// Delete 删除用户
func Delete() {

}