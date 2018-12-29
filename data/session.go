package data

import (
	"time"
)

// Session 会话信息
type Session struct {
	ID      int64
	UID     string
	Email   string
	Name    string
	Created time.Time
}

// SessionByEmail 通过email获取session
var sess map[string]*Session

// CreateSession 创建会话
func (u *User) CreateSession() *Session {

	if sess == nil {
		sess = make(map[string]*Session)
	}

	key := randNumStr()
	sess[key] = &Session{
		UID:     key,
		ID:      u.ID,
		Email:   u.Email,
		Name:    u.Name,
		Created: time.Now(),
	}
	return sess[key]
}

// Session 获取用户会话
func (u *User) Session(uid string) *Session {
	if sess == nil {
		return nil
	}

	if v, ok := sess[uid]; ok {
		return v
	}
	return nil
}

// IsLogin 判断是否登录
func IsLogin(uid string) bool {
	if sess == nil {
		return false
	}

	_, ok := sess[uid]
	return ok
}
