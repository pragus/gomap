package cache

import (
	"sync"
	"hash/fnv"
)

type CacheShard struct {
	mtx* sync.RWMutex
	items map[string]int
}

type Cache struct {
	segments int
	shards []*CacheShard
}


func New(scnt int) *Cache {
	c := Cache{segments: scnt, shards: make([]*CacheShard, scnt, scnt)}
	for i := 0; i < scnt; i++ {
		c.shards[i] = &CacheShard{
			mtx: &sync.RWMutex{},
			items: make(map[string]int),
		}

	}
	return &c
}

func (c *Cache) getShard(key string) *CacheShard {
	hasher := fnv.New64()
	hasher.Write([]byte(key))
	d := hasher.Sum(nil)
	idx := int(d[0])
	if c.segments > 256 {
	idx += 256*int((d[1]))
	}
	return c.shards[idx]

}

func (c *Cache) Set(key string, val int) {
	s := c.getShard(key)
	s.mtx.Lock()
	s.items[key] = val
	s.mtx.Unlock()
}

func (c *Cache) Get(key string) int {
	s := c.getShard(key)
	s.mtx.RLock()
	r := s.items[key]
	s.mtx.RUnlock()
	return r
}
