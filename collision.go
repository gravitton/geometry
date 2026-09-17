package geom

// CollisionRectangles checks if the given rectangles collide. Touching rectangles collide,
// within Epsilon of T, the same closed convention as Rectangle.Contains.
func CollisionRectangles[T Number](rect1 Rectangle[T], rect2 Rectangle[T]) bool {
	min1, max1 := rect1.Min(), rect1.Max()
	min2, max2 := rect2.Min(), rect2.Max()

	return LessOrEqual(min1.X, max2.X) && LessOrEqual(min2.X, max1.X) &&
		LessOrEqual(min1.Y, max2.Y) && LessOrEqual(min2.Y, max1.Y)
}

// CollisionCircles checks if the given circles collide. Touching circles collide,
// within Epsilon of T, the same closed convention as CollisionRectangles. The radii are
// summed in float64, so a narrow integer T cannot overflow the threshold.
func CollisionCircles[T Number](circle1 Circle[T], circle2 Circle[T]) bool {
	distance := circle1.Center.Subtract(circle2.Center).Length()
	threshold := float64(circle1.Radius) + float64(circle2.Radius)

	// LessOrEqual would take the threshold back into T, or apply the float64 tolerance to it;
	// the delta form compares two float64 values with the tolerance of T.
	return LessOrEqualDelta(distance, threshold, Epsilon[T]())
}

// CollisionRectangleCircle checks if the given rectangle and circle collide: the point of the
// rectangle closest to the circle center lies within the radius. Touching shapes collide,
// within Epsilon of T, and the rectangle bounds are the same Min and Max that Contains and
// CollisionRectangles use.
func CollisionRectangleCircle[T Number](rect Rectangle[T], circle Circle[T]) bool {
	closest := rect.Clamp(circle.Center)

	return closest.Subtract(circle.Center).LessOrEqual(circle.Radius)
}
