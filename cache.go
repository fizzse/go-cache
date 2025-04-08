package gocache

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"
)

type Item struct {
	Object     any
	Expiration int64
}

// Returns true if the item has expired.
func (item Item) Expired() bool {
	if item.Expiration == 0 {
		return false
	}
	return time.Now().Unix() > item.Expiration
}

const (
	NoExpiration time.Duration = -1
)

var (
	ErrorKeyNotExists = fmt.Errorf("key not exists")
	ErrorKeyExpired   = fmt.Errorf("key had expired")
	ErrorMismatch     = fmt.Errorf("val type mismatch")
)

type Cache struct {
	defaultExpiration time.Duration
	items             map[string]*Item
	mu                sync.RWMutex
}

var (
	entity = &Cache{items: make(map[string]*Item)}
)

func SetExpiration(duration time.Duration) {
	entity.defaultExpiration = duration
}

func Set(key string, value any, expire time.Duration) (err error) {
	entity.mu.Lock()
	defer entity.mu.Unlock()

	expiration := time.Now().Unix() + int64(expire.Seconds())
	if expire == NoExpiration {
		expiration = 0
	}

	entity.items[key] = &Item{Object: value, Expiration: expiration}
	return
}

func Del(key string) (err error) {
	entity.mu.Lock()
	defer entity.mu.Unlock()
	delete(entity.items, key)
	return
}

func Get(key string) (value any, err error) {
	entity.mu.RLock()
	defer entity.mu.RUnlock()

	item, exists := entity.items[key]
	if !exists {
		err = ErrorKeyNotExists
		return
	}

	if item.Expired() {
		err = ErrorKeyExpired
		return
	}

	value = item.Object
	return
}

func GetString(key string) (str string, err error) {
	value, err := Get(key)
	if err != nil {
		return
	}

	str, ok := value.(string)
	if !ok {
		err = ErrorMismatch
		return
	}

	return
}

func GetObject(key string, container any) (err error) {
	value, err := Get(key)
	if err != nil {
		return
	}

	err = ReflectVal(value, container)
	return
}

func ReflectVal(value any, container any) (err error) {
	// 获取目标容器的反射值
	dest := reflect.ValueOf(container)

	// 检查目标容器必须是指针类型
	if dest.Kind() != reflect.Ptr {
		return errors.New("container must be a pointer")
	}

	// 获取容器指针指向的元素
	destElem := dest.Elem()

	// 获取缓存值的反射值
	src := reflect.ValueOf(value)

	// 如果缓存值是指针，尝试解引用它
	if src.Kind() == reflect.Ptr && src.Type().Elem() == destElem.Type() {
		src = src.Elem()
	}

	// 类型匹配检查
	if !src.Type().AssignableTo(destElem.Type()) {
		return fmt.Errorf("type mismatch: cached value type is %v, but container expects %v", src.Type(), destElem.Type())
	}

	// 如果目标是指向指针的指针（例如 var u *User）,
	// 我们需要先通过反射为它分配一个新的实例
	if destElem.Kind() == reflect.Ptr && src.Kind() != reflect.Ptr {
		destElem.Set(reflect.New(destElem.Type().Elem()))
	}

	// 设置值
	destElem.Set(src)
	return
}
