package main

import "sync"

type Cache struct {
	mutex sync.RWMutex
	store map[string][]byte
}

func NewCache() *Cache {
	return &Cache{
		store: make(map[string][]byte),
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	value, ok := c.store[key]
	return value, ok
}

func (c *Cache) Set(key string, value []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.store[key] = value
}

func (c *Cache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.store, key)
}

func (c *Cache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.store = make(map[string][]byte)
}
