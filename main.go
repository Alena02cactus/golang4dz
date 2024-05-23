package main

import (
	"fmt"
	"time"
)

type Cache interface {
	Get(k string) (string, bool)
	Set(k, v string)
	Cleanup()
}

var _ Cache = (*cacheImpl)(nil)

// Доработает конструктор и методы кеша, так чтобы они соответствовали интерфейсу Cache
func newCacheImpl() *cacheImpl {
	return &cacheImpl{
		    data:       make(map[string]string),
		    expiration: make(map[string]time.Time),
	}
}

type cacheImpl struct {
	data       map[string]string
	expiration map[string]time.Time
}

func (c *cacheImpl) Get(k string) (string, bool) {
	// TODO implement me
	c.Cleanup() 
	v, ok := c.data[k]
	return v, ok
}

func (c *cacheImpl) Set(k, v string) {
	// TODO implement me
	c.data[k] = v
	c.expiration[k] = time.Now().Add(5 * time.Minute)
}

func (c *cacheImpl) Cleanup() {
	now := time.Now()
	for k, exp := range c.expiration {
		if now.After(exp) {
			delete(c.data, k)
			delete(c.expiration, k)
		}
	}
}


func newDbImpl(cache Cache) *dbImpl {
	return &dbImpl{cache: cache, dbs: map[string]string{"hello": "world", "test": "test"}}
}

type dbImpl struct {
	cache Cache
	dbs   map[string]string
}

func (d *dbImpl) Get(k string) (string, bool) {
	v, ok := d.cache.Get(k)
	if ok {
		return fmt.Sprintf("answer from cache: key: %s, val: %s", k, v), ok
	}

	v, ok = d.dbs[k]
	return fmt.Sprintf("answer from dbs: key: %s, val: %s", k, v), ok
}

func main() {
	c := newCacheImpl()
	db := newDbImpl(c)
	db.cache.Set("key1", "value1")
	db.cache.Set("key2", "value2")
	time.Sleep(6 * time.Minute)
	fmt.Println(db.Get("test"))
	fmt.Println(db.Get("hello"))
}
