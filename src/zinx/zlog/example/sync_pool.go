package main

import (
	"encoding/json"
	"fmt"
	"sync"
)

type Student struct {
	Name   string
	Age    int32
	Remark [1024]byte
}

// 转换成json
var buf, _ = json.Marshal(Student{
	Name:   "binwen",
	Age:    22,
	Remark: [1024]byte{},
})

// 转换成对象
func unmarsh() {
	stu := &Student{}
	json.Unmarshal(buf, stu)
	fmt.Println(stu)
}

// 创建一个Student对象
var studentPool = sync.Pool{
	New: func() interface{} {
		return new(Student)
	},
}

func main() {
	fmt.Println(buf)
	unmarsh()
}
