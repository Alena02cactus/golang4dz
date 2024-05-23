package main

import (
 "fmt"
 "time"
)

type Cache interface {
 Get(k string) (string, bool)
 Set(k, v string)
}

var _ Cache = (*cacheImpl)(nil)

// Доработаем конструктор и методы кеша, чтобы они соответствовали интерфейсу Cache
func newCacheImpl() *cacheImpl {
 return &cacheImpl{
  cache: make(map[string]cacheEntry),
 }
}

type cacheImpl struct {
 cache map[string]cacheEntry
}

type cacheEntry struct {
 value     string
 timestamp time.Time
}

func (c *cacheImpl) Get(k string) (string, bool) {
 entry, exists := c.cache[k]
 if exists {
  // Проверяем, не просрочен ли ключ
  if time.Now().Sub(entry.timestamp) > time.Minute {
   delete(c.cache, k)
   return "", false
  }
  return entry.value, true
 }
 return "", false
}

func (c *cacheImpl) Set(k, v string) {
 entry := cacheEntry{
  value:     v,
  timestamp: time.Now(),
 }
 c.cache[k] = entry
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
 fmt.Println(db.Get("test"))
 fmt.Println(db.Get("hello"))
}
