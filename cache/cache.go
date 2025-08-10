package cache

import (
	"sync"
	"time"
)

const (
	NoExpiration      time.Duration = -1
	DefaultExpiration time.Duration = 0
)

type Cache struct {
	*cache
}

type cache struct {
	defaultExpiration time.Duration
	items            map[string]Item
	mu               sync.RWMutex
	janitor          *janitor
}

type Item struct {
	Object     interface{}
	Expiration int64
}

func (item Item) Expired() bool {
	if item.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > item.Expiration
}

func New(defaultExpiration, cleanupInterval time.Duration) *Cache {
	items := make(map[string]Item)
	cache := newCache(defaultExpiration, items)
	
	// 如果清理时间间隔 > 0，则启动定时清理任务
	if cleanupInterval > 0 {
		runJanitor(cache, cleanupInterval)
	}
	
	return &Cache{cache}
}

func newCache(de time.Duration, m map[string]Item) *cache {
	if de == 0 {
		de = -1
	}
	c := &cache{
		defaultExpiration: de,
		items:            m,
	}
	return c
}

func (c *cache) Set(k string, x interface{}, d time.Duration) {
	var e int64
	if d == DefaultExpiration {
		d = c.defaultExpiration
	}
	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}
	c.mu.Lock()
	c.items[k] = Item{
		Object:     x,
		Expiration: e,
	}
	c.mu.Unlock()
}

func (c *cache) Get(k string) (interface{}, bool) {
	c.mu.RLock()
	item, found := c.items[k]
	if !found {
		c.mu.RUnlock()
		return nil, false
	}
	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			c.mu.RUnlock()
			return nil, false
		}
	}
	c.mu.RUnlock()
	return item.Object, true
}

func (c *cache) Delete(k string) {
	c.mu.Lock()
	delete(c.items, k)
	c.mu.Unlock()
}

type janitor struct {
	Interval time.Duration
	stop     chan bool
}

func (j *janitor) Run(c *cache) {
	ticker := time.NewTicker(j.Interval)
	for {
		select {
		case <-ticker.C:
			c.DeleteExpired()
		case <-j.stop:
			ticker.Stop()
			return
		}
	}
}

func runJanitor(c *cache, ci time.Duration) {
	j := &janitor{
		Interval: ci,
		stop:     make(chan bool),
	}
	c.janitor = j
	go j.Run(c)
}

func (c *cache) DeleteExpired() {
	now := time.Now().UnixNano()
	c.mu.Lock()
	for k, v := range c.items {
		if v.Expiration > 0 && now > v.Expiration {
			delete(c.items, k)
		}
	}
	c.mu.Unlock()
}
