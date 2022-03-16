package deque_test

import (
	"fmt"

	"github.com/tebeka/deque"
)

func Example() {
	// Create a new deque and put some numbers in it.
	dq := deque.New[int]()

	for i := 0; i < 5; i++ {
		dq.Append(i)
	}

	// Pop from the left
	val, _ := dq.PopLeft()
	fmt.Println(val) // 0

	// Get an item
	val, _ = dq.Get(2)
	fmt.Println(val) // 2

	// Set an item
	dq.Set(2, 9)

	// Rotate
	dq.Rotate(-2)

	// Print
	fmt.Println(dq)

	// Output:
	// 0
	// 3
	// Deque{9, 4, 1, 2}
}
