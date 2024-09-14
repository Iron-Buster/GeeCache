package lru

import (
	"container/list"
	"fmt"
)

// Cache -> LRU Cache
type Cache struct {
	maxBytes  int64 // 允许使用的最大内存
	usedBytes int64 // 当前已使用的内存
	ll        *list.List
	cache     map[string]*list.Element
	OnEvicted func(key string, value Value) // entry被清除时的回调函数
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int // 返回值所占用的内存大小
}

// New 实例化Cache
func New(maxBytes int64, oneEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: oneEvicted,
	}
}

// Get Cache获取值的方法
func (c *Cache) Get(key string) (value Value, ok bool) {
	if node, ok := c.cache[key]; ok {
		c.ll.MoveToFront(node)
		kv := node.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest 删除最近最少访问的节点
func (c *Cache) RemoveOldest() {
	node := c.ll.Back()
	if node != nil {
		c.ll.Remove(node)
		kv := node.Value.(*entry)
		delete(c.cache, kv.key)
		c.usedBytes -= int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Add 新增缓存
func (c *Cache) Add(key string, value Value) {
	if node, ok := c.cache[key]; ok {
		c.ll.MoveToFront(node)
		kv := node.Value.(*entry)
		c.usedBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		node := c.ll.PushFront(&entry{key, value})
		c.cache[key] = node
		c.usedBytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.usedBytes {
		c.RemoveOldest()
	}
}

// Len 返回缓存的节点个数
func (c *Cache) Len() int {
	return c.ll.Len()
}

// PrintAllEntry 打印所有的entry
func (c *Cache) PrintAllEntry() {
	for e := c.ll.Front(); e != nil; e = e.Next() {
		kv := e.Value.(*entry)
		fmt.Printf("Key: %s, Value: %v\n", kv.key, kv.value)
	}
}
