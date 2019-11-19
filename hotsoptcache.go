// Package hotspotcache stores values like sync.Map does.
// Meanwhile, it handles a aging queue: every access to a existing object in the cache,
// it will be lifted up to he top of the queue and marked as the newest element. When
// the size in cache exceeds preset limit size, the bottom elements of the queue would
// be removed.
package hotspotcache

import (
	"bytes"
	"container/list"
	"fmt"
	"runtime"
	"sync"
)

// Cache is the hotspot cache object.
type Cache struct {
	c *cache
}

type cache struct {
	maxSize   int
	values    sync.Map
	agingList *list.List // one in the front is the hottest one
	agingMap  map[interface{}]*list.Element
	access    chan interface{}
	stop      chan bool
}

// New returns a initialized cache with given size. When size of the cache exceeds maxSize, it will start cleaning.
func New(maxSize int) *Cache {
	if maxSize <= 0 {
		maxSize = 10240
	}

	c := newCache(maxSize)
	ret := Cache{c: c}

	go c.run()
	runtime.SetFinalizer(&ret, stopRunning)

	return &ret
}

// Load read a value in the cache. If value does not exist in cache, the return exist will be false, otherwise true. If value exists, its hotspot will also be risen to the top.
func (c *Cache) Load(key interface{}) (value interface{}, exist bool) {
	value, exist = c.c.values.Load(key)
	if exist {
		// log.Printf("Key %v not exists\n", key)
		c.c.access <- key
	}
	// log.Printf("Key %v exists\n", key)
	return
}

// Store saves a value by corresponding key.
func (c *Cache) Store(key, value interface{}) {
	c.c.values.Store(key, value)
	c.c.access <- key
	return
}

// MaxSize returns the max size of this cache
func (c *Cache) MaxSize() int {
	return c.c.maxSize
}

func newCache(maxSize int) *cache {
	c := cache{
		maxSize:   maxSize,
		agingList: list.New(),
		agingMap:  make(map[interface{}]*list.Element),
		access:    make(chan interface{}),
		stop:      make(chan bool),
	}
	return &c
}

func (c *Cache) dumpStatus() string {
	buff := bytes.Buffer{}
	buff.WriteString(fmt.Sprintf("\nMax size: %d", c.c.maxSize))
	buff.WriteString(fmt.Sprintf("\naging list length: %d", c.c.agingList.Len()))
	buff.WriteString(fmt.Sprintf("\nremain in channel: %d", len(c.c.access)))
	return buff.String()
}

func (c *cache) run() {
	for {
		select {
		case k := <-c.access:
			c.updateHotspot(k)
		case <-c.stop:
			return
		}
	}
}

func (c *cache) updateHotspot(k interface{}) {
	li := c.agingList
	e, exist := c.agingMap[k]

	// this key already in hot spotqueue
	if exist {
		// just move it to the front
		// log.Printf("upd hotspot %v\n", k)
		li.MoveToFront(e)
		return
	}

	// this key not exists in hotspot queue
	// log.Printf("new hotspot %v\n", k)
	e = li.PushFront(k)
	c.agingMap[k] = e

	// check if the queue is full
	if li.Len() <= c.maxSize {
		return
	}

	// remove the coldest one
	e = li.Back()
	k = e.Value
	c.values.Delete(k)
	li.Remove(e)
	delete(c.agingMap, k)
	// log.Printf("del hotspot %v\n", k)
	return
}

func stopRunning(c *Cache) {
	c.c.stop <- true
}
