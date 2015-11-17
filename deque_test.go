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

func TestAppend(t *testing.T) {
	q := New()
	for i := 0; i < nItems; i++ {
		q.Append(i)
	}
	if q.Len() != nItems {
		t.Fatalf("Size after %d Append is %d", nItems, q.Len())
	}
}

func TestAppendLeft(t *testing.T) {
	q := New()
	for i := 0; i < nItems; i++ {
		q.AppendLeft(i)
	}
	if q.Len() != nItems {
		t.Fatalf("Size after %d AppendLeft is %d", nItems, q.Len())
	}
}

func TestPop(t *testing.T) {
	dq := New()
	for i := 0; i < nItems; i++ {
		dq.Append(i)
	}

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
	dq := New()
	for i := 0; i < nItems; i++ {
		dq.Append(i)
	}

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
