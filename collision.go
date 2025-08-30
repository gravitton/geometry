package geom

// CollisionRectangles checks if the given rectangles collide.
func CollisionRectangles[T Number](r1 Rectangle[T], r2 Rectangle[T]) bool {
	// TODO: support rotation

	min1, max1 := r1.Min(), r1.Max()
	min2, max2 := r2.Min(), r2.Max()

	return min1.X <= max2.X && min2.X <= max1.X && min1.Y <= max2.Y && min2.Y <= max1.Y
}

// CollisionRectangleCircle checks if the given rectangle and circle collide.
func CollisionRectangleCircle[T Number](r Rectangle[T], c Circle[T]) bool {
	// TODO: support rotation

	distance := c.Center.Subtract(r.Center).Abs()
	extends := r.Size.Scale(0.5).ToVector()

	// circle center is more than size outside rectangle borders
	if distance.X > extends.X+c.Radius || distance.Y > extends.Y+c.Radius {
		return false
	}

	// circle center is in rectangle
	if distance.X <= extends.X || distance.Y <= extends.Y {
		return true
	}

	// circle center is less than its size outside nearest border
	return distance.Subtract(extends).Less(c.Radius)
}

// CollisionCircles checks if the given circles collide.
func CollisionCircles[T Number](c1 Circle[T], c2 Circle[T]) bool {
	distance := c1.Center.Subtract(c2.Center)
	threshold := c1.Radius + c2.Radius

	return distance.Less(threshold)
}
