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

func TestGetObject(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}

	_ = Set("user", User{Name: "Simon", Age: 30}, 1*time.Second)

	// 获取缓存到结构体变量
	var u User
	err := GetObject("user", &u)
	if err != nil {
		fmt.Println("Get error:", err)
	} else {
		fmt.Println("User from cache:", u)
	}
}

func TestList(t *testing.T) {
	key := "list"

	for i := 0; i < 5; i++ {
		LPush(key, i+1)
	}

	for i := 0; i < 6; i++ {
		fmt.Println(LPop(key))
	}
	for i := 0; i < 5; i++ {
		LPush(key, i+1)
	}

	for i := 0; i < 6; i++ {
		fmt.Println(RPop(key))
	}

	for i := 0; i < 5; i++ {
		RPush(key, i+1)
	}

	for i := 0; i < 6; i++ {
		fmt.Println(RPop(key))
	}
}
