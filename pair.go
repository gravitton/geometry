package geom

import (
	"encoding/json"
	"fmt"
	"iter"
)

// Pair holds a result of at most two items, Count of them in Items from the first, the rest left at
// the zero value so equal pairs compare equal with == where the items do. The items are a Point or
// a Segment.
//
// The items are read by ranging All, or as Items[:Count] of a pair held in a variable.
type Pair[E interface {
	Equal(E) bool
	String() string
}] struct {
	Items [2]E
	Count int
}

// All returns an iterator over the items in order.
func (p Pair[E]) All() iter.Seq[E] {
	return func(yield func(E) bool) {
		for i := range p.Count {
			if !yield(p.Items[i]) {
				return
			}
		}
	}
}

// Equal reports whether the pairs hold as many items, each equal by its own Equal to the one in its
// place. The order counts.
func (p Pair[E]) Equal(pair Pair[E]) bool {
	if p.Count != pair.Count {
		return false
	}

	for i := range p.Count {
		if !p.Items[i].Equal(pair.Items[i]) {
			return false
		}
	}

	return true
}

// IsZero reports whether the pair is the zero value, one that holds no item, so the omitzero option
// of encoding/json leaves it out.
func (p Pair[E]) IsZero() bool {
	return p.Count == 0
}

// IsEmpty reports whether the pair holds no item.
func (p Pair[E]) IsEmpty() bool {
	return p.Count == 0
}

// AppendTo appends the items to dst in order and returns the extended slice.
func (p Pair[E]) AppendTo(dst []E) []E {
	return append(dst, p.Items[:p.Count]...)
}

// String returns the pair as Pair(), Pair(item) or Pair(item;item), each item by its own String,
// which allocates beside the result.
func (p Pair[E]) String() string {
	buf := make([]byte, 0, stringSize)
	buf = append(buf, "Pair("...)
	for i := range p.Count {
		if i > 0 {
			buf = append(buf, ';')
		}

		buf = append(buf, p.Items[i].String()...)
	}

	buf = append(buf, ')')

	return string(buf)
}

// MarshalJSON implements json.Marshaler as the array of the items, empty for a pair of none.
func (p Pair[E]) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Items[:p.Count])
}

// UnmarshalJSON implements json.Unmarshaler from an array of at most two items, and an error that
// leaves the pair unchanged for a longer one.
func (p *Pair[E]) UnmarshalJSON(bytes []byte) error {
	var items []E
	if err := json.Unmarshal(bytes, &items); err != nil {
		return err
	}

	if len(items) > len(p.Items) {
		return fmt.Errorf("geom: a pair of %d items, at most %d", len(items), len(p.Items))
	}

	pair := Pair[E]{Count: len(items)}
	copy(pair.Items[:], items)
	*p = pair

	return nil
}
