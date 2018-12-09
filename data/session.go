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
var SessionByEmail map[string]*Session

// CreateSession 创建会话
func (u *User) CreateSession() {

	if SessionByEmail == nil {
		SessionByEmail = make(map[string]*Session)
	}

	sess := &Session{
		UID:     randNumStr(),
		ID:      u.ID,
		Email:   u.Email,
		Name:    u.Name,
		Created: time.Now(),
	}

	SessionByEmail[u.Email] = sess
}

// SessionUID 获取sessionuid
func (u *User) SessionUID() string {
	if v, ok := SessionByEmail[u.Email]; ok {
		return v.UID
	}
	return ""
}

// Session 获取用户会话
func (u *User) Session() *Session {
	if v, ok := SessionByEmail[u.Email]; ok {
		return v
	}
	return nil
}
