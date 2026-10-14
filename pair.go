package geom

import (
	"encoding/json"
	"fmt"
	"iter"
)

// Pair holds a result of at most two items, the points where two boundaries cross or the
// tangents of a circle, as a value that needs no allocation: Count of them, in Items from the
// first, the rest left at the zero value so two pairs of the same items compare equal with ==,
// where the items compare at all. An item is held once, so two results that are the same point
// give a Count of one. The items are a geometry of this package, a Point or a Segment, each
// with the Equal the pair compares them by and the String it prints them with.
//
// The items are read by ranging All, or as the slice Items[:Count] of a pair held in a
// variable.
type Pair[E interface {
	Equal(E) bool
	String() string
}] struct {
	Items [2]E
	Count int
}

// All returns an iterator over the items the pair holds, in order. Ranging it allocates
// nothing.
func (p Pair[E]) All() iter.Seq[E] {
	return func(yield func(E) bool) {
		for i := range p.Count {
			if !yield(p.Items[i]) {
				return
			}
		}
	}
}

// Equal reports whether the two pairs hold as many items and each is equal to the one in its
// place by its own Equal, within the tolerance that applies. The order counts: the same two
// points from the other end are another pair.
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

// IsZero reports whether the pair is the zero value, one that holds no item, the check the
// omitzero option of encoding/json makes to leave a field out.
func (p Pair[E]) IsZero() bool {
	return p.Count == 0
}

// IsEmpty reports whether the pair holds no item.
func (p Pair[E]) IsEmpty() bool {
	return p.Count == 0
}

// AppendTo appends the items the pair holds to dst, in order, and returns the extended slice,
// so a caller gathering the results of many calls into one slice allocates nothing once it has
// room. It is the slice Items[:Count] for a result that is in no variable.
func (p Pair[E]) AppendTo(dst []E) []E {
	return append(dst, p.Items[:p.Count]...)
}

// String returns the pair as the items it holds, each in its own form: Pair(), Pair(item) or
// Pair(item;item). The items are printed by their own String, so it allocates one string for
// each beside the one it returns.
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

// MarshalJSON implements json.Marshaler: the array of the items the pair holds, in order, an
// empty array for a pair of none.
func (p Pair[E]) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Items[:p.Count])
}

// UnmarshalJSON implements json.Unmarshaler: an array of at most two items, and an error for a
// longer one, which leaves the pair as it was.
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

// mapPair returns the pair of the items converted one by one, in their places, an item held
// once where the two convert to the same one, as two points rounding onto one lattice point do.
func mapPair[E interface {
	Equal(E) bool
	String() string
}, R interface {
	Equal(R) bool
	String() string
}](pair Pair[E], convert func(E) R) Pair[R] {
	mapped := Pair[R]{Count: pair.Count}
	for i := range pair.Count {
		mapped.Items[i] = convert(pair.Items[i])
	}

	if mapped.Count == 2 && mapped.Items[0].Equal(mapped.Items[1]) {
		var none R

		mapped.Items[1] = none
		mapped.Count = 1
	}

	return mapped
}
