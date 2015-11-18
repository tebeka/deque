package deque

import "testing"

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
