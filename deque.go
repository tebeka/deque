package deque

import (
	"bytes"
	"fmt"
)

const (
	blockLen    = 64
	blockCenter = (blockLen - 1) / 2
)

type block[T any] struct {
	left, right *block[T]
	data        []T
}

func newBlock[T any](left, right *block[T]) *block[T] {
	return &block[T]{
		left:  left,
		right: right,
		data:  make([]T, blockLen),
	}
}

// Deque is a double ended queue
type Deque[T any] struct {
	left, right       *block[T]
	leftIdx, rightIdx int /* in range(BLOCKLEN) */
	size              int
	maxSize           int
}

// New returns a new unbounded Deque
func New[T any]() *Deque[T] {
	block := newBlock[T](nil, nil)
	return &Deque[T]{
		right:    block,
		left:     block,
		leftIdx:  blockCenter + 1,
		rightIdx: blockCenter,
		size:     0,
		maxSize:  -1,
	}
}

// NewBounded returns a new bounded Deque
// A bounded Deque will not grow over maxSize items
func NewBounded[T any](maxSize int) (*Deque[T], error) {
	if maxSize <= 0 {
		return nil, fmt.Errorf("maxSize must be > 0 (got %d)", maxSize)
	}
	dq := New[T]()
	dq.maxSize = maxSize
	return dq, nil
}

// Len returns the number of items in the queue
func (dq *Deque[T]) Len() int {
	return dq.size
}

// Append appends an item to the right of the deque
func (dq *Deque[T]) Append(item T) {
	if dq.rightIdx == blockLen-1 {
		block := newBlock(dq.right, nil)
		dq.right.right = block
		dq.right = block
		dq.rightIdx = -1
	}
	dq.size++
	dq.rightIdx++
	dq.right.data[dq.rightIdx] = item
	if dq.maxSize != -1 && dq.Len() > dq.maxSize {
		dq.PopLeft() // nolint
	}
}

// AppendLeft appends an item to the left of the deque
func (dq *Deque[T]) AppendLeft(item T) {
	if dq.leftIdx == 0 {
		block := newBlock[T](nil, dq.left)
		dq.left.left = block
		dq.left = block
		dq.leftIdx = blockLen
	}
	dq.size++
	dq.leftIdx--
	dq.left.data[dq.leftIdx] = item
	if dq.maxSize != -1 && dq.Len() > dq.maxSize {
		dq.Pop() // nolint
	}
}

// Pop removes and return the rightmost element
func (dq *Deque[T]) Pop() (T, error) {
	if dq.Len() == 0 {
		var z T
		return z, fmt.Errorf("Pop from an empty Deque")
	}

	item := dq.right.data[dq.rightIdx]
	dq.rightIdx--
	dq.size--
	if dq.rightIdx == -1 {
		if dq.Len() == 0 {
			// re-center instead of freeing a block
			dq.leftIdx = blockCenter + 1
			dq.rightIdx = blockCenter
		} else {
			prev := dq.right.left
			prev.right = nil
			dq.right = prev
			dq.rightIdx = blockLen - 1
		}
	}
	return item, nil
}

// PopLeft removes and return the leftmost element.
func (dq *Deque[T]) PopLeft() (T, error) {
	if dq.Len() == 0 {
		var z T
		return z, fmt.Errorf("PopLeft from an empty Deque")
	}

	item := dq.left.data[dq.leftIdx]
	dq.leftIdx++
	dq.size--

	if dq.leftIdx == blockLen {
		if dq.Len() == 0 {
			// re-center instead of freeing a block
			dq.leftIdx = blockCenter + 1
			dq.rightIdx = blockCenter
		} else {
			prev := dq.left.right
			prev.left = nil
			dq.left = prev
			dq.leftIdx = 0
		}
	}
	return item, nil
}

func (dq *Deque[T]) locate(i int) (*block[T], int) {
	// first block
	firstSize := blockLen - dq.leftIdx
	if i < firstSize {
		return dq.left, dq.leftIdx + i
	}

	b := dq.left.right // 2nd block
	i -= firstSize

	for i >= blockLen {
		b = b.right
		i -= blockLen
	}
	return b, i
}

// Get return the item at position i
func (dq *Deque[T]) Get(i int) (T, error) {
	if i < 0 || i >= dq.Len() {
		var z T
		return z, fmt.Errorf("index %d out of range", i)
	}

	b, idx := dq.locate(i)

	return b.data[idx], nil
}

// Set sets the item at position i to val
func (dq *Deque[T]) Set(i int, val T) error {
	if i < 0 || i >= dq.Len() {
		return fmt.Errorf("index %d out of range", i)
	}

	b, idx := dq.locate(i)
	b.data[idx] = val
	return nil
}

// Rotate rotates the queue.
// If n is positive then rotate right n steps, otherwise rotate left -n steps
func (dq *Deque[T]) Rotate(n int) {
	if dq.Len() == 0 || n == 0 {
		return
	}

	var popfn func() (T, error)
	var appendfn func(T)

	if n > 0 {
		popfn = dq.Pop
		appendfn = dq.AppendLeft
	} else {
		popfn = dq.PopLeft
		appendfn = dq.Append
		n = -n
	}

	for i := 0; i < n; i++ {
		val, _ := popfn()
		appendfn(val)
	}
}

func (dq *Deque[T]) String() string {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "Deque{")
	n := dq.Len()
	chopped := false
	if n > 10 {
		n = 10
		chopped = true
	}
	for i := 0; i < n-1; i++ {
		val, _ := dq.Get(i)
		fmt.Fprintf(&buf, "%#v, ", val)
	}
	if chopped {
		fmt.Fprintf(&buf, "...")
	} else {
		val, _ := dq.Get(n - 1)
		fmt.Fprintf(&buf, "%#v", val)
	}
	fmt.Fprintf(&buf, "}")
	return buf.String()
}
