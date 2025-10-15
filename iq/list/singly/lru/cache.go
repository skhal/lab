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
	size     int
	capacity int
}

type node struct {
	k, v int
	next *node
}

func NewCache(capacity int) (*Cache, error) {
	if capacity < CapMinimum {
		return nil, ErrCapacity
	}
	c := &Cache{
		capacity: capacity,
	}
	return c, nil
}

func (c *Cache) Put(k, v int) {
	c.add(k, v)
	if c.size > c.capacity {
		c.dropLeastRecent()
	}
}

func (c *Cache) add(k, v int) {
	n := &node{
		k: k,
		v: v,
	}
	switch c.head {
	case nil:
		c.head = n
	default:
		c.tail.next = n
	}
	c.tail = n
	c.size++
}

func (c *Cache) dropLeastRecent() {
	next := c.head.next
	c.head.next = nil
	c.head = next
	c.size--
}

func (c *Cache) Get(k int) (v int, ok bool) {
	for prev, curr := (*node)(nil), c.head; curr != nil; curr = curr.next {
		if curr.k == k {
			v = curr.v
			ok = true
			c.makeMostRecent(prev)
			break
		}
		prev = curr
	}
	return
}

func (c *Cache) makeMostRecent(prev *node) {
	if c.size == 1 {
		return
	}
	if prev == nil {
		c.makeHeadMostRecent()
		return
	}
	curr := prev.next
	if curr == c.tail {
		return
	}
	prev.next = curr.next
	curr.next = nil
	c.tail.next = curr
	c.tail = curr
}

func (c *Cache) makeHeadMostRecent() {
	curr := c.head
	c.head = curr.next
	curr.next = nil
	c.tail.next = curr
	c.tail = curr
}

func (c *Cache) String() string {
	return fmt.Sprintf("%s", c.Items())
}

func (c *Cache) Items() []Item {
	var ii []Item
	for n := c.head; n != nil; n = n.next {
		ii = append(ii, Item{n.k, n.v})
	}
	return ii
}

type Item struct {
	K, V int
}

func (i Item) String() string {
	return fmt.Sprintf("%d:%d", i.K, i.V)
}
