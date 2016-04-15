// Package deque provides a double ended queue modeled after Python's
// collections.deque
package deque

import (
	"bytes"
	"fmt"
)

const (
	blockLen    = 64
	blockCenter = (blockLen - 1) / 2

	// Version is the package version
	Version = "1.0.0"
)

type block struct {
	left, right *block
	data        []interface{}
}

func newBlock(left, right *block) *block {
	return &block{
		left:  left,
		right: right,
		data:  make([]interface{}, blockLen),
	}
}

// Deque is a double ended queue
type Deque struct {
	left, right       *block
	leftIdx, rightIdx int /* in range(BLOCKLEN) */
	size              int
	maxSize           int
}

// New returns a new unbounded Deque
func New() *Deque {
	block := newBlock(nil, nil)
	return &Deque{
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
func NewBounded(maxSize int) (*Deque, error) {
	if maxSize <= 0 {
		return nil, fmt.Errorf("maxSize must be > 0 (got %d)", maxSize)
	}
	dq := New()
	dq.maxSize = maxSize
	return dq, nil
}

// Len returns the number of items in the queue
func (dq *Deque) Len() int {
	return dq.size
}

// Append appends an item to the right of the deque
func (dq *Deque) Append(item interface{}) {
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
		dq.PopLeft()
	}
}

// AppendLeft appends an item to the left of the deque
func (dq *Deque) AppendLeft(item interface{}) {
	if dq.leftIdx == 0 {
		block := newBlock(nil, dq.left)
		dq.left.left = block
		dq.left = block
		dq.leftIdx = blockLen
	}
	dq.size++
	dq.leftIdx--
	dq.left.data[dq.leftIdx] = item
	if dq.maxSize != -1 && dq.Len() > dq.maxSize {
		dq.Pop()
	}
}

// Pop removes and return the rightmost element
func (dq *Deque) Pop() (interface{}, error) {
	if dq.Len() == 0 {
		return nil, fmt.Errorf("Pop from an empty Deque")
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
func (dq *Deque) PopLeft() (interface{}, error) {
	if dq.Len() == 0 {
		return nil, fmt.Errorf("PopLeft from an empty Deque")
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

func (dq *Deque) locate(i int) (b *block, idx int) {
	if i == 0 {
		i = dq.leftIdx
		b = dq.left
	} else if i == dq.Len()-1 {
		i = dq.rightIdx
		b = dq.right
	} else {
		index := i
		i += dq.leftIdx
		n := i / blockLen
		i %= blockLen
		if index < (dq.Len() >> 1) {
			b = dq.right
			for ; n > 0; n-- {
				b = b.right
			}
		} else {
			n = (dq.leftIdx+dq.size-1)/blockLen - n
			b = dq.right
			for ; n > 0; n-- {
				b = b.left
			}
		}
	}
	return b, i
}

// Get return the item at position i
func (dq *Deque) Get(i int) (interface{}, error) {
	if i < 0 || i >= dq.Len() {
		return nil, fmt.Errorf("index %d out of range", i)
	}

	b, idx := dq.locate(i)

	return b.data[idx], nil
}

// Set sets the item at position i to val
func (dq *Deque) Set(i int, val interface{}) error {
	if i < 0 || i >= dq.Len() {
		return fmt.Errorf("index %d out of range", i)
	}

	b, idx := dq.locate(i)
	b.data[idx] = val
	return nil
}

// Rotate rotates the queue.
// If n is positive then rotate right n steps, otherwise rotate left -n steps
func (dq *Deque) Rotate(n int) {
	if dq.Len() == 0 || n == 0 {
		return
	}

	var popfn func() (interface{}, error)
	var appendfn func(interface{})

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

func (dq *Deque) String() string {
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
