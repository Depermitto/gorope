// Package gorope provides easy and error-free implementation of the Rope data
// structure. It is an alternative to []bytes, providing faster operations of
// Rope.Insert, Rope.Delete and Rope.Concat.
package gorope

import (
	"fmt"
)

// Rope is an alternative to string. This data structure consists of smaller
// chunks of []bytes connected by nodes. It has efficient operations for Insert, Delete
// and Split.
type Rope struct {
	value  []byte
	left   *Rope
	right  *Rope
	weight int
}

// NewWith creates a new Rope out a byte array and a maximum number
// of characters in each node. This is the preffered way of creating
// Ropes, especially if reading from file. NewWith has complexity of O(n).
func NewWith(value []byte, chunkSize int) *Rope {
	if len(value) <= chunkSize {
		return &Rope{
			left:   nil,
			right:  nil,
			weight: len(value),
			value:  value,
		}
	}

	mid := len(value) / 2
	left := value[:mid]
	right := value[mid:]
	return &Rope{
		left:   NewWith(left, chunkSize),
		right:  NewWith(right, chunkSize),
		weight: len(left),
		value:  nil,
	}
}

// New creates a new Rope with automatic chunk sizing. It has the same
// complexity as NewWith.
func New(value []byte) *Rope {
	return NewWith(value, max(len(value)/20, 5))
}

// FromStringWith creates a new Rope out of a string value and a maximum number
// of characters in each node. It has the same complexity as NewWith.
func FromStringWith(value string, chunkSize int) *Rope {
	return NewWith([]byte(value), chunkSize)
}

// FromString creates a new Rope out of a string value with automatic chunk
// sizing. It has the same complexity as NewWith.
func FromString(value string) *Rope {
	return New([]byte(value))
}

// Concat joins multiple instances of Rope and returns the root of the combination.
// This method simply assigns new pointers of the same data to the root, so it
// does not modify the caller or the callee. Concat has complexity of O(1).
func (r *Rope) Concat(other ...*Rope) *Rope {
	if other == nil {
		return r
	}

	rope := r
	for _, other := range other {
		rope = &Rope{
			left:   rope.Copy(),
			right:  other.Copy(),
			weight: rope.Len(),
			value:  nil,
		}
	}
	return rope
}

// Split splits the caller Rope at the specified position in two parts. The caller
// is modified to be the left part and the right part is returned. Incorrect position
// results in undefined behaviour. Split has complexity of O(log n).
func (r *Rope) Split(pos int) *Rope {
	orphans := &Rope{}
	if pos >= r.weight && r.right != nil {
		orphans = orphans.Concat(r.right.Split(pos - r.weight))
	} else if r.left != nil {
		orphans = orphans.Concat(r.left.Split(pos), r.right)
		r.right = nil
	} else {
		// We split the node
		pos = min(pos, r.weight)
		left := r.value[:pos]
		right := r.value[pos:]

		// Update current r
		*r = Rope{
			left:   &Rope{value: left, weight: len(left)},
			weight: len(left),
		}

		// Return the right child
		return &Rope{value: right, weight: len(right)}
	}

	return orphans
}

// Insert inserts a []byte value at the specified position. It's an efficient
// operation consisting of two Concat operations and a Split. Insert modifies
// the underlying Rope. Error is non-nil if position is incorrect. Insert has
// complexity of O(log n).
func (r *Rope) Insert(pos int, value []byte) error {
	n := &Rope{value: value, weight: len(value)}
	if pos == 0 {
		n = n.Concat(r)
		*r = *n
	} else if pos <= r.Len() {
		other := r.Split(pos)
		n = r.Concat(n)
		other = n.Concat(other)
		*r = *other
	} else {
		return fmt.Errorf("incorrect split position")
	}
	return nil
}

// Delete removes `n` characters starting from the specified position. This
// operation is made from two Split operations and a Concat. Delete has
// complexity of O(log n).
func (r *Rope) Delete(pos int, n int) error {
	lhs := r.Split(pos)

	n = min(n, lhs.Len())
	rhs := lhs.Split(n)

	*r = *r.Concat(rhs)
	return nil
}

// At Returns a byte at the specified position. Error is non-nil if the
// position is incorrect. At has complexity of O(log n).
func (r *Rope) At(pos int) (byte, error) {
	switch {
	case pos < 0:
		return 0, fmt.Errorf("cannot index negative value %v", pos)
	case pos >= r.weight && r.right != nil:
		return r.right.At(pos - r.weight)
	case r.left != nil:
		return r.left.At(pos)
	case pos >= len(r.value):
		return 0, fmt.Errorf("index %v out of bounds %v", pos, len(r.value))
	default:
		return r.value[pos], nil
	}
}

// Copy is a semi-deep copy on the Rope. It is not a deep-copy like Clone. Copy
// returns a new pointer to the same underlying data.
func (r *Rope) Copy() *Rope {
	if r == nil {
		return nil
	}

	temp := *r
	return &temp
}

// CloneWith collects all leaf nodes and creates a new Rope with
// a structure based on chunkSize. CloneWith has complexity of O(n).
func (r *Rope) CloneWith(chunkSize int) *Rope {
	return FromStringWith(r.String(), chunkSize)
}

// Clone deep-copies the entire Rope, allocating new memory with
// automatic chunk sizes. Clone has the same complexity as CloneWith.
func (r *Rope) Clone() *Rope {
	return FromString(r.String())
}

// Len calculates and returns the length of the rope (sum of all the
// characters). Len has complexity of O(log n).
func (r *Rope) Len() int {
	sum := r.weight
	if r.right != nil {
		sum += r.right.Len()
	}
	return sum
}

// String collects the leaves and returns the []byte held by Rope.
// String has complexity of O(n).
func (r *Rope) String() string {
	value := string(r.value)
	if r.left != nil {
		value += r.left.String()
	}

	if r.right != nil {
		value += r.right.String()
	}
	return value
}
