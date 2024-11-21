package cache

import (
	"errors"
	"log"
)

type Cacher interface {
	Get(key string) (string, error)
	Put(key string, val string) (string, error)
	Del(key string) (string, error)
	Size() int64
}

type Node struct {
	Key   string
	Left  *Node
	Val   string
	Right *Node
}

type Queue struct {
	Head *Node
	Tail *Node
	Len  int64
}

type Hash map[string]*Node

type Cache struct {
	queue Queue
	dict  Hash
	limit int64
}

func NewCache(limit int64) *Cache {
	return &Cache{
		queue: newQueue(),
		dict:  make(map[string]*Node),
		limit: limit,
	}
}

func newQueue() Queue {
	head := newNode("", "", nil, nil)
	tail := newNode("", "", nil, nil)

	head.Right = tail
	tail.Left = head

	return Queue{
		Head: head,
		Tail: tail,
		Len:  0,
	}
}

func newNode(key string, val string, left *Node, right *Node) *Node {
	return &Node{
		Key:   key,
		Val:   val,
		Left:  left,
		Right: right,
	}
}

func (l *Cache) Put(key string, val string) (string, error) {
	// check if already exists in our cache
	node, ok := l.dict[key]

	if ok {
		// remove that node from our cache
		node.Left.Right = node.Right
		node.Right.Left = node.Left
		delete(l.dict, key)
		l.queue.Len -= 1
	}

	if l.queue.Len == l.limit {
		log.Print("size full performing [DELETE] operation on TAIL")

		// Delete from tail
		nodeToDel := l.queue.Tail.Left

		nodeToDel.Left.Right = nodeToDel.Right
		nodeToDel.Right.Left = nodeToDel.Left
		delete(l.dict, nodeToDel.Key)
		l.queue.Len -= 1
	}

	// create new Node and add to right of head
	node = newNode(key, val, nil, nil)

	node.Right = l.queue.Head.Right
	node.Left = l.queue.Head

	l.queue.Head.Right.Left = node
	l.queue.Head.Right = node

	l.queue.Len += 1

	// Add to dict
	l.dict[key] = node

	return "OK", nil
}

func (l *Cache) Get(key string) (string, error) {
	node, ok := l.dict[key]

	if !ok {
		return "", errors.New("not found")
	}

	return node.Val, nil
}

func (l *Cache) Del(key string) (string, error) {
	node, ok := l.dict[key]

	if !ok {
		return "", errors.New("key not found")
	}

	node.Left.Right = node.Right
	node.Right.Left = node.Left

	l.queue.Len -= 1

	// Delete from dict
	delete(l.dict, key)

	return "OK", nil
}

func (l *Cache) Size() int64 {
	return l.queue.Len
}
