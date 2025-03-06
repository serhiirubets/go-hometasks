package profile

import (
	"cache/internal/order"
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	mu   sync.RWMutex
	data map[string]*CacheItem
	ttl  time.Duration
}

type CacheItem struct {
	profile  *Profile
	expireAt time.Time
}

func NewCache(ttl time.Duration) *Cache {
	cache := &Cache{
		data: make(map[string]*CacheItem),
		ttl:  ttl,
		mu:   sync.RWMutex{},
	}
	go cache.cleanup()
	return cache
}

func (c *Cache) Get(uuid string) Profile {
	c.mu.RLock()
	defer c.mu.RUnlock()

	p, ok := c.data[uuid]

	if !ok {
		return Profile{}
	}

	if time.Now().After(p.expireAt) {
		return Profile{}
	}

	return *p.profile
}

func (c *Cache) Set(uuid string, p *Profile) {
	c.mu.Lock()
	defer c.mu.Unlock()
	fmt.Println(p.Orders)

	profileCopy := *p

	profileCopy.Orders = make([]*order.Order, len(p.Orders))

	// Note: if inside orders pointers, it won't copy them deeply
	// Need to implement deep clone
	// Not in the scope of current task
	for i, ord := range p.Orders {
		if ord == nil {
			profileCopy.Orders[i] = nil
			continue
		}

		profileCopy.Orders[i] = &order.Order{
			UUID:      ord.UUID,
			Value:     ord.Value,
			CreatedAt: ord.CreatedAt,
			UpdatedAt: ord.UpdatedAt,
		}
	}

	c.data[uuid] = &CacheItem{
		profile:  &profileCopy,
		expireAt: time.Now().Add(c.ttl),
	}
}

func (c *Cache) cleanup() {
	ticker := time.NewTicker(c.ttl)
	defer ticker.Stop()

	for range ticker.C {
		var toDelete []string

		c.mu.RLock()
		for uuid, item := range c.data {
			if time.Now().After(item.expireAt) {
				toDelete = append(toDelete, uuid)
			}
		}
		c.mu.RUnlock()

		// Delete only expired profiles from cache, but takes additional memory
		// for storing them.
		c.mu.Lock()
		fmt.Println(toDelete)
		for _, uuid := range toDelete {
			delete(c.data, uuid)
		}
		c.mu.Unlock()
	}
}

func (c *Cache) Delete(uuid string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, uuid)
}
