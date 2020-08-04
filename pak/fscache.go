package pak

import (
	"hash/fnv"
	"sync"
)

type entry struct {
	key   string
	data  []byte
	mutex sync.RWMutex
}

type cache map[uint8]*entry

func makecache() cache {
	c := make(cache)
	for i := 0; i < 0x100; i++ {
		c[uint8(i)] = &entry{}
	}
	return c
}

func (c cache) set(key string, data []byte) {
	hash := fnv.New32()
	hash.Write([]byte(key))
	index := uint8(hash.Sum32() & 0xFF)
	c[index].mutex.Lock()
	defer c[index].mutex.Unlock()
	c[index].key = key
	c[index].data = data
}

func (c cache) get(key string) []byte {
	hash := fnv.New32()
	hash.Write([]byte(key))
	index := uint8(hash.Sum32() & 0xFF)
	c[index].mutex.RLock()
	defer c[index].mutex.RUnlock()
	if c[index].key != key {
		return nil
	}
	return c[index].data
}
