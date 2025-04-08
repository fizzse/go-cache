package gocache

import (
	"fmt"
	"testing"
	"time"
)

func TestSet(t *testing.T) {
	Set("name", "zhangsan", time.Second)
	a, _ := Get("name")
	fmt.Println(a)
	a, _ = GetString("name")
	fmt.Println(a)
}

func TestGetContainer(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}

	_ = Set("user", User{Name: "Simon", Age: 30}, 1*time.Second)

	// 获取缓存到结构体变量
	var u User
	err := GetContainer("user", &u)
	if err != nil {
		fmt.Println("Get error:", err)
	} else {
		fmt.Println("User from cache:", u)
	}
}
