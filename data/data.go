package data

import (
	"crypto/rand"
	"crypto/sha1"
	"database/sql"
	"fmt"
	"math/big"
	mr "math/rand"
	"time"

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

// 生成随机数字字符串
func randNumStr() string {
	max := new(big.Int).Lsh(big.NewInt(1), 128)
	s, _ := rand.Int(rand.Reader, max)
	return fmt.Sprintf("%s", s)
}

// randNum 生成随机数字
func randNum(max int) (n int) {
	mr.Seed(time.Now().Unix())
	return mr.Intn(max)
}
