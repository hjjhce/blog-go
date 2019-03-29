package data

import (
	"errors"
	"time"
)

// Post 内容结构
type Post struct {
	ID        int64  `json:"id" validate:"numeric"`
	Title     string `json:"title" validate:"required,alphanum"`
	Content   string `json:"content" validate:"required"`
	Status    int8   `json:"status" validate:"numeric,gte=0"`
	Created   string
	Updated   string
	StatusTag string
}

// StatusList post status
var StatusList = []string{"创建", "发布", "回收"}

// Posts list
func Posts() ([]*Post, error) {
	rows, err := Conn.Query("SELECT id,title,content,status,updated FROM posts")
	if err != nil {
		return nil, err
	}
	var res []*Post
	for rows.Next() {
		p := Post{}
		err = rows.Scan(&p.ID, &p.Title, &p.Content, &p.Status, &p.Updated)
		if err != nil {
			return nil, err
		}
		p.StatusTag = StatusList[p.Status]
		res = append(res, &p)
	}
	return res, nil
}

// GetPostRow 获取单条内容
func GetPostRow(id int64) (*Post, error) {

	if id <= 0 {
		return nil, errors.New("id is invalid")
	}

	post := &Post{}
	err := Conn.QueryRow("SELECT id,title,content,status FROM posts WHERE id = ?", id).Scan(&post.ID, &post.Title, &post.Content, &post.Status)
	if err != nil {
		return nil, err
	}
	return post, nil
}

// Create Post
func (p *Post) Create() (err error) {
	stmt, err := Conn.Prepare("INSERT INTO posts(`title`,`content`,`status`,`created`,`updated`) VALUES (?,?,?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	datetime := time.Now().Format("2006-01-02 15:04:05")

	res, err := stmt.Exec(p.Title, p.Content, 0, datetime, datetime)
	if err != nil {
		return err
	}

	p.ID, err = res.LastInsertId()
	if err != nil {
		return err
	}
	return
}

// Update Post
func (p *Post) Update() (err error) {

	stmt, err := Conn.Prepare("UPDATE posts SET title=?, content=?, status=? WHERE id=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(p.Title, p.Content, p.Status, p.ID)
	if err != nil {
		return err
	}
	return nil
}

// DeletePost delete post item
func DeletePost(id int64) (err error) {
	_, err = Conn.Exec("DELETE FROM posts WHERE id=?", id)
	if err != nil {
		return err
	}
	return
}
