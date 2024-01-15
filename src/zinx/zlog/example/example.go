package main

import (
	"log"
	"os"
	"zinx/zlog"
)

func main() {
	zlog.Info("std log")
	zlog.SetOptions(zlog.WithLevel(zlog.DebugLevel))
	zlog.Debug("change std log to debug level")
	zlog.SetOptions(zlog.WithFormatter(&zlog.JsonFormatter{IgnoreBasicFields: false}))
	zlog.Debug("log in json format")
	zlog.Info("another log in json format")

	// 输出到文件
	fd, err := os.OpenFile("test.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln("create file test.log failed")
	}
	defer fd.Close()

	l := zlog.New(zlog.WithLevel(zlog.InfoLevel),
		zlog.WithOutput(fd),
		zlog.WithFormatter(&zlog.JsonFormatter{IgnoreBasicFields: false}),
	)
	l.Info("custom log with json formatter")
}
