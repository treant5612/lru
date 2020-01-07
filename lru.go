package lru

import (
	"math/rand"
	"time"
)

type node struct {
	key   string
	value interface{}
	prev  *node
	next  *node

	setup  time.Time
	expire time.Duration
}

type Map struct {
	head  *node
	tail  *node
	cache map[string]*node

	size int
	cap  int
}

func NewMap(capacity int) *Map {
	return &Map{
		cache: make(map[string]*node),
		cap:   capacity,
	}
}

func (m *Map) Get(key string) (value interface{}, ok bool) {
	node, ok := m.cache[key]
	if !ok {
		return nil, false
	}
	if node.expire != -1 && time.Since(node.setup) > node.expire {
		m.Delete(key)
		return nil, false
	}
	m.refresh(key)
	return node.value, ok
}

func (m *Map) Delete(key string) {
	node, ok := m.cache[key]
	if !ok {
		return
	}
	delete(m.cache, key)
	m.size--
	if node == m.head || node == m.tail {
		if node == m.head {
			m.head = node.next
			if m.head != nil {
				m.head.prev = nil
			}
		}
		if node == m.tail {
			if node.prev != nil {
				node.prev.next = nil
			}
			m.tail = node.prev
		}
	} else {
		prev := node.prev
		next := node.next
		prev.next, next.prev = next, prev
	}
}

func (m *Map) Set(key string, value interface{}, expire time.Duration) {
	if n, ok := m.cache[key]; ok {
		n.setup = time.Now()
		n.expire = expire
		n.value = value
		m.refresh(key)
		return
	}

	newNode := &node{
		key:    key,
		value:  value,
		setup:  time.Now(),
		expire: expire,
	}

	m.insert(newNode)
	m.cache[key] = newNode
}

func (m *Map) insert(node *node) {
	m.size++
	if m.head == nil {
		m.head = node
		m.tail = node
		return
	}

	m.head.prev = node
	node.next = m.head
	m.head = node
	m.head.prev = nil

	m.trim()
}

func (m *Map) refresh(key string) {
	node, ok := m.cache[key]
	if !ok || node == nil {
		panic("unexpected operation")
	}
	if m.head == node {
		return
	}
	prev := node.prev
	next := node.next
	prev.next = next
	if next != nil {
		next.prev = prev
	}

}

func (m *Map) trim() {
	if m.cap <= 0 {
		return
	}
	for m.size >= m.cap {
		m.Delete(m.tail.key)
	}

	if rand.Float32()<.2{
		m.fullTrim()
	}
}

func init(){
	rand.Seed(time.Now().UnixNano())
}
func( m *Map) fullTrim(){
	for key,node :=range m.cache{
		if node.expire>0 && time.Since(node.setup)>node.expire{
			m.Delete(key)
		}
	}
}