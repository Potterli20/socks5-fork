package socks5

import (
	"context"
	"errors"
	"sync"
	"time"
)

var ErrNotFound = errors.New("not found")

type LocalCache interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string, value interface{}) error
	Delete(ctx context.Context, key string) error
	IsNotFoundError(err error) bool
}

type memoryCache struct {
	mu    sync.RWMutex
	items map[string]*cacheItem
}

type cacheItem struct {
	value     interface{}
	expiresAt time.Time
}

func NewLocalCache() LocalCache {
	c := &memoryCache{
		items: make(map[string]*cacheItem),
	}
	go c.cleanup()
	return c
}

func (c *memoryCache) Set(_ context.Context, key string, value interface{}, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = &cacheItem{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
	return nil
}

func (c *memoryCache) Get(_ context.Context, key string, value interface{}) error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, ok := c.items[key]
	if !ok {
		return ErrNotFound
	}
	if time.Now().After(item.expiresAt) {
		return ErrNotFound
	}
	ptr, ok := value.(*interface{})
	if !ok {
		return errors.New("value must be a pointer to interface{}")
	}
	*ptr = item.value
	return nil
}

func (c *memoryCache) Delete(_ context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
	return nil
}

func (c *memoryCache) IsNotFoundError(err error) bool {
	return err == ErrNotFound
}

func (c *memoryCache) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, item := range c.items {
			if now.After(item.expiresAt) {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}
}