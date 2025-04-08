# go-cache 内存缓存
支持 string list zset等常规redis操作

## String

```go
Set("hello", "world", time.Second)
val, _ = Get("name")
val, _ = GetString("name")

# 直接返回struct
type User struct {
	Name string
	Age  int
}
Set("user", User{Name: "Simon", Age: 30}, 1*time.Second)
var u User
err := GetObject("user", &u)


```

## List
```go
LPush("list",1)
val, _ = LPop("list")
```

## Zset TODO