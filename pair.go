package geom

import "iter"

// Pair holds a result of at most two items, the points where two boundaries cross or the
// tangents of a circle, as a value that needs no allocation: Count of them, in Items from the
// first, the rest left at the zero value so two pairs of equal items compare equal. An item is
// held once, so two results that are the same point give a Count of one.
//
// The items are read by ranging All, or as the slice Items[:Count] of a pair held in a
// variable.
type Pair[E any] struct {
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
