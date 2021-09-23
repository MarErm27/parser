package core

import (
	"errors"
	"sync"
)

var store = struct {
	sync.RWMutex
	m map[string]int
}{m: make(map[string]int)}

func Delete(key string) error {
	store.Lock()
	delete(store.m, key)
	store.Unlock()

	return nil
}

var ErrorNoSuchKey = errors.New("no such key")

func Get(key string) (int, error) {
	store.RLock()
	value, ok := store.m[key]
	store.RUnlock()

	if !ok {
		return 0, ErrorNoSuchKey
	}

	return value, nil
}

func GetAll() map[string]int {
	return store.m
}

func Put(key string, value int) error {
	store.Lock()
	store.m[key] = value
	store.Unlock()

	return nil
}

