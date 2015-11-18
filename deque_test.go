package deque

import "testing"

var (
	nItems = blockLen*9 + 7
)

func TestEmpty(t *testing.T) {
	q := New()
	if q.Len() != 0 {
		t.Fatalf("Empty queue with size: %d", q.Len())
	}
}

func newFilled(n int) *Deque {
	dq := New()
	for i := 0; i < nItems; i++ {
		dq.Append(i)
	}
	return dq
}

func TestAppend(t *testing.T) {
	dq := newFilled(nItems)
	if dq.Len() != nItems {
		t.Fatalf("Size after %d Append is %d", nItems, dq.Len())
	}
}

func TestAppendLeft(t *testing.T) {
	dq := New()
	for i := 0; i < nItems; i++ {
		dq.AppendLeft(i)
	}
	if dq.Len() != nItems {
		t.Fatalf("Size after %d AppendLeft is %d", nItems, dq.Len())
	}
}

func TestPop(t *testing.T) {
	dq := newFilled(nItems)

	count := 0
	for dq.Len() > 0 {
		count += 1
		item, err := dq.Pop()
		if err != nil {
			t.Fatalf("[%d] error Pop: %v", count, err)
		}
		val := item.(int)
		expected := nItems - count
		if val != expected {
			t.Fatalf("Pop %d, expected %d", val, expected)
		}
	}

	_, err := dq.Pop()
	if err == nil {
		t.Fatalf("Pop from empty Deque")
	}
	if count != nItems {
		t.Fatalf("Pop %d items, expected %d", count, nItems)
	}
}

func TestPopLeft(t *testing.T) {
	dq := newFilled(nItems)

	count := 0
	for dq.Len() > 0 {
		count += 1
		item, err := dq.PopLeft()
		if err != nil {
			t.Fatalf("[%d] error PopLeft: %v", count, err)
		}
		val := item.(int)
		expected := count - 1
		if val != expected {
			t.Fatalf("PopLeft %d, expected %d", val, expected)
		}
	}

	_, err := dq.PopLeft()
	if err == nil {
		t.Fatalf("PopLeft succeeded from empty Deque")
	}
	if count != nItems {
		t.Fatalf("PopLeft %d items, expected %d", count, nItems)
	}
}

func TestGet(t *testing.T) {
	dq := newFilled(nItems)

	for _, idx := range []int{0, dq.Len() / 2, dq.Len() - 1} {
		item, err := dq.Get(idx)
		if err != nil {
			t.Fatalf("Get(%d) failed on deque sized %d - %s", idx, dq.Len(), err)
		}
		i := item.(int)
		if i != idx {
			t.Fatalf("Get(%d) = %d, expected %d", idx, i, idx)
		}
	}

	idx := -7
	_, err := dq.Get(idx)
	if err == nil {
		t.Fatalf("Get(%d) didn't return error", idx)
	}

	idx = dq.Len() + 1
	_, err = dq.Get(idx)
	if err == nil {
		t.Fatalf("Get(%d) of queue with len %d didn't return error", idx, dq.Len())
	}
}

func TestSet(t *testing.T) {
	dq := newFilled(nItems)
	val := "hello"

	idx := dq.Len() / 2
	if err := dq.Set(idx, val); err != nil {
		t.Fatalf("Set(%d, %s) of deque size %d returned error - %s",
			idx, val, dq.Len(), err)
	}

	item, _ := dq.Get(idx)
	val2 := item.(string)
	if val != val2 {
		t.Fatalf("%s != %s at index %d", val, val2, idx)
	}
}

func TestBounded(t *testing.T) {
	size := 7
	dq, err := NewBounded(size)
	if err != nil {
		t.Fatalf("NewBounded(%d) failed - %s", size, err)
	}
	for i := 0; i < size; i++ {
		dq.Append(i)
		if dq.Len() != i+1 {
			t.Fatalf("bad size after %d append - %d", i+1, dq.Len())
		}
	}

	dq.Append(100)
	if dq.Len() != size {
		t.Fatalf("size = %d (should be %d)", dq.Len(), size)
	}
	val, _ := dq.Get(0)
	i := val.(int)
	if i != 1 {
		t.Fatalf("Get(0) -> %d (should be %d)", i, 1)
	}

	dq.AppendLeft(10)
	dq.AppendLeft(20)
	idx := dq.Len() - 1
	val, _ = dq.Get(idx)
	i = val.(int)
	if i != idx-1 {
		t.Fatalf("Get(%d) -> %d (should be %d)", idx, i, idx-1)
	}
}

type Point struct {
	X, Y int
}

func BenchmarkAppend(b *testing.B) {
	dq := New()
	for i := 0; i < b.N; i++ {
		dq.Append(&Point{i, i})
	}
}

func BenchmarkAppendLeft(b *testing.B) {
	dq := New()
	for i := 0; i < b.N; i++ {
		dq.AppendLeft(&Point{i, i})
	}
}

func BenchmarkPop(b *testing.B) {
	dq := New()
	for i := 0; i < b.N; i++ {
		dq.Append(&Point{i, i})
	}
	for dq.Len() > 0 {
		dq.Pop()
	}
}

func BenchmarkPopLeft(b *testing.B) {
	dq := New()
	for i := 0; i < b.N; i++ {
		dq.Append(&Point{i, i})
	}
	for dq.Len() > 0 {
		dq.PopLeft()
	}
}
