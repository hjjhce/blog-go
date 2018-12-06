package main

import (
	"os"
	"time"
)

func createFileName() string {
	t := time.Now().Format("2006-01-02")
	return t + ".log"
}

func createLogFile() *os.File {
	filename := createFileName()

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	return f
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
