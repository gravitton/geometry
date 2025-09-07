// Package geom provides a small, generic, and immutable 2D geometry toolkit.
//
// Design goals
//   - Generic: All core types are parameterized by a Number constraint so the
//     same API works with ints and floats.
//   - Immutable: Methods do not mutate receivers; they return new values.
//   - Practical: Focused on game/graphics use-cases with clear, minimal API.
package geom
