package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	Head   *ListItem
	Tail   *ListItem
	length int
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.Head
}

func (l *list) Back() *ListItem {
	return l.Tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	l.length++

	var item ListItem
	item.Value = v

	if l.Head == nil {
		l.Head = &item
		l.Tail = &item
		return &item
	}

	item.Next = l.Head
	l.Head.Prev = &item
	l.Head = &item

	return &item
}

func (l *list) PushBack(v interface{}) *ListItem {
	if l.Tail == nil {
		return l.PushFront(v)
	}

	l.length++

	var item ListItem
	item.Value = v
	item.Prev = l.Tail
	l.Tail.Next = &item
	l.Tail = &item

	return &item
}

func (l *list) Remove(i *ListItem) {
	if i == nil {
		return
	}

	l.length--

	if l.Len() == 0 {
		l.Head = nil
		l.Tail = nil
		return
	}

	if i == l.Head {
		l.Head = i.Next
		l.Head.Prev = nil
		return
	}

	if i == l.Tail {
		l.Tail = i.Prev
		l.Tail.Next = nil
		return
	}

	prevElem := i.Prev
	nextElem := i.Next

	if prevElem != nil {
		prevElem.Next = nextElem
	}
	if nextElem != nil {
		nextElem.Prev = prevElem
	}
}

func (l *list) MoveToFront(i *ListItem) {
	if i == l.Head {
		return
	}

	prevElem := i.Prev
	nextElem := i.Next

	if prevElem != nil {
		prevElem.Next = nextElem
	}
	if nextElem != nil {
		nextElem.Prev = prevElem
	}

	if i == l.Back() {
		l.Tail = prevElem
	}

	i.Prev = nil
	i.Next = l.Head

	l.Head.Prev = i
	l.Head = i
}

func NewList() List {
	return new(list)
}
