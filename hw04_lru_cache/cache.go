package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	sync.Mutex
}

type cacheItem struct {
	key   Key
	value interface{}
}

func (cache *lruCache) Set(key Key, value interface{}) bool {
	cache.Lock()

	item, ok := cache.items[key]
	if ok {
		cacheItem := item.Value.(cacheItem)
		cacheItem.value = value
		item.Value = cacheItem
		cache.queue.MoveToFront(item)
	} else {
		if cache.capacity == cache.queue.Len() {
			cacheItem := cache.queue.Back().Value.(cacheItem)
			cache.queue.Remove(cache.queue.Back())
			delete(cache.items, cacheItem.key)
		}
		cacheItem := cacheItem{key, value}
		listItem := cache.queue.PushFront(cacheItem)
		cache.items[key] = listItem
	}

	cache.Unlock()
	return ok
}

func (cache *lruCache) Get(key Key) (val interface{}, ok bool) {
	cache.Lock()

	item, ok := cache.items[key]
	if ok {
		cache.queue.MoveToFront(item)
		val = item.Value.(cacheItem).value
	} else {
		val = nil
	}

	cache.Unlock()

	return val, ok
}

func (cache *lruCache) Clear() {
	cache.Lock()
	cache.queue = NewList()
	cache.items = make(map[Key]*ListItem, cache.capacity)
	cache.Unlock()
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
