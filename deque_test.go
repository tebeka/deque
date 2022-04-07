package deque

import (
	"testing"
)

var (
	nItems = blockLen*9 + 7
)

func TestEmpty(t *testing.T) {
	q := New[int]()
	if q.Len() != 0 {
		t.Fatalf("Empty queue with size: %d", q.Len())
	}
}

func newFilled(n int) *Deque[int] {
	dq := New[int]()
	for i := 0; i < n; i++ {
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
	dq := New[int]()
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
		count++
		item, err := dq.Pop()
		if err != nil {
			t.Fatalf("[%d] error Pop: %v", count, err)
		}
		expected := nItems - count
		if item != expected {
			t.Fatalf("Pop %d, expected %d", item, expected)
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
		count++
		item, err := dq.PopLeft()
		if err != nil {
			t.Fatalf("[%d] error PopLeft: %v", count, err)
		}
		expected := count - 1
		if item != expected {
			t.Fatalf("PopLeft %d, expected %d", item, expected)
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
		if item != idx {
			t.Fatalf("Get(%d) = %d, expected %d", idx, item, idx)
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
	val := 17

	idx := dq.Len() / 2
	if err := dq.Set(idx, val); err != nil {
		t.Fatalf("Set(%d, %d) of deque size %d returned error - %s",
			idx, val, dq.Len(), err)
	}

	item, _ := dq.Get(idx)
	if val != item {
		t.Fatalf("%d != %d at index %d", val, item, idx)
	}
}

func TestBounded(t *testing.T) {
	size := 7
	dq, err := NewBounded[int](size)
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
	if val != 1 {
		t.Fatalf("Get(0) -> %d (should be %d)", val, 1)
	}

	dq.AppendLeft(10)
	dq.AppendLeft(20)
	idx := dq.Len() - 1
	val, _ = dq.Get(idx)
	if val != idx-1 {
		t.Fatalf("Get(%d) -> %d (should be %d)", idx, val, idx-1)
	}
}

func TestRotate(t *testing.T) {
	// [0, 1, 2, 3]
	dq := newFilled(4)
	dq.Rotate(2)
	// [2, 3, 0, 1]

	val, err := dq.Get(0)
	if err != nil {
		t.Fatalf("Get(0) error: %s", err)
	}
	if val != 2 {
		t.Fatalf("Get(0) -> %d, expected 2", val)
	}
	val, err = dq.Get(3)
	if err != nil {
		t.Fatalf("Get(3) error: %s", err)
	}
	if val != 1 {
		t.Fatalf("Get(3) -> %d, expected 1", val)
	}

	// [0, 1, 2, 3, 4]
	dq = newFilled(5)
	dq.Rotate(-2)
	// [2, 3, 4, 0, 1]
	val, err = dq.Get(0)
	if err != nil {
		t.Fatalf("Get(0) error: %s", err)
	}
	if val != 2 {
		t.Fatalf("Get(0) -> %d, expected 2", val)
	}
	val, err = dq.Get(4)
	if err != nil {
		t.Fatalf("Get(4) error: %s", err)
	}
	if val != 1 {
		t.Fatalf("Get(3) -> %d, expected 1", val)
	}
}

func TestString(t *testing.T) {
	dq := newFilled(3)
	s := dq.String()
	if s != "Deque{0, 1, 2}" {
		t.Fatalf("bad string: %s", s)
	}

	dq = newFilled(30)
	s = dq.String()
	if s != "Deque{0, 1, 2, 3, 4, 5, 6, 7, 8, ...}" {
		t.Fatalf("bad string: %s", s)
	}
}

func TestAppendAndGet(t *testing.T) {
	n := 128
	q := New[int]()

	for i := 0; i < n; i++ {
		q.Append(i + 1)
	}

	for i := 0; i < n; i++ {
		val, err := q.Get(i)
		if err != nil {
			t.Fatal(err)
		}

		if val != i+1 {
			t.Fatalf("get - %d: %v != %v", i, i+1, val)
		}
	}
}

// TODO: Add checks, currently we only try to crash
func FuzzDeque(f *testing.F) {
	f.Fuzz(func(t *testing.T, actions []byte) {
		d := New[int]()
		for i, action := range actions {
			switch action % 6 {
			case 0:
				d.Append(i)
			case 1:
				d.AppendLeft(i)
			case 2:
				d.Pop()
			case 3:
				d.PopLeft()
			case 4:
				d.Set(i, i)
			case 5:
				d.Rotate(i)
			}

			if n := d.Len(); n < 0 {
				t.Fatal(n)
			}
		}
	})
}
