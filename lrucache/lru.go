//go:build !solution

package lrucache

import "container/list"

func New(cap int) Cache {
	return &LRUCache{hashTable: make(map[int]*list.Element, cap), queue: list.New(), cap: cap}
}

func (l *LRUCache) Get(key int) (int, bool) {
	element, ok := l.hashTable[key]
	if ok {
		value := element.Value.(entry).value
		l.queue.MoveToFront(element)
		return value, ok
	}
	return 0, false
}

func (l *LRUCache) Set(key int, value int) {
	element, ok := l.hashTable[key]
	if ok {
		element.Value = entry{key, value}
		l.queue.MoveToFront(element)
		return
	}
	l.hashTable[key] = l.queue.PushFront(entry{key, value})
	if l.queue.Len() > l.cap {
		back := l.queue.Back()
		e := back.Value.(entry)
		delete(l.hashTable, e.key)
		l.queue.Remove(l.queue.Back())
	}
}

func (l *LRUCache) Range(f func(key, value int) bool) {
	element := l.queue.Back()
	for element != nil {
		if !f(element.Value.(entry).key, element.Value.(entry).value) {
			break
		}
		element = element.Prev()
	}
}

func (l *LRUCache) Clear() {
	element := l.queue.Back()
	for element != nil {
		delete(l.hashTable, element.Value.(entry).key)
		l.queue.Remove(element)
		element = l.queue.Back()
	}

}
