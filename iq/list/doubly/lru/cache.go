// Copyright 2025 Samvel Khalatyan. All rights reserved.

package lru

import (
	"errors"
	"fmt"
)

var ErrCapacity = errors.New("invalid capacity")

const CapMinimum = 1

type Cache struct {
	head     *node
	tail     *node
	nodes    map[int]*node
	size     int
	capacity int
}

type node struct {
	key  int
	val  int
	prev *node
	next *node
}

func NewCache(capacity int) (*Cache, error) {
	if capacity < CapMinimum {
		return nil, ErrCapacity
	}
	c := &Cache{
		nodes:    make(map[int]*node),
		capacity: capacity,
	}
	return c, nil
}

func (c *Cache) Put(k, v int) {
	if c.has(k) {
		c.overwrite(k, v)
		return
	}
	c.add(k, v)
}

func (c *Cache) has(k int) bool {
	_, ok := c.nodes[k]
	return ok
}

func (c *Cache) overwrite(k, v int) {
	n := c.nodes[k]
	n.val = v
	c.makeMostRecent(n)
}

func (c *Cache) add(k, v int) {
	if c.size == c.capacity {
		c.dropLeastRecent()
	} else {
		c.size++
	}
	n := &node{
		key: k,
		val: v,
	}
	c.nodes[k] = n
	switch c.head {
	case nil:
		c.head = n
	default:
		c.tail.next = n
		n.prev = c.tail
	}
	c.tail = n
}

func (c *Cache) dropLeastRecent() {
	n := c.head
	c.head = n.next
	n.next = nil
	c.head.prev = nil
	delete(c.nodes, n.key)
}

func (c *Cache) Get(k int) (v int, ok bool) {
	n, ok := c.nodes[k]
	if !ok {
		return
	}
	c.makeMostRecent(n)
	return n.val, true
}

func (c *Cache) makeMostRecent(n *node) {
	if c.size == 1 {
		return
	}
	if n == c.tail {
		return
	}
	if n == c.head {
		c.makeHeadMostRecent()
		return
	}
	prev := n.prev
	next := n.next
	prev.next = next
	next.prev = prev
	n.next = nil
	n.prev = c.tail
	c.tail.next = n
	c.tail = n
}

func (c *Cache) makeHeadMostRecent() {
	n := c.head
	c.head = n.next
	n.next = nil
	c.head.prev = nil
	c.tail.next = n
	n.prev = c.tail
	c.tail = n
}

func (c *Cache) String() string {
	return fmt.Sprintf("%s", c.Items())
}

func (c *Cache) Items() []Item {
	var ii []Item
	for n := c.head; n != nil; n = n.next {
		ii = append(ii, Item{n.key, n.val})
	}
	return ii
}

type Item struct {
	K, V int
}

func (i Item) String() string {
	return fmt.Sprintf("%d:%d", i.K, i.V)
}
