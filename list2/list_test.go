package list2

import (
	"reflect"
	"testing"
)

func checkListLen(t *testing.T, l *List, len int) bool {
	if n := l.Len(); n != len {
		t.Errorf("l.Len() = %d, want %d", n, len)
		return false
	}
	return true
}

func checkListPointers(t *testing.T, l *List, es []*Element) {

	root := l.root

	if !checkListLen(t, l, len(es)) {
		return
	}

	// zero length lists must be the zero value or properly initialized (sentinel circle)
	if len(es) == 0 {
		if l.root.next != 0 && l.root.next != root.index || l.root.prev != 0 && l.root.prev != root.index {
			t.Errorf("l.root.next = %d, l.root.prev = %d; both should both be 0 or %d", l.root.next, l.root.prev, root.index)
		}
		return
	}
	// len(es) > 0

	e1 := []*Element{}
	for e := l.Front(); e != nil; e = e.Next() {
		e1 = append(e1, e)
	}

	if !reflect.DeepEqual(e1, es) {
		t.Fatal("front not equal")
	}

	e2 := []*Element{}
	for e := l.Back(); e != nil; e = e.Prev() {
		e2 = append([]*Element{e}, e2...)
	}

	if !reflect.DeepEqual(e2, es) {
		t.Fatal("front not equal")
	}

	// check internal and external prev/next connections
	for i, e := range es {
		prev := root
		Prev := (*Element)(nil)
		if i > 0 {
			prev = es[i-1]
			Prev = prev
		}
		if p := e.prev; p != prev.index {
			t.Errorf("elt[%d](%p).prev = %d, want %p", i, e, p, prev)
		}
		if p := e.Prev(); p != Prev {
			t.Errorf("elt[%d](%p).Prev() = %p, want %p", i, e, p, Prev)
		}

		next := root
		Next := (*Element)(nil)
		if i < len(es)-1 {
			next = es[i+1]
			Next = next
		}
		if n := e.next; n != next.index {
			t.Errorf("elt[%d](%p).next = %d, want %p", i, e, n, next)
		}
		if n := e.Next(); n != Next {
			t.Errorf("elt[%d](%p).Next() = %p, want %p", i, e, n, Next)
		}
	}
}

func TestList(t *testing.T) {
	l := New()
	checkListPointers(t, l, []*Element{})

	// Single element list
	e := l.PushFront("a")
	checkListPointers(t, l, []*Element{e})
	l.MoveToFront(e)
	checkListPointers(t, l, []*Element{e})
	l.MoveToBack(e)
	checkListPointers(t, l, []*Element{e})
	l.Remove(e)
	checkListPointers(t, l, []*Element{})

	// Bigger list
	e2 := l.PushFront(2)
	e1 := l.PushFront(1)
	e3 := l.PushBack(3)
	e4 := l.PushBack("banana")
	checkListPointers(t, l, []*Element{e1, e2, e3, e4})

	l.Remove(e2)
	checkListPointers(t, l, []*Element{e1, e3, e4})

	l.MoveToFront(e3) // move from middle
	checkListPointers(t, l, []*Element{e3, e1, e4})

	l.MoveToFront(e1)
	l.MoveToBack(e3) // move from middle
	checkListPointers(t, l, []*Element{e1, e4, e3})

	l.MoveToFront(e3) // move from back
	checkListPointers(t, l, []*Element{e3, e1, e4})
	l.MoveToFront(e3) // should be no-op
	checkListPointers(t, l, []*Element{e3, e1, e4})

	l.MoveToBack(e3) // move from front
	checkListPointers(t, l, []*Element{e1, e4, e3})
	l.MoveToBack(e3) // should be no-op
	checkListPointers(t, l, []*Element{e1, e4, e3})

	e2 = l.InsertBefore(2, e1) // insert before front
	checkListPointers(t, l, []*Element{e2, e1, e4, e3})
	l.Remove(e2)
	e2 = l.InsertBefore(2, e4) // insert before middle
	checkListPointers(t, l, []*Element{e1, e2, e4, e3})
	l.Remove(e2)
	e2 = l.InsertBefore(2, e3) // insert before back
	checkListPointers(t, l, []*Element{e1, e4, e2, e3})
	l.Remove(e2)

	e2 = l.InsertAfter(2, e1) // insert after front
	checkListPointers(t, l, []*Element{e1, e2, e4, e3})
	l.Remove(e2)
	e2 = l.InsertAfter(2, e4) // insert after middle
	checkListPointers(t, l, []*Element{e1, e4, e2, e3})
	l.Remove(e2)
	e2 = l.InsertAfter(2, e3) // insert after back
	checkListPointers(t, l, []*Element{e1, e4, e3, e2})
	l.Remove(e2)

	// Check standard iteration.
	sum := 0
	for e := l.Front(); e != nil; e = e.Next() {
		if i, ok := e.Value.(int); ok {
			sum += i
		}
	}
	if sum != 4 {
		t.Errorf("sum over l = %d, want 4", sum)
	}

	// Clear all elements by iterating
	var next *Element
	for e := l.Front(); e != nil; e = next {
		next = e.Next()
		l.Remove(e)
	}
	checkListPointers(t, l, []*Element{})
}

func checkList(t *testing.T, l *List, es []interface{}) {
	if !checkListLen(t, l, len(es)) {
		return
	}

	i := 0
	for e := l.Front(); e != nil; e = e.Next() {
		le := e.Value.(int)
		if le != es[i] {
			t.Errorf("elt[%d].Value = %v, want %v", i, le, es[i])
		}
		i++
	}
}

func TestExtending(t *testing.T) {
	l1 := New()
	l2 := New()

	l1.PushBack(1)
	l1.PushBack(2)
	l1.PushBack(3)

	l2.PushBack(4)
	l2.PushBack(5)

	l3 := New()
	l3.PushBackList(l1)
	checkList(t, l3, []interface{}{1, 2, 3})
	l3.PushBackList(l2)
	checkList(t, l3, []interface{}{1, 2, 3, 4, 5})

	l3 = New()
	l3.PushFrontList(l2)
	checkList(t, l3, []interface{}{4, 5})
	l3.PushFrontList(l1)
	checkList(t, l3, []interface{}{1, 2, 3, 4, 5})

	checkList(t, l1, []interface{}{1, 2, 3})
	checkList(t, l2, []interface{}{4, 5})

	l3 = New()
	l3.PushBackList(l1)
	checkList(t, l3, []interface{}{1, 2, 3})
	l3.PushBackList(l3)
	checkList(t, l3, []interface{}{1, 2, 3, 1, 2, 3})

	l3 = New()
	l3.PushFrontList(l1)
	checkList(t, l3, []interface{}{1, 2, 3})
	l3.PushFrontList(l3)
	checkList(t, l3, []interface{}{1, 2, 3, 1, 2, 3})

	l3 = New()
	l1.PushBackList(l3)
	checkList(t, l1, []interface{}{1, 2, 3})
	l1.PushFrontList(l3)
	checkList(t, l1, []interface{}{1, 2, 3})
}

func TestRemove(t *testing.T) {
	l := New()
	e1 := l.PushBack(1)
	e2 := l.PushBack(2)
	checkListPointers(t, l, []*Element{e1, e2})
	e := l.Front()
	l.Remove(e)
	checkListPointers(t, l, []*Element{e2})
	l.Remove(e)
	checkListPointers(t, l, []*Element{e2})
}

func TestIssue4103(t *testing.T) {
	l1 := New()
	l1.PushBack(1)
	l1.PushBack(2)

	l2 := New()
	l2.PushBack(3)
	l2.PushBack(4)

	e := l1.Front()
	l2.Remove(e) // l2 should not change because e is not an element of l2
	if n := l2.Len(); n != 2 {
		t.Errorf("l2.Len() = %d, want 2", n)
	}

	l1.InsertBefore(8, e)
	if n := l1.Len(); n != 3 {
		t.Errorf("l1.Len() = %d, want 3", n)
	}
}

func TestIssue6349(t *testing.T) {
	l := New()
	l.PushBack(1)
	l.PushBack(2)

	e := l.Front()

	//because element may be reused later, we must clear value in Remove
	v := l.Remove(e)

	if v != 1 {
		t.Errorf("e.value = %d, want 1", e.Value)
	}

	if e.Next() != nil {
		t.Errorf("e.Next() != nil")
	}
	if e.Prev() != nil {
		t.Errorf("e.Prev() != nil")
	}
}

func TestMove(t *testing.T) {
	l := New()
	e1 := l.PushBack(1)
	e2 := l.PushBack(2)
	e3 := l.PushBack(3)
	e4 := l.PushBack(4)

	l.MoveAfter(e3, e3)
	checkListPointers(t, l, []*Element{e1, e2, e3, e4})
	l.MoveBefore(e2, e2)
	checkListPointers(t, l, []*Element{e1, e2, e3, e4})

	l.MoveAfter(e3, e2)
	checkListPointers(t, l, []*Element{e1, e2, e3, e4})
	l.MoveBefore(e2, e3)
	checkListPointers(t, l, []*Element{e1, e2, e3, e4})

	l.MoveBefore(e2, e4)
	checkListPointers(t, l, []*Element{e1, e3, e2, e4})
	e1, e2, e3, e4 = e1, e3, e2, e4

	l.MoveBefore(e4, e1)
	checkListPointers(t, l, []*Element{e4, e1, e2, e3})
	e1, e2, e3, e4 = e4, e1, e2, e3

	l.MoveAfter(e4, e1)
	checkListPointers(t, l, []*Element{e1, e4, e2, e3})
	e1, e2, e3, e4 = e1, e4, e2, e3

	l.MoveAfter(e2, e3)
	checkListPointers(t, l, []*Element{e1, e3, e2, e4})
	e1, e2, e3, e4 = e1, e3, e2, e4
}
