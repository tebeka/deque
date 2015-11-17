/*
Deque provides a double ended queue modeled after Python's collections.deque
*/
package deque

import "fmt"

const (
	blockLen    = 64
	blockCenter = (blockLen - 1) / 2
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
}

// NewDeque returns a new Deque
func NewDeque() *Deque {
	block := newBlock(nil, nil)
	return &Deque{
		right:    block,
		left:     block,
		leftIdx:  blockCenter + 1,
		rightIdx: blockCenter,
		size:     0,
	}
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
	dq.size += 1
	dq.rightIdx += 1
	dq.right.data[dq.rightIdx] = item
}

// AppendLeft appends an item to the left of the deque
func (dq *Deque) AppendLeft(item interface{}) {
	if dq.leftIdx == 0 {
		block := newBlock(nil, dq.left)
		dq.left.left = block
		dq.left = block
		dq.leftIdx = blockLen
	}
	dq.size += 1
	dq.leftIdx -= 1
	dq.left.data[dq.leftIdx] = item
}

// Pop removes and return the rightmost element
func (dq *Deque) Pop() (interface{}, error) {
	if dq.Len() == 0 {
		return nil, fmt.Errorf("Pop from an empty Deque")
	}

	item := dq.right.data[dq.rightIdx]
	dq.rightIdx -= 1
	dq.size -= 1
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
	dq.leftIdx += 1
	dq.size -= 1

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
