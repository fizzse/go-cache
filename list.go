package gocache

import (
	"container/list"
	"fmt"
)

// 使用list容器
func LPush(key string, val any) (err error) {
	entity.mu.Lock()
	defer entity.mu.Unlock()

	l := list.New()
	ok := false
	item, ok := entity.items[key]
	if ok {
		if item.Expired() {
			err = ErrorKeyExpired
			return
		}

		l, ok = item.Object.(*list.List)
		if !ok {
			err = ErrorMismatch
			return
		}
	}

	l.PushBack(val)
	entity.items[key] = &Item{Object: l}
	return
}

func RPush(key string, val any) (err error) {
	// 使用list容器
	entity.mu.Lock()
	defer entity.mu.Unlock()

	l := list.New()
	ok := false
	item, ok := entity.items[key]
	if ok {
		if item.Expired() {
			err = ErrorKeyExpired
			return
		}

		l, ok = item.Object.(*list.List)
		if !ok {
			err = ErrorMismatch
			return
		}
	}

	l.PushFront(val)
	entity.items[key] = &Item{Object: l}
	return
}

func LPop(key string) (val any, err error) {
	entity.mu.Lock()
	defer entity.mu.Unlock()
	item, ok := entity.items[key]
	if !ok {
		err = ErrorKeyNotExists
		return
	}

	if item.Expired() {
		err = ErrorKeyExpired
		return
	}

	l, ok := item.Object.(*list.List)
	if !ok {
		err = ErrorMismatch
		return
	}

	if l.Len() == 0 {
		err = fmt.Errorf("list is empty")
		return
	}

	val = l.Remove(l.Front())
	return
}

func RPop(key string) (val any, err error) {
	entity.mu.Lock()
	defer entity.mu.Unlock()
	item, ok := entity.items[key]
	if !ok {
		err = ErrorKeyNotExists
		return
	}

	if item.Expired() {
		err = ErrorKeyExpired
		return
	}

	l, ok := item.Object.(*list.List)
	if !ok {
		err = ErrorMismatch
		return
	}

	if l.Len() == 0 {
		err = fmt.Errorf("list is empty")
		return
	}

	val = l.Remove(l.Back())
	return
}
