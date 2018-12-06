package data

import (
	"crypto/sha1"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // 使用mysql驱动
)

// Conn 数据库句柄
var Conn *sql.DB

func init() {
	var err error
	Conn, err = sql.Open("mysql", "root:root@/blog")
	if err != nil {
		panic(err)
	}
}

// Encrypt sh1加密
func Encrypt(s string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(s)))
}
