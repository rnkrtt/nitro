// Copyright 2021-2022, Offchain Labs, Inc.
// For license information, see https://github.com/OffchainLabs/nitro/blob/master/LICENSE.md

package containers

import (
	"fmt"
	"math"
	"testing"

	"github.com/offchainlabs/nitro/util/testhelpers"
)

// This function calculates the number of elements needed to pop from the queue
// to cause a shrink, given a certain shrink ratio.
//
// A shrink occurs when the length * ratio <= capacity.
// Let length = l, ratio = r, and capacity = c.
//
// Then we need to find n such that
// (c-n) / (l-n) > r,
// because both the length and the capacity decrease by one when we assign
// `q.slice = q.slice[1:]`.
//
// Rearranging terms to solve for n, we get
// n > (r * l - c) / (r - 1).
//
// Take the ceiling of this value to find the number of elements needed to pop
// to force a shrink.
func calcNumElementsToPop(capacity, length, ratio int) int {
	return int(math.Ceil(float64((ratio*length)-capacity) / float64(ratio-1)))
}

func TestQueue(t *testing.T) {
	q := Queue[int]{}

	// Need enough elements that we can force a shrink.
	initNumElements := 10000
	for i := 0; i < initNumElements; i++ {
		q.Push(i)
	}

	// Save the capacity to calculate how many elements we need to pop.
	bigCap := cap(q.slice)
	if bigCap < initNumElements {
		testhelpers.FailImpl(t, fmt.Sprintf("Unexpected capacity %d<%d: ", bigCap, initNumElements))
	}

	// Pop elements up to the one that should cause shrink.
	popCount := calcNumElementsToPop(bigCap, initNumElements, shrinkRatio)
	for i := 0; i < popCount-1; i++ {
		got := q.Pop()
		if got != i {
			testhelpers.FailImpl(t, fmt.Sprintf("Unexpected element popped: want %d, got %d", i, got))
		}
	}

	// The next pop should cause the shrink.
	got := q.Pop()
	if got != popCount-1 {
		testhelpers.FailImpl(t, fmt.Sprintf("Unexpected element popped: want %d, got %d", popCount, got))
	}

	// After shrink, the capacity should be exactly twice the length.
	expectedNewCap := len(q.slice) * 2
	if cap(q.slice) != expectedNewCap {
		testhelpers.FailImpl(t, fmt.Sprintf("Unexpected post-shrink cap: want %d, got %d", expectedNewCap, cap(q.slice)))
	}

	// Pop the remaining elements.
	for i := popCount; i < initNumElements; i++ {
		got := q.Pop()
		if got != i {
			testhelpers.FailImpl(t, fmt.Sprintf("Unexpected element popped: want %d, got %d", i, got))
		}
	}

	// Assert that queue is empty.
	if len(q.slice) != 0 {
		testhelpers.FailImpl(t, fmt.Sprintf("Non-empty queue: len=%d", len(q.slice)))
	}

	// Pop on empty queue should return the default value of an int, which is 0.
	if got := q.Pop(); got != 0 {
		testhelpers.FailImpl(t, fmt.Sprintf("Unexpected element popped: want %d, got %d", 0, got))
	}
}
